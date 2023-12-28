package armotypes

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestGetControlIDsByRiskFactors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single Risk Factor",
			input:    "RiskFactorInternetFacing",
			expected: []string{"C-0256"},
		},
		{
			name:     "Multiple Risk Factors",
			input:    "RiskFactorPrivileged,RiskFactorSecretAccess",
			expected: []string{"C-0046", "C-0057", "C-0255"},
		},
		{
			name:     "No Risk Factors",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Invalid Risk Factor",
			input:    "RiskFactorNonExistent",
			expected: []string{},
		},
		{
			name:     "Duplicate Risk Factors",
			input:    "RiskFactorHostAccess,RiskFactorHostAccess",
			expected: []string{"C-0038", "C-0041", "C-0044", "C-0048"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetControlIDsByRiskFactors(tt.input)
			sort.Strings(result)
			sort.Strings(tt.expected)

			assert.Equal(t, tt.expected, result)
		})
	}
}

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
			expected: nil,
		},
		{
			name:     "nil controls list",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Single Risk Factor",
			input:    []string{"C-0256"},
			expected: []RiskFactor{RiskFactorInternetFacing},
		},
		{
			name:     "No Risk Factors",
			input:    []string{"C-0000"},
			expected: nil,
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

			assert.Equal(t, tt.expected, result, "GetRiskFactors(%v) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}
