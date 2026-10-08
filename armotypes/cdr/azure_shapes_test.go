package cdr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// eventHubShape is a real Activity Log record as delivered by the Event Hub /
// diagnostic-settings export (captured from a dev tenant; identity values
// replaced with placeholders). This is what the in-account collector consumes.
//
// Note what it does NOT have: a top-level "caller", "authorization" or "claims".
// They are nested under "identity". It also carries wrapper fields the struct
// deliberately ignores.
const eventHubShape = `{
  "RoleLocation": "East US",
  "Stamp": "FDWeb",
  "ReleaseVersion": "7.2026.30.7",
  "VmSku": "Standard_DS5_V2",
  "time": "2026-08-12T15:11:47.5329521Z",
  "resourceId": "/SUBSCRIPTIONS/AAAAAAAA-1111-2222-3333-444444444444/RESOURCEGROUPS/RG1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/SA1",
  "operationName": "MICROSOFT.STORAGE/STORAGEACCOUNTS/LISTKEYS/ACTION",
  "category": "Administrative",
  "resultType": "Success",
  "resultSignature": "Succeeded.OK",
  "durationMs": "61",
  "callerIpAddress": "203.0.113.10",
  "correlationId": "cccccccc-1111-2222-3333-444444444444",
  "identity": {
    "authorization": {
      "scope": "/subscriptions/aaaaaaaa/resourceGroups/RG1/providers/Microsoft.Storage/storageAccounts/SA1",
      "action": "Microsoft.Storage/storageAccounts/listKeys/action",
      "evidence": { "role": "Owner", "principalType": "User" }
    },
    "claims": {
      "name": "Example User",
      "upn": "user@example.com",
      "idtyp": "user",
      "ipaddr": "203.0.113.10",
      "http://schemas.microsoft.com/identity/claims/objectidentifier": "bbbbbbbb-1111-2222-3333-444444444444",
      "http://schemas.microsoft.com/identity/claims/tenantid": "dddddddd-1111-2222-3333-444444444444"
    }
  },
  "level": "Information",
  "location": "global",
  "properties": { "statusCode": "OK", "eventCategory": "Administrative" },
  "tenantId": "dddddddd-1111-2222-3333-444444444444"
}`

// restShape captures the Azure Monitor REST API's IDENTITY LAYOUT — caller,
// authorization and claims at the top level — which is the only part of that
// shape this struct models.
//
// It is deliberately NOT a verbatim REST response: real REST sends
// operationName and category as {value, localizedValue} objects, which the
// struct's string fields reject. Nothing decodes a REST response through this
// struct (the collector only ever sees the Event Hub shape), so normalizing
// those objects would be speculative. If a REST reader is ever added, this
// fixture is where the object form has to appear first.
const restShape = `{
  "time": "2026-08-12T15:11:47Z",
  "resourceId": "/subscriptions/aaaaaaaa/resourceGroups/RG1",
  "operationName": "Microsoft.Storage/storageAccounts/listKeys/action",
  "category": "Administrative",
  "resultType": "Success",
  "callerIpAddress": "203.0.113.10",
  "caller": "user@example.com",
  "authorization": { "scope": "/subscriptions/aaaaaaaa", "action": "Microsoft.Storage/storageAccounts/listKeys/action" },
  "claims": {
    "idtyp": "user",
    "http://schemas.microsoft.com/identity/claims/objectidentifier": "bbbbbbbb-1111-2222-3333-444444444444"
  },
  "tenantId": "dddddddd-1111-2222-3333-444444444444"
}`

func TestEventHubShape_IdentityIsReadThroughEffectiveAccessors(t *testing.T) {
	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(eventHubShape), &ev))

	// The flat fields are empty on this shape — reading them directly is the bug
	// these accessors exist to prevent.
	assert.Empty(t, ev.Caller)
	assert.Empty(t, ev.Claims)
	assert.Nil(t, ev.Authorization)

	require.NotNil(t, ev.Identity)
	assert.Equal(t, "user@example.com", ev.EffectiveCaller())
	assert.Equal(t, "user", ev.EffectiveClaims()["idtyp"])
	assert.Equal(t, "bbbbbbbb-1111-2222-3333-444444444444",
		ev.EffectiveClaims()["http://schemas.microsoft.com/identity/claims/objectidentifier"])

	auth := ev.EffectiveAuthorization()
	require.NotNil(t, auth)
	assert.Equal(t, "Microsoft.Storage/storageAccounts/listKeys/action", auth.Action)
	assert.Equal(t, "Owner", auth.Evidence["role"])
}

