package cdr

import (
	"bytes"
	"encoding/json"
	"strings"
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
// Two delivery paths deliver the same event in slightly different shapes, and
// this struct models BOTH (verified against captured samples of each):
//
//   - Event Hub / diagnostic-settings export — what the in-account collector
//     consumes. Nests the RBAC context and token claims under an "identity"
//     object, carries no top-level "caller", and adds wrapper fields
//     (RoleLocation, Stamp, ReleaseVersion, VmSku) that are intentionally not
//     modelled.
//   - Azure Monitor REST / management API — puts "authorization" and "claims"
//     at the top level alongside a "caller".
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
	// hierarchy, statusCode, responseBody, ...); shape varies by operation.
	Properties AzureProperties `json:"properties,omitempty"`
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
// event arrived in, preferring the Event Hub's nested identity. Never nil-maps a
// lookup: the zero value is an empty map, so indexing is always safe.
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

// AzureProperties is the operation-specific property bag. Activity Log delivers
// it either as a JSON object or, for some operations, as a STRING containing
// JSON. A plain map[string]any fails to unmarshal the second form, and because
// the whole record decode then fails, a detection built from that record is
// discarded — so this type accepts both and never fails the surrounding decode.
type AzureProperties map[string]any

// UnmarshalJSON accepts an object, a JSON-encoded string, or a bare string. A
// bare (non-JSON) string is preserved under the "message" key rather than
// dropped, since that is what Activity Log puts there.
func (p *AzureProperties) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*p = nil
		return nil
	}

	// The common case: a JSON object.
	if trimmed[0] == '{' {
		var m map[string]any
		if err := json.Unmarshal(trimmed, &m); err != nil {
			return err
		}
		*p = m
		return nil
	}

	// A JSON string. Its contents are usually themselves JSON.
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		inner := strings.TrimSpace(s)
		if strings.HasPrefix(inner, "{") {
			var m map[string]any
			if err := json.Unmarshal([]byte(inner), &m); err == nil {
				*p = m
				return nil
			}
		}
		*p = map[string]any{"message": s}
		return nil
	}

	// Anything else (array, number, bool): keep it rather than failing the record.
	var v any
	if err := json.Unmarshal(trimmed, &v); err != nil {
		return err
	}
	*p = map[string]any{"value": v}
	return nil
}

// AzureAuthorization is the RBAC authorization context of an Activity Log operation.
type AzureAuthorization struct {
	Scope    string                 `json:"scope,omitempty"`
	Action   string                 `json:"action,omitempty"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}
