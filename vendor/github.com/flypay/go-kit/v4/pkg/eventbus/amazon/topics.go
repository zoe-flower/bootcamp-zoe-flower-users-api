package amazon

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// createLocalTopics creates topics locally for all topics to which the service subscribes.
// Should only be used when running the service locally using `dev-tool run`.
func createLocalTopics(sess *session.Session, region, account string, topics []string) {
	svc := sns.New(sess)
	for _, t := range topics {
		if err := createLocalTopic(context.Background(), svc, FlytTopicARNStrategy(region, account, t)); err != nil {
			log.Error(err)
			continue
		}
	}
}

// createLocalTopic is used when running services locally and topics need to be created
func createLocalTopic(ctx context.Context, svc snsiface.SNSAPI, topicARN string) error {
	topicARN = extractTopicNameFromARN(topicARN)
	_, err := svc.CreateTopicWithContext(ctx, &sns.CreateTopicInput{
		Name: aws.String(topicARN),
	})
	return errors.Wrapf(err, "creating topic %s", topicARN)
}

func dynamicTopicCreation(cfg ProducerConfig) eventbus.ProducerMiddleware {
	// Which topics are created cache to prevent creating cache on every emit.
	// Using a mutex to prevent data races in the cache which will cause sequencial emits when
	// running locally. If you have a problem with slow concurrent emits this might be the reason.
	cache := map[string]struct{}{}
	var lock sync.Mutex

	svc := sns.New(cfg.Sessions.SQSSNSSession)

	return func(next eventbus.EmitterFunc) eventbus.EmitterFunc {
		return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
			topic := cfg.TopicName(event)
			if topic == "" {
				return next(ctx, event, headers...)
			}

			lock.Lock()
			defer lock.Unlock()

			if _, exists := cache[topic]; exists {
				return next(ctx, event, headers...)
			}

			if err := createLocalTopic(ctx, svc, extractTopicNameFromARN(cfg.TopicARN(cfg.Region, cfg.AccountID, topic))); err == nil {
				cache[topic] = struct{}{}
			}

			return next(ctx, event, headers...)
		}
	}
}
