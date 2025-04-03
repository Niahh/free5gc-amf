package ngap

import "github.com/prometheus/client_golang/prometheus"

const(
	ngapMsgCounterName = "ngap_msg_received_total"
	ngapMsgCounterDescription = "Total number of received NGAP message by the AMF "
)

var (
	NgapMsgCounter prometheus.CounterVec
)

func GetNgapHandlerMetrics(namespace string) []prometheus.Collector{
	var metrics []prometheus.Collector

	NgapMsgCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: ngapMsgCounterName,
			Help: ngapMsgCounterDescription,
		}, 
		[]string{"name", "status"},
	)

	metrics = append(metrics, NgapMsgCounter)

	return metrics
}
