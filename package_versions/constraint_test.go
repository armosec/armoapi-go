package packages_versions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionConstraint_Matches(t *testing.T) {
	tests := []struct {
		name      string
		constraint VersionConstraint
		candidate string
		pkgType   string
		want      bool
		wantErr   bool
	}{
		// Range match within bounds
		{
			name:      "deb: within range",
			constraint: VersionConstraint{Start: "1.0.0-1", End: "2.0.0-1"},
			candidate: "1.5.0-1",
			pkgType:   "deb",
			want:      true,
		},
		{
			name:      "rpm: within range",
			constraint: VersionConstraint{Start: "1.0.0-1", End: "2.0.0-1"},
			candidate: "1.5.0-1",
			pkgType:   "rpm",
			want:      true,
		},
		{
			name:      "npm: within range",
			constraint: VersionConstraint{Start: "1.0.0", End: "2.0.0"},
			candidate: "1.5.0",
			pkgType:   "npm",
			want:      true,
		},
		// Start inclusive
		{
			name:      "start is inclusive",
			constraint: VersionConstraint{Start: "1.0.0", End: "2.0.0"},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      true,
		},
		// End exclusive
		{
			name:      "end is exclusive",
			constraint: VersionConstraint{Start: "1.0.0", End: "2.0.0"},
			candidate: "2.0.0",
			pkgType:   "npm",
			want:      false,
		},
		// Below range
		{
			name:      "below range",
			constraint: VersionConstraint{Start: "1.0.0", End: "2.0.0"},
			candidate: "0.9.0",
			pkgType:   "npm",
			want:      false,
		},
		// Above range
		{
			name:      "above range",
			constraint: VersionConstraint{Start: "1.0.0", End: "2.0.0"},
			candidate: "3.0.0",
			pkgType:   "npm",
			want:      false,
		},
		// No lower bound (only upper bound)
		{
			name:      "no lower bound, within upper",
			constraint: VersionConstraint{End: "2.0.0"},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      true,
		},
		{
			name:      "no lower bound, at upper",
			constraint: VersionConstraint{End: "2.0.0"},
			candidate: "2.0.0",
			pkgType:   "npm",
			want:      false,
		},
		// No upper bound (only lower bound)
		{
			name:      "no upper bound, at lower",
			constraint: VersionConstraint{Start: "1.0.0"},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      true,
		},
		{
			name:      "no upper bound, above lower",
			constraint: VersionConstraint{Start: "1.0.0"},
			candidate: "99.0.0",
			pkgType:   "npm",
			want:      true,
		},
		{
			name:      "no upper bound, below lower",
			constraint: VersionConstraint{Start: "1.0.0"},
			candidate: "0.5.0",
			pkgType:   "npm",
			want:      false,
		},
		// Fully unbounded (no start, no end, no exact)
		{
			name:      "fully unbounded matches anything",
			constraint: VersionConstraint{},
			candidate: "5.0.0",
			pkgType:   "npm",
			want:      true,
		},
		// Exact match found
		{
			name:      "exact match found",
			constraint: VersionConstraint{ExactMatch: []string{"1.2.3", "4.5.6"}},
			candidate: "4.5.6",
			pkgType:   "npm",
			want:      true,
		},
		// Exact match not found
		{
			name:      "exact match not found",
			constraint: VersionConstraint{ExactMatch: []string{"1.2.3", "4.5.6"}},
			candidate: "7.8.9",
			pkgType:   "npm",
			want:      false,
		},
		// Exact match uses string equality (not version-aware comparison)
		{
			name:      "exact match is string equality, not version-aware",
			constraint: VersionConstraint{ExactMatch: []string{"1.0.0"}},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      true,
		},
		{
			name:      "exact match string mismatch despite version equivalence",
			constraint: VersionConstraint{ExactMatch: []string{"v1.0.0"}},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      false,
		},
		// Exact match takes precedence when both range and exact are set
		{
			name: "exact match takes precedence over range (match in exact, outside range)",
			constraint: VersionConstraint{
				Start:      "2.0.0",
				End:        "3.0.0",
				ExactMatch: []string{"1.0.0"},
			},
			candidate: "1.0.0",
			pkgType:   "npm",
			want:      true,
		},
		{
			name: "exact match takes precedence (not in exact, inside range)",
			constraint: VersionConstraint{
				Start:      "1.0.0",
				End:        "3.0.0",
				ExactMatch: []string{"5.0.0"},
			},
			candidate: "2.0.0",
			pkgType:   "npm",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.constraint.Matches(tt.candidate, tt.pkgType)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
