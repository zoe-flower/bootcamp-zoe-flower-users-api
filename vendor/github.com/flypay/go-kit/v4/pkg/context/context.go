package context

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type flytContextKey string

func (k flytContextKey) String() string {
	return fmt.Sprintf("Flyt context key %s", string(k))
}

const (
	flytRequestID      flytContextKey = "X-Flyt-Request-ID"
	flytJobID          flytContextKey = "X-Flyt-Job-ID"
	flytExecutionID    flytContextKey = "X-Execution-ID"
	flytExecutionNonce flytContextKey = "X-Execution-Nonce"
	flytLocationID     flytContextKey = "X-Flyt-Location-ID"
	flytEventDeadline  flytContextKey = "X-Event-Deadline"
	skipRestaurantID   flytContextKey = "X-Skip-Restaurant-ID"
)

// ExecutionNonce extracts an execution nonce from a context
func ExecutionNonce(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(flytExecutionNonce).(string)
	return rid, ok
}

// EncodeNonce returns an encoded nonce for capability test requests
func EncodeNonce(serviceName, commitHash, uuid string) string {
	return fmt.Sprintf("%s:%s:%s", serviceName, commitHash, uuid)
}

// DecodeNonce returns a slice containing all elements of a nonce: serviceName, commitHash, uuid
func DecodeNonce(nonce string) (serviceName, commitHash, uuid string, err error) {
	n := strings.Split(nonce, ":")
	if len(n) < 3 {
		return "", "", "", fmt.Errorf("Error decoding nonce. Format should be serviceName:commitHash:uuid, got %s", nonce)
	}
	return n[0], n[1], n[2], nil
}

// RequestID extracts a unique request identifier from a context
func RequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(flytRequestID).(string)
	return rid, ok
}

// JobID extracts a unique job identifier from a context
func JobID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(flytJobID).(string)
	return rid, ok
}

// WithRequestID returns a new context with the given requestID
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, flytRequestID, requestID)
}

// WithJobID returns a new context with the given jobID
func WithJobID(ctx context.Context, jobID string) context.Context {
	return context.WithValue(ctx, flytJobID, jobID)
}

// ExecutionID extracts a unique request identifier for this execution from a context
func ExecutionID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(flytExecutionID).(string)
	return rid, ok
}

// WithExecutionID returns a new context with the given executionID
func WithExecutionID(ctx context.Context, executionID string) context.Context {
	return context.WithValue(ctx, flytExecutionID, executionID)
}

// WithExecutionNonce returns a new context with the given flyt execution nonce
func WithExecutionNonce(ctx context.Context, executionNonce string) context.Context {
	return context.WithValue(ctx, flytExecutionNonce, executionNonce)
}

// FlytLocationID extracts the Flyt Location ID from a context
func FlytLocationID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(flytLocationID).(string)
	return id, ok
}

// WithFlytLocationID returns a new context with the given Flyt Location ID
func WithFlytLocationID(ctx context.Context, locationID string) context.Context {
	return context.WithValue(ctx, flytLocationID, locationID)
}

// EventDeadline extracts the event deadline time from a context
func EventDeadline(ctx context.Context) (time.Time, bool) {
	deadline, ok := ctx.Value(flytEventDeadline).(time.Time)
	return deadline, ok
}

// WithEventDeadline return a new context with the event deadline time
func WithEventDeadline(ctx context.Context, eventDeadline time.Time) context.Context {
	return context.WithValue(ctx, flytEventDeadline, eventDeadline)
}

// SkipRestaurantID extracts the Flyt Location ID from a context
func SkipRestaurantID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(skipRestaurantID).(string)
	return id, ok
}

// WithSkipRestaurantID returns a new context with the given Flyt Location ID
func WithSkipRestaurantID(ctx context.Context, skipID string) context.Context {
	return context.WithValue(ctx, skipRestaurantID, skipID)
}
