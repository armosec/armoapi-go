package cdr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// gcpSyncAuditLogSample is the real synchronous Admin Activity entry captured in
// the gcp-cdr-poc (FINDINGS.md §1d, redacted) — kept UNTRIMMED, including
// authorizationInfo[].resourceAttributes and requestMetadata.requestAttributes,
// so TestGcpAuditLogFieldCoverage can assert we model every field the real
// capture carries. Synchronous admin calls carry NO top-level "operation" block —
// the majority case the detection gate `!has(operation) || operation.last` must
// let through. Note resourceName is the PARENT project while the created service
// account is only in response.name, and status is {} on success.
const gcpSyncAuditLogSample = `{
  "insertId": "1ort6tqf2s9qz6",
  "logName": "projects/cdr-project-503613/logs/cloudaudit.googleapis.com%2Factivity",
  "severity": "NOTICE",
  "timestamp": "2026-07-26T13:42:27.024755819Z",
  "receiveTimestamp": "2026-07-26T13:42:28.100000000Z",
  "resource": { "type": "service_account", "labels": { "project_id": "cdr-project-503613" } },
  "protoPayload": {
    "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
    "serviceName": "iam.googleapis.com",
    "methodName": "google.iam.admin.v1.CreateServiceAccount",
    "resourceName": "projects/cdr-project-503613",
    "status": {},
    "authenticationInfo": {
      "principalEmail": "attacker@armosec.io",
      "principalSubject": "user:attacker@armosec.io",
      "oauthInfo": { "oauthClientId": "32555940559.apps.googleusercontent.com" }
    },
    "authorizationInfo": [
      { "granted": true, "permission": "iam.serviceAccounts.create",
        "permissionType": "ADMIN_WRITE", "resource": "projects/cdr-project-503613",
        "resourceAttributes": { "type": "iam.googleapis.com/ServiceAccount" } }
    ],
    "requestMetadata": {
      "callerIp": "199.203.132.136",
      "callerSuppliedUserAgent": "google-cloud-sdk gcloud/552.0.0",
      "requestAttributes": { "time": "2026-07-26T13:42:27.024755819Z" }
    },
    "request": {
      "@type": "type.googleapis.com/google.iam.admin.v1.CreateServiceAccountRequest",
      "account_id": "cdr-poc-ss1", "name": "projects/cdr-project-503613"
    },
    "response": {
      "@type": "type.googleapis.com/google.iam.admin.v1.ServiceAccount",
      "name": "projects/cdr-project-503613/serviceAccounts/cdr-poc-ss5@cdr-project-503613.iam.gserviceaccount.com"
    }
  }
}`

// gcpLongRunningCompletionSample is the completion entry of a long-running
// operation captured in the POC (FINDINGS.md §1b/§1c). It carries an "operation"
// block with last=true (and NO "first" key — booleans are omitted when false).
// This is the entry the gate keeps; it retains full authenticationInfo, and for
// networks.insert its resourceName IS the target resource.
const gcpLongRunningCompletionSample = `{
  "insertId": "3byqwkd4f8m",
  "logName": "projects/cdr-project-503613/logs/cloudaudit.googleapis.com%2Factivity",
  "severity": "NOTICE",
  "timestamp": "2026-07-26T13:33:09.401394Z",
  "resource": { "type": "gce_network", "labels": { "project_id": "cdr-project-503613" } },
  "operation": {
    "id": "operation-1785072775774-65783a4b2c062-c7f6a912-641fdc82",
    "producer": "compute.googleapis.com",
    "last": true
  },
  "protoPayload": {
    "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
    "serviceName": "compute.googleapis.com",
    "methodName": "v1.compute.networks.insert",
    "resourceName": "projects/cdr-project-503613/global/networks/cdr-poc-lro-net",
    "status": {},
    "authenticationInfo": {
      "principalEmail": "attacker@armosec.io",
      "principalSubject": "user:attacker@armosec.io"
    },
    "authorizationInfo": [
      { "granted": true, "permission": "compute.networks.create",
        "permissionType": "ADMIN_WRITE", "resource": "projects/cdr-project-503613/global/networks/cdr-poc-lro-net" }
    ]
  }
}`

