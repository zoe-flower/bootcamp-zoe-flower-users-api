package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/getsentry/sentry-go"
)

type SentryHub interface {
	RecoverWithContext(ctx context.Context, err interface{}) *sentry.EventID
}

type SentryMiddleware struct {
	hub SentryHub
}

func NewSentry(hub SentryHub) *SentryMiddleware {
	return &SentryMiddleware{
		hub: hub,
	}
}

func (m *SentryMiddleware) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) (err error) {
		defer func() {
			if r := recover(); r != nil {
				m.hub.RecoverWithContext(ctx, r)
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

		return next(ctx, event, headers...)
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
