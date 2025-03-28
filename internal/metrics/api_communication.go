package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	UEContextTransferSuccessfulCounterName = "communication_ue_context_transfer_success_total"
	UEContextTransferErrorCounterName = "communication_ue_context_transfer_error_total"
)

var (
	UEContextTransferSuccessfulCounter prometheus.Counter
	UEContextTransferErrorCounter prometheus.CounterVec
)

func GetCommunicationServiceMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	UEContextTransferSuccessfulCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: UEContextTransferSuccessfulCounterName,
			Help: "Show the number of successful UEContextTransfer this AMF has done",
		},
	)

	metrics = append(metrics, UEContextTransferSuccessfulCounter)

	UEContextTransferErrorCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name: UEContextTransferErrorCounterName,
			Help: "Show the number of UEContextTransfer this AMF could not carried out due to errors",
		}, []string{"code"},
	)

	metrics = append(metrics, UEContextTransferErrorCounter)

	return metrics
}
