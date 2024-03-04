package middleware

import (
	"context"
	"time"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/proto"
)

const headerEventDeadline = "X-Event-Deadline"

func NewDeadline() *DeadlineMiddleware {
	return &DeadlineMiddleware{}
}

type DeadlineMiddleware struct{}

func (mw *DeadlineMiddleware) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if deadline, ok := eventbus.HeaderValue(headers, headerEventDeadline); ok {
			if t, err := time.Parse(time.RFC3339Nano, deadline); err == nil {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, t)
				defer cancel()
			}
		}
		return next(ctx, event, headers...)
	}
}

func (mw *DeadlineMiddleware) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if deadline, ok := ctx.Deadline(); ok {
			headers = append(headers, eventbus.Header{
				Key:   headerEventDeadline,
				Value: deadline.Format(time.RFC3339Nano),
			})
		}
		return next(ctx, event, headers...)
	}
}
