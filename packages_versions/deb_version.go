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

	deb "github.com/knqyf263/go-deb-version"
)

type debVersion struct {
	obj deb.Version
}

func newDebVersion(raw string) (*debVersion, error) {
	ver, err := deb.NewVersion(raw)
	if err != nil {
		return nil, err
	}
	return &debVersion{
		obj: ver,
	}, nil
}

func (d *debVersion) Compare(other *Version) (int, error) {
	if other.Format != version.DebFormat {
		return -1, fmt.Errorf("unable to compare deb to given format: %s", other.Format)
	}
	if other.rich.debVer == nil {
		return -1, fmt.Errorf("given empty debVersion object")
	}

	return other.rich.debVer.obj.Compare(d.obj), nil
}
