package nas

import (
	"github.com/free5gc/amf/internal/metrics/utils"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
)

const (
	NasMsgRcvCounterName = "nas_msg_received_total"
	NasMsgRcvCounterDesc = "Total number of received NAS message by the AMF"

	NasMsgSentCounterName = "nas_msg_sent_total"
	NasMsgSentCounterDesc = "Total number of NAS message sent by the AMF"
)

// These values are tied to the metrics NAS message type
const (
	AUTHENTICATION_REQUEST                               = "AuthenticationRequest"
	AUTHENTICATION_RESPONSE                              = "AuthenticationResponse"
	AUTHENTICATION_RESULT                                = "AuthenticationResult"
	AUTHENTICATION_FAILURE                               = "AuthenticationFailure"
	AUTHENTICATION_REJECT                                = "AuthenticationReject"
	REGISTRATION_REQUEST                                 = "RegistrationRequest"
	REGISTRATION_ACCEPT                                  = "RegistrationAccept"
	REGISTRATION_ACCEPT_TIMER                            = "RegistrationAcceptTimer"
	REGISTRATION_COMPLETE                                = "RegistrationComplete"
	REGISTRATION_REJECT                                  = "RegistrationReject"
	UL_NAS_TRANSPORT                                     = "ULNASTransport"
	DL_NAS_TRANSPORT                                     = "DLNASTransport"
	DEREGISTRATION_REQUEST_UE_ORIGINATING_DEREGISTRATION = "DeregistrationRequestUEOriginatingDeregistration"
	DEREGISTRATION_ACCEPT_UE_ORIGINATING_DEREGISTRATION  = "DeregistrationAcceptUEOriginatingDeregistration"
	DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION  = "DeregistrationRequestUETerminatedDeregistration"
	DEREGISTRATION_ACCEPT_UE_TERMINATED_DEREGISTRATION   = "DeregistrationAcceptUETerminatedDeregistration"
	SERVICE_REQUEST                                      = "ServiceRequest"
	SERVICE_ACCEPT                                       = "ServiceAccept"
	SERVICE_REJECT                                       = "ServiceReject"
	CONFIGURATION_UPDATE_COMMAND                         = "ConfigurationUpdateCommand"
	CONFIGURATION_UPDATE_COMMAND_TIMER                   = "ConfigurationUpdateCommandTimer"
	CONFIGURATION_UPDATE_COMPLETE                        = "ConfigurationUpdateComplete"
	IDENTITY_REQUEST                                     = "IdentityRequest"
	IDENTITY_RESPONSE                                    = "IdentityResponse"
	NOTIFICATION                                         = "Notification"
	NOTIFICATION_TIMER                                   = "NotificationTimer"
	NOTIFICATION_RESPONSE                                = "NotificationResponse"
	SECURITY_MODE_COMMAND                                = "SecurityModeCommand"
	SECURITY_MODE_COMPLETE                               = "SecurityModeComplete"
	SECURITY_MODE_REJECT                                 = "SecurityModeReject"
	SECURITY_PROTECTED_5GS_NAS_MESSAGE                   = "SecurityProtected5GSNASMessage"
	STATUS_5GMM                                          = "Status5GMM"
)

// Additional error causes
const (
	RAN_UE_NIL_ERR      = "RanUe is nil"
	AMF_UE_NIL_ERR      = "AmfUe is nil"
	RAN_NIL_ERR         = "Ran is nil"
	NAS_PDU_NIL_ERR     = "nasPdu is nil"
	AUTH_CTX_UE_NIL_ERR = "Authentication Context of UE is nil"
	NAS_MSG_BUILD_ERR   = "Could not build NAS message"
	DECODE_NAS_MSG_ERR  = "Could not decode NAS message"
	AUSF_AUTH_ERR       = "Ausf Authentication Failure"
	HRES_AUTH_ERR       = "HRES* validation failure"
)

var (
	NasMsgRcvCounter  prometheus.CounterVec
	NasMsgSentCounter prometheus.CounterVec
)

var (
	suffixRe = regexp.MustCompile(`\s*\(\d+\)$`)
)

func GetNasHandlerMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	NasMsgRcvCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      NasMsgRcvCounterName,
			Help:      NasMsgRcvCounterDesc,
		},
		[]string{"name", "status", "cause"},
	)

	metrics = append(metrics, NasMsgRcvCounter)

	NasMsgSentCounter = *prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      NasMsgSentCounterName,
			Help:      NasMsgSentCounterDesc,
		},
		[]string{"name", "status", "cause"},
	)

	metrics = append(metrics, NasMsgSentCounter)

	return metrics
}

func removeDigitSuffix(s string) string {
	return suffixRe.ReplaceAllString(s, "")
}

func IncrMetricsRcvNasMsg(msg *nas.Message, isStatusSuccess *bool, cause *string) {

	nasMessageIe := getMessageStrFromGmmMessage(msg)
	metricCause := removeDigitSuffix(nasMessage.Cause5GMMToString(nasMessageIe.cause.Octet))
	metricStatus := utils.FailureMetric

	if cause != nil && *cause != "" {
		metricCause = *cause
	}

	if isStatusSuccess != nil && *isStatusSuccess {
		metricStatus = utils.SuccessMetric
	}
	NasMsgRcvCounter.With(prometheus.Labels{"name": nasMessageIe.nasMessageType, "status": metricStatus, "cause": metricCause}).Inc()
}

