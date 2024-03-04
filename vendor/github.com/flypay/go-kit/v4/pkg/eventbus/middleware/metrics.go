package middleware

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/eventbus"
	"github.com/flypay/go-kit/v4/pkg/eventbus/amazon"
	"github.com/flypay/go-kit/v4/pkg/metrics"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Metrics struct {
	consumed                         *prometheus.CounterVec
	consumerErrors                   *prometheus.CounterVec
	consumerPanics                   *prometheus.CounterVec
	consumerStatusCodes              *prometheus.CounterVec
	consumerDuration                 prometheus.ObserverVec
	consumerLag                      prometheus.ObserverVec
	consumerEventTimeout             *prometheus.CounterVec
	consumerEventTimeoutExceeded     prometheus.ObserverVec
	consumerInflightEvents           *prometheus.GaugeVec
	consumerInflightEventsPercentage *prometheus.GaugeVec
	produced                         *prometheus.CounterVec
	producerErrors                   *prometheus.CounterVec
	producerDuration                 prometheus.ObserverVec

	clock     Nower
	topicExt  protoreflect.ExtensionType
	queueName string

	consumerWorkerPoolSize float64
	inflightEvents         int64 // atomically incremented
}

// Nower returns the current time
type Nower interface {
	Now() time.Time
}

const (
	typeSNSSQS = "SNS_SQS"
)

// NewMetrics factories out new metrics for use in the consumer and producer middleware
func NewMetrics(clock Nower, topicExt protoreflect.ExtensionType, consumerWorkerPoolSize int) *Metrics {
	m := newMetrics(clock, consumerWorkerPoolSize)
	m.topicExt = topicExt
	return m
}

func NewMetricsForQueue(clock Nower, queueName string, consumerWorkerPoolSize int) *Metrics {
	m := newMetrics(clock, consumerWorkerPoolSize)
	m.queueName = queueName
	return m
}

func newMetrics(clock Nower, consumerWorkerPoolSize int) *Metrics {
	labels := prometheus.Labels{"eventbus_type": typeSNSSQS}

	consumed := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_total",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.CounterVec)

	consumerErrors := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_errors_total",
			ConstLabels: labels,
		}, []string{"event", "topic", "terminal"})).(*prometheus.CounterVec)

	consumerPanics := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_panics_total",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.CounterVec)

	consumerStatusCodes := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_status_codes_total",
			ConstLabels: labels,
		}, []string{"event", "topic", "status_code", "terminal"})).(*prometheus.CounterVec)

	consumerDuration := metrics.RegisterOrExisting(
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem:   "eventbus",
			Name:        "duration_of_event_processing_seconds",
			Buckets:     []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 40, 60, 80, 160},
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.HistogramVec)

	consumerLag := metrics.RegisterOrExisting(
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem:   "eventbus",
			Name:        "consumer_lag_milliseconds",
			Help:        "Amount of milliseconds the consumer is running behind the producer.",
			ConstLabels: labels,
			Buckets:     prometheus.ExponentialBuckets(10, 2, 12),
		}, []string{"event", "topic"})).(*prometheus.HistogramVec)

	consumerEventTimeout := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_event_timeout_total",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.CounterVec)

	consumerEventTimeoutExceeded := metrics.RegisterOrExisting(
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_event_timeout_exceeded_seconds",
			Buckets:     []float64{1, 2, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.HistogramVec)

	consumerInflightEvents := metrics.RegisterOrExisting(
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_events_inflight",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.GaugeVec)

	consumerInflightEventsPercentage := metrics.RegisterOrExisting(
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Subsystem:   "eventbus",
			Name:        "events_consumer_events_inflight_percentage",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.GaugeVec)

	produced := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_producer_total",
			Help:        "Number of events emitted to the producer.",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.CounterVec)

	producerErrors := metrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "eventbus",
			Name:        "events_producer_errors_total",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.CounterVec)

	producerDuration := metrics.RegisterOrExisting(
		prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem:   "eventbus",
			Name:        "duration_of_event_emitting_seconds",
			Help:        "Time spent emitting in seconds.",
			ConstLabels: labels,
		}, []string{"event", "topic"})).(*prometheus.HistogramVec)

	// Set an initial value so kubernetes HPA doesn't error and include
	// the topic name if available.
	consumerInflightEventsPercentage.WithLabelValues(
		"",
		runtime.DefaultConfig().EventSQSTopicName,
	).Set(0)

	return &Metrics{
		consumed:                         consumed,
		consumerErrors:                   consumerErrors,
		consumerPanics:                   consumerPanics,
		consumerStatusCodes:              consumerStatusCodes,
		consumerDuration:                 consumerDuration,
		consumerLag:                      consumerLag,
		consumerEventTimeout:             consumerEventTimeout,
		consumerEventTimeoutExceeded:     consumerEventTimeoutExceeded,
		consumerInflightEvents:           consumerInflightEvents,
		consumerInflightEventsPercentage: consumerInflightEventsPercentage,
		produced:                         produced,
		producerErrors:                   producerErrors,
		producerDuration:                 producerDuration,

		clock:                  clock,
		consumerWorkerPoolSize: float64(consumerWorkerPoolSize),
	}
}

func (m Metrics) ConsumerMiddleware(next eventbus.EventHandler) eventbus.EventHandler {
	next = m.consumerCounter(next)
	next = m.countConsumerErrors(next)
	next = m.countConsumerPanics(next)
	next = m.consumeDuration(next)
	next = m.measureConsumerLag(next, m.clock)
	next = m.countStatusCodes(next)
	next = m.consumerEventsInFlight(next)
	next = m.consumerEventTimeouts(next)

	return next
}

