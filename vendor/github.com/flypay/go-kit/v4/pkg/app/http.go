package app

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	// #nosec Add pprof for debugging on the internal http server
	_ "net/http/pprof"

	"github.com/flypay/go-kit/v4/pkg/safe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultInternalHTTPPort = 8081

// StartInternalHTTP brings up an HTTP server to listening on the internal
// interface used by kubernetes and service monitors
func StartInternalHTTP() Closeable {
	http.Handle("/liveness", healthHandler())
	http.Handle("/readiness", readyHandler())
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: net.JoinHostPort("", strconv.Itoa(defaultInternalHTTPPort)),
	}

	safe.Go(func() {
		_ = server.ListenAndServe()
	})

	return server
}

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "healthy")
	})
}

func readyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ready")
	})
}
