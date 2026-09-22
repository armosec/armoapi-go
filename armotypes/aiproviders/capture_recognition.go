package aiproviders

import (
	"net"
	"strings"
)

// This file holds the host-recognition helpers shared by event-ingester-service's
// AI-Sandbox reassembler (which tags a provider / gen_ai.system per transaction) and
// config-service's HTTP-capture write-time host-verification warning (SUB-8712). Keeping
// them here, next to ClassifyEndpoint, is the whole point of SUB-8712 item 5: the set of
// hosts config-service considers "recognized" is DERIVED FROM the exact catalog the
// reassembler tags from, so the two can never drift.

// hostForMatch strips an optional :port (a captured Host header / capture-rule host can
// carry one — an in-cluster gateway on :8080, or the standard :443), then lowercases and
// trims a trailing dot. Mirrors ClassifyEndpoint's own port handling so every recognizer
// in this package agrees that host:port and host are the same host.
//
// Surrounding whitespace is trimmed BEFORE the port split (a leading/trailing space would
// otherwise make net.SplitHostPort fail and leave the port on, so " api.openai.com:443 "
// would miss recognition).
func hostForMatch(host string) string {
	host = strings.TrimSpace(host)
	if hostOnly, _, err := net.SplitHostPort(host); err == nil {
		host = hostOnly
	}
	return normalizeHost(host)
}

// googleGenAIHostTokens are curated substrings naming Google's known generative-AI
// surfaces. Each is specific enough that it would not plausibly appear in an unrelated
// GCP service's hostname (checked against the standard service list in
// TestIsCoarseGoogleGenAIHost_DoesNotMatchUnrelatedGCPServices) — deliberately NOT a
// short/generic token like "ai" or "genai" alone, which risks an incidental substring
// hit.
var googleGenAIHostTokens = []string{
	"generativelanguage",
	"generativeai",
	"aiplatform",
	"vertexai",
	"cloudcode", // the daily-cloudcode-pa.googleapis.com motivating host (SUB-8712)
}

// IsCoarseGoogleGenAIHost reports whether an otherwise-unclassified host is nonetheless
// plausibly a Google generative-AI surface: it must be under the googleapis.com domain
// AND contain one of the curated Gen-AI naming tokens. Both conditions are required — a
// token match outside googleapis.com (some unrelated host that merely contains
// "aiplatform") does not count, and a googleapis.com host that matches no token does not
// count either.
//
// ClassifyEndpoint (the strict catalog) only recognizes an EXACT whitelist for Google
// (generativelanguage.googleapis.com, aiplatform.googleapis.com, and the
// *-aiplatform.googleapis.com Vertex convention) — unlike Bedrock/Azure/OpenAI, whose
// catalog rules already match by domain prefix/suffix across their whole family. This is
// the gap the SUB-8712 motivating case fell into: daily-cloudcode-pa.googleapis.com
// (Google's internal Code-Assist/Antigravity wire) is nowhere in that whitelist.
//
// This is a DELIBERATELY NARROW widening, NOT a blanket *.googleapis.com match — that
// domain also serves Cloud Storage, Compute Engine, BigQuery, IAM, Cloud Run, GKE and
// dozens of other unrelated services, and matching all of it would badge ordinary
// GCP-using workloads as AI clients (the same false-positive class SUB-8116's
// control-plane exclusion fought to eliminate for Bedrock). It is a WEAKER,
// lower-confidence signal than ClassifyEndpoint.
func IsCoarseGoogleGenAIHost(host string) bool {
	h := hostForMatch(host)
	if h == "" || !strings.HasSuffix(h, ".googleapis.com") {
		return false
	}
	for _, tok := range googleGenAIHostTokens {
		if strings.Contains(h, tok) {
			return true
		}
	}
	return false
}

// IsRecognizedInferenceHost reports whether a host is a STRICT-catalog-recognized AI
// inference endpoint — ClassifyEndpoint recognizes it AND its kind implies AI use
// (KindImpliesAIUse excludes control-plane management APIs). This is the confident "this
// is really an AI inference host" signal: config-service uses it as the anchor that marks
// a capture config as an "AI provider bucket" (SUB-8712 item 5), and it is the
// port-agnostic complement to ClassifyEndpoint's own gating.
func IsRecognizedInferenceHost(host string) bool {
	_, kind, _, ok := ClassifyEndpoint(host)
	return ok && KindImpliesAIUse(kind)
}

// RecognizedForCapture reports whether the AI-Sandbox derivation code will attach a
// provider / gen_ai.system to a transaction for this host — the strict inference catalog
// OR the coarse Google matcher (which the reassembler's fragment_derive.go host-only
// fallback also consults, SUB-8712 item 4). config-service uses this as the exact
// "already handled, do not warn" set, so its warning can never fire for a host the
// reassembler in fact tags.
func RecognizedForCapture(host string) bool {
	return IsRecognizedInferenceHost(host) || IsCoarseGoogleGenAIHost(host)
}

// HostDomainFamily returns a coarse registrable-domain key for grouping hosts by
// provider family — the last two dot-labels of the (port-stripped, lowercased) host, so
// generativelanguage.googleapis.com and daily-cloudcode-pa.googleapis.com share the key
// "googleapis.com". config-service's SUB-8712 warning uses it to fire only when an
// unrecognized full-capture host shares a family with a recognized AI host in the same
// config ("added to an existing provider bucket").
//
// A deliberate simplification of the public-suffix list (golang.org/x/net/publicsuffix):
// the last-two-labels rule is correct for every PUBLIC domain the catalog matches on
// (googleapis.com, amazonaws.com, openai.com, azure.com, anthropic.com, …), all of which
// are plain .com registrable domains, and this package intentionally stays dependency-free
// (stdlib only). A host with fewer than two labels returns itself.
//
// Kubernetes in-cluster service DNS is the one shape the last-two-labels rule gets wrong:
// vllm.serving.svc.cluster.local and payments.default.svc.cluster.local would both
// collapse to "cluster.local", so an unrelated in-cluster service would falsely share a
// family with a recognized in-cluster AI gateway (vllm/ollama/litellm ARE catalog-matched
// on their .svc.cluster.local component). For a *.svc.cluster.local host the family key is
// therefore the FULL host (service.namespace preserved), so two distinct in-cluster
// services never share a family — an in-cluster host only "shares a bucket" with itself.
func HostDomainFamily(host string) string {
	h := hostForMatch(host)
	if h == "" {
		return ""
	}
	// In-cluster service DNS: keep the whole host as its own family (see doc above).
	if strings.HasSuffix(h, ".svc.cluster.local") || strings.HasSuffix(h, ".cluster.local") {
		return h
	}
	labels := strings.Split(h, ".")
	if len(labels) < 2 {
		return h
	}
	return strings.Join(labels[len(labels)-2:], ".")
}
