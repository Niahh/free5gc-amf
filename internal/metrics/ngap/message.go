package ngap

import (
	metric_utils "github.com/free5gc/amf/internal/metrics/utils"
	"github.com/free5gc/ngap/ngapType"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ngapMsgRcvCounterName = "ngap_msg_received_total"
	ngapMsgRcvCounterDesc = "Total number of received NGAP message by the AMF "

	ngapMsgSentCounterName = "ngap_msg_sent_total"
	ngapMsgSentCounterDesc = "Total number of NGAP message sent by the AMF "
)

const (
	NG_SETUP_RESPONSE                          = "NGSetupResponse"
	NG_SETUP_FAILURE                           = "NGSetupFailure"
	NG_RESET                                   = "NGReset"
	NG_RESET_ACKNOWLEDGE                       = "NGResetAcknowledge"
	DOWNLINK_NAS_TRANSPORT                     = "DownlinkNasTransport"
	PDUSESSION_RESOURCE_RELEASE_COMMAND        = "PDUSessionResourceReleaseCommand"
	UE_CONTEXT_RELEASE_COMMAND                 = "UEContextReleaseCommand"
	ERROR_INDICATION                           = "ErrorIndication"
	UE_RADIO_CAPABILITY_CHECK_REQUEST          = "UERadioCapabilityCheckRequest"
	HANDOVER_CANCEL_ACKNOWLEDGE                = "HandoverCancelAcknowledge"
	PDUSESSION_RESOURCE_SETUP_REQUEST          = "PDUSessionResourceSetupRequest"
	PDUSESSION_RESOURCE_MODIFY_CONFIRM         = "PDUSessionResourceModifyConfirm"
	PDUSESSION_RESOURCE_MODIFY_REQUEST         = "PDUSessionResourceModifyRequest"
	INITIAL_CONTEXT_SETUP_REQUEST              = "InitialContextSetupRequest"
	UE_CONTEXT_MODIFICATION_REQUEST            = "UEContextModificationRequest"
	HANDOVER_COMMAND                           = "HandoverCommand"
	HANDOVER_PREPARATION_FAILURE               = "HandoverPreparationFailure"
	HANDOVER_REQUEST                           = "HandoverRequest"
	PATH_SWITCH_REQUEST_ACKNOWLEDGE            = "PathSwitchRequestAcknowledge"
	PATH_SWITCH_REQUEST_FAILURE                = "PathSwitchRequestFailure"
	DOWNLINK_RAN_STATUS_TRANSFER               = "DownlinkRanStatusTransfer"
	PAGING                                     = "Paging"
	REROUTE_NAS_REQUEST                        = "RerouteNasRequest"
	RAN_CONFIGURATION_UPDATE_ACKNOWLEDGE       = "RanConfigurationUpdateAcknowledge"
	RAN_CONFIGURATION_UPDATE_FAILURE           = "RanConfigurationUpdateFailure"
	AMF_STATUS_INDICATION                      = "AMFStatusIndication"
	OVERLOAD_START                             = "OverloadStart"
	OVERLOAD_STOP                              = "OverloadStop"
	DOWNLINK_RAN_CONFIGURATION_TRANSFER        = "DownlinkRanConfigurationTransfer"
	DOWNLINK_NON_UE_ASSOCIATED_NRPPA_TRANSPORT = "DownlinkNonUEAssociatedNRPPATransport"
	DEACTIVATE_TRACE                           = "DeactivateTrace"
	AMF_CONFIGURATION_UPDATE                   = "AMFConfigurationUpdate"
	DOWNLINK_UE_ASSOCIATED_NRPPA_TRANSPORT     = "DownlinkUEAssociatedNRPPaTransport"
	LOCATION_REPORTING_CONTROL                 = "LocationReportingControl"
	UE_TNLA_BINDING_RELEASE_REQUEST            = "UETNLABindingReleaseRequest"
)

// Additional error causes
const (
	RAN_UE_NIL_ERR                                = "RanUe is nil"
	AMF_UE_NIL_ERR                                = "AmfUe is nil"
	RAN_NIL_ERR                                   = "Ran is nil"
	GUAMI_LIST_OOR_ERR                            = "GUAMI List out of range"
	AMF_TRAFFIC_LOAD_REDUCTION_INDICATION_OOO_ERR = "AmfTrafficLoadReductionIndication out of range (should be 1 ~ 99)"
	NSSAI_LIST_OOR_ERR                            = "NSSAI List out of range"
	AOI_LIST_OOR_ERR                              = "AOI List out of range"
	LOCATION_REPORTING_REFERENCE_ID_OOR_ERR       = "LocationReportingReferenceIDToBeCancelled out of range (should be 1 ~ 64)"
	NRPPA_LEN_ZERO_ERR                            = "length of NRPPA-PDU is 0"
	NGAP_MSG_NIL_ERR                              = "Ngap Message is nil"
	NGAP_MSG_BUILD_ERR                            = "Could not build NAS message"
	CAUSE_NIL_ERR                                 = "Cause present is nil"
	SOURCE_UE_NIL_ERR                             = "SourceUe is nil"
	PDU_SESS_RESOURCE_SWITCH_OOO_ERR              = "Pdu Session Resource Switched List out of range"
	PDU_List_OOR_ERR                              = "Pdu List out of range"
	TARGET_RAN_NIL_ERR                            = "targetRan is nil"
	HANDOVER_REQUIRED_DUP_ERR                     = "Handover Required Duplicated"
	SRC_TO_TARGET_TRANSPARENT_CONTAINER_NIL_ERR   = "Source To Target TransparentContainer is nil"
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

func IncrMetricsRcvNgapMsg(msgType string, metricStatusSuccess *bool, syntaxCause *ngapType.Cause) {

	msgCause := ""

	if syntaxCause != nil && syntaxCause.Present != 0 {
		msgCause = GetCauseErrorStr(syntaxCause)
	}

	NgapMsgRcvCounter.With(prometheus.Labels{"name": msgType, "status": getStatus(metricStatusSuccess), "cause": msgCause}).Add(1)
}

func IncrMetricsSentNgapMsg(msgType string, metricStatusSuccess *bool, syntaxCause ngapType.Cause, otherCause *string) {

	msgCause := ""

	if syntaxCause.Present != 0 {
		msgCause = GetCauseErrorStr(&syntaxCause)
	} else if otherCause != nil {
		msgCause = *otherCause
	}

	NgapMsgSentCounter.With(prometheus.Labels{"name": msgType, "status": getStatus(metricStatusSuccess), "cause": msgCause}).Add(1)
}

func getStatus(metricStatusSuccess *bool) string {
	if metricStatusSuccess != nil && *metricStatusSuccess {
		return metric_utils.SuccessMetric
	}
	return metric_utils.FailureMetric
}
