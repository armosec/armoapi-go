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
	"github.com/anchore/grype/grype/version"

	mvnv "github.com/masahiro331/go-mvn-version"
)

type mavenVersion struct {
	raw     string
	version mvnv.Version
}

func newMavenVersion(raw string) (*mavenVersion, error) {
	ver, err := mvnv.NewVersion(raw)
	if err != nil {
		return nil, fmt.Errorf("could not generate new java version from: %s; %w", raw, err)
	}

	return &mavenVersion{
		raw:     raw,
		version: ver,
	}, nil
}

// Compare returns 0 if j2 == j, 1 if j2 > j, and -1 if j2 < j.
// If an error returns the int value is -1
func (j *mavenVersion) Compare(j2 *Version) (int, error) {
	if j2.Format != version.MavenFormat {
		return -1, fmt.Errorf("unable to compare java to given format: %s", j2.Format)
	}
	if j2.rich.mavenVer == nil {
		return -1, fmt.Errorf("given empty mavenVersion object")
	}

	submittedVersion := j2.rich.mavenVer.version
	if submittedVersion.Equal(j.version) {
		return 0, nil
	}
	if submittedVersion.LessThan(j.version) {
		return -1, nil
	}
	if submittedVersion.GreaterThan(j.version) {
		return 1, nil
	}

	return -1, fmt.Errorf(
		"could not compare java versions: %v with %v",
		submittedVersion.String(),
		j.version.String())
}
