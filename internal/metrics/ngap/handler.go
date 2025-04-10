package ngap

import "github.com/prometheus/client_golang/prometheus"

const (
	ngapMsgRcvCounterName = "ngap_msg_received_total"
	ngapMsgRcvCounterDesc = "Total number of received NGAP message by the AMF "

	ngapMsgSentCounterName = "ngap_msg_sent_total"
	ngapMsgSentCounterDesc = "Total number of NGAP message sent by the AMF "
)

var (
	NgapMsgRcvCounter  prometheus.CounterVec
	NgapMsgSentCounter prometheus.CounterVec
)

func GetNgapHandlerMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	NgapMsgRcvCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      ngapMsgRcvCounterName,
			Help:      ngapMsgRcvCounterDesc,
		},
		[]string{"name", "status", "cause"},
	)

	metrics = append(metrics, NgapMsgRcvCounter)

	NgapMsgSentCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      ngapMsgSentCounterName,
			Help:      ngapMsgSentCounterDesc,
		},
		[]string{"name", "status", "cause"},
	)

	metrics = append(metrics, NgapMsgSentCounter)

	return metrics
}
