package amazon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/flypay/go-kit/v4/pkg/objstore"
	"github.com/flypay/go-kit/v4/pkg/tracing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/inmemory"
	"github.com/flypay/go-kit/v4/pkg/eventbus/internal/duration"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/queues"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var (
	QueueReceiveWaitTimeSeconds int64 = 20

	// ErrMissingSubscriptionConfirmation is raised when a job handler receives a subscription cofirmation message
	// but there has not been a callback configured.
	ErrMissingSubscriptionConfirmation = errors.New("No subscription confirmation callback configured")
)

const (
	typeSubscriptionConfirmation = "SubscriptionConfirmation"
)

// Consumer is an SQS implementation of eventbus.Consumer
type Consumer struct {
	eventbus.ConsumerRegister

	largeMessageStorage objstore.URLObjectDownloader

	sqs sqsiface.SQSAPI
	cfg ConsumerConfig

	producer  *queues.Producer
	consumers []queues.Stoppable

	handlerTimeout    time.Duration
	visibilityTimeout time.Duration

	muListening sync.Mutex
	listening   bool

	// receiverMiddlewares are wrapped around the queue receiver. Commonly used to add logging, metrics etc.
	receiverMiddlewares []ReceiverMiddleware
}

func ResolveNewConsumer(config ConsumerConfig) (eventbus.Consumer, error) {
	return ResolveNewConsumerWithRuntimeConfig(config, runtime.DefaultConfig())
}

func ResolveNewConsumerWithRuntimeConfig(config ConsumerConfig, rcfg runtime.Config) (eventbus.Consumer, error) {
	if rcfg.EventInMemory {
		return inmemory.NewMemConsumer(inmemory.ConsumerConfig{
			TopicName:                  config.TopicName,
			DefaultEventHandlerTimeout: config.DefaultEventHandlerTimeout,
		}), nil
	}

	c, err := NewConsumer(config)
	if err != nil {
		return nil, err
	}

	largeMessageStorage, err := createLargeMessageStorage(config.Sessions.S3Session, rcfg)
	if err != nil {
		return nil, err
	}

	c.largeMessageStorage = largeMessageStorage

	if rcfg.EventLocalCreate {
		createLocalQueue(config.Sessions.SQSSNSSession, rcfg.AWSRegion, rcfg.AWSAccountID, rcfg.EventSQSQueueName, rcfg.EventTopics)
	}

	return c, nil
}

// NewConsumer builds a new amazon.Consumer
func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	svc := sqs.New(cfg.Sessions.SQSSNSSession)

	handlerTimeout, err := duration.GreaterThanZero(cfg.DefaultEventHandlerTimeout)
	if err != nil {
		return nil, err
	}

	visibilityTimeout, err := duration.GreaterThanZero(cfg.VisibilityTimeout)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		ConsumerRegister: eventbus.ConsumerRegister{
			TopicName:      cfg.TopicName,
			DefaultTimeout: handlerTimeout,
		},
		sqs:               svc,
		cfg:               cfg,
		handlerTimeout:    handlerTimeout,
		visibilityTimeout: visibilityTimeout,
	}, nil
}

// Listen binds workers in pools to the specified queue for the topic
func (c *Consumer) Listen() {
	c.muListening.Lock()
	defer c.muListening.Unlock()

	// Protect against calling the listen function twice
	if c.listening {
		return
	}

	if !c.HasSubscriptions() {
		log.Warn("SQS consumer is not starting, there are no handlers registered")
		return
	}

	// Setup defaults before starting the consumer
	if c.cfg.MessageSubject == nil {
		c.cfg.MessageSubject = func(msg SNSMessage) string { return msg.Subject }
	}

	if c.cfg.TopicARN == nil {
		c.cfg.TopicARN = DefaultTopicARNStrategy
	}

	var receiver queues.QueueReceiver = queues.NewQueue(c.cfg.QueueURL, c.sqs, QueueReceiveWaitTimeSeconds, c.visibilityTimeout)
	// TODO move the instrumentation to middleware
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

	// TODO this feels like it belongs in the Just Saying library rather than the generic SNS/SQS consumer code.
	// There is too much confusing configuration in the consumer now, and the clean Subscribe hook is polluted with the concept
	// of Just Eat tenants.
	if c.cfg.Subscribe != nil {
		for topic, p := range c.ConsumerRegister.Subscriptions() {
			if c.cfg.TopicMultipleARNS != nil && c.cfg.TopicTenants != nil {
				arns := c.cfg.TopicMultipleARNS(c.cfg.Region, c.cfg.AccountID, topic, c.cfg.TopicTenants(p))
				for _, arn := range arns {
					c.cfg.Subscribe(arn)
				}
			} else {
				c.cfg.Subscribe(c.cfg.TopicARN(c.cfg.Region, c.cfg.AccountID, topic))
			}
		}
	}

	c.listening = true
}

