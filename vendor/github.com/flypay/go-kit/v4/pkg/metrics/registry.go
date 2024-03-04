package metrics

import (
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Metric provides a wrapper around the description and raw metric
// returned from Prometheus.
type Metric struct {
	// Description of the metric name and labels as a string and may change
	// format in different version.
	// Unfortunately we only have access to prometheus.Desc.String() and no other
	// internals at the moment, see below issues:
	// - https://github.com/prometheus/client_golang/issues/516
	// - https://github.com/prometheus/client_golang/issues/222
	Description string

	// The metric returned from Prometheus
	RawMetric dto.Metric
}

// Registry implements prometheus.Registerer interface and provides
// a way to register and collect prometheus metrics
type Registry struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

func (p *Registry) Register(collector prometheus.Collector) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.collectors = append(p.collectors, collector)
	return nil
}

func (p *Registry) MustRegister(collectors ...prometheus.Collector) {
	for _, c := range collectors {
		if err := p.Register(c); err != nil {
			panic(err)
		}
	}
}

func (p *Registry) Unregister(prometheus.Collector) bool {
	return true
}

// Metrics return slightly more useful metric information than 'RawMetrics'
func (p *Registry) Metrics() ([]Metric, error) {
	rawMetrics := p.RawMetrics()

	metrics := make([]Metric, 0, len(rawMetrics))
	for _, rm := range rawMetrics {
		var dtoMetric dto.Metric
		if err := rm.Write(&dtoMetric); err != nil {
			return nil, err
		}
		metrics = append(metrics, Metric{
			Description: rm.Desc().String(),
			RawMetric:   dtoMetric,
		})
	}
	return metrics, nil
}

// RawMetrics collects all the current metrics from the registered collectors
func (p *Registry) RawMetrics() []prometheus.Metric {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan prometheus.Metric)
	go func() {
		for _, c := range p.collectors {
			c.Collect(ch)
		}
		close(ch)
	}()

	var metrics []prometheus.Metric
	for m := range ch {
		metrics = append(metrics, m)
	}
	return metrics
}

// MetricsByName will return the metrics with that name in metrics description
func (p *Registry) MetricsByName(name string) ([]Metric, error) {
	mtrcs, err := p.Metrics()
	if err != nil {
		return nil, err
	}

	var metrics []Metric
	for _, m := range mtrcs {
		if strings.Contains(m.Description, name) {
			metrics = append(metrics, m)
			continue
		}
	}
	return metrics, nil
}
