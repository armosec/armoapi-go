package containerscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBogusHash(t *testing.T) {
	tests := []struct {
		name       string
		attributes map[string]string
		want       string
	}{
		{
			name: "kube-proxy",
			attributes: map[string]string{
				"cluster":   "minikube",
				"namespace": "kube-system",
				"kind":      "daemonset",
				"name":      "kube-proxy",
			},
			want: "5485464254446115801",
		},
		{
			name: "coredns",
			attributes: map[string]string{
				"cluster":   "bez-longrun3",
				"namespace": "kube-system",
				"kind":      "deployment",
				"name":      "coredns",
			},
			want: "16878991644547272844",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GenerateBogusHash(tt.attributes), "GenerateWorkloadHash(%v)", tt.attributes)
		})
	}
}

func TestGenerateWorkloadHash(t *testing.T) {
	tests := []struct {
		name       string
		attributes map[string]string
		want       string
	}{
		{
			name: "kube-proxy",
			attributes: map[string]string{
				"cluster":   "minikube",
				"namespace": "kube-system",
				"kind":      "daemonset",
				"name":      "kube-proxy",
			},
			want: "8248822369989173472",
		},
		{
			name: "coredns",
			attributes: map[string]string{
				"cluster":   "bez-longrun3",
				"namespace": "kube-system",
				"kind":      "deployment",
				"name":      "coredns",
			},
			want: "12836087988784946749",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GenerateWorkloadHash(tt.attributes), "GenerateWorkloadHash(%v)", tt.attributes)
		})
	}
}
