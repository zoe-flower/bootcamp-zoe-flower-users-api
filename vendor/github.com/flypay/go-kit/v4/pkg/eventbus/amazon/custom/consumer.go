package custom

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/env"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/amazon"
	"github.com/flypay/go-kit/v4/pkg/eventbus/internal/duration"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/queues"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type CustomHandler func(context.Context, amazon.SNSMessage) (bool, error)

//counterfeiter:generate . QueueConsumer
type QueueConsumer interface {
	eventbus.Consumer
	RegisterHandler(CustomHandler)
	WithReceiver(ReceiverMiddleware)
}

var QueueReceiveWaitTimeSeconds int64 = 20

type ConsumerConfig struct {
	// The AWS session to use with the SNS or SQS clients
	Session *session.Session

	ServiceName string

	// QueueURL is used to know from which queue to receive messages
	QueueURL                   string
	DefaultEventHandlerTimeout string
	VisibilityTimeout          string

	// Queue is added to metrics to differentiate topics in receive and pool metrics.
	Queue string

	// How many events should be handled concurrently.
	WorkerPoolSize int
}

type Consumer struct {
	sqs sqsiface.SQSAPI
	cfg ConsumerConfig

	producer  *queues.Producer
	consumers []queues.Stoppable

	handlerTimeout    time.Duration
	visibilityTimeout time.Duration

	muListening sync.Mutex
	listening   bool

	// TODO: hack to continue using eventbus consumer middleware until we can
	// make the middleware and metrics not require a proto.
	handlerMiddleware []eventbus.ConsumerMiddleware

	handler CustomHandler

	// receiverMiddlewares are wrapped around the queue receiver. Commonly used to add logging, metrics etc.
	receiverMiddlewares []ReceiverMiddleware
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	svc := sqs.New(cfg.Session)

	handlerTimeout, err := duration.GreaterThanZero(cfg.DefaultEventHandlerTimeout)
	if err != nil {
		return nil, err
	}

	visibilityTimeout, err := duration.GreaterThanZero(cfg.VisibilityTimeout)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		sqs:               svc,
		cfg:               cfg,
		handlerTimeout:    handlerTimeout,
		visibilityTimeout: visibilityTimeout,
	}, nil
}

func (c *Consumer) RegisterHandler(h CustomHandler) {
	c.handler = h
}

func (c *Consumer) Use(mw eventbus.ConsumerMiddleware) {
	c.handlerMiddleware = append(c.handlerMiddleware, mw)
}

// Listen binds workers in pools to the specified queue for the topic
func (c *Consumer) Listen() {
	c.muListening.Lock()
	defer c.muListening.Unlock()

	// Protect against calling the listen function twice
	if c.listening {
		return
	}

	var receiver queues.QueueReceiver = queues.NewQueue(c.cfg.QueueURL, c.sqs, QueueReceiveWaitTimeSeconds, c.visibilityTimeout)
	receiver = &queues.InstrumentedQueueReceiver{
		QueueReceiver: receiver,
		Rate:          newQueueReceiveCounter(c.cfg.ServiceName, c.cfg.Queue),
		Quantity:      newQueueJobsReceivedHistogram(c.cfg.ServiceName, c.cfg.Queue),
	}

	for _, mw := range c.receiverMiddlewares {
		receiver = mw(receiver)
	}

	c.producer = queues.NewProducer(
		receiver,
		c.cfg.WorkerPoolSize,
		newQueueCapacityGauge(c.cfg.ServiceName, c.cfg.Queue),
	)

	for i := 0; i < c.cfg.WorkerPoolSize; i++ {
		consumer := queues.NewConsumer(c.producer.Pool, c.handle)
		c.consumers = append(c.consumers, consumer)
	}

	c.listening = true
}

