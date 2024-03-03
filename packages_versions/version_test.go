package packages_versions

import (
	"github.com/anchore/syft/syft/pkg"
	"testing"

	"github.com/anchore/grype/grype/version"
	"github.com/stretchr/testify/assert"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		format    version.Format
		wantError bool
	}{
		{"Valid Semantic Version", "1.0.0", version.SemanticFormat, false},
		{"Invalid Format", "1.0.0", version.ParseFormat("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVersion(tt.raw, tt.format)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewVersionFromPkgType(t *testing.T) {
	tests := []struct {
		name       string
		versionStr string
		pkgTypeStr string
		wantError  bool
	}{
		{"Valid Java Package", "1.0.0", string(pkg.JavaPkg), false},
		{"Invalid Package Type", "1.0.0", "invalid", true},
		{"Valid Python Package", "1.0.0", string(pkg.PythonPkg), false},
		{"Valid Golang Package", "v1.2.3", string(pkg.GoModulePkg), false},
		{"Invalid Version for Golang", "1.2ee", string(pkg.GoModulePkg), true},
		{"Valid Maven Package", "1.0.0-RELEASE", string(pkg.JavaPkg), false},
		{"Valid RPM Package", "1.0.0-1", string(pkg.RpmPkg), false},
		{"Valid Debian Package", "1.0.0-1ubuntu1", string(pkg.DebPkg), false},
		{"Valid APK Package", "1.0.0-r0", string(pkg.ApkPkg), false},
		{"Valid Semantic Version", "2.0.0", "semver", false},
		{"Invalid Semantic Version", "2.0", "semver", true},
		{"Valid Gem Package", "1.0.0", string(pkg.GemPkg), false},
		{"Valid Portage Package", "1.0.0-r1", string(pkg.PortagePkg), false},
		{"Valid KB Format", "KB123456", string(pkg.KbPkg), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVersionFromPkgType(tt.versionStr, tt.pkgTypeStr)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortVersions(t *testing.T) {
	tests := []struct {
		name        string
		pkgType     string
		versionStrs []string
		want        []string
		wantError   bool
	}{
		{
			name:        "Empty array",
			pkgType:     "java",
			versionStrs: []string{},
			want:        []string{},
			wantError:   false,
		},
		{
			name:        "Sort Semantic Versions",
			pkgType:     "java-archive",
			versionStrs: []string{"1.0.1", "1.0.0", "1.10.0"},
			want:        []string{"1.0.0", "1.0.1", "1.10.0"},
			wantError:   false,
		},
		{
			name:        "Sort Debian Versions",
			pkgType:     "deb",
			versionStrs: []string{"1:1.0.0-2", "2:1.0.0-1", "1:1.0.0-1"},
			want:        []string{"1:1.0.0-1", "1:1.0.0-2", "2:1.0.0-1"},
			wantError:   false,
		},
		{
			name:        "Sort APK Versions",
			pkgType:     "apk",
			versionStrs: []string{"1.0.0-r2", "1.0.0-r1", "1.0.1-r0"},
			want:        []string{"1.0.0-r1", "1.0.0-r2", "1.0.1-r0"},
			wantError:   false,
		},
		{
			name:        "Sort Python Versions",
			pkgType:     "python",
			versionStrs: []string{"1.0.1", "1.0.0", "1.0.0a1"},
			want:        []string{"1.0.0a1", "1.0.0", "1.0.1"},
			wantError:   false,
		},
		{
			name:        "Invalid Package Type",
			pkgType:     "unknown",
			versionStrs: []string{"1.0.0", "1.0.1"},
			wantError:   true,
		},
		{
			name:        "Empty Package Type",
			pkgType:     "",
			versionStrs: []string{"1.0.1", "1.0.0"},
			want:        []string{"1.0.1", "1.0.0"},
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortVersions(tt.pkgType, tt.versionStrs)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
