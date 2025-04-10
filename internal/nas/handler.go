package nas

import (
	"fmt"
	"github.com/free5gc/amf/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"

	amf_context "github.com/free5gc/amf/internal/context"
	gmm_common "github.com/free5gc/amf/internal/gmm/common"
	"github.com/free5gc/amf/internal/logger"
	nas_metrics "github.com/free5gc/amf/internal/metrics/nas"
	"github.com/free5gc/amf/internal/nas/nas_security"
	"github.com/free5gc/nas"
)

func incrMetrics(msg *nas.Message, metricStatusSuccess *bool) {
	nasMessageType := ""
	if msg.GmmMessage.AuthenticationRequest != nil {
		nasMessageType = "AuthenticationRequest"
	} else if msg.GmmMessage.AuthenticationResponse != nil {
		nasMessageType = "AuthenticationResponse"
	} else if msg.GmmMessage.AuthenticationResult != nil {
		nasMessageType = "AuthenticationResult"
	} else if msg.GmmMessage.AuthenticationFailure != nil {
		nasMessageType = "AuthenticationFailure"
	} else if msg.GmmMessage.AuthenticationReject != nil {
		nasMessageType = "AuthenticationReject"
	} else if msg.GmmMessage.RegistrationRequest != nil {
		nasMessageType = "RegistrationRequest"
	} else if msg.GmmMessage.RegistrationAccept != nil {
		nasMessageType = "RegistrationAccept"
	} else if msg.GmmMessage.RegistrationComplete != nil {
		nasMessageType = "RegistrationComplete"
	} else if msg.GmmMessage.RegistrationReject != nil {
		nasMessageType = "RegistrationReject"
	} else if msg.GmmMessage.ULNASTransport != nil {
		nasMessageType = "ULNASTransport"
	} else if msg.GmmMessage.DLNASTransport != nil {
		nasMessageType = "DLNASTransport"
	} else if msg.GmmMessage.DeregistrationRequestUEOriginatingDeregistration != nil {
		nasMessageType = "DeregistrationRequestUEOriginatingDeregistration"
	} else if msg.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration != nil {
		nasMessageType = "DeregistrationAcceptUEOriginatingDeregistration"
	} else if msg.GmmMessage.DeregistrationRequestUETerminatedDeregistration != nil {
		nasMessageType = "DeregistrationRequestUETerminatedDeregistration"
	} else if msg.GmmMessage.DeregistrationAcceptUETerminatedDeregistration != nil {
		nasMessageType = "DeregistrationAcceptUETerminatedDeregistration"
	} else if msg.GmmMessage.ServiceRequest != nil {
		nasMessageType = "ServiceRequest"
	} else if msg.GmmMessage.ServiceAccept != nil {
		nasMessageType = "ServiceAccept"
	} else if msg.GmmMessage.ServiceReject != nil {
		nasMessageType = "ServiceReject"
	} else if msg.GmmMessage.ConfigurationUpdateCommand != nil {
		nasMessageType = "ConfigurationUpdateCommand"
	} else if msg.GmmMessage.ConfigurationUpdateComplete != nil {
		nasMessageType = "ConfigurationUpdateComplete"
	} else if msg.GmmMessage.IdentityRequest != nil {
		nasMessageType = "IdentityRequest"
	} else if msg.GmmMessage.IdentityResponse != nil {
		nasMessageType = "IdentityResponse"
	} else if msg.GmmMessage.Notification != nil {
		nasMessageType = "Notification"
	} else if msg.GmmMessage.NotificationResponse != nil {
		nasMessageType = "NotificationResponse"
	} else if msg.GmmMessage.SecurityModeCommand != nil {
		nasMessageType = "SecurityModeCommand"
	} else if msg.GmmMessage.SecurityModeComplete != nil {
		nasMessageType = "SecurityModeComplete"
	} else if msg.GmmMessage.SecurityModeReject != nil {
		nasMessageType = "SecurityModeReject"
	} else if msg.GmmMessage.SecurityProtected5GSNASMessage != nil {
		nasMessageType = "SecurityProtected5GSNASMessage"
	} else if msg.GmmMessage.Status5GMM != nil {
		nasMessageType = "Status5GMM"
	}

	metricStatus := metrics.FailureMetric
	if metricStatusSuccess != nil && *metricStatusSuccess {
		metricStatus = metrics.SuccessMetric
	}
	nas_metrics.NasMsgRcvCounter.With(prometheus.Labels{"name": nasMessageType, "status": metricStatus}).Inc()
}

