package cdr

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// over a field no rule references. No Go code reads it; rules are evaluated
	// against the raw record JSON, not this struct. See AzureProperties for why the
	// raw form needs its own type rather than a bare json.RawMessage.
	Properties AzureProperties `json:"properties,omitempty" bson:"properties,omitempty"`
}

// AzureProperties holds an Activity Log record's "properties" bag verbatim as
// JSON, and survives a MongoDB round trip in both directions.
//
// Raw, because Activity Log delivers the bag as an object for most operations and
// as a STRING containing JSON for some. A typed field rejects the second form and
// — the failure not being local — fails the WHOLE record decode, discarding a
// detection over a field no rule references.
//
// Its own type, because this struct is not only a wire shape: config-service
// persists runtime incidents with it embedded, so the bag also crosses BSON. A
// bare json.RawMessage is a []byte there, and the driver's byte-slice codec
// accepts only BSON binary or string — never an embedded document, which is what
// every incident stored before this type existed holds. Reading one back failed
// the whole find with
//
//	error decoding key cdrevent.eventdata.azureactivitylog.properties:
//	cannot decode document into json.RawMessage
//
// and every Azure incident 500ed. Writing was as quietly wrong in the other
// direction: []byte marshals to BSON binary, so incidents stored during that
// window hold an opaque blob that no Mongo filter or index can reach.
//
// The methods below therefore accept every shape the collection now holds —
// document (before), binary (during), string, null — and store the bag in its
// natural BSON form, so a structured bag stays queryable as a document or array
// rather than becoming an opaque blob. Reading and writing are symmetric: every
// shape the writer can produce, the reader reads back. Decoding never fails: an
// unreadable bag yields nil rather than sinking the incident it belongs to,
// which is the same bargain the JSON side makes.
type AzureProperties json.RawMessage

// MarshalJSON emits the bag verbatim. A named type does not inherit
// json.RawMessage's methods, and without this encoding/json would base64 it.
func (p AzureProperties) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("null"), nil
	}
	return p, nil
}

// UnmarshalJSON keeps the bag verbatim, whatever JSON type it arrived as.
func (p *AzureProperties) UnmarshalJSON(data []byte) error {
	*p = append((*p)[:0], data...)
	return nil
}

// MarshalBSONValue stores the bag in its natural BSON form: an object becomes an
// embedded document and an array a BSON array, so a structured bag stays
// queryable and indexable rather than becoming an opaque blob; a top-level JSON
// scalar becomes the matching BSON scalar. A bag that is not valid JSON is stored
// as a string rather than dropped, and an empty or null one as BSON null.
//
// Every type this can emit is read back by UnmarshalBSONValue. Keep the two in
// step: a shape the writer emits and the reader does not recognise is silent
// data loss on the round trip.
func (p AzureProperties) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if len(p) == 0 {
		return bsontype.Null, nil, nil
	}
	var v interface{}
	if err := json.Unmarshal(p, &v); err != nil {
		return bson.MarshalValue(string(p))
	}
	if v == nil {
		// A bag that is literally JSON null. bson.MarshalValue(nil) errors
		// ("no encoder found for <nil>"), which would fail the whole incident.
		return bsontype.Null, nil, nil
	}
	return bson.MarshalValue(v)
}

// UnmarshalBSONValue reads back every shape the collection holds, and never
// fails: a bag it cannot read becomes nil instead of failing the surrounding
// incident decode.
//
// Documents reach JSON through jsonifyBSON, so BSON types with no JSON scalar
// form (dates, binary) come back as their driver representation rather than
// their original text. An Activity Log bag is plain JSON, so this is exact in
// practice; it is not a byte-exact round trip for arbitrary BSON.
func (p *AzureProperties) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}
	switch t {
	case bsontype.EmbeddedDocument, bsontype.Array,
		bsontype.Double, bsontype.Boolean, bsontype.Int32, bsontype.Int64:
		// Documents are what the original map[string]any field wrote. Scalars are
		// what MarshalBSONValue writes for a bag whose top-level JSON is a number
		// or a bool — rare in Activity Log traffic, which documents properties as
		// an object or a JSON-encoded string, but the writer emits them, so the
		// reader has to take them back or the round trip silently drops the bag.
		var v interface{}
		if err := raw.Unmarshal(&v); err != nil {
			*p = nil
			return nil
		}
		b, err := json.Marshal(jsonifyBSON(v))
		if err != nil {
			*p = nil
			return nil
		}
		*p = b
	case bsontype.Binary:
		// Written by the bare json.RawMessage field: []byte marshals to binary.
		if _, b, ok := raw.BinaryOK(); ok && json.Valid(b) {
			*p = append((*p)[:0], b...)
		} else {
			*p = nil
		}
	case bsontype.String:
		// A bag Azure sent as a JSON string, or one MarshalBSONValue could not
		// parse. Keep it as the JSON it is; quote it if it is not JSON at all.
		s, ok := raw.StringValueOK()
		if !ok {
			*p = nil
			return nil
		}
		if json.Valid([]byte(s)) {
			*p = []byte(s)
			return nil
		}
		b, err := json.Marshal(s)
		if err != nil {
			*p = nil
			return nil
		}
		*p = b
	default:
		// Null, Undefined, and BSON types this type never writes (dates, ObjectIDs,
		// Decimal128, ...). Nothing stores those in this field today; if that ever
		// changes, add them here rather than letting them fall through to nil.
		*p = nil
	}
	return nil
}

// jsonifyBSON rewrites the driver's document and array types into their JSON
// equivalents. The driver decodes a BSON document reached through an interface{}
// into a primitive.D — an ordered []{Key, Value} — which encoding/json would
// render as a list of {"Key":…,"Value":…} pairs rather than as an object. Scalars
// pass through: every one the driver produces already marshals to valid JSON.
func jsonifyBSON(v interface{}) interface{} {
	switch t := v.(type) {
	case primitive.D:
		m := make(map[string]interface{}, len(t))
		for _, e := range t {
			m[e.Key] = jsonifyBSON(e.Value)
		}
		return m
	case primitive.M:
		m := make(map[string]interface{}, len(t))
		for k, e := range t {
			m[k] = jsonifyBSON(e)
		}
		return m
	case primitive.A:
		a := make([]interface{}, len(t))
		for i, e := range t {
			a[i] = jsonifyBSON(e)
		}
		return a
	default:
		return v
	}
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
