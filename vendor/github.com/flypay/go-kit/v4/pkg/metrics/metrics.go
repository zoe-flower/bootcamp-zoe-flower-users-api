package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterOrExisting will register the metric or if it already exists return the existing metric
func RegisterOrExisting(col prometheus.Collector) prometheus.Collector {
	if err := prom.Register(col); err != nil {
		if regErr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return regErr.ExistingCollector
		}
		// Not ideal to panic but it's a bit like MustRegister
		panic(err)
	}

	return col
}

func Labels() prometheus.Labels {
	return emptyLabels
}

// ConstMetricLabels contains labels and values all metrics should contain such as application name
var (
	ConstMetricLabels prometheus.Labels
	emptyLabels       = prometheus.Labels{}
)

// prom is the prometheus registry to use, defaults to the
// prometheus global provided registry
var prom prometheus.Registerer

func init() {
	RegisterDefaultPrometheusRegisterer()
}

func RegisterDefaultPrometheusRegisterer() {
	RegisterRegistry(prometheus.DefaultRegisterer)
}

// RegisterRegistry will register a new prometheus registerer
func RegisterRegistry(registry prometheus.Registerer) {
	prom = registry
}
