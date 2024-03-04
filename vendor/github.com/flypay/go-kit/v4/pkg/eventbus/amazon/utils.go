package amazon

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	defaultRegion                   = "eu-west-1"
	protocolSQS                     = "sqs"
	subscriptionPendingConfirmation = "PendingConfirmation"

	partitionAWS = "aws"
	serviceSNS   = "sns"
	serviceSQS   = "sqs"

	largeMessageSize = 262144
)

func BuildARN(service, region, account, resource string) string {
	return arn.ARN{
		Partition: partitionAWS,
		Service:   service,
		Region:    region,
		AccountID: account,
		Resource:  resource,
	}.String()
}

// DefaultTopicARNStrategy defines the most common way of constructing a topic ARN
func DefaultTopicARNStrategy(region, account, topic string) string {
	return arn.ARN{
		Partition: partitionAWS,
		Service:   serviceSNS,
		Region:    region,
		AccountID: account,
		Resource:  SafeTopicName(topic),
	}.String()
}

// FlytTopicARNStrategy defines the most common way of constructing a topic ARN
func FlytTopicARNStrategy(region, account, topic string) string {
	return DefaultTopicARNStrategy(region, account, "event-bus-"+topic)
}

func extractTopicNameFromARN(topicARN string) string {
	if i := strings.LastIndex(topicARN, ":"); i >= 0 {
		return topicARN[i+1:]
	}
	return topicARN
}

// SafeTopicName returns the topic name but safe for AWS SQS queues
func SafeTopicName(topic string) string {
	return strings.ReplaceAll(topic, ".", "-")
}

// BuildQueueURL will build an AWS SQS url
func BuildQueueURL(region, accountID, queueName string) string {
	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, queueName)
}

// imports because it gets imported by internal packages which isn't good.
// HashedQueueName returns a fixed length queue name for any service & topic
func HashedQueueName(serviceName, topicName string) string {
	return fmt.Sprintf("event-bus-%x", sha256.Sum256([]byte(serviceName+"-"+topicName)))
}

// HashedDLQueueName returns a fixed length deadletter queue name for any service & topic
func HashedDLQueueName(serviceName, topicName string) string {
	return fmt.Sprintf("event-bus-dlq-%x", sha256.Sum256([]byte(serviceName+"-"+topicName)))
}

// MakeEnvAwareSession create a session based on the environment
func MakeEnvAwareSQSSNSSession() (*session.Session, error) {
	return MakeEnvAwareSQSSNSSessionWithRuntimeConfig(runtime.DefaultConfig())
}

func MakeEnvAwareSQSSNSSessionWithRuntimeConfig(cfg runtime.Config) (*session.Session, error) {
	if cfg.EventEndpoint != "" {
		return newCustomEndpointSQSSNSSession(cfg.EventEndpoint)
	}

	return newDefaultSession()
}

func MakeEnvAwareSFNSession() (*session.Session, error) {
	return MakeEnvAwareSFNSessionWithRuntimeConfig(runtime.DefaultConfig())
}

func MakeEnvAwareSFNSessionWithRuntimeConfig(cfg runtime.Config) (*session.Session, error) {
	if cfg.EventSFNEndpoint != "" {
		return newCustomEndpointSFNSession(cfg.EventSFNEndpoint)
	}

	return newDefaultSession()
}

func newDefaultSession() (*session.Session, error) {
	return session.NewSession(aws.NewConfig().WithRegion(defaultRegion))
}

func newCustomEndpointS3Session(endpoint string) (*session.Session, error) {
	return session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials("flyt", "flyt1234", "")).
		WithEndpoint(endpoint).
		WithRegion(defaultRegion).
		WithS3ForcePathStyle(true))
}

func newCustomEndpointSQSSNSSession(endpoint string) (*session.Session, error) {
	return session.NewSession(
		aws.NewConfig().
			WithCredentials(credentials.NewStaticCredentials("id", "secret", "token")).
			WithEndpoint(endpoint).
			WithRegion(defaultRegion),
	)
}

func newCustomEndpointSFNSession(endpoint string) (*session.Session, error) {
	return session.NewSession(
		aws.NewConfig().
			WithCredentials(credentials.NewStaticCredentials("id", "secret", "token")).
			WithEndpoint(endpoint).
			WithRegion(defaultRegion),
	)
}

