package middleware

import (
	"context"
	"errors"

	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/log"
	"google.golang.org/protobuf/proto"
)

type Logging struct{}

// NewLogging creates a type that returns consumer and producer logging middlewares
func NewLogging() Logging {
	return Logging{}
}

func (Logging) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		log.WithContext(ctx).Debugf("Handling %T event", event)

		if err := next(ctx, event, headers...); err != nil {
			hErr := &eventbus.HandlerError{}
			if errors.As(err, &hErr) {
				log.WithContext(ctx).Errorf("Error handling %T event (code: %s, terminal: %t): %+v",
					event, hErr.ErrorCode, hErr.Terminal, err)
			} else {
				log.WithContext(ctx).Errorf("Error handling %T event: %+v", event, err)
			}

			return err
		}

		return nil
	}
}

func (Logging) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		if err := next(ctx, event, headers...); err != nil {
			log.WithContext(ctx).Errorf("Error emitting %T event: %+v", event, err)
			return err
		}

		log.WithContext(ctx).Debugf("Emitted %T event", event)

		return nil
	}
}
