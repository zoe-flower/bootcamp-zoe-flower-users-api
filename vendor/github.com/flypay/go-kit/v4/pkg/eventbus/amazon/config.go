package amazon

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProducerSessions struct {
	SQSSNSSession *session.Session
	SFNSession    *session.Session
	S3Session     *session.Session
	Config        *aws.Config
}

type ConsumerSessions struct {
	SQSSNSSession *session.Session
	S3Session     *session.Session
}

// ProducerConfig holds configuration values for the SNS and the SFN producers
type ProducerConfig struct {
	Sessions ProducerSessions

	// AccountID AWS account ID and is used to build Topic ARNs
	AccountID string

	// Region is the AWS region to use when building Topic ARNs
	Region string

	// DelayStateMachineARN points at the state machine used to delay the emitting of event
	DelayStateMachineARN string

	// TopicName resolves a topic name from an event
	// The returned value should be just the topic name, before any alteration for AWS and also not an ARN.
	TopicName func(event protoreflect.ProtoMessage) string

	// TopicTenant resolves a tenant from an event
	TopicTenant func(event protoreflect.ProtoMessage) string

	// TopicARN allows configuration to mutate the way ARN strings are created for any topic.
	// If omitted the DefaultTopicARNStrategy is used. The returned string value should be a valid AWS ARN.
	TopicARN func(region, account, topic string) string

	TopicTenantedARN func(region, account, topic, tenant string) string

	// MarshallerOptions allows custom tuning of marshalling the encoding of the JSON message payload
	MarshallerOptions func(opt *protojson.MarshalOptions)

	// LargeMessageBucket is where messages over the Amazon threshold are stored
	LargeMessageBucket string
}

// ConsumerConfig holds configuration values for the SQS consumers
type ConsumerConfig struct {
	// The AWS session to use with the SNS or SQS clients
	Sessions ConsumerSessions

	Region      string
	AccountID   string
	ServiceName string

	// TopicName resolves a topic name from an event
	// The returned value should be just the topic name, before any alteration for AWS and also not an ARN.
	TopicName func(event protoreflect.ProtoMessage) string

	// TopicARN allows configuration to mutate the way ARN strings are created for any topic.
	// If omitted the DefaultTopicARNStrategy is used. The returned string value should be a valid AWS ARN.
	TopicARN func(region, account, topic string) string

	// TopicMultipleARNS allows for subscribing to a topic under multiple tenants
	TopicMultipleARNS func(region, account, topic string, tenants []string) []string

	// TopicTenants returns all topics from the event definition
	TopicTenants func(event protoreflect.ProtoMessage) []string

	// SubscriptionConfirmer is called when a confirmation message is received on the event queue
	SubscriptionConfirmer func(context.Context, ConfirmSubscriptionMessage) error

	// Subscribe is called during consumer.Listen for all events previously subscribed using `On` and `OnAll`.
	Subscribe func(topicARN string)

	// MessageSubject is an optional configuration taking a func that returns the topic name for a given
	// message recieved. This is useful when the SNS Message Subject field is not populated with the event name
	// as defined by the event proto options. This can be used to extract subject from elsewhere in the message payload
	// for example in the message attributes or message payload.
	// If omitted the default is to use the SNS Messages `Subject` attribute.
	MessageSubject func(SNSMessage) string

	// QueueURL is used to know from which queue to receive messages
	QueueURL                   string
	DefaultEventHandlerTimeout string
	VisibilityTimeout          string

	// Queue is added to metrics to differentiate topics in receive and pool metrics.
	Queue string

	// How many events should be handled concurrently.
	WorkerPoolSize int
}