// gcpLongRunningRequestSample is the request (first) entry of the same
// long-running operation — same operation.id, first=true, no "last" key. The
// gate skips it.
const gcpLongRunningRequestSample = `{
  "insertId": "1yg3nhe1wbpq",
  "logName": "projects/cdr-project-503613/logs/cloudaudit.googleapis.com%2Factivity",
  "timestamp": "2026-07-26T13:32:55.853384Z",
  "resource": { "type": "gce_network", "labels": { "project_id": "cdr-project-503613" } },
  "operation": {
    "id": "operation-1785072775774-65783a4b2c062-c7f6a912-641fdc82",
    "producer": "compute.googleapis.com",
    "first": true
  },
  "protoPayload": {
    "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
    "serviceName": "compute.googleapis.com",
    "methodName": "v1.compute.networks.insert",
    "resourceName": "projects/cdr-project-503613/global/networks/cdr-poc-lro-net",
    "authenticationInfo": {
      "principalEmail": "attacker@armosec.io",
      "principalSubject": "user:attacker@armosec.io"
    }
  }
}`

// gcpEventGate mirrors the HLD §3.4 detection gate `!has(operation) || operation.last`.
func gcpEventGate(e *GcpAuditLogEvent) bool {
	return e.Operation == nil || e.Operation.Last
}

// TestGcpAuditLogRoundTrip round-trips real synchronous and long-running samples
// from the POC and asserts the POC-measured corrections hold: operation is nil on
// synchronous events, non-nil with last=true on the completion entry, the
// detection gate behaves, response.name carries the create target that
// resourceName does not, and marshal/unmarshal is lossless.
func TestGcpAuditLogRoundTrip(t *testing.T) {
	tests := []struct {
		name             string
		sample           string
		wantHasOperation bool
		wantLast         bool
		wantGatePass     bool
		wantMethod       string
		wantService      string
		wantPrincipal    string
	}{
		{
			name:             "synchronous CreateServiceAccount (no operation block)",
			sample:           gcpSyncAuditLogSample,
			wantHasOperation: false,
			wantGatePass:     true, // !has(operation) => alerts
			wantMethod:       "google.iam.admin.v1.CreateServiceAccount",
			wantService:      "iam.googleapis.com",
			wantPrincipal:    "attacker@armosec.io",
		},
		{
			name:             "long-running networks.insert request entry (first)",
			sample:           gcpLongRunningRequestSample,
			wantHasOperation: true,
			wantLast:         false,
			wantGatePass:     false, // first entry is skipped
			wantMethod:       "v1.compute.networks.insert",
			wantService:      "compute.googleapis.com",
			wantPrincipal:    "attacker@armosec.io",
		},
		{
			name:             "long-running networks.insert completion entry (last)",
			sample:           gcpLongRunningCompletionSample,
			wantHasOperation: true,
			wantLast:         true,
			wantGatePass:     true, // completion entry alerts
			wantMethod:       "v1.compute.networks.insert",
			wantService:      "compute.googleapis.com",
			wantPrincipal:    "attacker@armosec.io",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e GcpAuditLogEvent
			require.NoError(t, json.Unmarshal([]byte(tc.sample), &e))

			require.NotNil(t, e.ProtoPayload)
			assert.Equal(t, tc.wantMethod, e.ProtoPayload.MethodName)
			assert.Equal(t, tc.wantService, e.ProtoPayload.ServiceName)

			require.NotNil(t, e.ProtoPayload.AuthenticationInfo, "completion/sync entries must retain identity")
			assert.Equal(t, tc.wantPrincipal, e.ProtoPayload.AuthenticationInfo.PrincipalEmail)
			assert.True(t, strings.HasPrefix(e.ProtoPayload.AuthenticationInfo.PrincipalSubject, "user:"),
				"principalSubject should carry the user:/serviceAccount: prefix")

			// operation is nilable — the load-bearing POC correction.
			if tc.wantHasOperation {
				require.NotNil(t, e.Operation)
				assert.Equal(t, tc.wantLast, e.Operation.Last)
				assert.NotEmpty(t, e.Operation.ID)
			} else {
				assert.Nil(t, e.Operation, "synchronous events carry no operation block")
			}

			// Detection gate `!has(operation) || operation.last`.
			assert.Equal(t, tc.wantGatePass, gcpEventGate(&e), "detection gate mismatch")

			// Round-trip back to JSON must not error.
			_, err := json.Marshal(e)
			require.NoError(t, err)
		})
	}
}

