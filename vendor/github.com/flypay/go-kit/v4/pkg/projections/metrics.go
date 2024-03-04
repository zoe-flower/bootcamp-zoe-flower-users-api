package projections

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/flypay/go-kit/v4/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type dynamoMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func registerDynamoMetrics() *dynamoMetrics {
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem:   "aws_dynamodb",
		Name:        "requests_total",
		ConstLabels: ConstMetricLabels,
	}, []string{"type", "status", "errCode"})
	c = metrics.RegisterOrExisting(c).(*prometheus.CounterVec)

	obs := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   "aws_dynamodb",
		Name:        "request_duration_milliseconds",
		Buckets:     []float64{0, 50, 100, 250, 500, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 15000},
		ConstLabels: ConstMetricLabels,
	}, []string{"type", "status"})
	obs = metrics.RegisterOrExisting(obs).(*prometheus.HistogramVec)

	return &dynamoMetrics{
		RequestsTotal:   c,
		RequestDuration: obs,
	}
}

func (m *dynamoMetrics) MetricsReq(startTime time.Time, reqType string, err error) {
	elapsed := time.Since(startTime).Milliseconds()
	status := "success"
	errCode := ""

	if err != nil {
		status = "failed"
		if aerr, ok := err.(awserr.Error); ok { // nolint:errorlint
			errCode = aerr.Code()
		}
	}

	m.RequestsTotal.WithLabelValues(reqType, status, errCode).Inc()
	m.RequestDuration.WithLabelValues(reqType, status).Observe(float64(elapsed))
}

// ConstMetricLabels contains labels and values all metrics should contain such as application name
var ConstMetricLabels prometheus.Labels
