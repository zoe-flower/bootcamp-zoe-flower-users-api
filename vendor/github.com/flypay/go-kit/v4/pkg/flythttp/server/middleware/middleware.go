package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/flypay/go-kit/v4/pkg/flythttp"
	"github.com/flypay/go-kit/v4/pkg/tracing"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// XRay adds distributed tracing to the request context.
// The name param is used as the segment name which should be
// the service name in most scenarios.
func XRay(name string) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return xray.Handler(xray.NewFixedSegmentNamer(name), h)
	})
}

// AddRequestID adds the request ID into the request context and xray tracing if available
func AddRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			r := c.Request()
			ctx := flythttp.WithRequestIDFromRequest(r.Context(), r)
			flythttp.AddRequestIDToTracing(ctx, r)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

// AddExecutionID adds the execution ID into the request context
func AddExecutionID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			r := c.Request()
			ctx := flythttp.WithExecutionIDFromRequest(r.Context(), r)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

// AddJobID adds the job ID into the request context if available
func AddJobID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			r := c.Request()
			ctx := flythttp.WithJobIDFromRequest(r.Context(), r)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

// AddTracingAnnotations adds annotations to the tracing segment
func AddTracingAnnotations(team, flow string, tier int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			tracing.AddAnnotations(
				c.Request().Context(),
				tracing.Annotation{Key: tracing.TeamKey, Value: team},
				tracing.Annotation{Key: tracing.FlowKey, Value: flow},
				tracing.Annotation{Key: tracing.TierKey, Value: tier},
			)
			return next(c)
		}
	}
}

// ExecutionTracing adds the execution nonce into the request context if available
func ExecutionTracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			r := c.Request()
			ctx := flythttp.WithExecutionNonceFromRequest(r.Context(), r)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

// Recover will recover from a panic during the handling of a request and send the
// panic to Sentry and wrap and return the error as a panic
func Recover(hub sentryHub, panicOnRecover bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					if panicOnRecover {
						panic(r)
					}
					hub.RecoverWithContext(c.Request().Context(), r)
					switch e := r.(type) {
					case string:
						err = Panic{errors.New(e)}
					case error:
						err = Panic{e}
					default:
						err = Panic{fmt.Errorf("%+v", e)}
					}
				}
			}()
			return next(c)
		}
	}
}

type IsPanic interface {
	IsPanic()
}

type Panic struct {
	error
}

func (Panic) IsPanic() {}

var _ IsPanic = Panic{}

type sentryHub interface {
	RecoverWithContext(ctx context.Context, err interface{}) *sentry.EventID
}