func IncrMetricsSentNasMsgs(msgType string, isStatusSuccess *bool, cause5GMM uint8, otherCause *string) {

	errCause := ""

	if cause5GMM != 0 {
		errCause = removeDigitSuffix(nasMessage.Cause5GMMToString(cause5GMM))
	} else if otherCause != nil {
		errCause = *otherCause
	}

	metricStatus := utils.FailureMetric
	if isStatusSuccess != nil && *isStatusSuccess {
		metricStatus = utils.SuccessMetric
	}
	NasMsgSentCounter.With(prometheus.Labels{"name": msgType, "status": metricStatus, "cause": errCause}).Inc()
}

type IeFromGmmMessage struct {
	nasMessageType string
	cause          nasType.Cause5GMM
}

func getMessageStrFromGmmMessage(msg *nas.Message) IeFromGmmMessage {

	ie := IeFromGmmMessage{nasMessageType: "Unknown gmm message"}

	if msg == nil || msg.GmmMessage == nil {
		return ie
	}

	if msg.GmmMessage.AuthenticationRequest != nil {
		ie.nasMessageType = AUTHENTICATION_REQUEST
	} else if msg.GmmMessage.AuthenticationResponse != nil {
		ie.nasMessageType = AUTHENTICATION_RESPONSE
	} else if msg.GmmMessage.AuthenticationResult != nil {
		ie.nasMessageType = AUTHENTICATION_RESULT
	} else if msg.GmmMessage.AuthenticationFailure != nil {
		ie.nasMessageType = AUTHENTICATION_FAILURE
		ie.cause = msg.GmmMessage.AuthenticationFailure.Cause5GMM
	} else if msg.GmmMessage.AuthenticationReject != nil {
		ie.nasMessageType = AUTHENTICATION_REJECT
	} else if msg.GmmMessage.RegistrationRequest != nil {
		ie.nasMessageType = REGISTRATION_REQUEST
	} else if msg.GmmMessage.RegistrationAccept != nil {
		ie.nasMessageType = REGISTRATION_ACCEPT
	} else if msg.GmmMessage.RegistrationComplete != nil {
		ie.nasMessageType = REGISTRATION_COMPLETE
	} else if msg.GmmMessage.RegistrationReject != nil {
		ie.nasMessageType = REGISTRATION_REJECT
		ie.cause = msg.GmmMessage.RegistrationReject.Cause5GMM
	} else if msg.GmmMessage.ULNASTransport != nil {
		ie.nasMessageType = UL_NAS_TRANSPORT
	} else if msg.GmmMessage.DLNASTransport != nil {
		ie.nasMessageType = DL_NAS_TRANSPORT
		ie.cause = *msg.GmmMessage.DLNASTransport.Cause5GMM
	} else if msg.GmmMessage.DeregistrationRequestUEOriginatingDeregistration != nil {
		ie.nasMessageType = DEREGISTRATION_REQUEST_UE_ORIGINATING_DEREGISTRATION
	} else if msg.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration != nil {
		ie.nasMessageType = DEREGISTRATION_ACCEPT_UE_ORIGINATING_DEREGISTRATION
	} else if msg.GmmMessage.DeregistrationRequestUETerminatedDeregistration != nil {
		ie.nasMessageType = DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION
		ie.cause = *msg.GmmMessage.DeregistrationRequestUETerminatedDeregistration.Cause5GMM
	} else if msg.GmmMessage.DeregistrationAcceptUETerminatedDeregistration != nil {
		ie.nasMessageType = DEREGISTRATION_ACCEPT_UE_TERMINATED_DEREGISTRATION
	} else if msg.GmmMessage.ServiceRequest != nil {
		ie.nasMessageType = SERVICE_REQUEST
	} else if msg.GmmMessage.ServiceAccept != nil {
		ie.nasMessageType = SERVICE_ACCEPT
	} else if msg.GmmMessage.ServiceReject != nil {
		ie.nasMessageType = SERVICE_REJECT
		ie.cause = msg.GmmMessage.ServiceReject.Cause5GMM
	} else if msg.GmmMessage.ConfigurationUpdateCommand != nil {
		ie.nasMessageType = CONFIGURATION_UPDATE_COMMAND
	} else if msg.GmmMessage.ConfigurationUpdateComplete != nil {
		ie.nasMessageType = CONFIGURATION_UPDATE_COMPLETE
	} else if msg.GmmMessage.IdentityRequest != nil {
		ie.nasMessageType = IDENTITY_REQUEST
	} else if msg.GmmMessage.IdentityResponse != nil {
		ie.nasMessageType = IDENTITY_RESPONSE
	} else if msg.GmmMessage.Notification != nil {
		ie.nasMessageType = NOTIFICATION
	} else if msg.GmmMessage.NotificationResponse != nil {
		ie.nasMessageType = NOTIFICATION_RESPONSE
	} else if msg.GmmMessage.SecurityModeCommand != nil {
		ie.nasMessageType = SECURITY_MODE_COMMAND
	} else if msg.GmmMessage.SecurityModeComplete != nil {
		ie.nasMessageType = SECURITY_MODE_COMPLETE
	} else if msg.GmmMessage.SecurityModeReject != nil {
		ie.nasMessageType = SECURITY_MODE_REJECT
		ie.cause = msg.GmmMessage.SecurityModeReject.Cause5GMM
	} else if msg.GmmMessage.SecurityProtected5GSNASMessage != nil {
		ie.nasMessageType = SECURITY_PROTECTED_5GS_NAS_MESSAGE
	} else if msg.GmmMessage.Status5GMM != nil {
		ie.nasMessageType = STATUS_5GMM
		ie.cause = msg.GmmMessage.Status5GMM.Cause5GMM
	}

	return ie
}
