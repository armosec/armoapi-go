package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name: "invalid severity",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "urgent",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
			wantErr: "invalid severity",
		},
		{
			name: "case insensitive severity",
			cve: HotCVE{
				CVEID:            "CVE-2024-3094",
				Severity:         "Critical",
				AffectedPackages: []HotCVEAffectedPackage{validPkg},
			},
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
