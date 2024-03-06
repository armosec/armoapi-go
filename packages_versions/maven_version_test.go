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
	"github.com/anchore/grype/grype/version"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_javaVersion_Compare(t *testing.T) {
	tests := []struct {
		name    string
		compare string
		want    int
	}{
		{
			name:    "1",
			compare: "2",
			want:    -1,
		},
		{
			name:    "1.8.0_282",
			compare: "1.8.0_282",
			want:    0,
		},
		{
			name:    "2.5",
			compare: "2.0",
			want:    1,
		},
		{
			name:    "2.414.2-cb-5",
			compare: "2.414.2",
			want:    1,
		},
		{
			name:    "5.2.25.RELEASE", // see https://mvnrepository.com/artifact/org.springframework/spring-web
			compare: "5.2.25",
			want:    0,
		},
		{
			name:    "5.2.25.release",
			compare: "5.2.25",
			want:    0,
		},
		{
			name:    "5.2.25.FINAL",
			compare: "5.2.25",
			want:    0,
		},
		{
			name:    "5.2.25.final",
			compare: "5.2.25",
			want:    0,
		},
		{
			name:    "5.2.25.GA",
			compare: "5.2.25",
			want:    0,
		},
		{
			name:    "5.2.25.ga",
			compare: "5.2.25",
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := NewVersion(tt.name, version.MavenFormat)
			assert.NoError(t, err)

			j2, err := NewVersion(tt.compare, version.MavenFormat)
			assert.NoError(t, err)

			if got, _ := j2.rich.mavenVer.Compare(j); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
