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
	"reflect"
)

type kbVersion struct {
	version string
}

func newKBVersion(raw string) kbVersion {
	// XXX Is this even useful/necessary?
	return kbVersion{
		version: raw,
	}
}

func (v *kbVersion) Compare(other *Version) (int, error) {
	if other.Format != version.KBFormat {
		return -1, fmt.Errorf("unable to compare kb to given format: %s", other.Format)
	}

	if other.rich.kbVer == nil {
		return -1, fmt.Errorf("given empty kbVersion object")
	}

	return other.rich.kbVer.compare(*v), nil
}

// Compare returns 0 if v == v2, 1 otherwise
func (v kbVersion) compare(v2 kbVersion) int {
	if reflect.DeepEqual(v, v2) {
		return 0
	}

	return 1
}

func (v kbVersion) String() string {
	return v.version
}
