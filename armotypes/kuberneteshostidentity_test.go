package armotypes

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const testHostMachineID = "0123456789abcdef0123456789abcdef"
const testHostFingerprint = "mf1-zt7egzdv74wbh5wcponay6w3s3rgw4cqae3pf3l2zbqj7l2futwq"
const testHostKey = "kh1-svueoocodoyh2esjs3npuw2ptxgxopn344ojyvpjxrc4nrwz2xca"

func testKubernetesHostIdentity() KubernetesHostIdentity {
	return KubernetesHostIdentity{Version: 1, ClusterUID: "cluster-uid", ClusterName: "cluster", NodeUID: "node-uid", NodeName: "node", MachineFingerprint: testHostFingerprint, Key: testHostKey}
}

func TestKubernetesHostIdentityGoldenVectors(t *testing.T) {
	// Expected digests were computed independently with Python hashlib, struct.pack
	// and base64.b32encode. The Unicode vector verifies byte rather than rune length.
	fingerprint, err := KubernetesHostMachineFingerprint(testHostMachineID)
	require.NoError(t, err)
	assert.Equal(t, testHostFingerprint, fingerprint)
	for _, tc := range []struct{ cluster, node, key string }{
		{"cluster-uid", "node-uid", testHostKey},
		{"ab", "c", "kh1-aaujg4b3jcfmt6eit6ztzsj6qj2dlz6p4sqrubveb6hsdbuzhqia"},
		{"a", "bc", "kh1-p5f5667kyyinhap4iss5pvdwsive7de4qxxdjozhhzimjlp5piiq"},
		{"集群", "节点", "kh1-7rhmefzugo6rs3yba5fayd2qgaxlly7alzazfdcpsgufrtp2d26a"},
	} {
		t.Run(tc.cluster+"/"+tc.node, func(t *testing.T) {
			key, err := KubernetesHostKey(tc.cluster, tc.node, fingerprint)
			require.NoError(t, err)
			assert.Equal(t, tc.key, key)
			assert.Len(t, key, 56)
			assert.Regexp(t, `^kh1-[a-z2-7]{52}$`, key)
		})
	}
}

func TestNormalizeKubernetesHostMachineID(t *testing.T) {
	for _, input := range []string{testHostMachineID, " \t\n\r\v\f01234567-89AB-CDEF-0123-456789ABCDEF\r\n"} {
		normalized, err := NormalizeKubernetesHostMachineID(input)
		require.NoError(t, err)
		assert.Equal(t, testHostMachineID, normalized)
		fingerprint, err := KubernetesHostMachineFingerprint(input)
		require.NoError(t, err)
		assert.Equal(t, testHostFingerprint, fingerprint)
	}
	for _, input := range []string{"", " ", strings.Repeat("0", 32), "00000000-0000-0000-0000-000000000000", "0123", testHostMachineID + "0", "g" + testHostMachineID[1:], "\u00a0" + testHostMachineID, testHostMachineID[:16] + " " + testHostMachineID[16:]} {
		_, err := NormalizeKubernetesHostMachineID(input)
		require.Error(t, err)
		_, err = KubernetesHostMachineFingerprint(input)
		require.Error(t, err)
		assert.NotContains(t, err.Error(), testHostMachineID)
	}
}

func TestKubernetesHostIdentityGeneration(t *testing.T) {
	original := testKubernetesHostIdentity()
	for _, change := range []struct {
		name   string
		mutate func(*KubernetesHostIdentity)
	}{
		{"agent restart", func(i *KubernetesHostIdentity) {}},
		{"host reboot", func(i *KubernetesHostIdentity) {}},
		{"node display name", func(i *KubernetesHostIdentity) { i.NodeName = "renamed" }},
		{"cluster display name", func(i *KubernetesHostIdentity) { i.ClusterName = "renamed" }},
		{"cloud enrichment", func(i *KubernetesHostIdentity) { i.ProviderID = "digitalocean://602659893" }},
	} {
		t.Run(change.name, func(t *testing.T) {
			identity := original
			change.mutate(&identity)
			require.NoError(t, identity.Validate())
			assert.Equal(t, original.Key, identity.Key)
		})
	}
	for _, change := range []struct {
		name   string
		mutate func(*KubernetesHostIdentity)
	}{
		{"replacement with same name", func(i *KubernetesHostIdentity) {
			i.MachineFingerprint, _ = KubernetesHostMachineFingerprint("1123456789abcdef0123456789abcdef")
		}},
		{"cloned machine ID on distinct Node", func(i *KubernetesHostIdentity) { i.NodeUID = "other-node" }},
		{"recreated Node on same hardware", func(i *KubernetesHostIdentity) { i.NodeUID = "recreated-node" }},
		{"different cluster", func(i *KubernetesHostIdentity) { i.ClusterUID = "other-cluster" }},
	} {
		t.Run(change.name, func(t *testing.T) {
			identity := original
			change.mutate(&identity)
			require.Error(t, identity.Validate())
			key, err := KubernetesHostKey(identity.ClusterUID, identity.NodeUID, identity.MachineFingerprint)
			require.NoError(t, err)
			assert.NotEqual(t, original.Key, key)
			identity.Key = key
			require.NoError(t, identity.Validate())
		})
	}
}

