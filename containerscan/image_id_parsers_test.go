package containerscan

import (
	"testing"
)

func TestGetRegistryFromImageID(t *testing.T) {
	tests := []struct {
		name     string
		imageID  string
		expected string
	}{
		{
			name:     "Valid image ID with hash",
			imageID:  "docker.io/library/alpine@sha256:2345",
			expected: "docker.io",
		},
		{
			name:     "Valid image ID with tag",
			imageID:  "quay.io/kubescape/gateway:v0.1.13",
			expected: "quay.io",
		},
		{
			name:     "No registry in image ID",
			imageID:  "library/alpine@sha256:2345",
			expected: "",
		},
		{
			name:     "Empty image ID",
			imageID:  "",
			expected: "",
		},
		{
			name:     "Image ID without hash",
			imageID:  "my.registry.io/library/alpine",
			expected: "my.registry.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRegistryFromImageID(tt.imageID)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetRepositoryFromImageID(t *testing.T) {
	tests := []struct {
		name     string
		imageID  string
		expected string
	}{
		{
			name:     "Valid image ID with hash",
			imageID:  "docker.io/library/alpine@sha256:2345",
			expected: "library/alpine",
		},
		{
			name:     "Valid image ID with tag",
			imageID:  "quay.io/kubescape/gateway:v0.1.13",
			expected: "kubescape/gateway",
		},
		{
			name:     "No registry in image ID",
			imageID:  "library/alpine@sha256:2345",
			expected: "library/alpine",
		},
		{
			name:     "Empty image ID",
			imageID:  "",
			expected: "",
		},
		{
			name:     "Image ID without hash",
			imageID:  "my.registry.io/library/alpine",
			expected: "library/alpine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRepositoryFromImageID(tt.imageID)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetImageTagFromImageID(t *testing.T) {
	tests := []struct {
		name     string
		imageID  string
		expected string
	}{
		{
			name:     "Valid image ID with hash",
			imageID:  "docker.io/library/alpine@sha256:2345",
			expected: "sha256:2345",
		},
		{
			name:     "Valid image ID with tag",
			imageID:  "quay.io/kubescape/gateway:v0.1.13",
			expected: "v0.1.13",
		},
		{
			name:     "No tag in image ID",
			imageID:  "docker.io/library/alpine",
			expected: "",
		},
		{
			name:     "Empty image ID",
			imageID:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getImageTagFromImageID(tt.imageID)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
