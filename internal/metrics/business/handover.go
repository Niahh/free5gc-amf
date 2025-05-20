package businness

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	handoverInProgressGauge *prometheus.GaugeVec
	handoverEventCounter    *prometheus.CounterVec
	handoverDuration        *prometheus.HistogramVec
)

func GetHandoverHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	handoverInProgressGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      HANDOVER_IN_PROGRESS_GAUGE_NAME,
			Help:      HANDOVER_IN_PROGRESS_GAUGE_DESC,
		},
		[]string{HANDOVER_TYPE_LABEL},
	)

	handoverEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      HANDOVER_EVENT_COUNTER_NAME,
			Help:      HANDOVER_EVENT_COUNTER_DESC,
		},
		[]string{HANDOVER_TYPE_LABEL, HANDOVER_EVENT_LABEL, HANDOVER_CAUSE_LABEL},
	)

	handoverDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      HANDOVER_DURATION_HISTOGRAM_NAME,
			Help:      HANDOVER_DURATION_HISTOGRAM_DESC,
			Buckets: []float64{
				0.0001,
				0.0050,
				0.0100,
				0.0200,
				0.0500,
				0.1000,
			},
		},
		[]string{HANDOVER_TYPE_LABEL, HANDOVER_EVENT_LABEL},
	)

	collectors = append(collectors, handoverInProgressGauge)
	collectors = append(collectors, handoverEventCounter)
	collectors = append(collectors, handoverDuration)

	return collectors
}

func IncrHoInProgressGauge(hoType string) {
	handoverInProgressGauge.With(prometheus.Labels{
		HANDOVER_TYPE_LABEL: hoType,
	}).Inc()
}

func DecrHoInProgressGauge(hoType string) {
	handoverInProgressGauge.With(prometheus.Labels{
		HANDOVER_TYPE_LABEL: hoType,
	}).Dec()
}

func IncrHoEventDurationCounter(handoverType string, handoverEvent string, hoStartTime time.Time) {
	duration := time.Since(hoStartTime).Seconds()

	handoverDuration.With(prometheus.Labels{
		HANDOVER_TYPE_LABEL:  handoverType,
		HANDOVER_EVENT_LABEL: handoverEvent,
	}).Observe(duration)
}

func IncrHoEventCounter(handoverType string, handoverEvent string, handoverCause string, hoStartTime time.Time) {

	//hoEvent := utils.ReadStringPtr(handoverEvent)
	//hoCause := utils.ReadStringPtr(handoverCause)

	handoverEventCounter.With(prometheus.Labels{
		HANDOVER_TYPE_LABEL:  handoverType,
		HANDOVER_EVENT_LABEL: handoverEvent,
		HANDOVER_CAUSE_LABEL: handoverCause,
	}).Inc()

	switch {
	case handoverEvent == HANDOVER_EVENT_ATTEMPT_VALUE:
		IncrHoInProgressGauge(handoverType)
	case handoverEvent == HANDOVER_EVENT_FAILURE_VALUE:
		fallthrough
	case handoverEvent == HANDOVER_EVENT_SUCCESS_VALUE:
		IncrHoEventDurationCounter(handoverType, handoverEvent, hoStartTime)
		DecrHoInProgressGauge(handoverType)
	}

}
