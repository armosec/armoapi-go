package cdr

import (
	"encoding/json"
	"time"
)

const (
	// CdrEventGcpProjectIDJsonPath is the JSON path to the GCP project ID on a GCP
	// Cloud Audit Log event — the account key the CDR alert is matched to. It is
	// the GCP equivalent of the AWS accountId / Azure subscriptionId
	// account-matching path.
	//
	// GCP exposes NO org-key equivalent of the AWS orgId / Azure tenantId path:
	// the audit record carries only project_id, with no org/folder ancestry, so
	// there is no field to point such a constant at. Org-level connections are
	// instead attributed by resolving this project_id to the org-connected
	// account via CADRGcpConfig (an armosec-infra concern), not from the event.
	CdrEventGcpProjectIDJsonPath = "cdrevent.eventData.gcpAuditLog.resource.labels.project_id"
)

// GcpAuditLogEvent is a single GCP Cloud Logging LogEntry carrying an Admin
// Activity audit record — the control-plane audit event that is the GCP
// equivalent of an AWS CloudTrail management event (see CloudTrailEvent in
// aws.go) and the Azure Activity Log record (see AzureActivityLogEvent in
// azure.go).
//
// The Pub/Sub push pipe delivers one LogEntry per message; unlike the Azure
// Event Hub path there is no batch envelope. The audit payload lives under
// ProtoPayload (google.cloud.audit.AuditLog), while Operation lives at the
// LogEntry level — deliberately, because the synchronous-vs-long-running
// distinction the detection gate depends on is a LogEntry property, not a
// payload one.
//
// The field set and its nil-safety were measured against real events in the
// gcp-cdr-poc (see its FINDINGS.md §1); the round-trip test in gcp_test.go is
// the guard. POC corrections baked into the shape below:
//   - Operation is a pointer: synchronous admin calls (create SA, set IAM
//     policy) carry NO operation block at all, so the detection gate is
//     `!has(operation) || operation.last`. Long-running calls emit two entries
//     sharing operation.id, flagged first / last.
//   - Operation booleans are OMITTED when false (not serialized as false), so
//     CEL rules must use has() guards; the Go zero-value (false) is safe.
//   - ResourceName is the PARENT for create-style methods (CreateServiceAccount
//     -> projects/X); the created resource is in response.name — model both.
//   - Status is `{}` on success, so status.code is absent rather than 0.
//   - AuthorizationInfo carries permission / permissionType / granted, enabling
//     denied-attempt (granted:false) detection rules.
//
// All optional object/slice/string fields are pointers or omitempty so CEL
// has()-guards behave and GCP's per-service omissions round-trip cleanly. The
// two time.Time fields follow the aws.go/azure.go convention (bare, non-pointer);
// like EventTime/Time there, an absent timestamp re-serializes as the zero time
// rather than staying omitted — harmless in practice, as Cloud Logging stamps
// both timestamp and receiveTimestamp on every delivered entry.
type GcpAuditLogEvent struct {
	// InsertID uniquely identifies the log entry; it is stable across Pub/Sub
	// redeliveries and forms half of the dedup key hash(insertId + rule_id).
	InsertID string `json:"insertId,omitempty"`
	// LogName is the fully-qualified log name, e.g.
	// projects/<project>/logs/cloudaudit.googleapis.com%2Factivity.
	LogName string `json:"logName,omitempty"`
	// Timestamp is when the logged action occurred.
	Timestamp time.Time `json:"timestamp"`
	// ReceiveTimestamp is when Cloud Logging received the entry.
	ReceiveTimestamp time.Time `json:"receiveTimestamp"`
	// Severity is the log severity (e.g. "NOTICE").
	Severity string `json:"severity,omitempty"`
	// Resource is the monitored resource the entry is about; resource.labels
	// carries the project_id used to match the alert to a GCP account.
	Resource *GcpMonitoredResource `json:"resource,omitempty"`
	// Operation ties together the entries of a single long-running operation.
	// ABSENT on synchronous events — the detection gate is
	// `!has(operation) || operation.last`.
	Operation *GcpLogEntryOperation `json:"operation,omitempty"`
	// ProtoPayload is the google.cloud.audit.AuditLog record.
	ProtoPayload *GcpAuditLogPayload `json:"protoPayload,omitempty"`
}

