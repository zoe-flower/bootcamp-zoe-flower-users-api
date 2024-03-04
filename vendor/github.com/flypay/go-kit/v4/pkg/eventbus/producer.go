package eventbus

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// EmitterFunc is a function responsible for sending an event.
type EmitterFunc func(ctx context.Context, event proto.Message, headers ...Header) error

// Emitter is the interface for raising a new event onto the eventbus
//
//counterfeiter:generate . Emitter
type Emitter interface {
	Emit(ctx context.Context, event proto.Message, headers ...Header) error
}

// DelayedEmitter is the interface for raising a new event onto the eventbus at a specific time
//
//counterfeiter:generate . DelayedEmitter
type DelayedEmitter interface {
	EmitAt(ctx context.Context, at time.Time, event proto.Message, headers ...Header) error
}

// The ProducerMiddleware allows a producer to be wrapped to add extra functionality.
type ProducerMiddleware func(next EmitterFunc) EmitterFunc

// Producer encapsulates the interface for emitting events in a service
//
//counterfeiter:generate . Producer
type Producer interface {
	Emitter
	Use(middleware ProducerMiddleware)
	Close() error
}
