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

	hashiVer "github.com/hashicorp/go-version"
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

// Compare checks semVer population rather than Format, because multiple formats
// populate rich.semVer (SemanticFormat, GemFormat, and the UnknownFormat fallback).
// The nil check is the real safety guard — Format alone would be too restrictive.
func (v *semanticVersion) Compare(other *Version) (int, error) {
	if other.rich.semVer == nil {
		return -1, fmt.Errorf("given Version has no semantic version populated")
	}

	return other.rich.semVer.verObj.Compare(v.verObj), nil
}
