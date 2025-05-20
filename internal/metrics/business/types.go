package businness

// Global metric information
const (
	SUBSYSTEM_NAME = "business"
)

// Collectors information
const (
	// States
	GMM_STATE_GAUGE_NAME = "ue_gmm_state_count"
	GMM_STATE_GAUGE_DESC = "Current number of UEs in each 5GMM state in the AMF"

	GMM_TRANSITION_COUNTER_NAME = "ue_gmm_transitions_total"
	GMM_TRANSITION_COUNTER_DESC = "Count of UE GMM state transitions in the AMF"

	GMM_DURATION_HISTOGRAM_NAME = "ue_gmm_state_duration_seconds"
	GMM_DURATION_HISTOGRAM_DESC = "Duration that UEs spend in a given GMM state before transitioning"

	// Connection Management
	UE_CM_STATE_GAUGE_NAME = "ue_cm_gmm_state_count"
	UE_CM_STATE_GAUGE_DESC = "Count of the UE in each Connection Management State (CM_IDLE, CM_CONNECTED) in the AMF"

	// Handover
	HANDOVER_IN_PROGRESS_GAUGE_NAME = "handover_current_count"
	HANDOVER_IN_PROGRESS_GAUGE_DESC = "Number of UEs currently in handover procedure (source AMF side)"

	HANDOVER_EVENT_COUNTER_NAME = "handover_events_total"
	HANDOVER_EVENT_COUNTER_DESC = "Count of handover events (attempts, successes, failures)"

	HANDOVER_DURATION_HISTOGRAM_NAME = "handover_duration_seconds"
	HANDOVER_DURATION_HISTOGRAM_DESC = "Histogram of the handover duration in seconds"

	// PDU Session
	PDU_SESSION_IN_PROGRESS_GAUGE_NAME = "pdu_session_current_count"
	PDU_SESSION_IN_PROGRESS_GAUGE_DESC = "Number of PDU Session currently in  active (source AMF side)"

	PDU_SESSION_EVENT_COUNTER_NAME = "pdu_session_events_total"
	PDU_SESSION_EVENT_COUNTER_DESC = "Count of pdu events (setup, release, modification)"

	PDU_SESSION_DURATION_HISTOGRAM_NAME = "pdu_session_duration_seconds"
	PDU_SESSION_DURATION_HISTOGRAM_DESC = "Histogram of the pdu session duration in seconds"
)

// Label names
const (
	// States
	GMM_STATE_ACCESS_LABEL     = "access"
	GMM_STATE_LABEL            = "gmm_state"
	GMM_STATE_FROM_STATE_LABEL = "from_state"
	GMM_STATE_TO_STATE_LABEL   = "to_state"

	// Connection Management
	UE_CM_STATE_LABEL        = "state"
	UE_CM_ACCESS_STATE_LABEL = "access_type"

	// Handover
	HANDOVER_TYPE_LABEL  = "type"
	HANDOVER_EVENT_LABEL = "event"
	HANDOVER_CAUSE_LABEL = "cause"

	// PDU
	PDU_SESSION_EVENT_LABEL  = "event"
	PDU_SESSION_STATUS_LABEL = "status"
	PDU_SESSION_CAUSE_LABEL  = "cause"
)

// Metrics Values
const (
	// Connection Management
	UE_CM_IDLE_VALUE      = "cm-idle"
	UE_CM_CONNECTED_VALUE = "cm-connected"

	// Handover
	HANDOVER_TYPE_XN_VALUE       = "xn"
	HANDOVER_TYPE_NGAP_VALUE     = "ngap"
	HANDOVER_EVENT_ATTEMPT_VALUE = "attempt"
	HANDOVER_EVENT_FAILURE_VALUE = "failure"
	HANDOVER_EVENT_SUCCESS_VALUE = "success"
)

// Potential Causes
const (
	HANDOVER_RAN_UE_MISSING_ERR                        = "ran ue missing"
	HANDOVER_AMF_UE_MISSING_ERR                        = "amf ue missing"
	HANDOVER_TARGET_UE_MISSING_ERR                     = "target ue is missing"
	HANDOVER_SECURITY_CONTEXT_MISSING_ERR              = "security context missing"
	HANDOVER_SWITCH_RAN_ERR                            = "ue could not switch ran"
	HANDOVER_TARGET_ID_NOT_SUPPORTED_ERR               = "target id type is not supported"
	HANDOVER_PDU_SESSION_RES_REL_LIST_ERR              = "some pdu session could not been release for handover"
	HANDOVER_BETWEEN_DIFFERENT_AMF_NOT_SUPPORTED       = "handover between different amf has not been implemented yet"
	HANDOVER_NOT_YET_IMPLEMENT_N2_HANDOVER_BETWEEN_AMF = "n2 Handover between amf has not been implemented yet"
	EMPTY_CAUSE                                        = ""
)
