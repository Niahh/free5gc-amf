package nas

import "github.com/prometheus/client_golang/prometheus"

const (
	nasMsgRcvCounterName = "nas_msg_received_total"
	nasMsgRcvCounterDesc = "Total number of received NAS message by the AMF"

	nasMsgSentCounterName = "nas_msg_sent_total"
	nasMsgSentCounterDesc = "Total number of NAS message sent by the AMF"
)

var (
	NasMsgRcvCounter  prometheus.CounterVec
	NasMsgSentCounter prometheus.CounterVec
)

func GetNasHandlerMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	NasMsgRcvCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      nasMsgRcvCounterName,
			Help:      nasMsgRcvCounterDesc,
		},
		[]string{"name", "status"},
	)

	metrics = append(metrics, NasMsgRcvCounter)

	NasMsgSentCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      nasMsgSentCounterName,
			Help:      nasMsgSentCounterDesc,
		},
		[]string{"name", "status"},
	)

	metrics = append(metrics, NasMsgSentCounter)

	return metrics
}
