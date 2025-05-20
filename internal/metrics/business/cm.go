package businness

import (
	"github.com/free5gc/openapi/models"
	"github.com/prometheus/client_golang/prometheus"
)

// ueCmStateGauge Connection Management different state (either cm-idle or cm-connected) Gauge
var (
	ueCmStateGauge *prometheus.GaugeVec
)

func GetUECMHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	ueCmStateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      UE_CM_STATE_GAUGE_NAME,
			Help:      UE_CM_STATE_GAUGE_DESC,
		},
		[]string{UE_CM_ACCESS_STATE_LABEL, UE_CM_STATE_LABEL},
	)

	collectors = append(collectors, ueCmStateGauge)

	return collectors
}

func IncrUeCmConnectedStateGauge(accessType models.AccessType) {
	ueCmStateGauge.With(prometheus.Labels{
		UE_CM_ACCESS_STATE_LABEL: string(accessType),
		UE_CM_STATE_LABEL:        UE_CM_CONNECTED_VALUE,
	}).Inc()
}

func IncrUeCmIdleStateGauge(accessType models.AccessType) {
	ueCmStateGauge.With(prometheus.Labels{
		UE_CM_ACCESS_STATE_LABEL: string(accessType),
		UE_CM_STATE_LABEL:        UE_CM_IDLE_VALUE,
	}).Inc()
}

func DecrUeCmIdleStateGauge(accessType models.AccessType) {
	ueCmStateGauge.With(prometheus.Labels{
		UE_CM_ACCESS_STATE_LABEL: string(accessType),
		UE_CM_STATE_LABEL:        UE_CM_IDLE_VALUE,
	}).Dec()
}

func DecrUeCmConnectedStateGauge(accessType models.AccessType) {
	ueCmStateGauge.With(prometheus.Labels{
		UE_CM_ACCESS_STATE_LABEL: string(accessType),
		UE_CM_STATE_LABEL:        UE_CM_CONNECTED_VALUE,
	}).Dec()
}
