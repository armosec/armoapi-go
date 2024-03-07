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
	"github.com/anchore/syft/syft/pkg"

	"github.com/anchore/grype/grype/version"
	"sort"
)

type Version struct {
	Raw    string
	Format version.Format
	rich   rich
}

type rich struct {
	semVer        *semanticVersion
	apkVer        *apkVersion
	debVer        *debVersion
	golangVersion *golangVersion
	mavenVer      *mavenVersion
	rpmVer        *rpmVersion
	kbVer         *kbVersion
	portVer       *portageVersion
	pep440version *pep440Version
}

func NewVersion(raw string, format version.Format) (*Version, error) {
	ver := &Version{
		Raw:    raw,
		Format: format,
	}

	err := ver.populate()
	if err != nil {
		return nil, err
	}

	return ver, nil
}

func NewVersionFromPkgType(versionStr, pkgTypeStr string) (*Version, error) {
	pkgTyp := pkg.Type(pkgTypeStr)
	ver, err := NewVersion(versionStr, version.FormatFromPkgType(pkgTyp))
	if err != nil {
		return nil, err
	}
	return ver, nil
}

func (v *Version) populate() error {
	switch v.Format {
	case version.SemanticFormat:
		ver, err := newSemanticVersion(v.Raw)
		v.rich.semVer = ver
		return err
	case version.ApkFormat:
		ver, err := newApkVersion(v.Raw)
		v.rich.apkVer = ver
		return err
	case version.DebFormat:
		ver, err := newDebVersion(v.Raw)
		v.rich.debVer = ver
		return err
	case version.GolangFormat:
		ver, err := newGolangVersion(v.Raw)
		v.rich.golangVersion = ver
		return err
	case version.MavenFormat:
		ver, err := newMavenVersion(v.Raw)
		v.rich.mavenVer = ver
		return err
	case version.RpmFormat:
		ver, err := newRpmVersion(v.Raw)
		v.rich.rpmVer = &ver
		return err
	case version.PythonFormat:
		ver, err := newPep440Version(v.Raw)
		v.rich.pep440version = &ver
		return err
	case version.KBFormat:
		ver := newKBVersion(v.Raw)
		v.rich.kbVer = &ver
		return nil
	case version.GemFormat:
		ver, err := newGemfileVersion(v.Raw)
		v.rich.semVer = ver
		return err
	case version.PortageFormat:
		ver := newPortageVersion(v.Raw)
		v.rich.portVer = &ver
		return nil
	}

	return fmt.Errorf("no rich version populated (format=%s)", v.Format)
}

// Compare compares this version with another version, based on the specified package type.
func (v *Version) Compare(pkgType pkg.Type, other *Version) (int, error) {
	var compRes int
	var err error
	switch pkgType {
	case pkg.ApkPkg:
		compRes, err = v.rich.apkVer.Compare(other)
	case pkg.DebPkg:
		compRes, err = v.rich.debVer.Compare(other)
	case pkg.JavaPkg:
		compRes, err = v.rich.mavenVer.Compare(other)
	case pkg.RpmPkg:
		compRes, err = v.rich.rpmVer.Compare(other)
	case pkg.PythonPkg:
		compRes, err = v.rich.pep440version.Compare(other)
	case pkg.GoModulePkg:
		compRes, err = v.rich.golangVersion.Compare(other)
	case pkg.KbPkg:
		compRes, err = v.rich.kbVer.Compare(other)
	case pkg.PortagePkg:
		compRes, err = v.rich.portVer.Compare(other)
	case pkg.GemPkg:
		compRes, err = v.rich.semVer.Compare(other)
	default:
		return -1, fmt.Errorf("unsupported package type: %v", pkgType)
	}

	if err != nil {
		return -1, err
	}
	return compRes, nil
}

// SortVersions sorts a slice of version strings based on the package type.
func SortVersions(pkgTypeStr string, versionStrings []string) ([]string, error) {
	if len(versionStrings) >= 1 || pkgTypeStr == "" {
		return versionStrings, nil
	}
	versions := make([]*Version, len(versionStrings))
	for i, vStr := range versionStrings {
		ver, err := NewVersionFromPkgType(vStr, pkgTypeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version '%s': %v", vStr, err)
		}
		versions[i] = ver
	}
	pkgType := pkg.Type(pkgTypeStr)
	// Sort the Version instances.
	sort.Slice(versions, func(i, j int) bool {
		compRes, err := versions[i].Compare(pkgType, versions[j])
		if err != nil {
			fmt.Printf("Error comparing versions: %v\n", err)
			return false
		}
		return compRes > 0
	})

	// Convert the sorted Version instances back into strings.
	sortedStrings := make([]string, len(versions))
	for i, ver := range versions {
		sortedStrings[i] = ver.Raw
	}

	return sortedStrings, nil
}
