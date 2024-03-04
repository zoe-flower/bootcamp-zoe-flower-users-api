package flythttp

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
)

// Client exists to wrap an existing *http.Client and send the X-Execution-Nonce
// header if the corresponding value exists in the request context.
func Client(c *http.Client) *http.Client {
	if c == nil {
		c = http.DefaultClient
	}

	transport := c.Transport

	// If the user provides a custom transport and it'll cast to an http.transport
	// Tune it to have connection pooling disabled to prevent DNS lookup problems
	// related issue here: https://github.com/golang/go/issues/23427
	if customTransport, ok := c.Transport.(*http.Transport); ok {
		customTransport.DisableKeepAlives = true
		customTransport.MaxIdleConnsPerHost = -1
	}

	// If there's no transport, set the default transport
	// Tune it to have connection pooling disabled to prevent DNS lookup problems
	// related issue here: https://github.com/golang/go/issues/23427
	if transport == nil {
		if t, ok := http.DefaultTransport.(*http.Transport); ok {
			tt := t.Clone()
			tt.DisableKeepAlives = true
			tt.MaxIdleConnsPerHost = -1
			transport = tt
		} else {
			// We are unable to set the transport properties for a custom default transport
			transport = http.DefaultTransport
		}
	}

	return &http.Client{
		Transport:     &roundTripper{base: transport},
		CheckRedirect: c.CheckRedirect,
		Jar:           c.Jar,
		Timeout:       c.Timeout,
	}
}

// NewRequestWithExecutionNonce is a convenience method to be able to construct a
// new request with an execution nonce.
func NewRequestWithExecutionNonce(ctx context.Context, nonce, method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(flytcontext.WithExecutionNonce(ctx, nonce))

	return request, err
}

// RoundTripper is a custom transport that pulls out the execution nonce from the
// context and passes it on via a header.
type roundTripper struct {
	base http.RoundTripper
}

// RoundTrip wraps a single HTTP transaction and add corresponding X-Execution-Nonce
func (rt *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if nonce, ok := flytcontext.ExecutionNonce(r.Context()); ok {
		r.Header.Add(FlytExecutionNonce, nonce)
	} else {
		// help other JET services identify the caller
		r.Header.Set("User-Agent", "JET-Connect")
	}

	return rt.base.RoundTrip(r)
}

// DoWithRetries sends an HTTP request and retries when a specified
// condition has not been satisfied.
//
// requestCall is a function that handles doing the HTTP request,
// which will be repeatedly executed until it's either successful
// or the retry limit has been reached.
//
// retryCondition is a function that evaluates the success of an
// HTTP request, in order to decide whether a retry should be attempted.
//
// retryEvery is the time interval between retry attempts.
//
// retryLimit is the maximum number of retries to be attempted when the
// retry condition is false before returning an error.
func DoWithRetries(requestCall func() (*http.Response, error), retryCondition func(request *http.Response) bool, retryEvery time.Duration, retryLimit int) (*http.Response, error) {
	response, err := requestCall()
	if err != nil {
		return response, err
	}

	if retryCondition(response) {
		retryLimit--

		if retryLimit < 0 {
			return response, errors.New("HTTP request failed. Maximum number of retries attempted")
		}

		time.Sleep(retryEvery)

		return DoWithRetries(requestCall, retryCondition, retryEvery, retryLimit)
	}

	return response, nil
}
