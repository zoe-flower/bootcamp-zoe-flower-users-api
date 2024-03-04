package middleware

import (
	"context"
	"reflect"

	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/tracing"
	"google.golang.org/protobuf/proto"
)

// headerServiceName is applied to all events emitted to set the service name
// whilst not technically tracing this does serve debugging.
const headerServiceName = "service_name"

type Tracing struct {
	name string
	team string
	flow string
	tier string
}

func NewTracing(name, team, flow, tier string) *Tracing {
	return &Tracing{
		name: name,
		team: team,
		flow: flow,
		tier: tier,
	}
}

func (t Tracing) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		xrayHeader := t.findTraceHeader(headers)

		var err error

		ctx1, seg := t.segmentForHeader(ctx, xrayHeader, t.name)
		defer func() {
			seg.Close(err)
		}()

		tracing.AddAnnotations(
			ctx1,
			tracing.Annotation{Key: tracing.TeamKey, Value: t.team},
			tracing.Annotation{Key: tracing.FlowKey, Value: t.flow},
			tracing.Annotation{Key: tracing.TierKey, Value: t.tier},
		)

		err = xray.Capture(ctx1, "eventbus.Handle", func(ctx2 context.Context) error {
			if eventType := reflect.TypeOf(event); eventType != nil {
				tracing.AddAnnotation(ctx2, "event", eventType.String())
			}
			return next(ctx2, event, headers...)
		})
		return err
	}
}

func (t Tracing) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		return xray.Capture(ctx, "eventbus.Emit", func(ctx context.Context) error {
			if seg := xray.GetSegment(ctx); seg != nil {
				headers = append(headers, eventbus.Header{
					Key:   xray.TraceIDHeaderKey,
					Value: seg.DownstreamHeader().String(),
				})
				tracing.AddAnnotations(
					ctx,
					tracing.Annotation{Key: tracing.FlowKey, Value: t.flow},
					tracing.Annotation{Key: tracing.TeamKey, Value: t.team},
					tracing.Annotation{Key: tracing.TierKey, Value: t.tier},
					tracing.Annotation{Key: "event", Value: reflect.TypeOf(event).String()},
				)
			}

			// Adds the service name header for debugging the source of an event. This was
			// extracted from the amazon.Producer.Emit code in buildMessageHeaders
			headers = append(headers, eventbus.Header{Key: headerServiceName, Value: t.name})

			return next(ctx, event, headers...)
		})
	}
}

func (t Tracing) findTraceHeader(headers []eventbus.Header) *header.Header {
	for _, h := range headers {
		if h.Key == xray.TraceIDHeaderKey {
			return header.FromString(h.Value)
		}
	}

	return nil
}

func (t Tracing) segmentForHeader(ctx context.Context, h *header.Header, name string) (context.Context, *xray.Segment) {
	if h != nil {
		return xray.NewSegmentFromHeader(ctx, name, nil, h)
	}

	return xray.BeginSegment(ctx, name)
}
