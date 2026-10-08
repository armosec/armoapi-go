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
		// Mega-cloud umbrella domains group by service label, so unrelated AWS services
		// don't collapse into one family with Bedrock (would false-warn on S3/DynamoDB).
		{"bedrock-runtime.us-east-1.amazonaws.com", "bedrock-runtime.amazonaws.com"},
		{"s3.us-east-1.amazonaws.com", "s3.amazonaws.com"},
		{"myresource.openai.azure.com", "myresource.azure.com"},
		{"bedrock-runtime.us-east-1.api.aws", "bedrock-runtime.api.aws"},
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

// An unrelated AWS service (S3) must NOT share a family with a Bedrock host — else it
// would false-warn as an "in the AI bucket" capture host (the mega-cloud umbrella
// collapse). Google stays coarse by design (that case is accepted / needed).
func TestHostDomainFamily_UnrelatedAWSServiceNotGroupedWithBedrock(t *testing.T) {
	bedrock := HostDomainFamily("bedrock-runtime.us-east-1.amazonaws.com")
	s3 := HostDomainFamily("s3.us-east-1.amazonaws.com")
	if bedrock == s3 {
		t.Errorf("S3 must not share Bedrock's family, both = %q", bedrock)
	}
	// Same Bedrock service across regions still shares a family (service, not region).
	if HostDomainFamily("bedrock-runtime.eu-west-1.amazonaws.com") != bedrock {
		t.Error("same AWS service across regions must share one family")
	}
}

// Whitespace around a recognized host+port must not defeat recognition.
func TestRecognizedForCapture_TrimsSurroundingWhitespace(t *testing.T) {
	if !RecognizedForCapture(" bedrock-runtime.us-east-1.amazonaws.com:443 ") {
		t.Error("a recognized host with surrounding whitespace + port must still be recognized")
	}
}
