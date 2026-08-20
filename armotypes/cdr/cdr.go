package cdr

import "github.com/armosec/armoapi-go/armotypes/common"

type CustomerDetails struct {
	// CustomerGUID is the unique identifier of the customer account
	CustomerGUID string `json:"customerGUID"`
	// AccessKey is the access key of the customer account
	AccessKey string `json:"accessKey"`
}

// Cloud services
type CloudService string

const (
	// CloudTrail is the cloudtrail service
	CloudTrail CloudService = "cloudtrail"
	// ActivityLogs is the Azure Activity Log service (control-plane audit log)
	ActivityLogs CloudService = "activitylogs"
	// CloudAuditLogs is the GCP Cloud Audit Logs service; the CDR pipe consumes
	// the Admin Activity control-plane audit log (the GCP equivalent of AWS
	// CloudTrail management events / the Azure Activity Log).
	CloudAuditLogs CloudService = "cloudauditlogs"
	// Add more cloud services here
)

// Cloud providers
type CloudProvider string

const (
	// AWS is the AWS cloud provider
	AWS CloudProvider = "aws"
	// Azure is the Microsoft Azure cloud provider
	Azure CloudProvider = "azure"
	// GCP is the Google Cloud Platform cloud provider
	GCP CloudProvider = "gcp"
	// Add more cloud providers here
)

type CloudMetadata struct {
	// Provider is the cloud provider
	Provider CloudProvider `json:"provider,omitempty"`
	// SourceService is the source service (e.g cloudtrail, cloudwatch, etc)
	SourceService CloudService `json:"sourceService,omitempty"`
}

// The types corresponds to the SourceService type
type EventData struct {
	// AWSCloudTrail cloudtrail event
	AWSCloudTrail *CloudTrailEvent `json:"awsCloudTrail,omitempty"`
	// AzureActivityLog azure activity log event
	AzureActivityLog *AzureActivityLogEvent `json:"azureActivityLog,omitempty"`
	// GcpAuditLog gcp cloud audit log (Admin Activity) event
	GcpAuditLog *GcpAuditLogEvent `json:"gcpAuditLog,omitempty"`
	// Target resource
	TargetResource string `json:"targetResource,omitempty"`
	// Identifiers of the alert
	Identifiers *common.Identifiers `json:"identifiers,omitempty"`
	// Add more cloud event data here
}

type CdrAlert struct {
	// CloudMetadata is the metadata of the cloud
	CloudMetadata `json:"cloudMetadata,omitempty"`
	// EventData is the event data
	EventData `json:"eventData,omitempty"`
	// RuleName is the name of the rule
	RuleName string `json:"ruleName,omitempty"`
	// RuleID is the unique identifier of the rule
	RuleID string `json:"ruleID,omitempty"`
	// Description is the description of the rule
	Description string `json:"description,omitempty"`
	// Priority is the severity of the rule
	Priority string `json:"priority,omitempty"`
	// Tags is the tags of the rule
	Tags []string `json:"tags,omitempty"`
	// Message is the failure message
	Message string `json:"message,omitempty"`
	// MitreTactic is the MITRE ATT&CK tactic
	MitreTactic string `json:"mitreTactic,omitempty"`
	// MitreTechnique is the MITRE ATT&CK technique
	MitreTechnique string `json:"mitreTechnique,omitempty"`
	// UniqueID is the unique identifier of the alert
	UniqueID string `json:"uniqueID,omitempty"`
}

// ConnectionLevel states whether a CDR batch belongs to a single-account connection or an
// organization/tenant-wide one. The ingester routes keep-alive on this stated intent rather
// than inferring it from OrgID presence; an empty value means "infer" — the legacy behavior
// existing AWS producers rely on. See docs/features/cdr-heartbeat-contract.md.
type ConnectionLevel string

const (
	// ConnectionLevelAccount is a single-account connection (AWS account / Azure
	// subscription / GCP project); keep-alive is keyed on CloudAccountID.
	ConnectionLevelAccount ConnectionLevel = "account"
	// ConnectionLevelOrganization is an organization/tenant-wide connection (AWS org /
	// Azure tenant / GCP org); keep-alive is keyed on OrgID.
	ConnectionLevelOrganization ConnectionLevel = "organization"
)

type CdrAlertBatch struct {
	// CustomerGUID is the unique identifier of the customer
	CustomerGUID string `json:"customerGUID,omitempty"`
	// CloudAccountID is the unique identifier of the cloud account
	CloudAccountID string `json:"cloudAccountID,omitempty"`
	// OrgID is the unique identifier of the organization
	OrgID string `json:"orgID,omitempty"`
	// Provider is the cloud provider
	Provider CloudProvider `json:"provider,omitempty"`
	// RuleFailures is the list of rule failures
	RuleFailures []CdrAlert `json:"ruleFailures,omitempty"`
	// IsHeartbeat marks a periodic liveness message (no RuleFailures); absent/false = a normal alert batch. See docs/features/cdr-heartbeat-contract.md.
	IsHeartbeat bool `json:"isHeartbeat,omitempty"`
	// ConnectionLevel states account- vs organization-level so the ingester routes keep-alive by intent, not by OrgID presence; empty = legacy inference. See docs/features/cdr-heartbeat-contract.md.
	ConnectionLevel ConnectionLevel `json:"connectionLevel,omitempty"`
	// LogsSeen counts audit records the collector received through its log pipe since start, so the backend can gate Pending -> Connected on evidence rather than liveness; nil = producer does not report it, 0 = reports it and has seen none. Read it with LogsSeenValue. See docs/features/cdr-heartbeat-contract.md.
	LogsSeen *uint64 `json:"logsSeen,omitempty"`
}

// LogsSeenValue returns the batch's logsSeen count and whether the producer
// reported it at all.
//
// Branch on reported rather than treating a missing count as zero. The two say
// different things, and conflating them breaks the gate in both directions: a
// producer that does not report the signal (AWS, Azure) has to fall back to
// liveness, while a producer reporting zero is explicitly saying "no log has
// traversed the pipe yet" and must not be treated as connected. This is why the
// field is a pointer — omitempty tests nil, not the pointee, so an explicit 0
// still reaches the wire.
//
// The count gates the Pending -> Connected transition ONLY, and that transition
// latches. The count is cumulative per collector instance and resets when the
// instance is replaced, so an already-Connected account will legitimately send
// logsSeen: 0 again after a restart; regressing it to Pending on that would flap
// the connection on every recycle.
func (b *CdrAlertBatch) LogsSeenValue() (count uint64, reported bool) {
	if b == nil || b.LogsSeen == nil {
		return 0, false
	}
	return *b.LogsSeen, true
}

// NewLogsSeen returns a pointer suitable for CdrAlertBatch.LogsSeen. Producers use
// it so that reporting a genuine zero is a one-liner and does not accidentally
// become "not reported".
func NewLogsSeen(count uint64) *uint64 {
	return &count
}
