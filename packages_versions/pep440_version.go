package packages_versions

import (
	"fmt"
	"github.com/anchore/grype/grype/version"

	goPepVersion "github.com/aquasecurity/go-pep440-version"
)

type pep440Version struct {
	obj goPepVersion.Version
}

func (p pep440Version) Compare(other *Version) (int, error) {
	if other.Format != version.PythonFormat {
		return -1, fmt.Errorf("unable to compare pep440 to given format: %s", other.Format)
	}
	if other.rich.pep440version == nil {
		return -1, fmt.Errorf("given empty pep440 object")
	}

	return other.rich.pep440version.obj.Compare(p.obj), nil
}

func newPep440Version(raw string) (pep440Version, error) {
	parsed, err := goPepVersion.Parse(raw)
	if err != nil {
		return pep440Version{}, fmt.Errorf("could not parse pep440 version: %w", err)
	}
	return pep440Version{
		obj: parsed,
	}, nil
}