// TestGcpAuditLogFieldCoverage decodes the untrimmed real capture with
// DisallowUnknownFields, so the struct is proven to model every field the real
// GCP event carries. This is the guard the plain round-trip lacks: encoding/json
// silently drops unmodeled fields, so without this a wire field we forgot to
// model would vanish undetected. If GCP (or a fixture update) introduces a field
// we don't model, this fails loudly — the signal to extend the struct.
func TestGcpAuditLogFieldCoverage(t *testing.T) {
	dec := json.NewDecoder(strings.NewReader(gcpSyncAuditLogSample))
	dec.DisallowUnknownFields()
	var e GcpAuditLogEvent
	require.NoError(t, dec.Decode(&e), "an unmodeled field is present in the real capture — extend GcpAuditLogEvent")

	// The two fields Ben's review flagged as silently dropped are now modeled.
	require.NotNil(t, e.ProtoPayload)
	require.Len(t, e.ProtoPayload.AuthorizationInfo, 1)
	require.NotNil(t, e.ProtoPayload.AuthorizationInfo[0].ResourceAttributes)
	assert.Equal(t, "iam.googleapis.com/ServiceAccount", e.ProtoPayload.AuthorizationInfo[0].ResourceAttributes.Type)
	require.NotNil(t, e.ProtoPayload.RequestMetadata)
	assert.NotEmpty(t, e.ProtoPayload.RequestMetadata.RequestAttributes["time"])
}

// gcpImpersonationSample is a synthetic (NOT POC-captured) Admin Activity entry
// modeling a service-account impersonation + static-key call — the canonical GCP
// privilege-escalation path. The POC authenticated as a human via gcloud so it
// never exercised these fields; this fixture guards that the struct models them.
const gcpImpersonationSample = `{
  "insertId": "imp1",
  "logName": "projects/cdr-project-503613/logs/cloudaudit.googleapis.com%2Factivity",
  "timestamp": "2026-07-26T14:00:00Z",
  "resource": { "type": "service_account", "labels": { "project_id": "cdr-project-503613" } },
  "protoPayload": {
    "@type": "type.googleapis.com/google.cloud.audit.AuditLog",
    "serviceName": "iamcredentials.googleapis.com",
    "methodName": "GenerateAccessToken",
    "resourceName": "projects/-/serviceAccounts/victim@cdr-project-503613.iam.gserviceaccount.com",
    "authenticationInfo": {
      "principalEmail": "victim@cdr-project-503613.iam.gserviceaccount.com",
      "principalSubject": "serviceAccount:victim@cdr-project-503613.iam.gserviceaccount.com",
      "serviceAccountKeyName": "projects/cdr-project-503613/serviceAccounts/attacker@cdr-project-503613.iam.gserviceaccount.com/keys/abc123",
      "serviceAccountDelegationInfo": [
        { "principalSubject": "user:attacker@armosec.io",
          "firstPartyPrincipal": { "principalEmail": "attacker@armosec.io" } }
      ]
    },
    "authorizationInfo": [
      { "granted": false, "permission": "iam.serviceAccounts.getAccessToken",
        "permissionType": "ADMIN_WRITE",
        "resource": "projects/-/serviceAccounts/victim@cdr-project-503613.iam.gserviceaccount.com" }
    ],
    "policyViolationInfo": { "orgPolicyViolationInfo": { "violationInfo": [] } }
  }
}`

