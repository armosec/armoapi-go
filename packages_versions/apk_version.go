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
