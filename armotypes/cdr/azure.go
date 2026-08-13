package cdr

import (
	"encoding/json"
	"time"
)

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
// The same record reaches us in two shapes that disagree on where the caller
// identity lives, and this struct models the identity of BOTH (verified against
// captured samples of each):
//
//   - Event Hub / diagnostic-settings export — what the in-account collector
//     consumes, and the only shape decoded through this struct today. Nests the
//     RBAC context and token claims under an "identity" object, carries no
//     top-level "caller", and adds wrapper fields (RoleLocation, Stamp,
//     ReleaseVersion, VmSku) that are intentionally not modelled.
//   - Azure Monitor REST / management API — puts "authorization" and "claims"
//     at the top level alongside a "caller".
//
// The REST shape is modelled for its IDENTITY LAYOUT only. A verbatim REST
// response does NOT decode through this struct: REST sends operationName and
// category as {value, localizedValue} objects, which the string fields below
// reject. That is deliberate — no code decodes a REST response here, and the
// collector's records always carry the string form. Add normalization when a
// REST reader actually exists, not before.
//
// Read the identity through EffectiveClaims / EffectiveAuthorization /
// EffectiveCaller rather than the fields directly, so callers work on both
// shapes. Reading Claims or Authorization directly silently yields nothing on
// the Event Hub shape, which is where the collector's alerts come from.
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
	// Location is the Azure region the operation acted in — the analog of AWS
	// awsRegion. Often "global" for control-plane operations (subscription /
	// resource-group scope), a region (e.g. "eastus") for regional resources.
	Location string `json:"location,omitempty"`
	// Channels is the Activity Log channel (e.g. "Operation").
	Channels string `json:"channels,omitempty"`
	// SubscriptionID / TenantID / ResourceGroupName are the account identifiers
	// (Azure's equivalents of the AWS accountId/orgId). Present on the
	// management/REST shape; on the Event Hub shape SubscriptionID may need to be
	// derived from ResourceID.
	SubscriptionID    string `json:"subscriptionId,omitempty"`
	TenantID          string `json:"tenantId,omitempty"`
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// Authorization is the RBAC context of the operation, on the REST/management
	// shape. Empty on the Event Hub shape — use EffectiveAuthorization.
	Authorization *AzureAuthorization `json:"authorization,omitempty"`
	// Claims is the caller's token claims (objectId, appid, ipaddr, tenantid, ...),
	// on the REST/management shape. Empty on the Event Hub shape — use
	// EffectiveClaims.
	Claims map[string]string `json:"claims,omitempty"`
	// Identity is where the Event Hub / diagnostic-settings export nests the
	// authorization and claims. Nil on the REST/management shape.
	Identity *AzureIdentity `json:"identity,omitempty"`
	// Properties is the operation-specific bag (eventCategory, entity, message,
	// hierarchy, statusCode, responseBody, ...). Kept raw on purpose: Activity Log
	// delivers it as an object for most operations and as a STRING containing JSON
	// for some, so a map[string]any field rejects the second form and — because the
	// failure is not local — fails the WHOLE record decode, discarding a detection
	// over a field no rule references. json.RawMessage cannot fail to decode and
	// round-trips byte-for-byte into the alert payload. No Go code reads it; rules
	// are evaluated against the raw record JSON, not this struct.
	Properties json.RawMessage `json:"properties,omitempty"`
}

// AzureIdentity is the "identity" object of an Event Hub Activity Log record,
// holding the RBAC context and the caller's token claims.
type AzureIdentity struct {
	Authorization *AzureAuthorization `json:"authorization,omitempty"`
	Claims        map[string]string   `json:"claims,omitempty"`
}

// Claim keys that carry the caller's human identity. Ordered by preference: the
// Event Hub claims bag commonly has "name" and "upn"; the WS-Fed URI form shows
// up on tokens issued through the federated flow.
const (
	claimName    = "name"
	claimUPN     = "upn"
	claimNameURI = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"
	claimUPNURI  = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"
)

// EffectiveClaims returns the caller's token claims from whichever shape the
// event arrived in, preferring the Event Hub's nested identity. It always
// returns a non-nil map, so ranging over the result is safe even when the event
// carries no claims at all.
func (e AzureActivityLogEvent) EffectiveClaims() map[string]string {
	if e.Identity != nil && len(e.Identity.Claims) > 0 {
		return e.Identity.Claims
	}
	if len(e.Claims) > 0 {
		return e.Claims
	}
	return map[string]string{}
}

// EffectiveAuthorization returns the operation's RBAC context from whichever
// shape the event arrived in, preferring the Event Hub's nested identity.
func (e AzureActivityLogEvent) EffectiveAuthorization() *AzureAuthorization {
	if e.Identity != nil && e.Identity.Authorization != nil {
		return e.Identity.Authorization
	}
	return e.Authorization
}

// EffectiveCaller returns the human identity that performed the operation. The
// REST shape carries it in "caller"; the Event Hub shape has no such field, so it
// falls back to the claims that hold the user's name or UPN.
func (e AzureActivityLogEvent) EffectiveCaller() string {
	if e.Caller != "" {
		return e.Caller
	}
	claims := e.EffectiveClaims()
	for _, key := range []string{claimUPN, claimName, claimUPNURI, claimNameURI} {
		if v := claims[key]; v != "" {
			return v
		}
	}
	return ""
}

// AzureAuthorization is the RBAC authorization context of an Activity Log operation.
type AzureAuthorization struct {
	Scope    string                 `json:"scope,omitempty"`
	Action   string                 `json:"action,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}
