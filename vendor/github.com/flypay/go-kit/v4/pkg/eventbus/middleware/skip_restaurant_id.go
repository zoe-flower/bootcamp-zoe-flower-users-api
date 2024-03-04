package middleware

import (
	"context"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/proto"
)

type SkipRestaurantIDMigration struct{}

func NewSkipRestaurantIDMigration() *SkipRestaurantIDMigration {
	return &SkipRestaurantIDMigration{}
}

// ConsumerMiddleware populates the context with the Skip Restaurant ID from the event headers
func (t SkipRestaurantIDMigration) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if value := t.getSkipRestaurantID(headers); value != "" {
			ctx = flytcontext.WithSkipRestaurantID(ctx, value)
		}
		return next(ctx, event, headers...)
	}
}

// ProducerMiddleware populates the headers with the Skip Restaurant ID in the context
func (t SkipRestaurantIDMigration) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if value, ok := flytcontext.SkipRestaurantID(ctx); ok {
			headers = append(headers, eventbus.Header{
				Key:   SkipRestaurantID,
				Value: value,
			})
		}
		return next(ctx, event, headers...)
	}
}

func (t SkipRestaurantIDMigration) getSkipRestaurantID(headers []eventbus.Header) string {
	for _, header := range headers {
		if header.Key == SkipRestaurantID {
			return header.Value
		}
	}

	return ""
}
