package aisandbox

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// TestIdentityKey_AllKinds is the exhaustive per-kind proof that the discriminated
// <tag>:<canonical-body> format handles every identity kind the system produces
// — recognized (K8s / AWS role+user / GCP / Azure), the total unknown fallback,
// and the genuinely-empty (unattributed) case — with correct normalization.
func TestIdentityKey_AllKinds(t *testing.T) {
	type build func() (string, bool)
	cases := []struct {
		name    string
		build   build
		wantKey string
		wantOK  bool
		wantTag string // expected IdentityKeyTag(wantKey); "" when wantOK is false
	}{
		// --- K8s ServiceAccount ---
		{"k8s ok", func() (string, bool) { return K8sServiceAccountKey("prod", "payments", "checkout-sa") },
			"k8s-sa:prod/payments/checkout-sa", true, "k8s-sa"},
		{"k8s trims", func() (string, bool) { return K8sServiceAccountKey("  prod ", " payments", "checkout-sa ") },
			"k8s-sa:prod/payments/checkout-sa", true, "k8s-sa"},
		{"k8s missing cluster", func() (string, bool) { return K8sServiceAccountKey("", "payments", "checkout-sa") },
			"", false, ""},
		{"k8s missing namespace", func() (string, bool) { return K8sServiceAccountKey("prod", "", "checkout-sa") },
			"", false, ""},
		{"k8s missing sa", func() (string, bool) { return K8sServiceAccountKey("prod", "payments", "") },
			"", false, ""},

		// --- AWS IAM role (ARN, case preserved, embedded colons) ---
		{"aws role ok", func() (string, bool) { return AWSRoleKey("arn:aws:iam::123456789012:role/Checkout") },
			"aws-role:arn:aws:iam::123456789012:role/Checkout", true, "aws-role"},
		{"aws role preserves case", func() (string, bool) { return AWSRoleKey("arn:aws:iam::1:role/MixedCase") },
			"aws-role:arn:aws:iam::1:role/MixedCase", true, "aws-role"},
		{"aws role empty", func() (string, bool) { return AWSRoleKey("   ") }, "", false, ""},

		// --- AWS IAM user ---
		{"aws user ok", func() (string, bool) { return AWSUserKey("arn:aws:iam::123456789012:user/svc") },
			"aws-user:arn:aws:iam::123456789012:user/svc", true, "aws-user"},

		// --- GCP service account (email, lower-cased) ---
		{"gcp ok", func() (string, bool) { return GCPServiceAccountKey("Checkout@Proj.iam.gserviceaccount.com") },
			"gcp-sa:checkout@proj.iam.gserviceaccount.com", true, "gcp-sa"},
		{"gcp empty", func() (string, bool) { return GCPServiceAccountKey("") }, "", false, ""},

		// --- Azure managed identity / SP (GUID lower-cased, tenant-qualified) ---
		{"azure with tenant", func() (string, bool) { return AzureIdentityKey("TENANT-1", "ABCD-EF") },
			"azure-id:tenant-1/abcd-ef", true, "azure-id"},
		{"azure no tenant", func() (string, bool) { return AzureIdentityKey("", "ABCD-EF") },
			"azure-id:abcd-ef", true, "azure-id"},
		{"azure missing clientID", func() (string, bool) { return AzureIdentityKey("tenant-1", "") },
			"", false, ""},

		// --- unknown / opaque (total fallback) ---
		{"opaque ok", func() (string, bool) { return OpaqueIdentityKey("some-federated-principal") },
			"unknown:some-federated-principal", true, "unknown"},
		{"opaque trims", func() (string, bool) { return OpaqueIdentityKey("  weird//id  ") },
			"unknown:weird//id", true, "unknown"},
		{"opaque empty -> unattributed", func() (string, bool) { return OpaqueIdentityKey("") }, "", false, ""},
		{"opaque whitespace -> unattributed", func() (string, bool) { return OpaqueIdentityKey("   \t ") }, "", false, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotKey, gotOK := tc.build()
			if gotOK != tc.wantOK {
				t.Fatalf("ok = %v, want %v (key=%q)", gotOK, tc.wantOK, gotKey)
			}
			if gotKey != tc.wantKey {
				t.Fatalf("key = %q, want %q", gotKey, tc.wantKey)
			}
			if tag := IdentityKeyTag(gotKey); tag != tc.wantTag {
				t.Fatalf("IdentityKeyTag(%q) = %q, want %q", gotKey, tag, tc.wantTag)
			}
		})
	}
}

