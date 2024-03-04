package eventbus

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"
)

type (
	// EventHandler is a function responsible for processing an event.
	EventHandler func(ctx context.Context, event interface{}, headers ...Header) error

	// ConsumerMiddleware allows a consumer to be wrapped to add extra functionality.
	ConsumerMiddleware func(next EventHandler) EventHandler
)

// Consumer is the interface for consumption of messages from the eventbus.
//
//counterfeiter:generate . Consumer
type Consumer interface {
	// On binds an event handler to a given event type.
	// opt is a variadic list of options that allow you to specify custom configuration for the event handler (such as
	// a timeout)
	On(event proto.Message, callback EventHandler, opt ...HandlerOption) error
	// OnAll binds the handlers to the events to which they are keyed. If any of the bindings returns an error the
	// binding is halted and the error returned.
	OnAll(map[proto.Message]EventHandler) error
	Use(middleware ConsumerMiddleware)
	Listen()
	Close() error
}

// HandlerBinding encapsulates the event type and handler.
// The event is used to unmarshal the received raw bytes into a known
// protocol buffer message type.
// HandlerOptions can be associated with a event handler binding in this type.
type HandlerBinding struct {
	Event   proto.Message
	Handle  EventHandler
	Options HandlerOptions
	Pool    *sync.Pool
}

// HeaderValue finds a given header by key and returns the value and a bool if it was found
func HeaderValue(headers []Header, key string) (string, bool) {
	for _, h := range headers {
		if h.Key == key {
			return h.Value, true
		}
	}
	return "", false
}

type HandlerBindings map[string]HandlerBinding

func (hb HandlerBindings) ApplyMiddleware(middleware []ConsumerMiddleware) {
	for k, h := range hb {
		for _, middleware := range middleware {
			h.Handle = middleware(h.Handle)
		}
		hb[k] = h
	}
}

func (hb HandlerBindings) BindingExists(key string) bool {
	_, ok := hb[key]
	return ok
}
