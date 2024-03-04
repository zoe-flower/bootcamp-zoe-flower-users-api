package middleware

import (
	"context"

	"github.com/google/uuid"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
)

type ExecutionTracing struct{}

func NewExecutionTracing() *ExecutionTracing {
	return &ExecutionTracing{}
}

// ConsumerMiddleware populates the context with the executionID from the event headers
func (t ExecutionTracing) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if _, ok := flytcontext.ExecutionID(ctx); !ok {
			executionUUID := uuid.NewString()
			ctx = flytcontext.WithExecutionID(ctx, executionUUID)
		}

		return next(ctx, event, headers...)
	}
}