// TestOpaqueIdentityKey_HashesPathologicalInput proves the unknown body is bounded:
// a principal longer than identityKeyMaxBodyLen is replaced by its sha256, so
// identity_key can never grow unbounded from an adversarial producer.
func TestOpaqueIdentityKey_HashesPathologicalInput(t *testing.T) {
	huge := strings.Repeat("x", identityKeyMaxBodyLen+1)
	key, ok := OpaqueIdentityKey(huge)
	if !ok {
		t.Fatal("expected ok for non-empty input")
	}
	sum := sha256.Sum256([]byte(huge))
	want := "unknown:sha256:" + hex.EncodeToString(sum[:])
	if key != want {
		t.Fatalf("key = %q, want %q", key, want)
	}
	if IdentityKeyTag(key) != AiSandboxIdentityTagUnknown {
		t.Fatalf("tag = %q, want %q", IdentityKeyTag(key), AiSandboxIdentityTagUnknown)
	}
	// At-threshold length is kept verbatim (not hashed).
	atLimit := strings.Repeat("y", identityKeyMaxBodyLen)
	key2, _ := OpaqueIdentityKey(atLimit)
	if key2 != "unknown:"+atLimit {
		t.Fatalf("at-limit key = %q, want verbatim", key2)
	}
}

// TestIdentityKey_CrossKindUniqueness proves the tag prefix prevents a collision
// between different kinds that happen to share the same underlying identifier.
func TestIdentityKey_CrossKindUniqueness(t *testing.T) {
	same := "shared-id"
	role, _ := AWSRoleKey(same)
	user, _ := AWSUserKey(same)
	opaque, _ := OpaqueIdentityKey(same)
	if role == user || role == opaque || user == opaque {
		t.Fatalf("keys collided across kinds: role=%q user=%q opaque=%q", role, user, opaque)
	}
}

// TestIdentityKey_Deterministic proves the same inputs always yield the same key
// (so independent producers converge on one row).
func TestIdentityKey_Deterministic(t *testing.T) {
	a, _ := K8sServiceAccountKey("prod", "ns", "sa")
	b, _ := K8sServiceAccountKey("prod", "ns", "sa")
	if a != b {
		t.Fatalf("non-deterministic: %q vs %q", a, b)
	}
}

// TestIdentityKeyTag_Malformed guards the tag extractor on bad input.
func TestIdentityKeyTag_Malformed(t *testing.T) {
	for _, k := range []string{"", "no-colon", ":leading-colon"} {
		if tag := IdentityKeyTag(k); tag != "" {
			t.Fatalf("IdentityKeyTag(%q) = %q, want empty", k, tag)
		}
	}
}

// BearerSessionKey requires BOTH coordinates (a missing half falls back to the
// per-fingerprint key at the caller — never a shared bucket), lower-cases the
// host, and round-trips through tag extraction.
func TestBearerSessionKey(t *testing.T) {
	key, ok := BearerSessionKey("rh-1", "Management.Azure.COM")
	if !ok || key != "bearer-session:rh-1/management.azure.com" {
		t.Fatalf("BearerSessionKey = %q, %v", key, ok)
	}
	if IdentityKeyTag(key) != AiSandboxIdentityTagBearerSession {
		t.Fatalf("tag = %q", IdentityKeyTag(key))
	}
	if _, ok := BearerSessionKey("", "host"); ok {
		t.Fatal("empty resourceHash must not form a key")
	}
	if _, ok := BearerSessionKey("rh-1", " "); ok {
		t.Fatal("empty host must not form a key")
	}
}
