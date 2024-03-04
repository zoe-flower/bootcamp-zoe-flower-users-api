package amazon

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flypay/go-kit/v4/pkg/objstore"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sfn/sfniface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/inmemory"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrProducerClosed = errors.New("producer closed")
	ErrNotConnected   = errors.New("no connection available to perform request")
)

const (
	// immediateEmitThreshold is the duration from time.Now() at which events are immediately
	// emitted to the SNS topic even though the EmitAt function was called (bypassing the step function).
	immediateEmitThreshold = 1 * time.Second

	// AWS message attribute data type for string values.
	dataTypeString = "String"

	// AWS message attribute data type for number values.
	dataTypeNumber = "Number"
)

// Producer is an implementation of a message producer that uses SNS and SFN
type Producer struct {
	sns                 snsiface.SNSAPI
	sm                  sfniface.SFNAPI
	largeMessageStorage objstore.Uploader
	middleware          []eventbus.ProducerMiddleware

	cfg ProducerConfig

	// quit is a channel that will be closed when calling
	// 'Close()' causing an error to be emitted if 'Emit()'
	// is called at a later time.
	quit chan struct{}
}

// ResolveNewProducer resolves a different producer based on the environment
func ResolveNewProducer(config ProducerConfig) (eventbus.Producer, error) {
	return ResolveNewProducerWithRuntimeConfig(config, runtime.DefaultConfig())
}

func ResolveNewProducerWithRuntimeConfig(config ProducerConfig, rcfg runtime.Config) (eventbus.Producer, error) {
	if rcfg.EventInMemory {
		return inmemory.NewMemProducer(), nil
	}

	p := NewProducer(config)

	largeMessageStorage, err := createLargeMessageStorage(config.Sessions.S3Session, rcfg)
	if err != nil {
		return nil, err
	}

	p.largeMessageStorage = largeMessageStorage

	if rcfg.EventLocalCreate {
		createLocalTopics(config.Sessions.SQSSNSSession, config.Region, config.AccountID, rcfg.EventTopics)
		p.cfg.DelayStateMachineARN = createDelayedEmitterStepFunction(config.Sessions.SFNSession, config.Region, config.AccountID, rcfg.EventDynamicTaskARN)
		p.Use(dynamicTopicCreation(config))
	}

	return p, nil
}

// NewProducer returns a new producer with an SNS and an SFN clients initialized based on the ProducerConfig.Sessions
func NewProducer(cfg ProducerConfig) *Producer {
	if cfg.AccountID == "" {
		panic("sns producer not starting. CONFIG_AWS_ACCOUNT_ID is empty")
	}

	svc := sns.New(cfg.Sessions.SQSSNSSession, cfg.Sessions.Config)
	xray.AWS(svc.Client)

	if cfg.TopicName == nil {
		cfg.TopicName = func(event protoreflect.ProtoMessage) string { return "" }
	}

	if cfg.TopicTenant == nil {
		cfg.TopicTenant = func(event protoreflect.ProtoMessage) string { return "" }
	}

	if cfg.TopicARN == nil {
		cfg.TopicARN = DefaultTopicARNStrategy
	}

	p := Producer{
		sns:  svc,
		cfg:  cfg,
		quit: make(chan struct{}),
	}

	if cfg.Sessions.SFNSession != nil {
		sm := sfn.New(cfg.Sessions.SFNSession, cfg.Sessions.Config)
		xray.AWS(sm.Client)
		p.sm = sm
	}

	return &p
}

// Emit will produce a message to the SNS topic for the event descriptor
func (p *Producer) Emit(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
	return p.EmitAt(ctx, time.Now(), event, headers...)
}

// EmitAt will produce a message (containing an event) to the dynamic task state machine, that will be published to SNS
// at the specified time
func (p *Producer) EmitAt(ctx context.Context, at time.Time, event proto.Message, headers ...eventbus.Header) error {
	select {
	case <-p.quit:
		return ErrProducerClosed
	default:
		// don't block and fallthrough
	}

	next := p.selectEmitStrategy(at)

	for _, middleware := range p.middleware {
		next = middleware(next)
	}

	return next(ctx, event, headers...)
}

func (p *Producer) selectEmitStrategy(emitAt time.Time) eventbus.EmitterFunc {
	if emitAt.Before(time.Now().Add(immediateEmitThreshold)) {
		return p.emit
	}

	// A new emitter function that passes the emitAt value
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		return p.emitAt(ctx, emitAt, event, headers...)
	}
}

func (p *Producer) messagePointer(ctx context.Context, event string) (string, error) {
	filename := uuid.NewString() + ".json"

	if p.cfg.LargeMessageBucket == "" {
		return "", fmt.Errorf("asked to upload a large payload to object storage, but ProducerConfig.LargeMessageBucket was empty")
	}

	messageURL, err := p.largeMessageStorage.Upload(
		ctx, objstore.Bucket(p.cfg.LargeMessageBucket), filename, strings.NewReader(event),
	)

	return messageURL, err
}

