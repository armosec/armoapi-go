package armotypes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	hashiVer "github.com/hashicorp/go-version"
)

// bareSemverPattern is deliberately stricter than hashiVer.NewVersion's own grammar: exactly
// three dot-separated numeric segments, optional "v" prefix, no prerelease/build metadata
// suffix. hashiVer accepts "1.2" (normalizes to 1.2.0), "1.2.3.4", and orders "1.2.3-rc1" below
// "1.2.3" - all of which would let ValidateVariants accept a MinAgentVersion that rulelibrary's
// Python-side twin (_parse_bare_semver in yaml_to_mongodb.py) rejects, or orders differently,
// producing a document one writer considers validly mirrored and the other doesn't. Rather than
// keep the two grammars in sync by convention, both are pinned to this same strict pattern.
var bareSemverPattern = regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

// parsedVariant pairs a RuleVariant with its once-parsed MinAgentVersion, so a single sort can
// serve both LowestVariant and ResolveVariantExpressions without re-parsing on every comparison.
type parsedVariant struct {
	variant RuleVariant
	version *hashiVer.Version // nil if MinAgentVersion didn't parse
}

// sortedParsedVariants parses every variant's MinAgentVersion once and sorts ascending;
// variants with an unparsable MinAgentVersion sort last (relative order among themselves
// preserved) since they can't be placed in the total order this design requires -
// ValidateVariants should be used at write time so this case shouldn't occur on persisted data.
func sortedParsedVariants(variants []RuleVariant) []parsedVariant {
	parsed := make([]parsedVariant, len(variants))
	for i, v := range variants {
		version, err := hashiVer.NewVersion(v.MinAgentVersion)
		if err != nil {
			parsed[i] = parsedVariant{variant: v}
			continue
		}
		parsed[i] = parsedVariant{variant: v, version: version}
	}
	sort.SliceStable(parsed, func(i, j int) bool {
		if parsed[i].version == nil || parsed[j].version == nil {
			return parsed[i].version != nil && parsed[j].version == nil
		}
		return parsed[i].version.LessThan(parsed[j].version)
	})
	return parsed
}

// SortedVariants returns a copy of variants sorted by MinAgentVersion ascending. See
// sortedParsedVariants for the unparsable-entry ordering rule.
func SortedVariants(variants []RuleVariant) []RuleVariant {
	parsed := sortedParsedVariants(variants)
	sorted := make([]RuleVariant, len(parsed))
	for i, p := range parsed {
		sorted[i] = p.variant
	}
	return sorted
}

// LowestVariant returns the variant with the lowest MinAgentVersion - the variant every writer
// must derive the RuntimeRule's top-level Expressions from, so version-unaware readers degrade
// to a working base body. Returns false if there are no variants.
func LowestVariant(variants []RuleVariant) (RuleVariant, bool) {
	if len(variants) == 0 {
		return RuleVariant{}, false
	}
	return sortedParsedVariants(variants)[0].variant, true
}

// ResolveVariantExpressions picks the RuleExpressions of the highest-MinAgentVersion variant
// satisfied by agentVersion (i.e. the highest variant with MinAgentVersion <= agentVersion).
//
// It falls back to the lowest-MinAgentVersion variant - never the newest - whenever agentVersion
// is empty, unparsable, or below every variant's floor. That prevents the failure this design
// exists to avoid (an old agent silently jumping to a body using a capability it doesn't have),
// but it is NOT by itself a guarantee that the returned body compiles on the requesting agent:
// if every variant's floor is above the agent's real version (e.g. variants authored at 5.0.0
// and 6.0.0 only, requested by a 1.0.0 agent), the lowest variant is still returned, and nothing
// in this package enforces that a rule's lowest variant covers the oldest agent version the
// fleet actually runs. That floor-coverage discipline is a write-time authoring convention this
// function cannot see or enforce - callers that need the stronger guarantee must enforce it
// themselves (e.g. requiring a variant at or below a known fleet floor before accepting a rule).
// Returns false if variants is empty (caller should keep using the RuntimeRule's plain top-level
// Expressions in that case).
func ResolveVariantExpressions(variants []RuleVariant, agentVersion string) (RuleExpressions, bool) {
	sorted := sortedParsedVariants(variants)
	if len(sorted) == 0 {
		return RuleExpressions{}, false
	}
	best := sorted[0].variant // lowest - the fallback

	requested, err := hashiVer.NewVersion(agentVersion)
	if err != nil {
		return best.Expressions, true
	}

	for _, p := range sorted {
		if p.version == nil {
			continue
		}
		if requested.LessThan(p.version) {
			break
		}
		best = p.variant
	}
	return best.Expressions, true
}

// ValidateVariants checks that every variant's MinAgentVersion matches bareSemverPattern - a
// bare, three-segment semver with no prerelease/build metadata and no constraint syntax (unlike
// AgentVersionRequirement) - and that no two variants share a MinAgentVersion. Both are required
// for the total order LowestVariant and ResolveVariantExpressions rely on:
//
//   - A duplicate floor makes the two functions disagree about which variant is "lowest" versus
//     "the one an agent at that exact floor resolves to" - SortedVariants is stable, so
//     LowestVariant keeps the first of the tied variants in input order, while
//     ResolveVariantExpressions' resolution loop keeps overwriting best on ties and so returns
//     the last one. A version-unaware reader and an agent at that exact floor would then see
//     different bodies for what the schema calls "the same" variant.
//   - Prerelease/build metadata (e.g. "1.2.3-rc1", "1.2.3+build5") is exactly the input where
//     this package's ordering (hashiVer, prerelease sorts below release) and rulelibrary's
//     Python twin (strips prerelease/build before comparing, so "1.2.3-rc1" == "1.2.3") can
//     silently disagree on ordering. Rejecting it here removes the ambiguity rather than relying
//     on two independent implementations to agree on how to order it.
func ValidateVariants(variants []RuleVariant) error {
	seen := make(map[string]int, len(variants)) // normalized version -> first index using it
	for i, v := range variants {
		if !bareSemverPattern.MatchString(v.MinAgentVersion) {
			return fmt.Errorf("variant %d: minAgentVersion %q is not a valid bare semver (e.g. \"1.2.3\") - constraint strings and prerelease/build metadata are not allowed", i, v.MinAgentVersion)
		}
		normalized := strings.TrimPrefix(v.MinAgentVersion, "v")
		if first, duplicate := seen[normalized]; duplicate {
			return fmt.Errorf("variant %d: minAgentVersion %q duplicates variant %d's", i, v.MinAgentVersion, first)
		}
		seen[normalized] = i
	}
	return nil
}
