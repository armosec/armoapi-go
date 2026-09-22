package aiproviders

import "testing"

func TestIsCoarseGoogleGenAIHost_MatchesKnownGenAITokens(t *testing.T) {
	for _, host := range []string{
		"generativelanguage.googleapis.com",
		"aiplatform.googleapis.com",
		"us-central1-aiplatform.googleapis.com",
		"daily-cloudcode-pa.googleapis.com",     // the SUB-8712 motivating host
		"CLOUDCODE-PA.GOOGLEAPIS.COM",           // case-insensitive
		"daily-cloudcode-pa.googleapis.com:443", // port stripped before matching
	} {
		if !IsCoarseGoogleGenAIHost(host) {
			t.Errorf("IsCoarseGoogleGenAIHost(%q) = false, want true", host)
		}
	}
}

// The false-positive boundary — unrelated GCP services under the same googleapis.com
// domain must never match; a blanket domain match would badge ordinary GCP traffic.
func TestIsCoarseGoogleGenAIHost_DoesNotMatchUnrelatedGCPServices(t *testing.T) {
	for _, host := range []string{
		"storage.googleapis.com",
		"compute.googleapis.com",
		"bigquery.googleapis.com",
		"pubsub.googleapis.com",
		"iam.googleapis.com",
		"cloudcode.evil-mirror.com", // token, wrong domain
		"aiplatform.example.org",    // token, wrong domain
		"",
	} {
		if IsCoarseGoogleGenAIHost(host) {
			t.Errorf("IsCoarseGoogleGenAIHost(%q) = true, want false", host)
		}
	}
}

func TestIsRecognizedInferenceHost(t *testing.T) {
	cases := []struct {
		host string
		want bool
	}{
		{"generativelanguage.googleapis.com", true},           // Google inference (strict)
		{"bedrock-runtime.us-east-1.amazonaws.com", true},     // Bedrock data plane
		{"bedrock-runtime.us-east-1.amazonaws.com:443", true}, // port-agnostic
		{"bedrock.us-east-1.amazonaws.com", false},            // control-plane — not AI use
		{"daily-cloudcode-pa.googleapis.com", false},          // coarse-only, not strict
		{"storage.googleapis.com", false},
		{"shop.example.com", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := IsRecognizedInferenceHost(tc.host); got != tc.want {
			t.Errorf("IsRecognizedInferenceHost(%q) = %v, want %v", tc.host, got, tc.want)
		}
	}
}

func TestRecognizedForCapture(t *testing.T) {
	cases := []struct {
		host string
		want bool
	}{
		{"generativelanguage.googleapis.com", true}, // strict inference
		{"daily-cloudcode-pa.googleapis.com", true}, // coarse Google
		{"bedrock-runtime.us-east-1.amazonaws.com", true},
		{"bedrock.us-east-1.amazonaws.com", false},     // control-plane: neither strict-inference nor coarse
		{"storage.googleapis.com", false},              // unrelated GCP
		{"daily-antigravity-pa.googleapis.com", false}, // novel Google host, no token yet
		{"shop.example.com", false},
	}
	for _, tc := range cases {
		if got := RecognizedForCapture(tc.host); got != tc.want {
			t.Errorf("RecognizedForCapture(%q) = %v, want %v", tc.host, got, tc.want)
		}
	}
}

func TestHostDomainFamily(t *testing.T) {
	cases := []struct {
		host string
		want string
	}{
		{"generativelanguage.googleapis.com", "googleapis.com"},
		{"daily-antigravity-pa.googleapis.com", "googleapis.com"},
		{"daily-cloudcode-pa.googleapis.com:443", "googleapis.com"}, // port stripped
		{"bedrock-runtime.us-east-1.amazonaws.com", "amazonaws.com"},
		{"api.openai.com", "openai.com"},
		{" api.openai.com:443 ", "openai.com"}, // surrounding whitespace trimmed before port split
		{"shop.example.com", "example.com"},
		{"localhost", "localhost"}, // single label returns itself
		{"", ""},
		// Kubernetes in-cluster DNS: distinct services must NOT collapse to cluster.local,
		// else an unrelated in-cluster host would falsely share a family with a recognized
		// in-cluster AI gateway. The full host is its own family.
		{"vllm.serving.svc.cluster.local", "vllm.serving.svc.cluster.local"},
		{"payments.default.svc.cluster.local", "payments.default.svc.cluster.local"},
	}
	for _, tc := range cases {
		if got := HostDomainFamily(tc.host); got != tc.want {
			t.Errorf("HostDomainFamily(%q) = %q, want %q", tc.host, got, tc.want)
		}
	}
}

// Two distinct in-cluster services must not be grouped into the same family (the
// cluster.local collapse a naive last-two-labels rule would produce).
func TestHostDomainFamily_InClusterServicesAreDistinct(t *testing.T) {
	a := HostDomainFamily("vllm.serving.svc.cluster.local")
	b := HostDomainFamily("payments.default.svc.cluster.local")
	if a == b {
		t.Errorf("distinct in-cluster services must not share a family, both = %q", a)
	}
}

// Whitespace around a recognized host+port must not defeat recognition.
func TestRecognizedForCapture_TrimsSurroundingWhitespace(t *testing.T) {
	if !RecognizedForCapture(" bedrock-runtime.us-east-1.amazonaws.com:443 ") {
		t.Error("a recognized host with surrounding whitespace + port must still be recognized")
	}
}
