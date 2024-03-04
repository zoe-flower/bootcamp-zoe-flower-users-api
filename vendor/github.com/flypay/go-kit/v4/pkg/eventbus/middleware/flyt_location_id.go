package middleware

import (
	"context"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/proto"
)

type FlytLocationIDMigration struct{}

func NewFlytLocationIDMigration() *FlytLocationIDMigration {
	return &FlytLocationIDMigration{}
}

// ConsumerMiddleware populates the context with the Flyt Location ID from the event headers
func (t FlytLocationIDMigration) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if value := t.getFlytLocationID(headers); value != "" {
			ctx = flytcontext.WithFlytLocationID(ctx, value)
		}
		return next(ctx, event, headers...)
	}
}

// ProducerMiddleware populates the headers with the Flyt Location ID in the context
func (t FlytLocationIDMigration) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if value, ok := flytcontext.FlytLocationID(ctx); ok {
			headers = append(headers, eventbus.Header{
				Key:   FlytLocationID,
				Value: value,
			})
		}
		return next(ctx, event, headers...)
	}
}

func (t FlytLocationIDMigration) getFlytLocationID(headers []eventbus.Header) string {
	for _, header := range headers {
		if header.Key == FlytLocationID {
			return header.Value
		}
	}

	return ""
}
