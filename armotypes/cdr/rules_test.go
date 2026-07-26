package cdr

import (
	"encoding/json"
	"testing"
	"time"
)

// TestCdrRuleBundleRoundTrip verifies the rule-delivery response survives a JSON round-trip
// unchanged — the collector must decode exactly what the endpoint encodes.
func TestCdrRuleBundleRoundTrip(t *testing.T) {
	original := CdrRuleBundle{
		Version:     "sha256:abc123",
		Provider:    Azure,
		GeneratedAt: time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC),
		Rules: []CdrRule{
			{
				RuleID:         "azure-role-assignment-created",
				Name:           "Privileged role assignment created",
				Service:        ActivityLogs,
				Expression:     `event.operationName == "Microsoft.Authorization/roleAssignments/write"`,
				Priority:       "high",
				MitreTactic:    "Privilege Escalation",
				MitreTechnique: "T1098",
				Tags:           []string{"iam", "azure"},
				Message:        "Role assignment created",
				Origin:         CdrRuleOriginManaged,
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded CdrRuleBundle
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Version != original.Version || decoded.Provider != original.Provider ||
		!decoded.GeneratedAt.Equal(original.GeneratedAt) || len(decoded.Rules) != 1 {
		t.Fatalf("bundle mismatch: %+v", decoded)
	}
	got, want := decoded.Rules[0], original.Rules[0]
	if got.RuleID != want.RuleID || got.Service != want.Service || got.Expression != want.Expression ||
		got.Priority != want.Priority || got.Origin != want.Origin || len(got.Tags) != len(want.Tags) {
		t.Errorf("rule mismatch:\n got  %+v\n want %+v", got, want)
	}
}

// TestOmitEmptyOptionalFields verifies empty optional fields drop off the wire.
func TestOmitEmptyOptionalFields(t *testing.T) {
	data, err := json.Marshal(CdrRule{RuleID: "r1", Name: "n", Service: CloudTrail, Expression: "true"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{"origin", "priority", "tags", "message", "description"} {
		if _, present := m[k]; present {
			t.Errorf("expected %q omitted when empty, got: %s", k, data)
		}
	}
}
