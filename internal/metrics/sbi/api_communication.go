package sbi

import "github.com/prometheus/client_golang/prometheus"

const (
	UEContextTransferCounterName = "communication_ue_context_transfer_handled_total"
)

var (
	UEContextTransferCounter prometheus.CounterVec
)

func GetCommunicationServiceMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	UEContextTransferCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: UEContextTransferCounterName,
			Help: "Show the total number of UEContextTransfer calls handled by the AMF, could be filtered by StatusCode",
		}, []string{"StatusCode"},
	)

	metrics = append(metrics, UEContextTransferCounter)

	return metrics
}
