package nas

import (
	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	metrics_nas "github.com/free5gc/amf/internal/metrics/nas"
	metrics_ngap "github.com/free5gc/amf/internal/metrics/ngap"
	metrics_utils "github.com/free5gc/amf/internal/metrics/utils"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"testing"
)

func getTai() models.Tai {
	tai := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "1",
	}
	return tai
}

func getMockRanUeCtx() *amf_context.RanUe {
	amfSelf := amf_context.GetSelf()
	amfSelf.ServedGuamiList = []models.Guami{{
		PlmnId: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
		AmfId: "cafe00",
	}}

	ue := new(amf_context.RanUe)
	ue.Ran = new(amf_context.AmfRan)
	ue.Ran.AnType = models.AccessType__3_GPP_ACCESS
	ue.Ran.Log = logger.NasLog
	ue.Log = logger.NasLog
	ue.Tai = getTai()
	ue.AmfUe = amfSelf.NewAmfUe("")
	ue.AmfUe.State[models.AccessType__3_GPP_ACCESS].Set(amf_context.Deregistered)

	return ue
}

func getServiceRequestNasMsg(b *testing.B) []byte {
	msg := nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeServiceRequest)
	msg.GmmMessage.ServiceRequest = nasMessage.NewServiceRequest(nas.MsgTypeServiceRequest)
	sr := msg.GmmMessage.ServiceRequest
	sr.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	sr.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	sr.ServiceRequestMessageIdentity.SetMessageType(nas.MsgTypeServiceRequest)
	sr.ServiceTypeAndNgksi.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	sr.ServiceTypeAndNgksi.SetNasKeySetIdentifiler(0)
	sr.ServiceTypeAndNgksi.SetServiceTypeValue(nasMessage.ServiceTypeSignalling)
	sr.TMSI5GS.SetLen(7)

	buf, err := msg.PlainNasEncode()
	require.NoError(b, err)

	buf = append([]uint8{
		nasMessage.Epd5GSMobilityManagementMessage,
		nas.SecurityHeaderTypeIntegrityProtected,
		0, 0, 0, 0, 0,
	},
		buf...)

	return buf
}

func BenchmarkHandlerAuthenticationRequest(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		HandleNAS(getMockRanUeCtx(), ngapType.ProcedureCodeUplinkNASTransport, getServiceRequestNasMsg(b), true)
	}

	rcvValue, errRcv := metrics_utils.GetCounterVecValue(
		metrics_nas.NasMsgRcvCounterName,
		metrics_nas.NasMsgRcvCounter,
		prometheus.Labels{"name": "ServiceRequest", "status": "successful", "cause": ""})

	require.Nilf(b, errRcv, "Could not retrieve %s counter", metrics_ngap.NgapMsgRcvCounterName)

	sentValue, errSent := metrics_utils.GetCounterVecValue(
		metrics_nas.NasMsgSentCounterName,
		metrics_nas.NasMsgSentCounter,
		prometheus.Labels{"name": "ServiceReject", "status": "failure", "cause": "UE identity cannot be derived by the network"})

	require.Nilf(b, errSent, "Could not retrieve %s counter", metrics_ngap.NgapMsgSentCounterName)

	logger.NasLog.Printf("%s %d\n", metrics_ngap.NgapMsgRcvCounterName, int(rcvValue))
	logger.NasLog.Printf("%s %d\n", metrics_ngap.NgapMsgSentCounterName, int(sentValue))
}
