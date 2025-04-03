package ngap

import "github.com/prometheus/client_golang/prometheus"

const (
	HandoverRequestCounterName = "handover_request_received_total"
	HandoverRequestSuccessfulCounterName = "handover_request_acknowledge_total"
	PathSwitchRequestFailureCounterName = "path_switch_request_failure_total"
)

// To do, rename te metrics to -> xn handover ... (req)
// To do, rename te metrics to -> xn handover ... (succ)
// To do, rename te metrics to -> xn handover ... (failed)

// Add metrics for n2
// n2 handover req
// n2 handover failed ... (check nokia's metric sheet)

var (
	PathSwitchRequestCounter prometheus.Counter
	PathSwitchRequestAcknowledgeCounter prometheus.Counter
	PathSwitchRequestFailureCounter prometheus.Counter
)

func GetHandoverRequestCounter(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	PathSwitchRequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: HandoverRequestCounterName,
			Help: "Show the total number of PathSwitchRequest NGAP received by the AMF",
		},
	)

	metrics = append(metrics, PathSwitchRequestCounter)

	PathSwitchRequestAcknowledgeCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: HandoverRequestSuccessfulCounterName,
			Help: "Show the total number of PathSwitchRequest NGAP Acknowledged by the AMF",
		},
	)

	metrics = append(metrics, PathSwitchRequestAcknowledgeCounter)

	PathSwitchRequestFailureCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: PathSwitchRequestFailureCounterName,
			Help: "Show the total number of PathSwitchRequest NGAP that did not succeed handled by the AMF",
		},
	)

	metrics = append(metrics, PathSwitchRequestFailureCounter)

	return metrics
}