func TestKubernetesHostIdentityRejectInvalid(t *testing.T) {
	for name, mutate := range map[string]func(*KubernetesHostIdentity){
		"unknown version":                       func(i *KubernetesHostIdentity) { i.Version = 2 },
		"missing version":                       func(i *KubernetesHostIdentity) { i.Version = 0 },
		"missing cluster UID":                   func(i *KubernetesHostIdentity) { i.ClusterUID = "" },
		"missing node UID":                      func(i *KubernetesHostIdentity) { i.NodeUID = "" },
		"missing cluster name":                  func(i *KubernetesHostIdentity) { i.ClusterName = "" },
		"missing node name":                     func(i *KubernetesHostIdentity) { i.NodeName = "" },
		"missing fingerprint":                   func(i *KubernetesHostIdentity) { i.MachineFingerprint = "" },
		"uppercase fingerprint":                 func(i *KubernetesHostIdentity) { i.MachineFingerprint = strings.ToUpper(i.MachineFingerprint) },
		"noncanonical fingerprint padding bits": func(i *KubernetesHostIdentity) { i.MachineFingerprint = testHostFingerprint[:55] + "r" },
		"wrong fingerprint domain":              func(i *KubernetesHostIdentity) { i.MachineFingerprint = "kh1-" + testHostFingerprint[4:] },
		"invalid fingerprint alphabet":          func(i *KubernetesHostIdentity) { i.MachineFingerprint = testHostFingerprint[:55] + "1" },
		"missing key":                           func(i *KubernetesHostIdentity) { i.Key = "" },
		"wrong key":                             func(i *KubernetesHostIdentity) { i.Key = testHostKey[:55] + "b" },
		"invalid UID UTF8":                      func(i *KubernetesHostIdentity) { i.NodeUID = "\xff" },
		"invalid name UTF8":                     func(i *KubernetesHostIdentity) { i.ClusterName = "\xff" },
		"invalid provider UTF8":                 func(i *KubernetesHostIdentity) { i.ProviderID = "\xff" },
	} {
		t.Run(name, func(t *testing.T) {
			identity := testKubernetesHostIdentity()
			mutate(&identity)
			require.Error(t, identity.Validate())
			canonical, err := identity.CanonicalJSON()
			require.Error(t, err)
			assert.Empty(t, canonical)
		})
	}
}

func TestKubernetesHostIdentityWireCompatibility(t *testing.T) {
	const legacy = `{"agent_version":"v1","sensor_updated":"2026-09-22T00:00:00Z","timestamp":"2026-09-22T01:00:00Z","cloudMetadata":{"host_type":"ec2","instance_id":"i-123","machine_id":"legacy-machine"}}`
	var report HealthReport
	require.NoError(t, json.Unmarshal([]byte(legacy), &report))
	assert.Nil(t, report.CloudMetadata.KubernetesHostIdentity)
	data, err := json.Marshal(report)
	require.NoError(t, err)
	assert.JSONEq(t, legacy, string(data))
	legacyBSON, err := bson.Marshal(report.CloudMetadata)
	require.NoError(t, err)
	var fields bson.M
	require.NoError(t, bson.Unmarshal(legacyBSON, &fields))
	assert.NotContains(t, fields, KubernetesHostIdentityMetadataKey)
	var restored CloudMetadata
	require.NoError(t, bson.Unmarshal(legacyBSON, &restored))
	assert.Equal(t, report.CloudMetadata, restored)
	identity := testKubernetesHostIdentity()
	report.CloudMetadata.KubernetesHostIdentity = &identity
	for _, provider := range []string{"", "digitalocean://602659893", "azure:///subscriptions/s/resourceGroups/r/providers/Microsoft.Compute/virtualMachineScaleSets/a/virtualMachines/0"} {
		identity.ProviderID = provider
		canonical, err := identity.CanonicalJSON()
		require.NoError(t, err)
		expected := `{"version":1,"cluster_uid":"cluster-uid","cluster_name":"cluster","node_uid":"node-uid","node_name":"node","machine_fingerprint":"` + testHostFingerprint + `","key":"` + testHostKey + `"`
		if provider != "" {
			expected += `,"provider_id":"` + provider + `"`
		}
		expected += "}"
		assert.Equal(t, expected, canonical)
		assert.NotContains(t, canonical, testHostMachineID)
		data, err = json.Marshal(report)
		require.NoError(t, err)
		var result HealthReport
		require.NoError(t, json.Unmarshal(data, &result))
		assert.Equal(t, report, result)
		data, err = bson.Marshal(report.CloudMetadata)
		require.NoError(t, err)
		require.NoError(t, bson.Unmarshal(data, &restored))
		assert.Equal(t, report.CloudMetadata, restored)
		require.NoError(t, bson.Unmarshal(data, &fields))
		nested, ok := fields[KubernetesHostIdentityMetadataKey].(bson.M)
		require.True(t, ok)
		var jsonFields map[string]any
		require.NoError(t, json.Unmarshal([]byte(canonical), &jsonFields))
		assert.Len(t, nested, len(jsonFields))
		for key := range jsonFields {
			assert.Contains(t, nested, key)
		}
	}
}
