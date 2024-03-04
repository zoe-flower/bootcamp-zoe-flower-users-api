package flythttp

import (
	"net/http"
)

const (
	FlytRequestID         string = "X-Flyt-Request-ID"
	FlytJobID             string = "X-Flyt-Job-ID"
	FlytRequestIDFallback string = "X-Request-ID"
	FlytExecutionNonce    string = "X-Execution-Nonce"
)

// ExecutionNonce extracts an execution nonce from a request
func ExecutionNonce(r *http.Request) string {
	return r.Header.Get(FlytExecutionNonce)
}

// RequestID extracts a unique request identifier from a request
func RequestID(r *http.Request) string {
	if rid := r.Header.Get(FlytRequestID); rid != "" {
		return rid
	}
	return r.Header.Get(FlytRequestIDFallback)
}

// RequestID extracts a unique request identifier from a request
func JobID(r *http.Request) string {
	if rid := r.Header.Get(FlytJobID); rid != "" {
		return rid
	}
	return ""
}
