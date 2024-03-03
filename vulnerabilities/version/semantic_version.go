package version

import (
	"fmt"
	"github.com/anchore/grype/grype/version"
	"strings"

	hashiVer "github.com/anchore/go-version"
)

var normalizer = strings.NewReplacer(".alpha", "-alpha", ".beta", "-beta", ".rc", "-rc")

type semanticVersion struct {
	verObj *hashiVer.Version
}

func newSemanticVersion(raw string) (*semanticVersion, error) {
	verObj, err := hashiVer.NewVersion(normalizer.Replace(raw))
	if err != nil {
		return nil, fmt.Errorf("unable to create semver obj: %w", err)
	}
	return &semanticVersion{
		verObj: verObj,
	}, nil
}

func (v *semanticVersion) Compare(other *Version) (int, error) {
	if other.Format != version.SemanticFormat {
		return -1, fmt.Errorf("unable to compare semantic version to given format: %s", other.Format)
	}
	if other.rich.semVer == nil {
		return -1, fmt.Errorf("given empty semanticVersion object")
	}

	return other.rich.semVer.verObj.Compare(v.verObj), nil
}