// GcpMonitoredResource is the Cloud Logging monitored-resource descriptor.
// For an Admin Activity entry, Type is e.g. "service_account" / "gce_network"
// and Labels carries "project_id" (and often "location", etc.).
type GcpMonitoredResource struct {
	Type   string            `json:"type,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// GcpLogEntryOperation is the LogEntry.operation block. It is present only on
// long-running operations, whose request and completion entries share the same
// ID and are flagged First / Last respectively.
//
// NOTE: First / Last are omitted from the wire when false (they are not
// serialized as false). Go unmarshals the missing key to the false zero-value
// (safe), but CEL rule authors must use has() guards rather than comparing to
// false. The chosen gate reads only operation.last, so it is has()-safe.
type GcpLogEntryOperation struct {
	ID       string `json:"id,omitempty"`
	Producer string `json:"producer,omitempty"`
	First    bool   `json:"first,omitempty"`
	Last     bool   `json:"last,omitempty"`
}

// GcpAuditLogPayload models the google.cloud.audit.AuditLog protoPayload of an
// Admin Activity entry — the fields the CDR detection and identifier
// normalization lean on. Optional bags (Request / Response) are kept generic;
// Response in particular carries the real target of create-style methods
// (response.name) that ResourceName does not.
type GcpAuditLogPayload struct {
	// Type is the payload @type discriminator
	// ("type.googleapis.com/google.cloud.audit.AuditLog").
	Type string `json:"@type,omitempty"`
	// ServiceName is the API that produced the event (e.g. "iam.googleapis.com").
	ServiceName string `json:"serviceName,omitempty"`
	// MethodName is the API method invoked (e.g.
	// "google.iam.admin.v1.CreateServiceAccount").
	MethodName string `json:"methodName,omitempty"`
	// ResourceName is the resource the request targets. For create-style methods
	// it is the PARENT (e.g. "projects/X"), not the created resource — the real
	// target is in Response["name"]. Do NOT assume ResourceName == affected
	// resource (POC FINDINGS §1c②).
	ResourceName string `json:"resourceName,omitempty"`
	// AuthenticationInfo is the caller identity (principalEmail / principalSubject).
	AuthenticationInfo *GcpAuthenticationInfo `json:"authenticationInfo,omitempty"`
	// AuthorizationInfo is the per-permission authorization decisions for the
	// call; granted:false entries enable denied/attempted-action detection.
	AuthorizationInfo []GcpAuthorizationInfo `json:"authorizationInfo,omitempty"`
	// RequestMetadata carries the caller IP and user agent.
	RequestMetadata *GcpRequestMetadata `json:"requestMetadata,omitempty"`
	// Status is the operation status. Present but EMPTY ({}) on success, so
	// status.code is absent rather than 0 — guard with has() in CEL.
	Status *GcpStatus `json:"status,omitempty"`
	// PolicyViolationInfo carries org-policy / VPC Service Controls violation
	// details when a request is denied by policy — a detection class of its own
	// (org-policy and perimeter denials). Shape varies; kept generic.
	PolicyViolationInfo map[string]interface{} `json:"policyViolationInfo,omitempty"`
	// ResourceLocation is where the resource is/was located (current vs original),
	// e.g. for data-residency-aware rules.
	ResourceLocation *GcpResourceLocation `json:"resourceLocation,omitempty"`
	// NumResponseItems is the number of items returned by a list/query method.
	// protojson encodes int64 as a quoted JSON string ("3"), but the protobuf-JSON
	// spec also accepts a bare number (3) on parse. json.Number accepts BOTH forms;
	// a plain string would fail the whole decode on a bare number — a dropped event
	// for a field nothing depends on.
	NumResponseItems json.Number `json:"numResponseItems,omitempty"`
	// Request is the operation-specific request bag; shape varies by method.
	Request map[string]interface{} `json:"request,omitempty"`
	// Response is the operation-specific response bag; for create-style methods
	// Response["name"] is the created resource (the real target).
	Response map[string]interface{} `json:"response,omitempty"`
	// Metadata is where several services (GKE, BigQuery, Cloud SQL, ...) put their
	// audit detail instead of Request / Response; shape varies by service.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// GcpResourceLocation is the AuditLog.resourceLocation — the resource's current
// and (for moves) original locations.
type GcpResourceLocation struct {
	CurrentLocations  []string `json:"currentLocations,omitempty"`
	OriginalLocations []string `json:"originalLocations,omitempty"`
}

// GcpAuthenticationInfo is the AuditLog.authenticationInfo — the caller identity.
type GcpAuthenticationInfo struct {
	// PrincipalEmail is the caller's email (user or service account).
	PrincipalEmail string `json:"principalEmail,omitempty"`
	// PrincipalSubject is the prefixed principal ("user:..." / "serviceAccount:..."),
	// which cleanly distinguishes user from service-account callers for
	// common.Identifiers normalization (POC FINDINGS §1c⑤).
	PrincipalSubject string `json:"principalSubject,omitempty"`
	// ServiceAccountKeyName is the resource name of the service-account key used
	// to authenticate, when a long-lived key was used rather than short-lived
	// credentials — a standing detection / compliance signal.
	ServiceAccountKeyName string `json:"serviceAccountKeyName,omitempty"`
	// ServiceAccountDelegationInfo is the identity-delegation (impersonation)
	// chain. Service-account impersonation (iam.serviceAccounts.getAccessToken /
	// actAs) is the canonical GCP privilege-escalation path; this is what lets an
	// alert say who impersonated whom. Absent unless impersonation occurred (so
	// not seen in the human-auth POC run — its absence there is not evidence of
	// rarity in the field).
	ServiceAccountDelegationInfo []GcpServiceAccountDelegationInfo `json:"serviceAccountDelegationInfo,omitempty"`
	// OAuthInfo carries the OAuth client that made the call, when present.
	OAuthInfo *GcpOAuthInfo `json:"oauthInfo,omitempty"`
}

// GcpServiceAccountDelegationInfo is one link in the impersonation chain
// (AuditLog.authenticationInfo.serviceAccountDelegationInfo[]).
type GcpServiceAccountDelegationInfo struct {
	// PrincipalSubject is the delegated principal at this link.
	PrincipalSubject string `json:"principalSubject,omitempty"`
	// FirstPartyPrincipal is set when a first-party (Google) identity delegated.
	FirstPartyPrincipal *GcpFirstPartyPrincipal `json:"firstPartyPrincipal,omitempty"`
	// ThirdPartyPrincipal is set for third-party (external) delegation; shape
	// varies, kept generic.
	ThirdPartyPrincipal map[string]interface{} `json:"thirdPartyPrincipal,omitempty"`
}

// GcpFirstPartyPrincipal is the first-party identity in a delegation link.
type GcpFirstPartyPrincipal struct {
	PrincipalEmail string `json:"principalEmail,omitempty"`
}

// GcpOAuthInfo is the AuditLog.authenticationInfo.oauthInfo block.
type GcpOAuthInfo struct {
	OAuthClientID string `json:"oauthClientId,omitempty"`
}

// GcpAuthorizationInfo is one AuditLog.authorizationInfo[] entry — the
// authorization decision for a single permission. PermissionType (e.g.
// "ADMIN_WRITE") is a clean control-plane-write discriminator, and Granted:false
// surfaces denied/attempted privileged actions (recon, privilege probing).
//
// NOTE: like the operation booleans, GCP omits Granted from the wire when false,
// so it unmarshals to the false zero-value; CEL rules must has()-guard it.
type GcpAuthorizationInfo struct {
	Permission     string `json:"permission,omitempty"`
	PermissionType string `json:"permissionType,omitempty"`
	Granted        bool   `json:"granted,omitempty"`
	Resource       string `json:"resource,omitempty"`
	// ResourceAttributes describes the resource the permission was checked
	// against. Its Type (e.g. "iam.googleapis.com/ServiceAccount") is the kind of
	// thing touched — the field to fall back on when ResourceName is the parent
	// rather than the target (POC FINDINGS §1c②); mirrors AWS Resource.ResourceType.
	ResourceAttributes *GcpResourceAttributes `json:"resourceAttributes,omitempty"`
}

// GcpResourceAttributes is the AuthorizationInfo.resourceAttributes — the typed
// resource the authorization decision was made against.
type GcpResourceAttributes struct {
	Service string `json:"service,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
}

// GcpRequestMetadata is the AuditLog.requestMetadata — the caller's network
// context, useful for SourceInformation normalization.
type GcpRequestMetadata struct {
	CallerIP                string `json:"callerIp,omitempty"`
	CallerSuppliedUserAgent string `json:"callerSuppliedUserAgent,omitempty"`
	// CallerNetwork is the VPC network the call originated from, when applicable.
	CallerNetwork string `json:"callerNetwork,omitempty"`
	// RequestAttributes is the request context (time, auth, method, ...); shape
	// varies, kept generic.
	RequestAttributes map[string]interface{} `json:"requestAttributes,omitempty"`
	// DestinationAttributes is the request destination context; shape varies,
	// kept generic.
	DestinationAttributes map[string]interface{} `json:"destinationAttributes,omitempty"`
}

// GcpStatus is the AuditLog.status (google.rpc.Status). Empty ({}) on success,
// so Code is absent rather than 0.
type GcpStatus struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