// The returned error is discarded by the queue consumer so we must handle it here
func (c *Consumer) handle(job queues.Job) (ok bool, err error) {
	logger := log.WithContext(job.Context())

	var msg SNSMessage
	err = json.NewDecoder(job.Payload()).Decode(&msg)
	if err != nil {
		logger.Errorf("Unable to decode SNS message: %s, err: %+v", job.ID(), err)
		return false, nil
	}

	switch msg.Type {
	case typeSubscriptionConfirmation:
		ok, err = c.handleSNSSubscriptionConfirmation(job.Context(), msg)
	default:
		ok, err = c.handleSNSNotification(job.Context(), msg)
	}
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

func (c *Consumer) handleSNSSubscriptionConfirmation(ctx context.Context, msg SNSMessage) (bool, error) {
	if c.cfg.SubscriptionConfirmer == nil {
		return false, errors.Wrapf(ErrMissingSubscriptionConfirmation, "no subscription confirmation callback configured for message: %s", msg.MessageID)
	}

	log.WithContext(ctx).Debugf("Subscription confirmation received for topic: %s", msg.TopicArn)

	cmsg := ConfirmSubscriptionMessage{TopicARN: msg.TopicArn, Token: msg.Token}

	if err := c.cfg.SubscriptionConfirmer(ctx, cmsg); err != nil {
		return false, errors.Wrapf(err, "handling subscription confirmation for job %s on topic %s", msg.MessageID, msg.TopicArn)
	}

	return true, nil
}

type resetter interface {
	Reset()
}

func (c *Consumer) resolveMessagePointer(ctx context.Context, messageURL string) (string, error) {
	storageBackedMessageURL, err := url.Parse(messageURL)
	if err != nil {
		return "", fmt.Errorf("url: unable to parse: %q, %w", messageURL, err)
	}

	ctx, span := tracing.StartSpan(ctx, "consumer.large-messageUrl")
	defer span.Close(nil)

	originalEv, err := c.largeMessageStorage.DownloadObjectURL(ctx, storageBackedMessageURL)
	if err != nil {
		return "", fmt.Errorf("storage: unable to download object url: %q, %w", messageURL, err)
	}

	return string(originalEv), nil
}

func (c *Consumer) handleSNSNotification(ctx context.Context, msg SNSMessage) (bool, error) {
	msgSubject := c.cfg.MessageSubject(msg)

	// Look up the handler for the topic (subject of the message)
	handler, ok := c.TopicHandler(msgSubject)
	if !ok {
		return true, fmt.Errorf("no handler found for topic %q (%q), message ID %q, message: %s",
			msgSubject, msg.Subject, msg.MessageID, msg.Message)
	}

	// Clone the handler.Event so we don't use the same pointer
	// for every handler.
	// event := proto.Clone(handler.Event)
	event := handler.Pool.Get().(proto.Message)
	defer func() {
		if ev, ok := event.(resetter); ok {
			ev.Reset()
		}
		handler.Pool.Put(event)
	}()

	if _, ok := msg.MessageAttributes["ExtendedPayloadSize"]; ok {
		resolvedMessage, err := c.resolveMessagePointer(ctx, msg.Message)
		if err != nil {
			return false, err
		}

		msg.Message = resolvedMessage
	}

	if err := decodeMessage(msg.Message, event); err != nil {
		return true, errors.Wrapf(err, "unmarshalling event for message: %s", msg.MessageID)
	}

	if err := c.Handle(handler, event, headers(msg)...); err != nil {
		// Check to see if we have a terminal error. Terminal errors
		// should not be retried via SQS but should still reported in metrics etc.
		var eventbusErr *eventbus.HandlerError
		if errors.As(err, &eventbusErr) && eventbusErr.Terminal {
			// Return true to delete the message as we do not want it to be redelivered (retried)
			return true, err
		}

		// Discard this error as the eventbus middlewares will capture and log it
		// Prevents double logging of errors in kibana.
		return false, nil
	}

	return true, nil
}

// Close stops the context and shuts down all producers and consumers
func (c *Consumer) Close() error {
	log.Info("Consumer stopping...")

	if c.producer != nil {
		queues.StopAll(c.producer)
	}
	queues.StopAll(c.consumers...)

	log.Info("Consumer stopped")

	return nil
}

// WithReceiver adds a middleware to the queue receiver for instrumentation etc.
func (c *Consumer) WithReceiver(mw ReceiverMiddleware) *Consumer {
	c.receiverMiddlewares = append(c.receiverMiddlewares, mw)
	return c
}

type ReceiverMiddleware func(next queues.QueueReceiver) queues.QueueReceiver
