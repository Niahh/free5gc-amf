// Package metrics sets and initializes Prometheus metrics.
package metrics

import (
	"github.com/free5gc/amf/pkg/factory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Init initializes all Prometheus metrics
func Init(cfg *factory.Config) *prometheus.Registry {
	reg := prometheus.NewRegistry()

	reg.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	var amfMetrics []prometheus.Collector

	// Append here the collector you want to register to the prometheus registry

	initMetric(amfMetrics, reg)

	return reg
}

func initMetric(metrics []prometheus.Collector, reg *prometheus.Registry) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
