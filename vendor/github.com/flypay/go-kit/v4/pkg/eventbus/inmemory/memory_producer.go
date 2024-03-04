package inmemory

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"google.golang.org/protobuf/proto"
)

// Message encapsulates the event and headers that is sent through the event bus.
type Message struct {
	Event proto.Message
	// Done is called when the event has been handled (either with or without error)
	Done    func()
	Headers []eventbus.Header
}

func (msg *Message) incrementAttempts() {
	for i := range msg.Headers {
		if msg.Headers[i].Key == eventbus.HeaderAttempt {
			prev, _ := strconv.Atoi(msg.Headers[i].Value)
			msg.Headers[i].Value = strconv.Itoa(prev + 1)
			return
		}
	}

	msg.Headers = append(msg.Headers, eventbus.Header{
		Key:   eventbus.HeaderAttempt,
		Value: "1",
	})
}

type MemProducer struct {
	Messages   chan *Message
	middleware []eventbus.ProducerMiddleware
	wg         sync.WaitGroup
}

func NewMemProducer() *MemProducer {
	return &MemProducer{
		wg:         sync.WaitGroup{},
		middleware: []eventbus.ProducerMiddleware{},
		Messages:   NewMessageChannel(),
	}
}

func NewBufferedMemProducer(bufferSize int) *MemProducer {
	return &MemProducer{
		wg:         sync.WaitGroup{},
		middleware: []eventbus.ProducerMiddleware{},
		Messages:   make(chan *Message, bufferSize),
	}
}

func (p *MemProducer) Emit(ctx context.Context, message proto.Message, headers ...eventbus.Header) error {
	next := p.emit

	for _, middleware := range p.middleware {
		next = middleware(next)
	}

	log.Debugf("Emitting event %v", message)

	return next(ctx, message, headers...)
}

func (p *MemProducer) EmitAt(ctx context.Context, at time.Time, event proto.Message, headers ...eventbus.Header) error {
	p.wg.Add(1)
	time.AfterFunc(time.Until(at), func() {
		defer p.wg.Done()
		if err := p.Emit(ctx, event, headers...); err != nil {
			log.Errorf("Emitting event at %s: %+v", at, err)
		}
	})

	return nil
}

func (p *MemProducer) emit(ctx context.Context, message proto.Message, headers ...eventbus.Header) error {
	msg := &Message{
		Event:   proto.Clone(message),
		Headers: headers,
	}

	msg.Headers = eventbus.WithEmittedAt(time.Now(), headers...)

	// Do a background write to the unbuffered channel so we don't block
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.Messages <- msg
	}()

	return nil
}

func (p *MemProducer) Use(middleware eventbus.ProducerMiddleware) {
	p.middleware = append(p.middleware, middleware)
}

func (p *MemProducer) Close() error {
	p.wg.Wait()

	return nil
}
