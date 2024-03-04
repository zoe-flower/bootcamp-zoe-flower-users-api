package flythttp

import (
	"context"
	"net/http"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/tracing"
	"github.com/google/uuid"
)

// WithRequestIDFromRequest populates the context with the requestID discovered in the request
func WithRequestIDFromRequest(ctx context.Context, r *http.Request) context.Context {
	if rid := RequestID(r); rid != "" {
		return flytcontext.WithRequestID(ctx, rid)
	}
	return ctx
}

// WithJobIDFromRequest populates the context with the jobID discovered in the request
func WithJobIDFromRequest(ctx context.Context, r *http.Request) context.Context {
	if rid := JobID(r); rid != "" {
		return flytcontext.WithJobID(ctx, rid)
	}
	return ctx
}

// WithExecutionIDFromRequest populates the context with the executionID discovered in the request
func WithExecutionIDFromRequest(ctx context.Context, r *http.Request) context.Context {
	executionUUID := uuid.NewString()
	return flytcontext.WithExecutionID(ctx, executionUUID)
}

// AddRequestIDToTracing adds the requestID discovered in the request to the xray segment
func AddRequestIDToTracing(ctx context.Context, r *http.Request) {
	if rid := RequestID(r); rid != "" {
		tracing.AddAnnotation(ctx, tracing.RequestIDKey, rid)
	}
}

// WithExecutionNonceFromRequest populates the context with the execution nonce discovered in the request
func WithExecutionNonceFromRequest(ctx context.Context, r *http.Request) context.Context {
	if nonce := ExecutionNonce(r); nonce != "" {
		return flytcontext.WithExecutionNonce(ctx, nonce)
	}
	return ctx
}
