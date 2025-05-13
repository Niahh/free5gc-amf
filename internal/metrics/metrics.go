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

	globalLabels := prometheus.Labels{
		NF_TYPE_LABEL: NF_TYPE_VALUE,
	}

	wrappedReg := prometheus.WrapRegistererWith(globalLabels, reg)

	wrappedReg.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	var amfMetrics []prometheus.Collector

	// Append here the collector you want to register to the prometheus registry
	amfMetrics = append(amfMetrics, nas.GetNasHandlerMetrics(namespace)...)
	amfMetrics = append(amfMetrics, ngap.GetNgapHandlerMetrics(namespace)...)

	amfMetrics = append(amfMetrics, sbi.GetSbiOutboundMetrics(namespace)...)
	amfMetrics = append(amfMetrics, sbi.GetSbiInboundMetrics(namespace)...)

	initMetric(amfMetrics, wrappedReg)

	return reg
}

func initMetric(metrics []prometheus.Collector, reg prometheus.Registerer) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
