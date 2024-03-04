package custom

import (
	flytmetrics "github.com/flypay/go-kit/v4/pkg/metrics"
	"github.com/go-kit/kit/metrics"
	promkit "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNamespace = "flyt"

	// LabelSuccess for messages that were successful
	LabelSuccess = "success"
	// LabelFailed for messages that failed
	LabelFailed = "failed"
	// LabelError for messages that were encountered an error
	LabelError = "error"
)

func newQueueReceiveCounter(serviceName, queueName string) metrics.Counter {
	labels := prometheus.Labels{"service": serviceName, "queue_name": queueName}

	m := flytmetrics.RegisterOrExisting(
		prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   metricNamespace,
			Name:        "queue_receives_total",
			Help:        "The total number of receives on the queue",
			ConstLabels: labels,
		}, []string{"status"})).(*prometheus.CounterVec)

	c := promkit.NewCounter(m)
	c.With("status", LabelSuccess).Add(0)
	c.With("status", LabelError).Add(0)
	return c
}

func newQueueJobsReceivedHistogram(serviceName, queueName string) metrics.Histogram {
	labels := prometheus.Labels{"service": serviceName, "queue_name": queueName}

	m := flytmetrics.RegisterOrExisting(prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   metricNamespace,
		Name:        "queue_jobs_received",
		Help:        "The number of jobs received from the queue",
		ConstLabels: labels,
		Buckets:     []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}, []string{"status"})).(*prometheus.HistogramVec)

	return promkit.NewHistogram(m)
}

func newQueueCapacityGauge(serviceName, queueName string) metrics.Gauge {
	labels := prometheus.Labels{"service": serviceName, "queue_name": queueName}

	m := flytmetrics.RegisterOrExisting(prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "job_pool_capacity",
		Help:        "The number of consumers ready to receive jobs",
		ConstLabels: labels,
	}, []string{})).(*prometheus.GaugeVec)

	return promkit.NewGauge(m)
}
