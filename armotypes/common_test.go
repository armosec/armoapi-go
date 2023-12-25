package armotypes

import (
	"reflect"
	"sort"
	"testing"
)

// TestGetRiskFactors tests the GetRiskFactors function.
func TestGetRiskFactors(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []RiskFactor
	}{
		{
			name:     "Multiple Risk Factors",
			input:    []string{"C-0256", "C-0046", "C-0057", "C-0255"},
			expected: []RiskFactor{RiskFactorInternetFacing, RiskFactorPrivileged, RiskFactorSecretAccess},
		},
		{
			name:     "Empty controls list",
			input:    []string{},
			expected: []RiskFactor{},
		},
		{
			name:     "nil controls list",
			input:    nil,
			expected: []RiskFactor{},
		},
		{
			name:     "Single Risk Factor",
			input:    []string{"C-0256"},
			expected: []RiskFactor{RiskFactorInternetFacing},
		},
		{
			name:     "No Risk Factors",
			input:    []string{"C-0000"},
			expected: []RiskFactor{},
		},
		{
			name:     "Duplicate Risk Factors",
			input:    []string{"C-0046", "C-0046"},
			expected: []RiskFactor{RiskFactorPrivileged},
		},
		{
			name:     "Mixed Valid and Invalid IDs",
			input:    []string{"C-0046", "C-9999"},
			expected: []RiskFactor{RiskFactorPrivileged},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRiskFactors(tt.input)
			sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
			sort.Slice(tt.expected, func(i, j int) bool { return tt.expected[i] < tt.expected[j] })

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetRiskFactors(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
