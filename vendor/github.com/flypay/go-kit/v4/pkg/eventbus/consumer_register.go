package eventbus

import (
	"context"
	"strings"
	"sync"
	"time"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ConsumerRegister holds handler bindings for events
// and middlewares.
type ConsumerRegister struct {
	// TopicName returns a topic string for a given message
	// The returned topic name string is used for handler lookups when events are received.
	TopicName func(protoreflect.ProtoMessage) string

	// DefaultTimeout is the time after which the handler function returns an error.
	DefaultTimeout time.Duration

	middleware []ConsumerMiddleware

	handlers map[string]*HandlerBinding

	// handlersCache is a cache of handlers wrapped with the registered middlewares.
	// This cache exists so that we don't need to apply middleware every time a handler
	// is requested for an event.
	handlersCache map[string]*HandlerBinding

	// mu for locking when using the maps to protect against concurrent access
	mu sync.RWMutex

	// multiEventHandlers is used when registering multiple handlers
	// for a single event
	multiEventHandlers multiEventHandlers
}

// On implements Consumer.On
func (c *ConsumerRegister) On(event proto.Message, callback EventHandler, opt ...HandlerOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	topic := c.TopicName(event)
	if topic == "" {
		return ErrTopicNameNotDefined
	}

	// consume independently of case
	topic = strings.ToLower(topic)

	if c.handlers == nil {
		c.handlers = make(map[string]*HandlerBinding)
	}

	opts := NewDefaultHandlerOptions(c.DefaultTimeout)
	opts.Apply(opt...)

	if _, ok := c.handlers[topic]; !ok {
		c.handlers[topic] = &HandlerBinding{
			Event:   event,
			Handle:  callback,
			Options: opts,
			Pool: &sync.Pool{
				New: func() interface{} {
					return proto.Clone(event)
				},
			},
		}
		return nil
	}

	if c.multiEventHandlers == nil {
		c.multiEventHandlers = map[string]*multiEventHandler{}
	}

	// If a topic is registered more than one, create a multi handler binding
	// that will allow handling multiple registered handlers for the same
	// event.
	hb, err := c.multiEventHandlers.binding(
		topic,
		c.handlers[topic],
		callback,
		opts,
	)
	if err != nil {
		return err
	}
	c.handlers[topic] = hb
	return nil
}

// OnAll implements Consumer.OnAll
func (c *ConsumerRegister) OnAll(events map[proto.Message]EventHandler) error {
	for event, handler := range events {
		if err := c.On(event, handler); err != nil {
			return err
		}
	}

	return nil
}

// Use implements Consumer.Use
func (c *ConsumerRegister) Use(middleware ConsumerMiddleware) {
	c.middleware = append(c.middleware, middleware)

	// Clear handler cache
	c.handlersCache = nil
}

// HasSubscriptions returns true if any event handlers have been registered using `On` or `OnAll`
func (c *ConsumerRegister) HasSubscriptions() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.handlers) > 0
}

// Subscriptions returns a map of topic names for which the consumer register has handlers.
func (c *ConsumerRegister) Subscriptions() map[string]protoreflect.ProtoMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	topics := make(map[string]protoreflect.ProtoMessage, len(c.handlers))
	for t, h := range c.handlers {
		topics[t] = h.Event
	}
	return topics
}

// EventHandler returns the handler for an event
func (c *ConsumerRegister) EventHandler(event proto.Message) (HandlerBinding, bool) {
	topic := c.TopicName(event)
	if topic == "" {
		log.Warnf("no topic configured for event: %T", event)
		return HandlerBinding{}, false
	}

	return c.handlerWithMiddleware(topic)
}

// TopicHandler returns the handler for a topic name
func (c *ConsumerRegister) TopicHandler(topic string) (HandlerBinding, bool) {
	// consume independently of case
	return c.handlerWithMiddleware(strings.ToLower(topic))
}

func (c *ConsumerRegister) handlerWithMiddleware(topic string) (HandlerBinding, bool) {
	if !c.isRegistered(topic) {
		log.Warnf("no handler registered for event %q", topic)
		return HandlerBinding{}, false
	}

	// Is the handler cached
	if handler, ok := c.cachedHandler(topic); ok {
		return *handler, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Init cache
	if c.handlersCache == nil {
		c.handlersCache = make(map[string]*HandlerBinding)
	}

	// Populate cache
	h := c.wrapWithMiddleware(c.handlers[topic])
	c.handlersCache[topic] = h
	return *h, true
}

func (c *ConsumerRegister) isRegistered(topic string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.handlers[topic]
	return ok
}

func (c *ConsumerRegister) cachedHandler(topic string) (*HandlerBinding, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if handler, ok := c.handlersCache[topic]; ok {
		return handler, true
	}

	return nil, false
}

// wrapWithMiddleware applies all registered middleware to the handler binding
func (c *ConsumerRegister) wrapWithMiddleware(h *HandlerBinding) *HandlerBinding {
	for i := len(c.middleware) - 1; i >= 0; i-- {
		m := c.middleware[i]
		h.Handle = m(h.Handle)
	}
	return h
}

func (c *ConsumerRegister) Handle(handler HandlerBinding, msg proto.Message, headers ...Header) error {
	deadline := time.Now().Add(handler.Options.Timeout)
	ctx := flytcontext.WithEventDeadline(context.Background(), deadline)

	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	// Run the event handler synchronously to ensure we don't have "rogue" event handlers
	// running in the background.
	if err := handler.Handle(ctx, msg, headers...); err != nil {
		return err
	}

	// Return the context error if there was one. This is likely to happen when
	// the context deadline exceeds and the handler doesn't return an error.
	return ctx.Err()
}