func createLocalQueue(sess *session.Session, region, accountID, serviceQueueName string, topics []string) {
	sqsClient := sqs.New(sess)
	snsClient := sns.New(sess)

	deadLetterQueueName := serviceQueueName + "-deadletter"

	if _, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(deadLetterQueueName),
	}); err != nil {
		log.Errorf("Could not create dead letter queue %q for service %s: %v", deadLetterQueueName, serviceQueueName, err)
	}

	redrivePolicy := fmt.Sprintf(`{"maxReceiveCount": "2", "deadLetterTargetArn":"arn:aws:sqs::000000000000:%s"}`, deadLetterQueueName)
	_, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(serviceQueueName),
		Attributes: map[string]*string{
			"RedrivePolicy": &redrivePolicy,
		},
	})
	if err != nil {
		log.Errorf("Could not create queue for service %s: %v", serviceQueueName, err)
	}

	createLocalTopics(sess, region, accountID, topics)

	for _, t := range topics {
		topic := FlytTopicARNStrategy(region, accountID, t)
		_, err := snsClient.Subscribe(&sns.SubscribeInput{
			Protocol: aws.String(protocolSQS),
			Endpoint: aws.String(BuildARN(serviceSQS, region, accountID, serviceQueueName)),
			TopicArn: aws.String(topic),
		})
		if err != nil {
			log.Errorf("Could not subscribe topic [%s] to queue %s: %v", topic, serviceQueueName, err)
		}
	}
}

// encodeMessage is used to encode the proto.Message into a format
// that is compatible with AWS.
func encodeMessage(message proto.Message, opts ...func(opt *protojson.MarshalOptions)) (string, error) {
	mo := protojson.MarshalOptions{}

	for _, f := range opts {
		if f == nil {
			continue
		}
		f(&mo)
	}

	b, err := mo.Marshal(message)

	return string(b), errors.Wrap(err, "unable to protojson encode message")
}

// decodeMessage will decode an encoded message that was sent to AWS.
func decodeMessage(message string, protoMessage proto.Message) error {
	marshaller := protojson.UnmarshalOptions{DiscardUnknown: true}

	err := marshaller.Unmarshal([]byte(message), protoMessage)
	return errors.Wrap(err, "unable to protojson decode message")
}

func headers(msg SNSMessage) []eventbus.Header {
	var headers []eventbus.Header

	for key, val := range msg.MessageAttributes {
		headers = append(headers, eventbus.Header{
			Key:   key,
			Value: val.Value,
		})
	}

	return headers
}

func isLargePayload(size int) bool {
	// 256 KiB in B = 262144
	return size >= largeMessageSize
}

// SNSSubscriptionConfirmer returns an SubscriptionConfirmer function that confirms the subscription with SNS
func SNSSubscriptionConfirmer(svc snsiface.SNSAPI) func(context.Context, ConfirmSubscriptionMessage) error {
	return func(ctx context.Context, msg ConfirmSubscriptionMessage) error {
		_, err := svc.ConfirmSubscriptionWithContext(ctx, &sns.ConfirmSubscriptionInput{
			Token:    &msg.Token,
			TopicArn: &msg.TopicARN,
		})
		return err
	}
}

// SNSSubscriber will attempt to subscribe to the topics registered with consumer.On and consumer.OnAll if the subscription for the
// given queue does not already exist. This can be wired up in the consumer.Subscriber field before calling consumer.Listen.
func SNSSubscriber(svc snsiface.SNSAPI, queueARN string) func(topic string) {
	return func(topicARN string) {
		var next *string
		for {
			res, err := svc.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
				TopicArn:  aws.String(topicARN),
				NextToken: next,
			})
			if err != nil {
				log.Warnf("Listing subscriptions for topic %q err: %+v", topicARN, err)
				res = &sns.ListSubscriptionsByTopicOutput{Subscriptions: nil}
			}

			log.Debugf("Topic %q subscriptions: %+v", topicARN, res.Subscriptions)

			// Are we already subscribed?
			for _, sub := range res.Subscriptions {
				if *sub.Endpoint == queueARN && *sub.SubscriptionArn != subscriptionPendingConfirmation {
					log.Infof("Found existing subscription %+v", *sub)
					return
				}
			}

			// No more pages of subscriptions
			if res.NextToken == nil || *res.NextToken == "" {
				break
			}

			next = res.NextToken
		}

		log.Infof("Subscribing queue %s to topic %s", queueARN, topicARN)

		_, err := svc.Subscribe(&sns.SubscribeInput{
			Endpoint: aws.String(queueARN),
			Protocol: aws.String(protocolSQS),
			TopicArn: aws.String(topicARN),
		})
		if err != nil {
			log.Errorf("Subscribing queue %s to topic %q err: %+v", queueARN, topicARN, err)
		}
	}
}