func HandleNAS(ranUe *amf_context.RanUe, procedureCode int64, nasPdu []byte, initialMessage bool) {
	amfSelf := amf_context.GetSelf()

	if ranUe == nil {
		logger.NasLog.Error("RanUe is nil")
		return
	}

	if nasPdu == nil {
		ranUe.Log.Error("nasPdu is nil")
		return
	}

	if ranUe.AmfUe == nil {
		// Only the New created RanUE will have no AmfUe in it

		if ranUe.HoldingAmfUe != nil && !ranUe.HoldingAmfUe.CmConnect(ranUe.Ran.AnType) {
			// If the UE is CM-IDLE, there is no RanUE in AmfUe, so here we attach new RanUe to AmfUe.
			gmm_common.AttachRanUeToAmfUeAndReleaseOldIfAny(ranUe.HoldingAmfUe, ranUe)
			ranUe.HoldingAmfUe = nil
		} else {
			// Assume we have an existing UE context in CM-CONNECTED state. (RanUe <-> AmfUe)
			// We will release it if the new UE context has a valid security context(Authenticated) in line 50.
			ranUe.AmfUe = amfSelf.NewAmfUe("")
			gmm_common.AttachRanUeToAmfUeAndReleaseOldIfAny(ranUe.AmfUe, ranUe)
		}
	}

	msg, integrityProtected, err := nas_security.Decode(ranUe.AmfUe, ranUe.Ran.AnType, nasPdu, initialMessage)
	if err != nil {
		ranUe.AmfUe.NASLog.Errorln(err)
		return
	}

	// [todo] Here I can retrieve the msg object, that contains the list of possible nas type message as pointers
	// I could make a function to retrieve the nas message type from it.

	metricStatusOk := true
	defer incrMetrics(msg, &metricStatusOk)

	ranUe.AmfUe.NasPduValue = nasPdu
	ranUe.AmfUe.MacFailed = !integrityProtected

	if ranUe.AmfUe.SecurityContextIsValid() && ranUe.HoldingAmfUe != nil {
		gmm_common.ClearHoldingRanUe(ranUe.HoldingAmfUe.RanUe[ranUe.Ran.AnType])
		ranUe.HoldingAmfUe = nil
	}

	if errDispatch := Dispatch(ranUe.AmfUe, ranUe.Ran.AnType, procedureCode, msg); errDispatch != nil {
		ranUe.AmfUe.NASLog.Errorf("Handle NAS Error: %v", errDispatch)
		metricStatusOk = false
	}
}

// Get5GSMobileIdentityFromNASPDU is used to find MobileIdentity from plain nas
// return value is: mobileId, mobileIdType, err
func GetNas5GSMobileIdentity(gmmMessage *nas.GmmMessage) (string, string, error) {
	var err error
	var mobileId, mobileIdType string

	if gmmMessage.GmmHeader.GetMessageType() == nas.MsgTypeRegistrationRequest {
		mobileId, mobileIdType, err = gmmMessage.RegistrationRequest.MobileIdentity5GS.GetMobileIdentity()
	} else if gmmMessage.GmmHeader.GetMessageType() == nas.MsgTypeServiceRequest {
		mobileId, mobileIdType, err = gmmMessage.ServiceRequest.TMSI5GS.Get5GSTMSI()
	} else {
		err = fmt.Errorf("gmmMessageType: [%d] is not RegistrationRequest or ServiceRequest",
			gmmMessage.GmmHeader.GetMessageType())
	}
	return mobileId, mobileIdType, err
}
