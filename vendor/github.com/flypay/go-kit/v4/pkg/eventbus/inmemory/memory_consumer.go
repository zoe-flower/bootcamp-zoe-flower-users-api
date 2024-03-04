package inmemory

import (
	"context"
	"strconv"
	"time"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/internal/duration"
	"github.com/flypay/go-kit/v4/pkg/log"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MemConsumer struct {
	eventbus.ConsumerRegister

	Messages chan *Message
	// DeadMessages is buffered to allow non-blocking writes. Should a lot of messages be issued
	// to the dead channel then it may start blocking and goroutines will start to become blocked until
	// you start pulling some messages from the channel.
	DeadMessages chan *Message

	stop context.CancelFunc
}

const (
	defaultEventHandlerTimeout = 5 * time.Second
)

type ConsumerConfig struct {
	// TopicName is used to lookup a topic name for an event
	TopicName func(protoreflect.ProtoMessage) string

	// DefaultEventHandlerTimeout it's the maximum duration for event handling, before a timeout error will be returned,
	// passed as a duration string. It can be overridden per single handle basis, passing `eventbus.WithTimeout` on
	// the handler registration
	DefaultEventHandlerTimeout string
}

func NewMemConsumer(cfg ConsumerConfig) *MemConsumer {
	handlerTimeout, err := duration.GreaterThanZero(cfg.DefaultEventHandlerTimeout)
	if err != nil {
		handlerTimeout = defaultEventHandlerTimeout
	}
	log.Infof("Default event timeout is set to %s", handlerTimeout)

	return &MemConsumer{
		ConsumerRegister: eventbus.ConsumerRegister{
			TopicName:      cfg.TopicName,
			DefaultTimeout: handlerTimeout,
		},

		Messages:     NewMessageChannel(),
		DeadMessages: deadletterChannel(),
	}
}

func (c *MemConsumer) Listen() {
	ctx, cancel := context.WithCancel(context.Background())

	c.stop = cancel

	go c.listen(ctx)
}

func done(msg *Message) {
	if msg.Done != nil {
		msg.Done()
	}
}

func (c *MemConsumer) listen(ctx context.Context) {
	for {
		select {
		case message := <-c.Messages:
			message.incrementAttempts()

			if handler, ok := c.EventHandler(message.Event); ok {
				err := c.Handle(handler, message.Event, message.Headers...)
				if err != nil {
					c.handleError(message, err, handler.Options)
					continue
				}

				done(message)
				continue
			}

			log.Warnf("did not handle event %T", message.Event)
			done(message)

		case <-ctx.Done():
			return
		}
	}
}

func (c *MemConsumer) handleError(message *Message, err error, opts eventbus.HandlerOptions) {
	if opts.MaxAttempts > 1 {
		attempt := 1
		if n, ok := eventbus.HeaderValue(message.Headers, eventbus.HeaderAttempt); ok {
			attempt, _ = strconv.Atoi(n)
		}

		if attempt < opts.MaxAttempts {
			log.Infof("Message attempt %d will be retried due to error: %+v", attempt, err)
			go func(msg *Message) {
				c.Messages <- msg
			}(message)
			return
		}

		log.Warnf("max event attempts %d reached with error: %+v", opts.MaxAttempts, err)
	}

	log.Errorf("Emitting event to deadletter: %+v", err)
	go func(msg *Message) {
		c.DeadMessages <- msg
		done(msg)
	}(message)
}

func (c *MemConsumer) Close() error {
	if c.stop != nil {
		c.stop()
	}

	return nil
}
