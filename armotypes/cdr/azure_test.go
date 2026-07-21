package cdr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// azureActivityLogSample is a trimmed real Azure Activity Log record captured in
// the POC (Event Hub `records[]` envelope; identity claims abbreviated). See the
// azure-cdr-poc FINDINGS.md for the full sample.
const azureActivityLogSample = `{
  "records": [
    {
      "time": "2026-07-19T14:14:22Z",
      "resourceId": "/SUBSCRIPTIONS/8FC00071/RESOURCEGROUPS/CDR-POC-RG",
      "operationName": "Microsoft.Resources/subscriptions/resourcegroups/write",
      "category": "Administrative",
      "resultType": "Success",
      "resultSignature": "Succeeded.Created",
      "correlationId": "abc-123",
      "caller": "benm@armosec.io",
      "level": "Information",
      "channels": "Operation",
      "tenantId": "50a70646-52e3-4e46-911e-6ca1b46afba3",
      "authorization": {
        "action": "Microsoft.Resources/subscriptions/resourcegroups/write",
        "scope": "/subscriptions/8fc00071/resourcegroups/cdr-poc-rg"
      },
      "claims": {
        "idtyp": "user",
        "appid": "04b07795-8ddb-461a-bbee-02f9e1bf7b46",
        "ipaddr": "199.203.132.136",
        "http://schemas.microsoft.com/identity/claims/objectidentifier": "0dfe9580"
      },
      "properties": {
        "statusCode": "Created",
        "eventCategory": "Administrative",
        "entity": "/subscriptions/8fc00071/resourcegroups/cdr-poc-rg",
        "message": "Microsoft.Resources/subscriptions/resourcegroups/write"
      }
    }
  ]
}`

func TestAzureActivityLogRoundTrip(t *testing.T) {
	var batch AzureActivityLogBatch
	if err := json.Unmarshal([]byte(azureActivityLogSample), &batch); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(batch.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(batch.Records))
	}

	e := batch.Records[0]
	if e.OperationName != "Microsoft.Resources/subscriptions/resourcegroups/write" {
		t.Errorf("operationName = %q", e.OperationName)
	}
	if e.Category != "Administrative" {
		t.Errorf("category = %q", e.Category)
	}
	if e.Caller != "benm@armosec.io" {
		t.Errorf("caller = %q", e.Caller)
	}
	if e.TenantID != "50a70646-52e3-4e46-911e-6ca1b46afba3" {
		t.Errorf("tenantId = %q", e.TenantID)
	}
	if e.Authorization == nil || e.Authorization.Action == "" {
		t.Fatalf("authorization not parsed: %+v", e.Authorization)
	}
	if e.Claims["idtyp"] != "user" {
		t.Errorf("claims.idtyp = %q", e.Claims["idtyp"])
	}
	if e.Properties["message"] != "Microsoft.Resources/subscriptions/resourcegroups/write" {
		t.Errorf("properties.message = %v", e.Properties["message"])
	}

	// Round-trip back to JSON.
	if _, err := json.Marshal(e); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

// TestAzureEventDataEmbedding verifies the Azure event embeds in the shared
// CdrAlert contract without disturbing the AWS path.
func TestAzureEventDataEmbedding(t *testing.T) {
	alert := CdrAlert{
		CloudMetadata: CloudMetadata{Provider: Azure, SourceService: ActivityLogs},
		EventData:     EventData{AzureActivityLog: &AzureActivityLogEvent{OperationName: "test"}},
	}
	b, err := json.Marshal(alert)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got CdrAlert
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Provider != Azure || got.SourceService != ActivityLogs {
		t.Errorf("cloud metadata = %+v", got.CloudMetadata)
	}
	if got.AzureActivityLog == nil || got.AzureActivityLog.OperationName != "test" {
		t.Errorf("azure event data = %+v", got.AzureActivityLog)
	}
	if got.AWSCloudTrail != nil {
		t.Errorf("AWS path should be nil, got %+v", got.AWSCloudTrail)
	}
}

// TestCdrAlertAzureJsonPath verifies the Azure subscription/tenant JSON-path
// constants match the json tags on CdrAlert, so config-service lookups don't
// drift from the wire contract (mirrors TestCdrAlertJsonPath for AWS).
func TestCdrAlertAzureJsonPath(t *testing.T) {
	tests := []struct {
		Name          string
		Path          string
		CdrAlert      CdrAlert
		ExpectedValue string
	}{
		{
			Name: "Test CdrEventAzureSubscriptionIDJsonPath",
			Path: CdrEventAzureSubscriptionIDJsonPath,
			CdrAlert: CdrAlert{
				EventData: EventData{
					AzureActivityLog: &AzureActivityLogEvent{
						SubscriptionID: "8fc00071-75e6-4b3f-82d4-844e7865bab3",
					},
				},
			},
			ExpectedValue: "8fc00071-75e6-4b3f-82d4-844e7865bab3",
		},
		{
			Name: "Test CdrEventAzureTenantIDJsonPath",
			Path: CdrEventAzureTenantIDJsonPath,
			CdrAlert: CdrAlert{
				EventData: EventData{
					AzureActivityLog: &AzureActivityLogEvent{
						TenantID: "50a70646-52e3-4e46-911e-6ca1b46afba3",
					},
				},
			},
			ExpectedValue: "50a70646-52e3-4e46-911e-6ca1b46afba3",
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
