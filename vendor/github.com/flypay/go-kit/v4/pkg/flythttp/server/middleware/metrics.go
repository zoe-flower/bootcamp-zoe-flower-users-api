package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/flypay/go-kit/v4/pkg/metrics"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() echo.MiddlewareFunc {
	requestsTotal := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
		}, []string{"code", "method", "path"})).(*prometheus.CounterVec)

	requestDuration := metrics.RegisterOrExisting(
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "request_duration_seconds",
		}, []string{"code", "method", "path"})).(*prometheus.HistogramVec)

	requestsInFlight := metrics.RegisterOrExisting(
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Subsystem: "http",
			Name:      "requests_in_flight",
		}, []string{"method", "path"})).(*prometheus.GaugeVec)

	requestPanics := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_panic",
		}, []string{"method", "path"})).(*prometheus.CounterVec)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		next = instrumentHandlerCounter(requestsTotal, next)
		next = instrumentHandlerDuration(requestDuration, next)
		next = instrumentHandlerInFlight(requestsInFlight, next)
		next = instrumentHandlerPanics(requestPanics, next)
		return next
	}
}

func instrumentHandlerCounter(requestsTotalMetric *prometheus.CounterVec, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		wrapper := newResponseWrapper(c.Response())
		c.SetResponse(echo.NewResponse(wrapper, c.Echo()))

		err := next(c)

		labels := prometheus.Labels{
			"method": c.Request().Method,
			"path":   c.Path(),
			"code":   wrapper.StatusCode(),
		}

		requestsTotalMetric.With(labels).Inc()

		return err
	}
}

func instrumentHandlerDuration(requestDurationMetric prometheus.ObserverVec, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		wrapper := newResponseWrapper(c.Response())
		c.SetResponse(echo.NewResponse(wrapper, c.Echo()))

		startTime := time.Now()

		err := next(c)

		labels := prometheus.Labels{
			"method": c.Request().Method,
			"path":   c.Path(),
			"code":   wrapper.StatusCode(),
		}

		requestDurationMetric.With(labels).Observe(time.Since(startTime).Seconds())

		return err
	}
}

func instrumentHandlerInFlight(requestInFlightMetric *prometheus.GaugeVec, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		labelsInFlight := prometheus.Labels{
			"method": c.Request().Method,
			"path":   c.Path(),
		}

		requestInFlightMetric.With(labelsInFlight).Inc()
		defer requestInFlightMetric.With(labelsInFlight).Dec()

		return next(c)
	}
}

func instrumentHandlerPanics(requestPanicsMetric *prometheus.CounterVec, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		var p IsPanic
		if err != nil && errors.As(err, &p) {
			labels := prometheus.Labels{
				"method": c.Request().Method,
				"path":   c.Path(),
			}
			requestPanicsMetric.With(labels).Inc()
		}
		return err
	}
}

type responseWrapper struct {
	delegate   http.ResponseWriter
	statusCode int
}

func newResponseWrapper(response *echo.Response) *responseWrapper {
	return &responseWrapper{
		delegate: response.Writer,
	}
}

// affirm at compilation time
var (
	_ http.Flusher        = (*responseWrapper)(nil)
	_ http.ResponseWriter = (*responseWrapper)(nil)
)

func (w *responseWrapper) Header() http.Header {
	return w.delegate.Header()
}

func (w *responseWrapper) Write(d []byte) (int, error) {
	return w.delegate.Write(d)
}

func (w *responseWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.delegate.WriteHeader(statusCode)
}

func (w *responseWrapper) StatusCode() string {
	return strconv.Itoa(w.statusCode)
}

func (w *responseWrapper) Flush() {
	if fl, ok := w.delegate.(http.Flusher); ok {
		fl.Flush()
	}
}
