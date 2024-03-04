package middleware

import (
	"context"
	"os"
	"strings"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/env"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"google.golang.org/protobuf/proto"
)

type Testing struct{}

func NewTesting() *Testing {
	return &Testing{}
}

func (t Testing) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		serviceName := os.Getenv("CONFIG_APP_NAME") // can be flyt-monolith-sender or pr-{short_hash}

		if value := t.getNonce(headers); value != "" {
			ctx = flytcontext.WithExecutionNonce(ctx, value)

			deployingServiceName, commitHash, _, err := flytcontext.DecodeNonce(value)
			if err != nil {
				if strings.HasPrefix(serviceName, "pr-") {
					// When capability tests run on master, existing pr- services should ignore events
					return nil
				}

				return next(ctx, event, headers...)
			}

			// When capability tests are running on a PR the branch will be deployed to staging.
			// The following ensures only the newly deployed version of the service consumes events.

			// The master version of the service should not consume events
			if serviceName == deployingServiceName {
				return nil
			}

			// pr- versions deployed with the wrong hash should not consume events
			if strings.HasPrefix(serviceName, "pr-") &&
				!strings.Contains(serviceName, commitHash) {
				return nil
			}
			return next(ctx, event, headers...)
		}

		if env.IsStaging() && strings.HasPrefix(serviceName, "pr-") {
			// pr- services should ignore all events that don't have a nonce and are not part of capability tests
			return nil
		}

		return next(ctx, event, headers...)
	}
}

func (t Testing) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if value, ok := flytcontext.ExecutionNonce(ctx); ok {
			headers = append(headers, eventbus.Header{
				Key:   FlytExecutionNonce,
				Value: value,
			})
		}

		return next(ctx, event, headers...)
	}
}

func (t Testing) getNonce(headers []eventbus.Header) string {
	for _, header := range headers {
		if header.Key == FlytExecutionNonce {
			return header.Value
		}
	}

	return ""
}
