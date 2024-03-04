package middleware

import (
	"context"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/proto"
)

type RequestTracing struct{}

func NewRequestTracing() *RequestTracing {
	return &RequestTracing{}
}

// ConsumerMiddleware populates the context with the requestID from the event headers
func (t RequestTracing) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if reqID := t.getRequestID(headers); reqID != "" {
			ctx = flytcontext.WithRequestID(ctx, reqID)
		}

		if jobID := t.getJobID(headers); jobID != "" {
			ctx = flytcontext.WithJobID(ctx, jobID)
		}

		return next(ctx, event, headers...)
	}
}

// ProducerMiddleware populates the headers with the requestID in the context
func (t RequestTracing) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if reqID, ok := flytcontext.RequestID(ctx); ok {
			headers = append(headers, eventbus.Header{
				Key:   FlytRequestID,
				Value: reqID,
			})
		}

		if jobID, ok := flytcontext.JobID(ctx); ok {
			headers = append(headers, eventbus.Header{
				Key:   FlytJobID,
				Value: jobID,
			})
		}

		return next(ctx, event, headers...)
	}
}

func (t RequestTracing) getRequestID(headers []eventbus.Header) string {
	for _, header := range headers {
		if header.Key == FlytRequestID {
			return header.Value
		}
	}

	return ""
}

func (t RequestTracing) getJobID(headers []eventbus.Header) string {
	for _, header := range headers {
		if header.Key == FlytJobID {
			return header.Value
		}
	}

	return ""
}
