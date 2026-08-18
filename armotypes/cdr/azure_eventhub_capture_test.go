package cdr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// azureEventHubCapture is a REAL Azure Activity Log record as delivered to Event
// Hub by a subscription-level diagnostic setting — the only shape the in-account
// collector ever decodes. Captured 2026-08-18 from a live Event Hub (339 records
// off subscription-level Administrative/Security/Policy logs); identifiers,
// names and IPs replaced, opaque token blobs (aio/rh/uti/sid/puid) dropped.
//
// This closes SUB-7951, which could not be resolved from Microsoft's docs alone.
// The POC fixture in azure_test.go carries authorization/claims/caller at the TOP
// LEVEL; that shape does not occur in real Event Hub output. Measured over the
// 339-record capture:
//
//	top-level "authorization" :   0/339
//	top-level "claims"        :   0/339
//	top-level "caller"        :   9/339  (only Microsoft.Sql/servers/* records)
//	"identity.authorization"  : 330/339
//	"identity.claims"         : 330/339
//
// So identity IS nested, as the public streaming schema says, and reading
// Claims/Authorization directly yields nothing on real data. Every consumer must
// go through EffectiveClaims / EffectiveAuthorization / EffectiveCaller.
const azureEventHubCapture = `{
  "records": [
    {
      "RoleLocation": "Central India",
      "Stamp": "FDWorker",
      "ReleaseVersion": "7.2026.30.7+a8879f0.release_2026w30",
      "VmSku": "Standard_DS5_V2",
      "time": "2026-08-17T12:44:11.0994607Z",
      "resourceId": "/SUBSCRIPTIONS/AAAAAAAA-1111-2222-3333-444444444444/RESOURCEGROUPS/ARMO-CDR/PROVIDERS/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/ARMO-CDR-COLLECTOR",
      "operationName": "MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/WRITE",
      "category": "Administrative",
      "resultType": "Start",
      "resultSignature": "Started.",
      "durationMs": "0",
      "callerIpAddress": "203.0.113.10",
      "correlationId": "CCCCCCCC-1111-2222-3333-444444444444",
      "identity": {
        "authorization": {
          "scope": "/subscriptions/AAAAAAAA-1111-2222-3333-444444444444/resourcegroups/armo-cdr/providers/Microsoft.ManagedIdentity/userAssignedIdentities/armo-cdr-collector",
          "action": "Microsoft.ManagedIdentity/userAssignedIdentities/write",
          "evidence": {
            "role": "Owner",
            "roleAssignmentScope": "/subscriptions/AAAAAAAA-1111-2222-3333-444444444444",
            "roleAssignmentId": "4310a6a019ae492eb1e9bdb1488ca37b",
            "roleDefinitionId": "8e3af657a8ff443ca75c2fe8c4bcb635",
            "principalId": "BBBBBBBB111122223333444444444444",
            "principalType": "User"
          }
        },
        "claims": {
          "aud": "https://management.core.windows.net/",
          "iss": "https://sts.windows.net/DDDDDDDD-1111-2222-3333-444444444444/",
          "iat": "1786970122",
          "nbf": "1786970122",
          "exp": "1786975519",
          "http://schemas.microsoft.com/claims/authnclassreference": "1",
          "acrs": "p1",
          "http://schemas.microsoft.com/claims/authnmethodsreferences": "pwd,mfa",
          "appid": "04b07795-8ddb-461a-bbee-02f9e1bf7b46",
          "appidacr": "0",
          "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname": "User",
          "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname": "Example",
          "idtyp": "user",
          "ipaddr": "203.0.113.10",
          "name": "Example User",
          "http://schemas.microsoft.com/identity/claims/objectidentifier": "BBBBBBBB-1111-2222-3333-444444444444",
          "http://schemas.microsoft.com/identity/claims/scope": "user_impersonation",
          "http://schemas.microsoft.com/identity/claims/tenantid": "DDDDDDDD-1111-2222-3333-444444444444",
          "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name": "user@example.com",
          "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn": "user@example.com",
          "ver": "1.0",
          "wids": "b79fbf4d-3ef9-4689-8143-76b194e85509",
          "xms_act_fct": "5 3",
          "xms_idrel": "1 26",
          "xms_sub_fct": "4 3",
          "xms_tcdt": "1511180936"
        }
      },
      "level": "Information",
      "properties": {
        "eventCategory": "Administrative",
        "entity": "/subscriptions/AAAAAAAA-1111-2222-3333-444444444444/resourcegroups/armo-cdr/providers/Microsoft.ManagedIdentity/userAssignedIdentities/armo-cdr-collector",
        "message": "Microsoft.ManagedIdentity/userAssignedIdentities/write",
        "hierarchy": "DDDDDDDD-1111-2222-3333-444444444444/AAAAAAAA-1111-2222-3333-444444444444"
      },
      "tenantId": "DDDDDDDD-1111-2222-3333-444444444444"
    }
  ]
}`

