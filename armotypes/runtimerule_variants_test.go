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

func TestValidateVariants_RejectsNonThreeSegment(t *testing.T) {
	// hashiVer.NewVersion alone would accept both of these (normalizing "1.2" to "1.2.0", and
	// happily parsing four segments) - the stricter bareSemverPattern rejects them so this
	// grammar matches rulelibrary's Python-side twin (_parse_bare_semver), which always
	// required exactly three segments.
	err := ValidateVariants([]RuleVariant{{MinAgentVersion: "1.2"}})
	require.Error(t, err, "two-segment version must be rejected, not silently normalized")

	err = ValidateVariants([]RuleVariant{{MinAgentVersion: "1.2.3.4"}})
	require.Error(t, err, "four-segment version must be rejected")
}

func TestValidateVariants_RejectsPrereleaseAndBuildMetadata(t *testing.T) {
	// Prerelease/build metadata is where this package's ordering (hashiVer: prerelease sorts
	// below release) and rulelibrary's Python twin (strips the suffix before comparing, so
	// "1.2.3-rc1" == "1.2.3") could silently disagree - rejecting it outright removes the
	// ambiguity instead of relying on two implementations to order it the same way.
	err := ValidateVariants([]RuleVariant{{MinAgentVersion: "1.2.3-rc1"}})
	require.Error(t, err, "prerelease suffix must be rejected")

	err = ValidateVariants([]RuleVariant{{MinAgentVersion: "1.2.3+build5"}})
	require.Error(t, err, "build metadata suffix must be rejected")
}

func TestValidateVariants_RejectsDuplicateMinAgentVersion(t *testing.T) {
	err := ValidateVariants([]RuleVariant{
		{MinAgentVersion: "1.0.0", Expressions: expr("a")},
		{MinAgentVersion: "1.0.0", Expressions: expr("b")},
	})
	require.Error(t, err, "two variants sharing a floor must be rejected")
}

func TestValidateVariants_RejectsDuplicateMinAgentVersion_VPrefixNormalized(t *testing.T) {
	// "1.0.0" and "v1.0.0" parse to the same hashiVer.Version, so LowestVariant and
	// ResolveVariantExpressions would treat them as a tie too - the duplicate check must catch
	// this normalized form, not just a literal string match.
	err := ValidateVariants([]RuleVariant{
		{MinAgentVersion: "1.0.0", Expressions: expr("a")},
		{MinAgentVersion: "v1.0.0", Expressions: expr("b")},
	})
	require.Error(t, err)
}

func TestLowestVariant_And_ResolveVariantExpressions_AgreeOnDuplicateFloors(t *testing.T) {
	// Before the duplicate check existed, LowestVariant (stable sort, keeps first of a tie) and
	// ResolveVariantExpressions (resolution loop overwrites best on ties, keeps last) disagreed
	// about which variant a tied floor actually meant - a version-unaware reader would see one
	// body and an agent at exactly that floor would receive another. ValidateVariants now
	// rejects this input at write time; this test pins that the two resolution functions would
	// still agree even if a duplicate somehow reached them (defense in depth, not a substitute
	// for the write-time check).
	variants := []RuleVariant{
		{MinAgentVersion: "1.0.0", Expressions: expr("a")},
		{MinAgentVersion: "1.0.0", Expressions: expr("a")}, // identical, not just same floor
	}
	lowest, ok := LowestVariant(variants)
	require.True(t, ok)
	resolved, ok := ResolveVariantExpressions(variants, "1.0.0")
	require.True(t, ok)
	assert.Equal(t, lowest.Expressions, resolved, "identical duplicate floors must resolve identically regardless of tie-break direction")
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
