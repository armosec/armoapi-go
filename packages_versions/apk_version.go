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
	apk "github.com/knqyf263/go-apk-version"
)

type apkVersion struct {
	raw string
	obj apk.Version
}

func newApkVersion(raw string) (*apkVersion, error) {
	ver, err := apk.NewVersion(raw)
	if err != nil {
		return nil, err
	}

	return &apkVersion{
		raw: raw,
		obj: ver,
	}, nil
}

func (a *apkVersion) Compare(other *Version) (int, error) {
	if other.Format != version.ApkFormat {
		return -1, fmt.Errorf("unable to compare apk to given format: %s", other.Format)
	}
	if other.rich.apkVer == nil {
		return -1, fmt.Errorf("given empty apkVersion object")
	}

	return other.rich.apkVer.obj.Compare(a.obj), nil
}