// TestAzureEventHubCapture_IdentityIsNested pins the finding from SUB-7951: on the
// real Event Hub shape the RBAC context and token claims live under "identity",
// and the flat fields the POC fixture used are empty. If Azure ever moves them
// back to the top level this fails loudly rather than silently emptying the actor
// on every Azure alert.
func TestAzureEventHubCapture_IdentityIsNested(t *testing.T) {
	var batch AzureActivityLogBatch
	require.NoError(t, json.Unmarshal([]byte(azureEventHubCapture), &batch))
	require.Len(t, batch.Records, 1)
	e := batch.Records[0]

	// Nested, not flat.
	require.NotNil(t, e.Identity, "real Event Hub records carry an identity object")
	assert.NotNil(t, e.Identity.Authorization)
	assert.NotEmpty(t, e.Identity.Claims)

	// The flat fields the POC fixture relied on are absent on this shape. These
	// assertions are the whole point: reading them directly is the bug.
	assert.Nil(t, e.Authorization, "top-level authorization does not exist on the Event Hub shape")
	assert.Empty(t, e.Claims, "top-level claims does not exist on the Event Hub shape")
	assert.Empty(t, e.Caller, "top-level caller does not exist on the Event Hub shape")
}

// TestAzureEventHubCapture_EffectiveAccessors proves the accessors recover the
// actor from the nested shape — the failure this struct exists to prevent is an
// Azure alert shipping with no principal identity.
func TestAzureEventHubCapture_EffectiveAccessors(t *testing.T) {
	var batch AzureActivityLogBatch
	require.NoError(t, json.Unmarshal([]byte(azureEventHubCapture), &batch))
	e := batch.Records[0]

	assert.Equal(t, "Example User", e.EffectiveCaller())
	require.NotNil(t, e.EffectiveAuthorization())
	assert.Equal(t, "Microsoft.ManagedIdentity/userAssignedIdentities/write", e.EffectiveAuthorization().Action)
	assert.NotEmpty(t, e.EffectiveAuthorization().Evidence, "RBAC evidence (role, principalId) must survive")

	claims := e.EffectiveClaims()
	assert.Equal(t, "user", claims["idtyp"])
	assert.Equal(t, "BBBBBBBB-1111-2222-3333-444444444444",
		claims["http://schemas.microsoft.com/identity/claims/objectidentifier"])

	// EffectiveCaller prefers the short "upn" claim, which was present in 0/339
	// real records — resolution happens via "name". Kept as the first preference
	// because it is harmless, but do not assume it ever fires on Event Hub data.
	assert.Empty(t, claims["upn"], "short-form upn claim is absent on real Event Hub records")
	assert.NotEmpty(t, claims["name"], "name is what EffectiveCaller actually resolves through")
}

// TestAzureEventHubCapture_SubscriptionMustBeDerived documents that the Event Hub
// shape carries no top-level subscriptionId (empty on 330/339 captured records),
// so batching has to recover it from resourceId — which was parseable on 339/339.
func TestAzureEventHubCapture_SubscriptionMustBeDerived(t *testing.T) {
	var batch AzureActivityLogBatch
	require.NoError(t, json.Unmarshal([]byte(azureEventHubCapture), &batch))
	e := batch.Records[0]

	assert.Empty(t, e.SubscriptionID, "Event Hub shape has no subscriptionId; derive it from resourceId")
	assert.Contains(t, e.ResourceID, "/SUBSCRIPTIONS/", "resourceId is the only subscription source")

	// Modelled but never delivered on this shape - see SUB-7951 findings.
	assert.Empty(t, e.Location, "Activity Log records carry no resource region at all")
	assert.Empty(t, e.Channels, "channels is absent on the Event Hub shape")
	assert.Empty(t, e.ResourceGroupName, "resourceGroupName is absent; Azure sends ResourceGroup=DummyValue on the few records that have it")
}

// TestAzureEventHubCapture_PropertiesSurviveRaw guards the deliberate choice to
// keep Properties as json.RawMessage: Activity Log sends it as an object for most
// operations and as a JSON STRING for some, and a typed map would fail the whole
// record decode over a field no rule references.
func TestAzureEventHubCapture_PropertiesSurviveRaw(t *testing.T) {
	var batch AzureActivityLogBatch
	require.NoError(t, json.Unmarshal([]byte(azureEventHubCapture), &batch))
	assert.NotEmpty(t, batch.Records[0].Properties)

	// The string form must decode too, not just the object form.
	const asString = `{"records":[{"operationName":"X","properties":"{\"statusCode\":\"OK\"}"}]}`
	var b2 AzureActivityLogBatch
	require.NoError(t, json.Unmarshal([]byte(asString), &b2), "properties-as-string must not fail the record")
	assert.NotEmpty(t, b2.Records[0].Properties)
}
