package eventbus

import (
	"context"
	"sync"

	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/safe"
	"google.golang.org/protobuf/proto"
)

const maxHandlersPerEvent = 2

// multiEventHandlers is a cache of the topic -> multiEventHandler
type multiEventHandlers map[string]*multiEventHandler

// binding will return a HandlerBinding with the handler for a multEventHandler
// to allow for calling multiple event handlers for the same event.
// For options, the longest timeout and largest max attempts will be used if multiple
// are set for event handlers.
func (h multiEventHandlers) binding(topic string, hb *HandlerBinding, callback EventHandler, options HandlerOptions) (*HandlerBinding, error) {
	if mh, ok := h[topic]; ok {
		if len(mh.handlers) == maxHandlersPerEvent {
			return nil, ErrTooManyHandlersForEvent
		}
		hb.Options = optsMerge(hb.Options, options)

		// Add the callback to the event handlers for the multiEventHandler
		mh.handlers = append(mh.handlers, callback)
		return hb, nil
	}

	// Create a new multiEventHandler using current HandlerBinding.Handle and
	// the callback as the initial handlers.
	mh := &multiEventHandler{
		HandlerBinding: *hb,
		handlers:       []EventHandler{hb.Handle, callback},
	}

	// Override the options
	hb.Options = optsMerge(hb.Options, options)

	// Override the Handle method with the one for the multiEventHandler
	hb.Handle = mh.Handle

	h[topic] = mh
	return hb, nil
}

// multiEventHandler stores handlers for a single event so we can concurrently
// run multiple handlers for an event
type multiEventHandler struct {
	HandlerBinding
	handlers []EventHandler
}

// Handle will concurrently call all the handlers registered for an event.
// Only the first error will be returned, the rest will be logged.
// NOTE: event handler timeouts are handled by the caller
func (m *multiEventHandler) Handle(ctx context.Context, event interface{}, headers ...Header) error {
	msg := event.(proto.Message)

	errs := make(chan error, len(m.handlers))
	wg := sync.WaitGroup{}

	for _, handler := range m.handlers {
		wg.Add(1)
		// Copy the handler into the scope for the goroutine
		handler := handler
		safe.Go(func() {
			defer wg.Done()
			if err := handler(ctx, proto.Clone(msg), headers...); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)

	// Always return the first error and log any others.
	var firstErr error
	for err := range errs {
		if firstErr == nil {
			firstErr = err
			continue
		}
		log.WithContext(ctx).Errorf("error handling %T event: %v", event, err)
	}
	return firstErr
}

func optsMerge(opts1, opts2 HandlerOptions) HandlerOptions {
	// Pick the largest timeout
	if opts1.Timeout < opts2.Timeout {
		opts1.Timeout = opts2.Timeout
	}
	return opts1
}