// TestGcpImpersonationFields verifies the privilege-escalation fields Ben's review
// added (serviceAccountDelegationInfo, serviceAccountKeyName) parse, and that a
// denied attempt (granted:false, omitted on the wire) reads as false.
func TestGcpImpersonationFields(t *testing.T) {
	var e GcpAuditLogEvent
	require.NoError(t, json.Unmarshal([]byte(gcpImpersonationSample), &e))
	require.NotNil(t, e.ProtoPayload)

	ai := e.ProtoPayload.AuthenticationInfo
	require.NotNil(t, ai)
	assert.Equal(t, "projects/cdr-project-503613/serviceAccounts/attacker@cdr-project-503613.iam.gserviceaccount.com/keys/abc123", ai.ServiceAccountKeyName)
	require.Len(t, ai.ServiceAccountDelegationInfo, 1)
	assert.Equal(t, "user:attacker@armosec.io", ai.ServiceAccountDelegationInfo[0].PrincipalSubject)
	require.NotNil(t, ai.ServiceAccountDelegationInfo[0].FirstPartyPrincipal)
	assert.Equal(t, "attacker@armosec.io", ai.ServiceAccountDelegationInfo[0].FirstPartyPrincipal.PrincipalEmail)

	// Denied attempt: granted omitted-when-false => reads false (denied-attempt rule signal).
	require.Len(t, e.ProtoPayload.AuthorizationInfo, 1)
	assert.False(t, e.ProtoPayload.AuthorizationInfo[0].Granted)

	// policyViolationInfo present (org-policy/VPC-SC denial detail).
	assert.NotNil(t, e.ProtoPayload.PolicyViolationInfo)

	// Round-trips without error.
	_, err := json.Marshal(e)
	require.NoError(t, err)
}

// TestGcpNumResponseItemsBothForms guards that numResponseItems decodes from both
// wire forms — protojson's quoted string ("3") and a bare number (3, spec-legal on
// parse). A plain string type would fail the whole-event decode on the bare form.
func TestGcpNumResponseItemsBothForms(t *testing.T) {
	for _, in := range []string{
		`{"protoPayload":{"numResponseItems":"3"}}`,
		`{"protoPayload":{"numResponseItems":3}}`,
	} {
		var e GcpAuditLogEvent
		require.NoError(t, json.Unmarshal([]byte(in), &e), "decode failed for %s", in)
		require.NotNil(t, e.ProtoPayload)
		assert.Equal(t, "3", e.ProtoPayload.NumResponseItems.String())
	}
}

// TestGcpAuditLogSyncSpecifics pins the two GCP quirks that are easy to regress:
// resourceName is the parent (real target in response.name) for create-style
// methods, and status is empty ({}) on success so status.code is absent.
func TestGcpAuditLogSyncSpecifics(t *testing.T) {
	var e GcpAuditLogEvent
	require.NoError(t, json.Unmarshal([]byte(gcpSyncAuditLogSample), &e))
	require.NotNil(t, e.ProtoPayload)

	// resourceName is the PARENT; the created SA is only in response.name.
	assert.Equal(t, "projects/cdr-project-503613", e.ProtoPayload.ResourceName)
	require.NotNil(t, e.ProtoPayload.Response)
	assert.Contains(t, e.ProtoPayload.Response["name"], "serviceAccounts/")

	// status {} on success => present but code absent (Go zero-value, nil-safe).
	require.NotNil(t, e.ProtoPayload.Status)
	assert.Equal(t, 0, e.ProtoPayload.Status.Code)

	// authorizationInfo enables denied-attempt rules; ADMIN_WRITE + granted.
	require.Len(t, e.ProtoPayload.AuthorizationInfo, 1)
	assert.Equal(t, "ADMIN_WRITE", e.ProtoPayload.AuthorizationInfo[0].PermissionType)
	assert.True(t, e.ProtoPayload.AuthorizationInfo[0].Granted)
}

// TestGcpNilSafety verifies CEL-style has()-guarded access is panic-free when the
// per-service optional fields GCP omits are absent (empty event, and the minimal
// long-running request entry).
func TestGcpNilSafety(t *testing.T) {
	// A zero event: nothing set. Guarded access must not panic.
	var empty GcpAuditLogEvent
	assert.True(t, gcpEventGate(&empty), "empty event has no operation => gate passes")
	assert.Nil(t, empty.ProtoPayload)
	assert.Nil(t, empty.Resource)

	// The minimal request entry has operation but no authorizationInfo / status /
	// requestMetadata / response.
	var e GcpAuditLogEvent
	require.NoError(t, json.Unmarshal([]byte(gcpLongRunningRequestSample), &e))
	require.NotNil(t, e.Operation)
	assert.True(t, e.Operation.First)
	assert.False(t, e.Operation.Last)
	require.NotNil(t, e.ProtoPayload)
	assert.Nil(t, e.ProtoPayload.Status)
	assert.Nil(t, e.ProtoPayload.RequestMetadata)
	assert.Empty(t, e.ProtoPayload.AuthorizationInfo)
}

