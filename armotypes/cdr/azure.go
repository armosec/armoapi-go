package cdr

import "time"

const (
	// CdrEventAzureSubscriptionIDJsonPath is the JSON path to the subscription ID
	// (the Azure equivalent of the AWS account ID) on an Azure Activity Log event.
	CdrEventAzureSubscriptionIDJsonPath = "cdrevent.eventData.azureActivityLog.subscriptionId"
	// CdrEventAzureTenantIDJsonPath is the JSON path to the tenant ID
	// (the Azure equivalent of the AWS org ID) on an Azure Activity Log event.
	CdrEventAzureTenantIDJsonPath = "cdrevent.eventData.azureActivityLog.tenantId"
)

// AzureActivityLogBatch is the envelope Azure Monitor delivers to Event Hub: one
// Event Hub message carries many Activity Log records. The in-account collector
// unwraps this and embeds each AzureActivityLogEvent into a CdrAlert's EventData.
type AzureActivityLogBatch struct {
	Records []AzureActivityLogEvent `json:"records"`
}

// AzureActivityLogEvent is a single Azure Activity Log record — the control-plane
// audit event that is the Azure equivalent of an AWS CloudTrail management event
// (see CloudTrailEvent in aws.go).
//
// The field set below covers the common Activity Log fields observed in the POC.
// NOTE: the exact shape differs slightly between the two delivery paths — the
// Event Hub stream vs. the Azure Monitor REST API (e.g. some callers nest
// authorization/claims under an "identity" object, and the Event Hub body adds
// wrapper fields). Reconcile/extend against real samples of both; the
// marshal/unmarshal round-trip test is the guard (see the POC FINDINGS.md for a
// captured sample).
type AzureActivityLogEvent struct {
	Time            time.Time `json:"time"`
	ResourceID      string    `json:"resourceId,omitempty"`
	OperationName   string    `json:"operationName,omitempty"`
	Category        string    `json:"category,omitempty"`
	ResultType      string    `json:"resultType,omitempty"`
	ResultSignature string    `json:"resultSignature,omitempty"`
	CorrelationID   string    `json:"correlationId,omitempty"`
	CallerIPAddress string    `json:"callerIpAddress,omitempty"`
	// Caller is the identity that performed the operation (UPN or object ID).
	Caller string `json:"caller,omitempty"`
	Level  string `json:"level,omitempty"`
	// Channels is the Activity Log channel (e.g. "Operation").
	Channels string `json:"channels,omitempty"`
	// SubscriptionID / TenantID / ResourceGroupName are the account identifiers
	// (Azure's equivalents of the AWS accountId/orgId). Present on the
	// management/REST shape; on the Event Hub shape SubscriptionID may need to be
	// derived from ResourceID.
	SubscriptionID    string `json:"subscriptionId,omitempty"`
	TenantID          string `json:"tenantId,omitempty"`
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// Authorization is the RBAC context of the operation.
	Authorization *AzureAuthorization `json:"authorization,omitempty"`
	// Claims is the caller's token claims (objectId, appid, ipaddr, tenantid, ...).
	Claims map[string]string `json:"claims,omitempty"`
	// Properties is the operation-specific bag (eventCategory, entity, message,
	// hierarchy, statusCode, responseBody, ...); shape varies by operation.
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// AzureAuthorization is the RBAC authorization context of an Activity Log operation.
type AzureAuthorization struct {
	Scope    string                 `json:"scope,omitempty"`
	Action   string                 `json:"action,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}
