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
	assert.Equal(t, "critical", hotCVEs[0].Severity)
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
				Severity:         "critical",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "valid with status active",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "high",
				Status:           HotCVEStatusActive,
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "valid with status inactive",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "medium",
				Status:           HotCVEStatusInactive,
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
		},
		{
			name: "missing cveId",
			cve: HotCVE{
				Severity:         "critical",
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
				Severity:         "critical",
				AffectedPackages: []HotCVEAffectedPackage{},
			},
			wantErr: "affectedPackages is required",
		},
		{
			name: "nil affected packages",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "critical",
			},
			wantErr: "affectedPackages is required",
		},
		{
			name: "invalid status",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "critical",
				Status:           "pending",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
			wantErr: "invalid status",
		},
		{
			name: "invalid affected package — missing packageName",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "critical",
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
				Severity: "critical",
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
				Severity: "critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{PackageName: "xz-utils", VersionExact: []string{"5.6.1"}},
				},
			},
		},
		{
			name: "valid with start only (no upper bound)",
			cve: HotCVE{
				CVEID:    "CVE-2024-3094",
				Severity: "critical",
				AffectedPackages: []HotCVEAffectedPackage{
					{PackageName: "xz-utils", VersionStart: "5.6.0"},
				},
			},
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
