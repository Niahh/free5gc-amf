package businness

import "github.com/prometheus/client_golang/prometheus"

var (
	pduSessionInProgressGauge prometheus.Gauge
	pduSessionEventCounter    *prometheus.CounterVec
)

func GetPDUHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	pduSessionInProgressGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      PDU_SESSION_IN_PROGRESS_GAUGE_NAME,
			Help:      PDU_SESSION_IN_PROGRESS_GAUGE_DESC,
		},
	)

	pduSessionEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      PDU_SESSION_EVENT_COUNTER_NAME,
			Help:      PDU_SESSION_EVENT_COUNTER_DESC,
		},
		[]string{PDU_SESSION_EVENT_LABEL, PDU_SESSION_STATUS_LABEL, PDU_SESSION_CAUSE_LABEL},
	)

	collectors = append(collectors, pduSessionInProgressGauge)
	collectors = append(collectors, pduSessionEventCounter)

	return collectors
}
