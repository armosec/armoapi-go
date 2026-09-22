package armotypes

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

const (
	KubernetesHostIdentityVersion     = 1
	KubernetesHostKeyLabel            = "kubescape.io/k8s-host-key"
	KubernetesHostIdentityAnnotation  = "kubescape.io/k8s-host-identity"
	KubernetesHostIdentityMetadataKey = "kubernetes_host_identity"
)

// KubernetesHostIdentity identifies a machine generation of a Kubernetes Node.
// It is nested wire metadata, not a Portal resource. Field order defines its
// canonical JSON representation. It never contains the raw Linux machine-id.
type KubernetesHostIdentity struct {
	Version            int    `json:"version" bson:"version"`
	ClusterUID         string `json:"cluster_uid" bson:"cluster_uid"`
	ClusterName        string `json:"cluster_name" bson:"cluster_name"`
	NodeUID            string `json:"node_uid" bson:"node_uid"`
	NodeName           string `json:"node_name" bson:"node_name"`
	MachineFingerprint string `json:"machine_fingerprint" bson:"machine_fingerprint"`
	Key                string `json:"key" bson:"key"`
	ProviderID         string `json:"provider_id,omitempty" bson:"provider_id,omitempty"`
}

// NormalizeKubernetesHostMachineID canonicalizes a Linux machine-id, accepting
// UUID hyphens and surrounding ASCII whitespace. The zero ID is invalid.
func NormalizeKubernetesHostMachineID(machineID string) (string, error) {
	normalized := strings.ToLower(strings.ReplaceAll(strings.Trim(machineID, " \t\n\r\v\f"), "-", ""))
	if len(normalized) != 32 {
		return "", fmt.Errorf("machine-id must contain 32 hexadecimal digits")
	}
	if _, err := hex.DecodeString(normalized); err != nil {
		return "", fmt.Errorf("machine-id must contain only hexadecimal digits")
	}
	if normalized == strings.Repeat("0", 32) {
		return "", fmt.Errorf("machine-id must not be zero")
	}
	return normalized, nil
}

// KubernetesHostMachineFingerprint hashes the normalized machine-id using the
// domain-separated v1 length-framed encoding.
func KubernetesHostMachineFingerprint(machineID string) (string, error) {
	normalized, err := NormalizeKubernetesHostMachineID(machineID)
	if err != nil {
		return "", err
	}
	return kubernetesHostDigest("mf1-", "kubernetes-host-machine-v1", normalized)
}

// KubernetesHostKey derives a generation key from exact API UID strings and a
// canonical machine fingerprint. Names and cloud metadata do not enter the key.
func KubernetesHostKey(clusterUID, nodeUID, machineFingerprint string) (string, error) {
	if clusterUID == "" || nodeUID == "" {
		return "", fmt.Errorf("cluster UID and node UID are required")
	}
	if !validKubernetesHostDigest(machineFingerprint, "mf1-") {
		return "", fmt.Errorf("invalid Kubernetes host machine fingerprint")
	}
	return kubernetesHostDigest("kh1-", "kubernetes-host-key-v1", clusterUID, nodeUID, machineFingerprint)
}

// Validate checks v1 syntax and key consistency. It does not authenticate the
// envelope or verify its values against a registered cluster and synced Node.
func (identity KubernetesHostIdentity) Validate() error {
	if identity.Version != KubernetesHostIdentityVersion {
		return fmt.Errorf("unsupported Kubernetes host identity version: %d", identity.Version)
	}
	if identity.ClusterName == "" || identity.NodeName == "" {
		return fmt.Errorf("cluster name and node name are required")
	}
	if !utf8.ValidString(identity.ClusterName) || !utf8.ValidString(identity.NodeName) || !utf8.ValidString(identity.ProviderID) {
		return fmt.Errorf("Kubernetes host identity metadata must be valid UTF-8")
	}
	key, err := KubernetesHostKey(identity.ClusterUID, identity.NodeUID, identity.MachineFingerprint)
	if err != nil {
		return err
	}
	if identity.Key != key {
		return fmt.Errorf("Kubernetes host identity key does not match its generation")
	}
	return nil
}

// CanonicalJSON returns the validated envelope in v1 field order without
// whitespace, suitable for profile annotations and string-valued metadata.
func (identity KubernetesHostIdentity) CanonicalJSON() (string, error) {
	if err := identity.Validate(); err != nil {
		return "", err
	}
	data, err := json.Marshal(identity)
	if err != nil {
		return "", fmt.Errorf("marshal Kubernetes host identity: %w", err)
	}
	return string(data), nil
}

func kubernetesHostDigest(prefix string, values ...string) (string, error) {
	hash := sha256.New()
	var length [4]byte
	for _, value := range values {
		if !utf8.ValidString(value) || uint64(len(value)) > math.MaxUint32 {
			return "", fmt.Errorf("Kubernetes host identity value must be UTF-8 with a uint32 byte length")
		}
		binary.BigEndian.PutUint32(length[:], uint32(len(value)))
		hash.Write(length[:])
		hash.Write([]byte(value))
	}
	return prefix + strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash.Sum(nil))), nil
}

func validKubernetesHostDigest(value, prefix string) bool {
	if len(value) != 56 || !strings.HasPrefix(value, prefix) || value != strings.ToLower(value) {
		return false
	}
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	digest, err := encoding.DecodeString(strings.ToUpper(value[len(prefix):]))
	return err == nil && len(digest) == sha256.Size && strings.ToLower(encoding.EncodeToString(digest)) == value[len(prefix):]
}
