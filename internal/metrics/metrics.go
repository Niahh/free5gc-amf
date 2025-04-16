// Package metrics sets and initializes Prometheus metrics.
package metrics

import (
	"github.com/free5gc/amf/internal/metrics/nas"
	"github.com/free5gc/amf/internal/metrics/ngap"
	"github.com/free5gc/amf/internal/metrics/sbi"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Init initializes all Prometheus metrics
func Init(cfg *factory.Config) *prometheus.Registry {
	reg := prometheus.NewRegistry()

	namespace := cfg.GetMetricsNamespace()
	reg.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	var amfMetrics []prometheus.Collector

	amfMetrics = append(amfMetrics, sbi.GetCommunicationServiceMetrics(namespace)...)
	amfMetrics = append(amfMetrics, ngap.GetHandoverRequestCounter(namespace)...)
	amfMetrics = append(amfMetrics, ngap.GetNgapHandlerMetrics(namespace)...)
	amfMetrics = append(amfMetrics, nas.GetNasHandlerMetrics(namespace)...)

	initMetric(amfMetrics, reg)

	return reg
}

func initMetric(metrics []prometheus.Collector, reg *prometheus.Registry) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
