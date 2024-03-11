/*
Copyright [2024] [anchore/grype]
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package packages_versions

import (
	"fmt"
	"strings"

	"github.com/anchore/grype/grype/version"

	hashiVer "github.com/hashicorp/go-version"
)

type golangVersion struct {
	raw    string
	semVer *hashiVer.Version
}

func (g golangVersion) Compare(other *Version) (int, error) {
	if other.Format != version.GolangFormat {
		return -1, fmt.Errorf("cannot compare %v to golang version", other.Format)
	}
	if other.rich.golangVersion == nil {
		return -1, fmt.Errorf("cannot compare version with nil golang version to golang version")
	}
	if other.rich.golangVersion.raw == g.raw {
		return 0, nil
	}
	if other.rich.golangVersion.raw == "(devel)" {
		return -1, fmt.Errorf("cannot compare %s with %s", g.raw, other.rich.golangVersion.raw)
	}

	return other.rich.golangVersion.compare(g), nil
}

func (g golangVersion) compare(o golangVersion) int {
	switch {
	case g.semVer != nil && o.semVer != nil:
		return g.semVer.Compare(o.semVer)
	case g.semVer != nil && o.semVer == nil:
		return 1
	case g.semVer == nil && o.semVer != nil:
		return -1
	default:
		return strings.Compare(g.raw, o.raw)
	}
}

func newGolangVersion(v string) (*golangVersion, error) {
	// go stdlib is reported by syft as a go package with version like "go1.24.1"
	// other versions have "v" as a prefix, which the semver lib handles automatically
	semver, err := hashiVer.NewSemver(strings.TrimPrefix(v, "go"))
	if err != nil {
		return nil, err
	}
	return &golangVersion{
		raw:    v,
		semVer: semver,
	}, nil
}
