package businness

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	GmmStateGauge        *prometheus.GaugeVec
	GmmTransitionCounter *prometheus.CounterVec
	GmmStateDurationHist *prometheus.HistogramVec
)

func GetGMMStatesHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	// Gauge for number of UEs in each GMM state, labeled by state name (and access type)
	GmmStateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: GMM_STATE_GAUGE_NAME,
			Help: GMM_STATE_GAUGE_DESC,
		},
		[]string{GMM_STATE_ACCESS_LABEL, GMM_STATE_LABEL},
	)

	collectors = append(collectors, GmmStateGauge)

	// Counter for GMM state transitions, labeled by from->to states
	GmmTransitionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: GMM_TRANSITION_COUNTER_NAME,
			Help: GMM_TRANSITION_COUNTER_DESC,
		},
		[]string{GMM_STATE_FROM_STATE_LABEL, GMM_STATE_TO_STATE_LABEL},
	)

	collectors = append(collectors, GmmTransitionCounter)

	// Histogram for time spent in a GMM state (seconds)
	GmmStateDurationHist = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: GMM_DURATION_HISTOGRAM_NAME,
			Help: GMM_DURATION_HISTOGRAM_DESC,
			Buckets: []float64{
				0.0001,
				0.0050,
				0.0100,
				0.0200,
				0.0250,
				0.0500,
			},
		},
		[]string{GMM_STATE_ACCESS_LABEL, GMM_STATE_LABEL},
	)

	collectors = append(collectors, GmmStateDurationHist)

	return collectors
}

func IncrGmmTransitionCounter(fromState string, toState string) {
	GmmTransitionCounter.With(prometheus.Labels{
		GMM_STATE_FROM_STATE_LABEL: fromState,
		GMM_STATE_TO_STATE_LABEL:   toState,
	}).Add(1)
}

func IncrGmmStateGauge(accessType string, state string) {
	GmmStateGauge.With(prometheus.Labels{
		GMM_STATE_ACCESS_LABEL: accessType,
		GMM_STATE_LABEL:        state,
	}).Add(1)
}

func DecrGmmStateGauge(accessType string, state string, enterTime time.Time) {
	dur := time.Since(enterTime).Seconds()

	GmmStateGauge.With(prometheus.Labels{
		GMM_STATE_ACCESS_LABEL: accessType,
		GMM_STATE_LABEL:        state,
	}).Dec()

	ObserveGmmTransitionDuration(accessType, state, dur)
}

func ObserveGmmTransitionDuration(accessType string, gmmState string, duration float64) {
	GmmStateDurationHist.With(prometheus.Labels{
		GMM_STATE_ACCESS_LABEL: accessType,
		GMM_STATE_LABEL:        gmmState,
	}).Observe(duration)
}