func TestEventHubShape_FlatFieldsDecodeAsStrings(t *testing.T) {
	// Regression guard: operationName and category are plain strings on this
	// delivery path. The REST API's {value, localizedValue} objects would fail
	// this decode, and CEL rules comparing them as strings depend on it.
	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(eventHubShape), &ev))

	assert.Equal(t, "MICROSOFT.STORAGE/STORAGEACCOUNTS/LISTKEYS/ACTION", ev.OperationName)
	assert.Equal(t, "Administrative", ev.Category)
	assert.Equal(t, "Success", ev.ResultType)
	assert.Equal(t, "cccccccc-1111-2222-3333-444444444444", ev.CorrelationID)
	assert.Equal(t, "dddddddd-1111-2222-3333-444444444444", ev.TenantID)
}

func TestRESTShape_FlatIdentityLayoutStillResolves(t *testing.T) {
	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(restShape), &ev))

	assert.Nil(t, ev.Identity)
	assert.Equal(t, "user@example.com", ev.EffectiveCaller())
	assert.Equal(t, "user", ev.EffectiveClaims()["idtyp"])
	require.NotNil(t, ev.EffectiveAuthorization())
	assert.Equal(t, "Microsoft.Storage/storageAccounts/listKeys/action", ev.EffectiveAuthorization().Action)
}

func TestEffectiveClaims_SafeToIndexWhenAbsent(t *testing.T) {
	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(`{"category":"Administrative"}`), &ev))
	assert.Empty(t, ev.EffectiveCaller())
	assert.Empty(t, ev.EffectiveClaims()["idtyp"]) // must not panic
	assert.Nil(t, ev.EffectiveAuthorization())
}

func TestEffectiveCaller_PrefersTopLevelThenUPNThenName(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"top-level caller wins", `{"caller":"flat@example.com","identity":{"claims":{"upn":"nested@example.com"}}}`, "flat@example.com"},
		{"upn preferred over name", `{"identity":{"claims":{"name":"Example User","upn":"user@example.com"}}}`, "user@example.com"},
		{"falls back to name", `{"identity":{"claims":{"name":"Example User"}}}`, "Example User"},
		{"falls back to WS-Fed URI", `{"identity":{"claims":{"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name":"URI User"}}}`, "URI User"},
		{"nothing to report", `{"identity":{"claims":{"idtyp":"app"}}}`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ev AzureActivityLogEvent
			require.NoError(t, json.Unmarshal([]byte(tc.raw), &ev))
			assert.Equal(t, tc.want, ev.EffectiveCaller())
		})
	}
}

func TestProperties_NeverFailsTheRecordAndSurvivesVerbatim(t *testing.T) {
	// A record whose "properties" arrives as a JSON string used to fail the whole
	// decode, which discarded the detection built from it. Holding the field raw
	// makes that impossible for any well-formed JSON value, and forwards exactly
	// what Azure sent instead of a rewritten approximation of it.
	tests := []struct {
		name string
		raw  string
		want string // what "properties" must look like after a decode/encode cycle
	}{
		{"object", `{"properties":{"statusCode":"OK"}}`, `{"statusCode":"OK"}`},
		{"JSON encoded as a string", `{"properties":"{\"statusCode\":\"OK\"}"}`, `"{\"statusCode\":\"OK\"}"`},
		{"bare string", `{"properties":"something happened"}`, `"something happened"`},
		{"array", `{"properties":[1,2]}`, `[1,2]`},
		{"null", `{"properties":null}`, `null`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ev AzureActivityLogEvent
			require.NoError(t, json.Unmarshal([]byte(tc.raw), &ev), "decode must not fail")

			data, err := json.Marshal(ev)
			require.NoError(t, err)
			var out struct {
				Properties json.RawMessage `json:"properties"`
			}
			require.NoError(t, json.Unmarshal(data, &out))
			assert.JSONEq(t, tc.want, string(out.Properties),
				"properties must reach the backend exactly as Azure sent it")
		})
	}

	t.Run("absent stays absent", func(t *testing.T) {
		var ev AzureActivityLogEvent
		require.NoError(t, json.Unmarshal([]byte(`{"category":"Administrative"}`), &ev))
		assert.Nil(t, ev.Properties)

		data, err := json.Marshal(ev)
		require.NoError(t, err)
		assert.NotContains(t, string(data), "properties")
	})
}

func TestEventHubShape_RoundTripsThroughTheAlertPayload(t *testing.T) {
	// The collector embeds the typed event into the alert, so identity has to
	// survive marshal/unmarshal or the backend and UI lose the actor.
	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(eventHubShape), &ev))

	data, err := json.Marshal(ev)
	require.NoError(t, err)

	var out AzureActivityLogEvent
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, "user@example.com", out.EffectiveCaller())
	assert.Equal(t, "Administrative", out.Category)
	require.NotNil(t, out.EffectiveAuthorization())
	assert.Equal(t, "Owner", out.EffectiveAuthorization().Evidence["role"])
}
