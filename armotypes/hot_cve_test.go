package armotypes

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHotCVE_Validate_FromJSON(t *testing.T) {
	data, err := os.ReadFile("testdata/hot_cve_example.json")
	require.NoError(t, err, "failed to read example JSON")

	var hotCVEs []HotCVE
	err = json.Unmarshal(data, &hotCVEs)
	require.NoError(t, err, "failed to unmarshal example JSON")
	require.Len(t, hotCVEs, 2, "expected 2 hot CVEs in example")

	for i, cve := range hotCVEs {
		t.Run(cve.CVEID, func(t *testing.T) {
			err := cve.Validate()
			assert.NoError(t, err, "hot CVE at index %d should be valid", i)
			assert.NotEmpty(t, cve.CVEID)
			assert.NotEmpty(t, cve.Severity)
			assert.NotEmpty(t, cve.AffectedPackages)
			for _, pkg := range cve.AffectedPackages {
				assert.NotEmpty(t, pkg.PackageName)
			}
		})
	}

	// Verify specific fields from the example
	assert.Equal(t, "CVE-2024-3094", hotCVEs[0].CVEID)
	assert.Equal(t, "xz/liblzma backdoor", hotCVEs[0].Title)
	assert.Equal(t, "Critical", hotCVEs[0].Severity)
	assert.Equal(t, 900, hotCVEs[0].SeverityScore)
	assert.Len(t, hotCVEs[0].AffectedPackages, 4)
	assert.Equal(t, "xz-utils", hotCVEs[0].AffectedPackages[0].PackageName)
	assert.Equal(t, []string{"deb"}, hotCVEs[0].AffectedPackages[0].PackageTypes)
	assert.Equal(t, "5.6.0", hotCVEs[0].AffectedPackages[0].VersionStart)
	assert.Equal(t, "5.6.2", hotCVEs[0].AffectedPackages[0].VersionEnd)
	assert.Equal(t, "5.6.2", hotCVEs[0].AffectedPackages[0].FixedVersion)
}

