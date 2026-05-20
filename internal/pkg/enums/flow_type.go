package enums

type FlowType string

var (
	FlowTypeUnknown         FlowType = "UNKNOWN"
	FlowTypeRegistration    FlowType = "REGISTRATION"
	FlowTypeRestorePassword FlowType = "RESTORE_PASSWORD"
)