// TestGcpOperationBooleanOmittedWhenFalse guards the POC finding that operation
// booleans are omitted from the wire when false (not serialized as false), so
// re-marshalling the completion entry must NOT emit "first".
func TestGcpOperationBooleanOmittedWhenFalse(t *testing.T) {
	var e GcpAuditLogEvent
	require.NoError(t, json.Unmarshal([]byte(gcpLongRunningCompletionSample), &e))
	require.NotNil(t, e.Operation)

	b, err := json.Marshal(e.Operation)
	require.NoError(t, err)
	assert.NotContains(t, string(b), `"first"`, "false booleans must be omitted, not serialized")
	assert.Contains(t, string(b), `"last":true`)
}

// TestGcpEventDataEmbedding verifies the GCP event embeds in the shared CdrAlert
// contract without disturbing the AWS or Azure paths (regression guard).
func TestGcpEventDataEmbedding(t *testing.T) {
	alert := CdrAlert{
		CloudMetadata: CloudMetadata{Provider: GCP, SourceService: CloudAuditLogs},
		EventData: EventData{GcpAuditLog: &GcpAuditLogEvent{
			ProtoPayload: &GcpAuditLogPayload{MethodName: "test"},
		}},
	}
	b, err := json.Marshal(alert)
	require.NoError(t, err)

	var got CdrAlert
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, GCP, got.Provider)
	assert.Equal(t, CloudAuditLogs, got.SourceService)
	require.NotNil(t, got.GcpAuditLog)
	require.NotNil(t, got.GcpAuditLog.ProtoPayload)
	assert.Equal(t, "test", got.GcpAuditLog.ProtoPayload.MethodName)
	assert.Nil(t, got.AWSCloudTrail, "AWS path must stay nil")
	assert.Nil(t, got.AzureActivityLog, "Azure path must stay nil")
}

// TestGcpProviderEnumsUnchanged is a regression guard: adding GCP must not alter
// the existing AWS/Azure enum values or the GCP additions themselves.
func TestGcpProviderEnumsUnchanged(t *testing.T) {
	assert.Equal(t, CloudProvider("aws"), AWS)
	assert.Equal(t, CloudProvider("azure"), Azure)
	assert.Equal(t, CloudProvider("gcp"), GCP)
	assert.Equal(t, CloudService("cloudtrail"), CloudTrail)
	assert.Equal(t, CloudService("activitylogs"), ActivityLogs)
	assert.Equal(t, CloudService("cloudauditlogs"), CloudAuditLogs)
}

// TestCdrAlertGcpJsonPath verifies the GCP project-ID JSON-path constant matches
// the json tags on CdrAlert, so config-service account matching doesn't drift
// from the wire contract (mirrors TestCdrAlertJsonPath / TestCdrAlertAzureJsonPath).
func TestCdrAlertGcpJsonPath(t *testing.T) {
	tests := []struct {
		Name          string
		Path          string
		CdrAlert      CdrAlert
		ExpectedValue string
	}{
		{
			Name: "Test CdrEventGcpProjectIDJsonPath",
			Path: CdrEventGcpProjectIDJsonPath,
			CdrAlert: CdrAlert{
				EventData: EventData{
					GcpAuditLog: &GcpAuditLogEvent{
						Resource: &GcpMonitoredResource{
							Labels: map[string]string{"project_id": "cdr-project-503613"},
						},
					},
				},
			},
			ExpectedValue: "cdr-project-503613",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// 1. Marshal to JSON
			data, err := json.Marshal(test.CdrAlert)
			require.NoError(t, err)

			// 2. Unmarshal into a generic map
			var genericCdrAlert map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &genericCdrAlert))

			// 3. Traverse the path (shared helper lives in aws_test.go)
			pathWithoutPrefix := strings.TrimPrefix(test.Path, "cdrevent.")
			val, ok := getValueAtJsonPath(genericCdrAlert, pathWithoutPrefix)
			require.True(t, ok)
			assert.Equal(t, test.ExpectedValue, val)
		})
	}
}
