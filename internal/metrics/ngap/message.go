package ngap

import (
	metric_utils "github.com/free5gc/amf/internal/metrics/utils"
	"github.com/free5gc/ngap/ngapType"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MsgRcvCounter  *prometheus.CounterVec
	MsgSentCounter *prometheus.CounterVec
)

func GetNgapHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	MsgRcvCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      MsgRcvCounterName,
			Help:      MsgRcvCounterDesc,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, MsgRcvCounter)

	MsgSentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      MsgSentCounterName,
			Help:      MsgSentCounterDesc,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, MsgSentCounter)

	return collectors
}

func IncrMetricsRcvMsg(msgType string, metricStatusSuccess *bool, syntaxCause *ngapType.Cause) {

	msgCause := ""

	if syntaxCause != nil && syntaxCause.Present != 0 {
		msgCause = GetCauseErrorStr(syntaxCause)
	}

	MsgRcvCounter.With(prometheus.Labels{
		NAME_LABEL:   msgType,
		STATUS_LABEL: getStatus(metricStatusSuccess),
		CAUSE_LABEL:  msgCause,
	}).Add(1)
}

func IncrMetricsSentMsg(msgType string, metricStatusSuccess *bool, syntaxCause ngapType.Cause, otherCause *string) {

	msgCause := ""

	if syntaxCause.Present != 0 {
		msgCause = GetCauseErrorStr(&syntaxCause)
	} else if otherCause != nil {
		msgCause = *otherCause
	}

	MsgSentCounter.With(prometheus.Labels{
		NAME_LABEL:   msgType,
		STATUS_LABEL: getStatus(metricStatusSuccess),
		CAUSE_LABEL:  msgCause,
	}).Add(1)
}

func getStatus(metricStatusSuccess *bool) string {
	if metricStatusSuccess != nil && *metricStatusSuccess {
		return metric_utils.SuccessMetric
	}
	return metric_utils.FailureMetric
}
