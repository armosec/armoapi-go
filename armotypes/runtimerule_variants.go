package armotypes

import (
	"fmt"
	"sort"

	hashiVer "github.com/hashicorp/go-version"
)

// SortedVariants returns a copy of variants sorted by MinAgentVersion ascending. Variants whose
// MinAgentVersion doesn't parse as a semver sort after every parsable one (relative order among
// themselves preserved), since they can't be placed in the total order this design requires;
// ValidateVariants should be used at write time so this case shouldn't occur on persisted data.
func SortedVariants(variants []RuleVariant) []RuleVariant {
	sorted := make([]RuleVariant, len(variants))
	copy(sorted, variants)
	sort.SliceStable(sorted, func(i, j int) bool {
		vi, erri := hashiVer.NewVersion(sorted[i].MinAgentVersion)
		vj, errj := hashiVer.NewVersion(sorted[j].MinAgentVersion)
		if erri != nil || errj != nil {
			return erri == nil && errj != nil
		}
		return vi.LessThan(vj)
	})
	return sorted
}

// LowestVariant returns the variant with the lowest MinAgentVersion - the variant every writer
// must derive the RuntimeRule's top-level Expressions from, so version-unaware readers degrade
// to a working base body. Returns false if there are no variants.
func LowestVariant(variants []RuleVariant) (RuleVariant, bool) {
	if len(variants) == 0 {
		return RuleVariant{}, false
	}
	return SortedVariants(variants)[0], true
}

// ResolveVariantExpressions picks the RuleExpressions of the highest-MinAgentVersion variant
// satisfied by agentVersion (i.e. the highest variant with MinAgentVersion <= agentVersion).
// Per the fail-safe/conservative-default requirement, it falls back to the lowest-MinAgentVersion
// variant - never the newest - whenever agentVersion is empty, unparsable, or below every
// variant's floor: an agent whose capability can't be confirmed must never be handed a body it
// might not be able to compile. Returns false if variants is empty (caller should keep using the
// RuntimeRule's plain top-level Expressions in that case).
func ResolveVariantExpressions(variants []RuleVariant, agentVersion string) (RuleExpressions, bool) {
	lowest, ok := LowestVariant(variants)
	if !ok {
		return RuleExpressions{}, false
	}

	requested, err := hashiVer.NewVersion(agentVersion)
	if err != nil {
		return lowest.Expressions, true
	}

	best := lowest
	for _, v := range SortedVariants(variants) {
		floor, err := hashiVer.NewVersion(v.MinAgentVersion)
		if err != nil {
			continue
		}
		if requested.LessThan(floor) {
			break
		}
		best = v
	}
	return best.Expressions, true
}

// ValidateVariants checks that every variant's MinAgentVersion is a well-formed bare semver, not
// a constraint string (unlike AgentVersionRequirement) - required for the total order that
// ResolveVariantExpressions and LowestVariant rely on.
func ValidateVariants(variants []RuleVariant) error {
	for i, v := range variants {
		if _, err := hashiVer.NewVersion(v.MinAgentVersion); err != nil {
			return fmt.Errorf("variant %d: minAgentVersion %q is not a valid semver: %w", i, v.MinAgentVersion, err)
		}
	}
	return nil
}
