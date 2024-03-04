package eventbus

import (
	"github.com/flypay/go-kit/v4/pkg/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetTopicFromEvent returns the topic name for an event.
func GetTopicFromEvent(message proto.Message, topic protoreflect.ExtensionType) (string, error) {
	options := message.ProtoReflect().Descriptor().Options()
	value := proto.GetExtension(options, topic)

	if s, ok := value.(string); ok && s != "" {
		return s, nil
	}

	return "", ErrTopicNameNotDefined
}

// GetTenantsFromEvent returns all tenants for an event from the defined proto extension type
func GetTenantsFromEvent(message proto.Message, topic protoreflect.ExtensionType) []string {
	options := message.ProtoReflect().Descriptor().Options()
	value := proto.GetExtension(options, topic)

	if s, ok := value.([]string); ok && len(s) > 0 {
		return s
	}

	return []string{"all"}
}

// ProtoTopicNamingStrategy returns the name of the topic from the defined proto extension type
func ProtoTopicNamingStrategy(ext protoreflect.ExtensionType) func(protoreflect.ProtoMessage) string {
	return func(event protoreflect.ProtoMessage) string {
		topic, _ := GetTopicFromEvent(event, ext)
		return topic
	}
}

// ProtoTopicTenantsStrategy returns all tenants of the topic from the defined proto extension type
func ProtoTopicTenantsStrategy(ext protoreflect.ExtensionType) func(protoreflect.ProtoMessage) []string {
	return func(event protoreflect.ProtoMessage) []string {
		return GetTenantsFromEvent(event, ext)
	}
}

// GracefulClose is a type that ensures that the consumers have closed before
// closing the producers. This prevents event handlers from trying to emit
// events on a closed producer.
type GracefulClose struct {
	Consumers []Consumer
	Producers []Producer
}

// Close will close the consumers, wait and then close the producers.
func (t GracefulClose) Close() error {
	var errConsumer error
	for i, c := range t.Consumers {
		if err := c.Close(); err != nil {
			log.Errorf("Unable to close consumer %d: %v", i+1, err)
			errConsumer = err
		}
	}

	var errProducer error
	for i, p := range t.Producers {
		if err := p.Close(); err != nil {
			log.Errorf("Unable to close producer %d: %v", i+1, err)
			errProducer = err
		}
	}

	if errConsumer != nil {
		return errConsumer
	}

	return errProducer
}
