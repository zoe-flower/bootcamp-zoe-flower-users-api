package flythttp

import (
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
)

const defaultTimeout = 10 * time.Second

// NewClientWithXray retuns a flythttp client with Xray wrapper and default timeout
func NewClientWithXray() *http.Client {
	return NewClientWithXrayTimeout(defaultTimeout)
}

// NewClientWithXrayTimeout retuns a flythttp client with Xray wrapper and set timeout
func NewClientWithXrayTimeout(timeout time.Duration) *http.Client {
	return xray.Client(Client(
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout: timeout,
				}).Dial,
				TLSHandshakeTimeout:   timeout,
				ResponseHeaderTimeout: timeout,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	))
}
