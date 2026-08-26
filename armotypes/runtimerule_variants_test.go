package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func expr(msg string) RuleExpressions {
	return RuleExpressions{Message: msg, UniqueID: "''", RuleExpression: []RuleExpression{{EventType: "http", Expression: "true"}}}
}

func TestLowestVariant_PicksLowest(t *testing.T) {
	variants := []RuleVariant{
		{MinAgentVersion: "1.2.0", Expressions: expr("v1.2.0")},
		{MinAgentVersion: "0.9.0", Expressions: expr("v0.9.0")},
		{MinAgentVersion: "2.0.0", Expressions: expr("v2.0.0")},
	}
	lowest, ok := LowestVariant(variants)
	require.True(t, ok)
	assert.Equal(t, "0.9.0", lowest.MinAgentVersion)
}

func TestLowestVariant_Empty(t *testing.T) {
	_, ok := LowestVariant(nil)
	assert.False(t, ok)
}

func TestResolveVariantExpressions_PicksHighestSatisfied(t *testing.T) {
	variants := []RuleVariant{
		{MinAgentVersion: "0.9.0", Expressions: expr("base")},
		{MinAgentVersion: "1.5.0", Expressions: expr("mid")},
		{MinAgentVersion: "2.0.0", Expressions: expr("new")},
	}

	tests := []struct {
		name         string
		agentVersion string
		wantMessage  string
	}{
		{"exact match on a floor", "1.5.0", "mid"},
		{"between floors picks the lower", "1.9.9", "mid"},
		{"above every floor picks newest", "3.0.0", "new"},
		{"below every floor falls back to lowest, not error", "0.1.0", "base"},
		{"empty version falls back to lowest", "", "base"},
		{"unparsable version falls back to lowest", "not-a-version", "base"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ResolveVariantExpressions(variants, tt.agentVersion)
			require.True(t, ok)
			assert.Equal(t, tt.wantMessage, got.Message)
		})
	}
}

func TestResolveVariantExpressions_NeverPicksNewestOnUnknownVersion(t *testing.T) {
	// Fail-safe requirement: unknown version info must never loosen coverage by jumping to a
	// body the requester might not be able to compile.
	variants := []RuleVariant{
		{MinAgentVersion: "0.1.0", Expressions: expr("base")},
		{MinAgentVersion: "9.9.9", Expressions: expr("newest")},
	}
	got, ok := ResolveVariantExpressions(variants, "")
	require.True(t, ok)
	assert.Equal(t, "base", got.Message, "unknown agent version must resolve to the lowest variant, never the newest")
}

func TestResolveVariantExpressions_NoVariants(t *testing.T) {
	_, ok := ResolveVariantExpressions(nil, "1.0.0")
	assert.False(t, ok)
}

func TestSortedVariants_UnparsableSortsLast(t *testing.T) {
	variants := []RuleVariant{
		{MinAgentVersion: "bogus"},
		{MinAgentVersion: "1.0.0"},
	}
	sorted := SortedVariants(variants)
	assert.Equal(t, "1.0.0", sorted[0].MinAgentVersion)
	assert.Equal(t, "bogus", sorted[1].MinAgentVersion)
}

func TestValidateVariants(t *testing.T) {
	assert.NoError(t, ValidateVariants([]RuleVariant{{MinAgentVersion: "1.2.3"}}))
	assert.NoError(t, ValidateVariants([]RuleVariant{{MinAgentVersion: "v1.2.3"}}))

	err := ValidateVariants([]RuleVariant{{MinAgentVersion: ">=1.2.3"}})
	require.Error(t, err, "MinAgentVersion is a bare semver, not a constraint string")

	err = ValidateVariants([]RuleVariant{{MinAgentVersion: "not-a-version"}})
	require.Error(t, err)
}

func TestRuntimeRule_Variants_YAMLRoundTrip(t *testing.T) {
	rule := RuntimeRule{
		ID:          "R9999",
		Name:        "test rule",
		Expressions: expr("base"),
		Variants: []RuleVariant{
			{MinAgentVersion: "0.9.0", Expressions: expr("base")},
			{MinAgentVersion: "2.0.0", Expressions: expr("new")},
		},
	}

	b, err := yaml.Marshal(rule)
	require.NoError(t, err)

	var roundTripped RuntimeRule
	require.NoError(t, yaml.Unmarshal(b, &roundTripped))
	assert.Equal(t, rule.Variants, roundTripped.Variants)
}