func (m Metrics) ProducerMiddleware(next eventbus.EmitterFunc) eventbus.EmitterFunc {
	next = m.countEmitted(next)
	next = m.countEmitErrors(next)
	next = m.emitDuration(next)

	return next
}

func (m Metrics) topic(event interface{}) string {
	if m.queueName != "" {
		return m.queueName
	}
	var topic string
	if _, ok := event.(proto.Message); !ok {
		return topic
	}

	if t, err := eventbus.GetTopicFromEvent(event.(proto.Message), m.topicExt); err == nil {
		// TODO remove dependency on the amazon version of the topic by migrating all
		// events options to the amazon format.
		topic = amazon.SafeTopicName(t)
	}
	return topic
}

func (m Metrics) event(event interface{}) string {
	if m.queueName != "" {
		return m.queueName
	}
	return reflect.TypeOf(event).String()
}

func (m Metrics) emitDuration(emitter eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		start := time.Now()
		err := emitter(ctx, event, headers...)

		m.producerDuration.
			WithLabelValues(m.event(event), m.topic(event)).
			Observe(time.Since(start).Seconds())

		return err
	}
}

func (m Metrics) countEmitErrors(emit eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		err := emit(ctx, event, headers...)
		if err != nil {
			m.producerErrors.
				WithLabelValues(m.event(event), m.topic(event)).
				Inc()
		}
		return err
	}
}

func (m Metrics) countEmitted(emit eventbus.EmitterFunc) eventbus.EmitterFunc {
	return func(ctx context.Context, event proto.Message, headers ...eventbus.Header) error {
		err := emit(ctx, event, headers...)
		m.produced.
			WithLabelValues(m.event(event), m.topic(event)).
			Inc()
		return err
	}
}

func (m Metrics) consumeDuration(handler eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		start := time.Now()
		err := handler(ctx, event, headers...)
		m.consumerDuration.
			WithLabelValues(m.event(event), m.topic(event)).
			Observe(time.Since(start).Seconds())
		return err
	}
}

func (m Metrics) countConsumerErrors(handler eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		err := handler(ctx, event, headers...)
		ctxErr := ctx.Err()
		if err != nil || ctxErr != nil {
			ev := m.event(event)
			topic := m.topic(event)
			terminal := false
			var ebErr *eventbus.HandlerError
			if errors.As(err, &ebErr) {
				terminal = ebErr.Terminal
			}
			m.consumerErrors.WithLabelValues(ev, topic, strconv.FormatBool(terminal)).Inc()
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctxErr, context.DeadlineExceeded) {
				m.consumerEventTimeout.WithLabelValues(ev, topic).Inc()
			}
		}
		return err
	}
}

func (m Metrics) countConsumerPanics(next eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		err := next(ctx, event, headers...)
		var p IsPanic
		if err != nil && errors.As(err, &p) {
			ev := m.event(event)
			topic := m.topic(event)
			m.consumerPanics.WithLabelValues(ev, topic).Inc()
		}
		return err
	}
}

func (m Metrics) countStatusCodes(handler eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		statusCode := eventbus.StatusOK
		terminal := false
		err := handler(ctx, event, headers...)
		if err != nil {
			statusCode = eventbus.StatusInternalError
			var eventbusErr *eventbus.HandlerError
			if errors.As(err, &eventbusErr) {
				statusCode = eventbusErr.ErrorCode
				terminal = eventbusErr.Terminal
			}
		}
		m.consumerStatusCodes.
			WithLabelValues(
				m.event(event),
				m.topic(event),
				statusCode.String(),
				strconv.FormatBool(terminal),
			).
			Inc()
		return err
	}
}

func (m Metrics) consumerCounter(handler eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		err := handler(ctx, event, headers...)
		m.consumed.
			WithLabelValues(m.event(event), m.topic(event)).
			Inc()
		return err
	}
}

func (m Metrics) measureConsumerLag(handler eventbus.EventHandler, clock Nower) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		if then, ok := eventbus.EmittedAt(headers...); ok {
			now := clock.Now()
			m.consumerLag.
				WithLabelValues(m.event(event), m.topic(event)).
				Observe(float64(now.Sub(then).Milliseconds()))
		}
		return handler(ctx, event, headers...)
	}
}

func (m *Metrics) consumerEventsInFlight(handler eventbus.EventHandler) eventbus.EventHandler {
	const makePercent = 100

	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		ev := m.event(event)
		topic := m.topic(event)

		inflight := m.consumerInflightEvents.WithLabelValues(ev, topic)
		percentage := m.consumerInflightEventsPercentage.WithLabelValues(ev, topic)

		measure := func(delta int64) {
			current := float64(atomic.AddInt64(&m.inflightEvents, delta))
			inflight.Set(current)
			percentage.Set((current / m.consumerWorkerPoolSize) * makePercent)
		}

		measure(1)
		defer measure(-1)

		return handler(ctx, event, headers...)
	}
}

func (m Metrics) consumerEventTimeouts(handler eventbus.EventHandler) eventbus.EventHandler {
	return func(ctx context.Context, event interface{}, headers ...eventbus.Header) error {
		err := handler(ctx, event, headers...)

		// Check the context error for a deadline exceeded. This happens when the
		// event handler timeout is exceeded.
		ctxErr := ctx.Err()
		if errors.Is(ctxErr, context.DeadlineExceeded) {
			ev := m.event(event)
			topic := m.topic(event)
			m.consumerEventTimeout.WithLabelValues(ev, topic).Inc()

			if deadline, ok := flytcontext.EventDeadline(ctx); ok {
				excess := time.Since(deadline).Seconds()
				m.consumerEventTimeoutExceeded.WithLabelValues(ev, topic).Observe(excess)
			}
		}
		return err
	}
}
