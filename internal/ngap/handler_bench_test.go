package ngap

import (
	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/metrics"
	metrics_ngap "github.com/free5gc/amf/internal/metrics/ngap"
	metrics_utils "github.com/free5gc/amf/internal/metrics/utils"
	nastesting "github.com/free5gc/amf/internal/nas/testing"
	ngaptesting "github.com/free5gc/amf/internal/ngap/testing"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/ngap/ngapType"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	benchRan *amf_context.AmfRan
	benchPDU ngapType.NGAPPDU
	connStub *ngaptesting.SctpConnStub
)

func getMockConfiguration() *factory.Config {
	cfg := &factory.Config{
		Configuration: &factory.Configuration{
			AmfName: "amf",
		},
	}
	return cfg
}

func init() {

	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	metrics.Init(&factory.Config{})

	// 1) fake SCTP connection
	connStub = new(ngaptesting.SctpConnStub)
	benchRan = NewAmfRan(connStub)

	// 2) init global AMF context
	amfSelf := amf_context.GetSelf()
	NewAmfContext(amfSelf)

	// 3) build a single “InitialUEMessage” PDU
	var nasId nasType.MobileIdentity5GS
	nasId.Len = 12
	nasId.Buffer = []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0, 0, 0, 0, 0x47, 0x78}
	benchPDU = BuildInitialUEMessage(
		1,
		nastesting.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration, nasId, nil, nil, nil, nil, nil),
		"",
	)

	factory.AmfConfig = getMockConfiguration()
}

func BenchmarkHandlerInitialUEMessage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		handlerInitialUEMessage(benchRan, &benchPDU, benchPDU.InitiatingMessage)
	}

	rcvValue, errRcv := metrics_utils.GetCounterVecValue(
		metrics_ngap.MsgRcvCounterName,
		metrics_ngap.MsgRcvCounter,
		prometheus.Labels{"name": "InitialUEMessage", "status": "successful", "cause": ""})

	require.Nilf(b, errRcv, "Could not retrieve %s counter", metrics_ngap.MsgRcvCounterName)

	sentValue, errSent := metrics_utils.GetCounterVecValue(
		metrics_ngap.MsgSentCounterName,
		metrics_ngap.MsgSentCounter,
		prometheus.Labels{"name": "DownlinkNasTransport", "status": "successful", "cause": ""})

	require.Nilf(b, errSent, "Could not retrieve %s counter", metrics_ngap.MsgSentCounterName)

	logger.NgapLog.Printf("%s %d\n", metrics_ngap.MsgRcvCounterName, int(rcvValue))
	logger.NgapLog.Printf("%s %d\n", metrics_ngap.MsgSentCounterName, int(sentValue))
}