func (p *Producer) emit(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
	topic, err := p.topicName(event)
	if err != nil {
		return err
	}

	topicARN := p.cfg.TopicARN(p.cfg.Region, p.cfg.AccountID, topic)
	if p.cfg.TopicTenantedARN != nil {
		tenant := p.cfg.TopicTenant(event)
		topicARN = p.cfg.TopicTenantedARN(p.cfg.Region, p.cfg.AccountID, topic, tenant)
	}

	headers = eventbus.WithEmittedAt(time.Now(), headers...)

	encodedMessage, err := encodeMessage(event, p.cfg.MarshallerOptions)
	if err != nil {
		return errors.Wrapf(err, "could not marshal proto.Message to publish to SNS")
	}

	publishInput := &sns.PublishInput{
		Message:           aws.String(encodedMessage),
		TopicArn:          aws.String(topicARN),
		Subject:           aws.String(topic),
		MessageAttributes: buildMessageHeaders(headers),
	}

	if isLargePayload(len(encodedMessage)) {
		publishInput.MessageAttributes = p.withExtendedPayloadSize(publishInput.MessageAttributes, len(encodedMessage))

		messagePointer, err := p.messagePointer(ctx, encodedMessage)
		if err != nil {
			return err
		}

		publishInput.SetMessage(messagePointer)
	}

	if _, err := p.sns.PublishWithContext(ctx, publishInput); err != nil {
		return errors.Wrapf(err, "could not publish message to SNS topic [%s]", topicARN)
	}

	return nil
}

func (p *Producer) emitAt(ctx context.Context, at time.Time, event proto.Message, headers ...eventbus.Header) error {
	if p.sm == nil {
		return ErrNotConnected
	}

	topic, err := p.topicName(event)
	if err != nil {
		return err
	}

	encodedMessage, err := encodeMessage(event)
	if err != nil {
		return errors.Wrapf(err, "could not marshal proto.Message to publish to SNS")
	}

	topicARN := p.cfg.TopicARN(p.cfg.Region, p.cfg.AccountID, topic)

	headers = eventbus.WithEmittedAt(at, headers...)

	input := sfnDynamicTaskInput{
		EventBusTopicArn:  topicARN,
		EventName:         topic,
		EventMessage:      encodedMessage,
		MessageAttributes: buildMessageHeaders(headers),
		EmitDateTime:      at.Format(time.RFC3339),
	}

	if isLargePayload(len(encodedMessage)) {
		input.MessageAttributes = p.withExtendedPayloadSize(input.MessageAttributes, len(encodedMessage))

		messagePointer, err := p.messagePointer(ctx, encodedMessage)
		if err != nil {
			return err
		}

		input.EventMessage = messagePointer
	}

	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	log.WithContext(ctx).Debugf("Event %T will be delayed and emitted at %s", event, at.Format(time.RFC3339))

	if _, err := p.sm.StartExecutionWithContext(ctx, &sfn.StartExecutionInput{
		Name:            aws.String(uuid.New().String()),
		StateMachineArn: aws.String(p.cfg.DelayStateMachineARN),
		Input:           aws.String(string(b)),
	}); err != nil {
		return err
	}

	return nil
}

func (p *Producer) withExtendedPayloadSize(attributes map[string]*sns.MessageAttributeValue, size int) map[string]*sns.MessageAttributeValue {
	attributes["ExtendedPayloadSize"] = &sns.MessageAttributeValue{
		StringValue: aws.String(strconv.Itoa(size)),
		DataType:    aws.String(dataTypeNumber),
	}

	return attributes
}

// topicName builds a topic name for a given event by using the naming strategy given
// in the producer configuration or falling back to the flyt version.
func (p *Producer) topicName(event protoreflect.ProtoMessage) (topic string, err error) {
	topic = p.cfg.TopicName(event)
	if topic == "" {
		err = eventbus.ErrTopicNameNotDefined
	}

	return
}

func buildMessageHeaders(headers []eventbus.Header) map[string]*sns.MessageAttributeValue {
	hdrs := make(map[string]*sns.MessageAttributeValue)

	for _, header := range headers {
		if header.Key == "" {
			continue
		}

		hdrs[header.Key] = &sns.MessageAttributeValue{
			StringValue: aws.String(header.Value),
			DataType:    aws.String(dataTypeString),
		}
	}

	return hdrs
}

// Use places middleware into the middleware stack for the producer
func (p *Producer) Use(middleware eventbus.ProducerMiddleware) {
	p.middleware = append(p.middleware, middleware)
}

// Close closes things
func (p *Producer) Close() error {
	close(p.quit)
	return nil
}