// The returned error is discarded by the queue consumer so we must handle it here
func (c *Consumer) handle(job queues.Job) (bool, error) {
	var msg amazon.SNSMessage
	if err := json.NewDecoder(job.Payload()).Decode(&msg); err != nil {
		return true, errors.Wrapf(err, "unable to decode sns message: %s", job.ID())
	}

	logger := log.WithContext(job.Context())
	logger.Debugf("handling SNS %s with message_id=%s", msg.Type, msg.MessageID)

	ok, err := c.handleMsg(msg)
	if err != nil {
		logger.Error(err)
		// No return here as sometimes we return ok=true with an error to ensure the message is not redelivered.
	}

	if ok {
		// Delete the message ONLY if the event was processed successfully do we remove it from the queue.
		// Otherwise it will reach the visibility timeout, retry or end up in the DLQ.
		job.Delete()
	}

	return ok, nil
}

func (c *Consumer) handleMsg(msg amazon.SNSMessage) (bool, error) {
	deadline := time.Now().Add(c.handlerTimeout)
	ctx := flytcontext.WithEventDeadline(context.Background(), deadline)
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	var seg *xray.Segment
	if traceID, ok := msg.MessageAttributes[xray.TraceIDHeaderKey]; ok {
		traceHeader := header.FromString(traceID.Value)
		ctx, seg = xray.NewSegmentFromHeader(ctx, env.GetAppName(), nil, traceHeader)
	}

	const flytRequestID = "X-Flyt-Request-ID"
	if rid, ok := msg.MessageAttributes[flytRequestID]; ok {
		ctx = flytcontext.WithRequestID(ctx, rid.Value)
	} else {
		// Fallback to using the message identifier
		ctx = flytcontext.WithRequestID(ctx, msg.MessageID)
	}

	const headerTimestamp = "timestamp"
	// If no header timestamp is populated in the sns message headers, add the "Timestamp" from
	// the SNS message.
	if _, ok := msg.MessageAttributes[headerTimestamp]; !ok {
		ts := msg.Timestamp.Format(time.RFC3339Nano)
		if msg.MessageAttributes == nil {
			msg.MessageAttributes = make(map[string]amazon.MessageAttribute)
		}
		msg.MessageAttributes[headerTimestamp] = amazon.MessageAttribute{Value: ts}
	}

	// Capture the 'ok' value since the middleware doesn't support this yet
	var ok bool

	// NOTE: This is a bit of a hack to continue using the current eventbus middleware.
	// We need to update the eventbus middleware to not require proto definitions to record
	// metrics etc.
	next := func(ctx context.Context, _ interface{}, _ ...eventbus.Header) error {
		var err error
		if ok, err = c.handler(ctx, msg); err == nil {
			err = ctx.Err()
		}
		return err
	}

	for _, mw := range c.handlerMiddleware {
		next = mw(next)
	}

	err := next(ctx, nil, headers(msg)...)
	if seg != nil {
		seg.Close(err)
	}
	return ok, err
}

// Close stops the context and shuts down all producers and consumers
func (c *Consumer) Close() error {
	log.Info("custom.Consumer stopping...")

	if c.producer != nil {
		queues.StopAll(c.producer)
	}
	queues.StopAll(c.consumers...)

	log.Info("custom.Consumer stopped")
	return nil
}

// Implement the rest of the eventbus.Consumer interface to continue
// working with closables and other consumer-sepcific code.

func (c *Consumer) On(event proto.Message, callback eventbus.EventHandler, opt ...eventbus.HandlerOption) error {
	return nil
}
func (c *Consumer) OnAll(map[proto.Message]eventbus.EventHandler) error { return nil }

func headers(msg amazon.SNSMessage) []eventbus.Header {
	var headers []eventbus.Header

	for key, val := range msg.MessageAttributes {
		headers = append(headers, eventbus.Header{
			Key:   key,
			Value: val.Value,
		})
	}

	return headers
}

// WithReceiver adds a middleware to the queue receiver for instrumentation etc.
func (c *Consumer) WithReceiver(mw ReceiverMiddleware) {
	c.receiverMiddlewares = append(c.receiverMiddlewares, mw)
}

type ReceiverMiddleware func(next queues.QueueReceiver) queues.QueueReceiver
