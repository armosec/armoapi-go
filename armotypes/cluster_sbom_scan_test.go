package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Round-trip a fully-populated message and assert no field is lost on
// the wire. Guards against accidental json-tag typos and against
// removing fields the consumer relies on.
func TestClusterSBOMScanMessage_RoundTrip(t *testing.T) {
	orig := ClusterSBOMScanMessage{
		CustomerGUID:         "11111111-1111-1111-1111-111111111111",
		Cluster:              "prod-east-1",
		WorkloadKind:         "Deployment",
		WorkloadName:         "checkout-api",
		WorkloadNamespace:    "shop",
		WorkloadResourceHash: "wlhash-deadbeef",
		ContainerProfileName: "cp-checkout-api-main",
		ImageDigest:          "sha256:abc123",
		SyftVersion:          "1.20.0",
		SBOMObjectRef:        `{"bucket":"sboms","key":"11111111/abc123/1.20.0/sbom.json"}`,
	}

	wire, err := json.Marshal(orig)
	require.NoError(t, err)

	var got ClusterSBOMScanMessage
	require.NoError(t, json.Unmarshal(wire, &got))
	assert.Equal(t, orig, got)
}

// The consumer reads exact json tag names, so freeze them — a rename
// here would silently break the Vulnerability Scanner's deserialiser.
func TestClusterSBOMScanMessage_JSONFieldNames(t *testing.T) {
	wire, err := json.Marshal(ClusterSBOMScanMessage{
		CustomerGUID:         "c",
		Cluster:              "cl",
		WorkloadKind:         "wk",
		WorkloadName:         "wn",
		WorkloadNamespace:    "wns",
		WorkloadResourceHash: "wrh",
		ContainerProfileName: "cpn",
		ImageDigest:          "id",
		SyftVersion:          "sv",
		SBOMObjectRef:        "sor",
	})
	require.NoError(t, err)

	var asMap map[string]string
	require.NoError(t, json.Unmarshal(wire, &asMap))
	assert.Equal(t, map[string]string{
		"customerGUID":         "c",
		"cluster":              "cl",
		"workloadKind":         "wk",
		"workloadName":         "wn",
		"workloadNamespace":    "wns",
		"workloadResourceHash": "wrh",
		"containerProfileName": "cpn",
		"imageDigest":          "id",
		"syftVersion":          "sv",
		"sbomObjectRef":        "sor",
	}, asMap)
}
