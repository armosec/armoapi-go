package packages_versions

import (
	"slices"

	syftPkg "github.com/anchore/syft/syft/pkg"
)

// VersionConstraint defines a version range and/or exact match list for matching
// against candidate versions. Used by the hot CVE evaluator to determine if an
// SBOM component falls within an affected version range.
type VersionConstraint struct {
	Start      string   // inclusive lower bound (empty = unbounded)
	End        string   // exclusive upper bound (empty = unbounded)
	ExactMatch []string // exact version strings (plain string equality)
}

// Matches checks whether candidateVersion satisfies this constraint for the given pkgType.
// If ExactMatch is non-empty, only exact string matches are checked (range is ignored).
// Otherwise, the range [Start, End) is evaluated using version-aware comparison.
// If all fields are empty, the constraint matches any version.
func (vc *VersionConstraint) Matches(candidateVersion, pkgType string) (bool, error) {
	// ExactMatch takes precedence when set
	if len(vc.ExactMatch) > 0 {
		return slices.Contains(vc.ExactMatch, candidateVersion), nil
	}

	// Fully unbounded: matches everything
	if vc.Start == "" && vc.End == "" {
		return true, nil
	}

	candidate, err := NewVersionFromPkgType(candidateVersion, pkgType)
	if err != nil {
		return false, err
	}

	pkg := syftPkg.Type(pkgType)

	// Check lower bound (inclusive): candidate >= Start
	if vc.Start != "" {
		startVer, err := NewVersionFromPkgType(vc.Start, pkgType)
		if err != nil {
			return false, err
		}
		cmp, err := candidate.Compare(pkg, startVer)
		if err != nil {
			return false, err
		}
		if cmp > 0 { // candidate < start
			return false, nil
		}
	}

	// Check upper bound (exclusive): candidate < End
	if vc.End != "" {
		endVer, err := NewVersionFromPkgType(vc.End, pkgType)
		if err != nil {
			return false, err
		}
		cmp, err := candidate.Compare(pkg, endVer)
		if err != nil {
			return false, err
		}
		if cmp <= 0 { // candidate >= end (end is exclusive)
			return false, nil
		}
	}

	return true, nil
}