func TestHotCVE_Validate(t *testing.T) {
	validPkg := HotCVEAffectedPackage{
		PackageName:  "xz-utils",
		VersionStart: "5.6.0",
		VersionEnd:   "5.6.2",
	}

	tests := []struct {
		name    string
		cve     HotCVE
		wantErr string
	}{
		{
			name: "valid hot CVE",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "Critical",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "valid with status active",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "High",
				Status:           HotCVEStatusActive,
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "valid with status inactive",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "Medium",
				Status:           HotCVEStatusInactive,
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "missing cveId",
			cve: HotCVE{
				Severity:         "Critical",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
			wantErr: "cveId is required",
		},
		{
			name: "missing severity",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
			wantErr: "severity is required",
		},
		{
			name: "empty affected packages",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "Critical",
				AffectedPackages: []HotCVEAffectedPackage{},
			},
			wantErr: "affectedPackages is required",
		},
		{
			name: "nil affected packages",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "Critical",
			},
			wantErr: "affectedPackages is required",
		},
		{
			name: "invalid status",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "Critical",
				Status:           "pending",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
			wantErr: "invalid status",
		},
		{
			name: "invalid affected package — missing packageName",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "Critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{VersionStart: "1.0.0"},
				},
			},
			wantErr: "affectedPackages[0]: packageName is required",
		},
		{
			name: "invalid affected package — no version constraint",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "Critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{PackageName: "xz-utils"},
				},
			},
			wantErr: "at least one version constraint",
		},
		{
			name: "valid with exact version only",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "Critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{PackageName: "xz-utils", VersionExact: []string{"5.6.1"}},
				},
			},
		},
		{
			name: "valid with start only (no upper bound)",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "Critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{PackageName: "xz-utils", VersionStart: "5.6.0"},
				},
			},
		},

		// Severity whitelist — the postgres-connector view upsert silently
		// filters out anything other than titlecase Critical/High/Medium/Low/
		// Unknown/Negligible (see hotCVEValidSeverities doc). These cases are
		// the regression guards for the SUB-7201 bug where lowercase "critical"
		// was accepted and then silently dropped by the SQL join, resulting in
		// zero is_hot_cve=true rows across 260k vulnerabilities_cves rows in
		// dev until this validator was tightened.
		{
			name:    "rejects lowercase critical",
			cve:     HotCVE{CVEID: "CVE-2024-3094", Severity: "critical", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
			wantErr: "invalid severity",
		},
		{
			name:    "rejects titlecase-but-nonmember Informational",
			cve:     HotCVE{CVEID: "CVE-2024-3094", Severity: "Informational", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
			wantErr: "invalid severity",
		},
		{
			name:    "rejects whitespace-padded Critical",
			cve:     HotCVE{CVEID: "CVE-2024-3094", Severity: " Critical", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
			wantErr: "invalid severity",
		},
		{
			name: "accepts Negligible (scanners produce this)",
			cve:  HotCVE{CVEID: "CVE-2024-3094", Severity: "Negligible", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
		},
		{
			name: "accepts Unknown",
			cve:  HotCVE{CVEID: "CVE-2024-3094", Severity: "Unknown", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
		},
		{
			name: "accepts Low",
			cve:  HotCVE{CVEID: "CVE-2024-3094", Severity: "Low", AffectedPackages: []HotCVEAffectedPackage{validPkg}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cve.Validate()
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestHotCVEOnFinishedMessage_RoundTrip verifies a message with all current
// fields (including the SUB-7201 follow-up additions: Title, References,
// AffectedWorkloads) serializes and deserializes losslessly. The UNS
// dispatcher reads References[0] for CVELink and AffectedWorkloads for
// per-workload Slack/Teams context — any silent JSON tag rename here
// breaks the production hot-CVE alert rendering.
func TestHotCVEOnFinishedMessage_RoundTrip(t *testing.T) {
	original := HotCVEOnFinishedMessage{
		CustomerGUID: "tenant-1",
		CVEID:        "CVE-2024-9999",
		Severity:     "Critical",
		Title:        "Demo hot CVE",
		References: []string{
			"https://nvd.nist.gov/vuln/detail/CVE-2024-9999",
			"https://example.com/advisory",
		},
		AffectedWorkloads: []HotCVEAffectedWorkload{
			{Cluster: "prod-eu", Namespace: "checkout", WorkloadKind: "Deployment", WorkloadName: "api", ImageTag: "registry.example/app:v1"},
			{Cluster: "prod-us", Namespace: "checkout", WorkloadKind: "Deployment", WorkloadName: "api", ImageTag: "registry.example/app:v1"},
		},
	}
	b, err := json.Marshal(original)
	require.NoError(t, err)

	var got HotCVEOnFinishedMessage
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, original, got)
}

// TestHotCVEOnFinishedMessage_BackwardCompatible_DecodesOldPayload guards
// the contract that a message produced by an older publisher (before Title /
// References / AffectedWorkloads were added) still decodes correctly.
// Otherwise rolling out the new struct breaks readers of in-flight messages
// from un-upgraded publishers.
func TestHotCVEOnFinishedMessage_BackwardCompatible_DecodesOldPayload(t *testing.T) {
	old := []byte(`{"customerGUID":"tenant-1","cveId":"CVE-2024-9999","severity":"Critical"}`)
	var got HotCVEOnFinishedMessage
	require.NoError(t, json.Unmarshal(old, &got))
	assert.Equal(t, "tenant-1", got.CustomerGUID)
	assert.Equal(t, "CVE-2024-9999", got.CVEID)
	assert.Equal(t, "Critical", got.Severity)
	assert.Empty(t, got.Title)
	assert.Empty(t, got.References)
	assert.Empty(t, got.AffectedWorkloads)
}

// TestHotCVEOnFinishedMessage_OmitemptyOnAddedFields guards the JSON tags:
// the three additions are `omitempty` so a publisher that doesn't populate
// them produces a message indistinguishable from the pre-enrichment shape.
// Drops in `omitempty` here would force every consumer to handle
// `"title":""` / `"references":null` / `"affectedWorkloads":null`.
func TestHotCVEOnFinishedMessage_OmitemptyOnAddedFields(t *testing.T) {
	minimal := HotCVEOnFinishedMessage{
		CustomerGUID: "tenant-1",
		CVEID:        "CVE-2024-9999",
		Severity:     "Critical",
	}
	b, err := json.Marshal(minimal)
	require.NoError(t, err)
	body := string(b)
	assert.NotContains(t, body, `"title"`)
	assert.NotContains(t, body, `"references"`)
	assert.NotContains(t, body, `"affectedWorkloads"`)
}
