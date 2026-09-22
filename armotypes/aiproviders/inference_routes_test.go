package aiproviders

import (
	"regexp"
	"testing"
)

// inferenceRouteMatchCases is the shared corpus both directions of the R6 drift
// guard use: MatchInferenceRoute behaviour here, and the generated-SQL-regex
// equivalence below. Paths are silver-shaped (query-free, percent-decoded) unless a
// case exercises the defensive query strip.
var inferenceRouteMatchCases = []struct {
	name       string
	path       string
	wantFamily string
	wantOK     bool
}{
	// --- OpenAI Responses API (adapter landing: PLAN-codex-responses-adapter) ------
	// codex 0.149.x speaks the Responses API over a WebSocket upgrade to
	// chatgpt.com. Its path ends /codex/responses, NOT /v1/responses, so the
	// pattern has to admit both spellings or codex keeps folding as `utility`.
	{name: "openai responses, codex websocket path", path: "/backend-api/codex/responses", wantFamily: RouteFamilyOpenAIResponses, wantOK: true},
	{name: "openai responses, canonical http path", path: "/v1/responses", wantFamily: RouteFamilyOpenAIResponses, wantOK: true},
	// A bare /responses suffix is deliberately NOT admitted: matching ignores the
	// host, so a survey or forms API's /api/responses would otherwise classify as
	// LLM inference and mint conversation cards out of nothing.
	{name: "unrelated responses path is not inference", path: "/api/survey/responses", wantFamily: "", wantOK: false},
	// --- Anthropic Messages ---
	{"anthropic messages", "/v1/messages", RouteFamilyAnthropicMessages, true},
	{"anthropic behind gateway prefix", "/anthropic/v1/messages", RouteFamilyAnthropicMessages, true},
	{"messages without v1 segment does not match", "/messages", "", false},
	{"embedded v1/messages without slash boundary does not match", "/av1/messages", "", false},

	// --- Bedrock runtime (the four adapter-backed invocation routes) ---
	{"bedrock invoke", "/model/anthropic.claude-3-5-sonnet-20241022-v2:0/invoke", RouteFamilyBedrockRuntime, true},
	{"bedrock invoke-with-response-stream", "/model/anthropic.claude-haiku-4-5/invoke-with-response-stream", RouteFamilyBedrockRuntime, true},
	{"bedrock converse", "/model/amazon.nova-lite-v1:0/converse", RouteFamilyBedrockRuntime, true},
	{"bedrock converse-stream", "/model/amazon.nova-lite-v1:0/converse-stream", RouteFamilyBedrockRuntime, true},
	// silver stores req.URL.Path — percent-DECODED, so an inference-profile ARN id
	// carries literal slashes and must still match.
	{"bedrock decoded ARN model id", "/model/arn:aws:bedrock:us-east-1:015253967648:inference-profile/us.anthropic.claude-haiku-4-5-20251001-v1:0/invoke", RouteFamilyBedrockRuntime, true},
	{"bedrock control-plane-ish path does not match", "/foundation-models", "", false},
	{"bedrock unknown verb does not match", "/model/anthropic.claude/apply-guardrail", "", false},
	{"bedrock verb without model segment does not match", "/invoke", "", false},

	// --- OpenAI-compatible chat completions (suffix — gateways included) ---
	{"openai chat completions", "/v1/chat/completions", RouteFamilyOpenAIChat, true},
	{"gateway chat completions", "/api/v1/chat/completions", RouteFamilyOpenAIChat, true},
	{"azure deployment chat completions", "/openai/deployments/gpt-4o/chat/completions", RouteFamilyOpenAIChat, true},
	// Deferred until their adapters exist (design v1 rule): the fold has no
	// normalizer/session extractor for these, and an (E1) match without an adapter
	// folds an unnamed promptless card.
	{"legacy text completions NOT registered", "/v1/completions", "", false},
	{"embeddings does not match", "/v1/embeddings", "", false},
	{"models listing does not match", "/v1/models", "", false},

	// --- Gemini / Vertex generate ---
	{"gemini generateContent", "/v1beta/models/gemini-2.5-flash:generateContent", RouteFamilyGeminiGenerate, true},
	{"gemini streamGenerateContent", "/v1beta/models/gemini-2.5-flash:streamGenerateContent", RouteFamilyGeminiGenerate, true},
	{"vertex publisher model generateContent", "/v1/projects/p/locations/us-central1/publishers/google/models/gemini-1.5-pro:generateContent", RouteFamilyGeminiGenerate, true},
	// Derive-side compat: isGeminiGeneratePath accepts a raw capture path with a
	// query tail; the registry's defensive strip must accept it too.
	{"gemini stream with query tail", "/v1beta/models/gemini-2.5-flash:streamGenerateContent?alt=sse", RouteFamilyGeminiGenerate, true},
	{"gemini embedContent does not match", "/v1beta/models/gemini-2.5-flash:embedContent", "", false},
	{"gemini countTokens does not match", "/v1beta/models/gemini-2.5-flash:countTokens", "", false},

	// --- degenerate inputs ---
	{"empty path", "", "", false},
	{"bare slash", "/", "", false},
	{"query-only path", "?alt=sse", "", false},
}

func TestMatchInferenceRoute(t *testing.T) {
	for _, tc := range inferenceRouteMatchCases {
		t.Run(tc.name, func(t *testing.T) {
			family, ok := MatchInferenceRoute(tc.path)
			if ok != tc.wantOK || family != tc.wantFamily {
				t.Errorf("MatchInferenceRoute(%q) = (%q, %v), want (%q, %v)",
					tc.path, family, ok, tc.wantFamily, tc.wantOK)
			}
		})
	}
}

// TestInferenceRouteRegistry_AdapterBackedOnly is direction (b) of the design's R6
// two-direction drift guard: the registry registers NOTHING the Conversations fold
// lacks an adapter for. The family set is pinned CLOSED — adding a family here is a
// deliberate act that must ship with its adapter (normalizer + stream reassembler +
// session extractor), never a drive-by registration. (Direction (a) — the registry
// covers every path shape the derive-side helpers accept — lives next to those
// helpers: ingesters/ai_sandbox_ingester/inference_route_parity_test.go.)
func TestInferenceRouteRegistry_AdapterBackedOnly(t *testing.T) {
	wantFamilies := map[string]bool{
		RouteFamilyAnthropicMessages: true,
		RouteFamilyBedrockRuntime:    true,
		RouteFamilyOpenAIChat:        true,
		RouteFamilyGeminiGenerate:    true,
		RouteFamilyOpenAIResponses:   true,
	}
	if len(inferenceRoutes) != len(wantFamilies) {
		t.Fatalf("registry has %d routes, want exactly the %d adapter-backed families", len(inferenceRoutes), len(wantFamilies))
	}
	seen := map[string]bool{}
	for _, r := range inferenceRoutes {
		if !wantFamilies[r.family] {
			t.Errorf("route family %q is not in the adapter-backed v1 set — register it only WITH its adapter", r.family)
		}
		if seen[r.family] {
			t.Errorf("route family %q registered twice", r.family)
		}
		seen[r.family] = true
	}
	// The remaining explicitly-deferred OpenAI shape stays out until its adapter
	// exists. `/v1/responses` used to sit here and graduated when openai_responses.go
	// landed; `/v1/completions` still has no adapter.
	for _, deferred := range []string{"/v1/completions"} {
		if _, ok := MatchInferenceRoute(deferred); ok {
			t.Errorf("deferred route %q must NOT match until its adapter exists (design v1 rule)", deferred)
		}
	}
}

// TestInferenceRouteSuffixRegex_AgreesWithGoMatcher pins the generated SQL prefilter
// regex to the Go matcher over the whole corpus: for every case, the single combined
// regex (compiled by Go's RE2 — the shared syntax subset with Trino's regexp_like)
// must reach the same verdict as MatchInferenceRoute. This is the Go-side half of the
// single-source guarantee; the SQL-side half (the literal in silver_transactions.sql
// equals this generated string) is pinned in lakequery, and the executing Trino
// sqltest exercises the literal against real rows.
func TestInferenceRouteSuffixRegex_AgreesWithGoMatcher(t *testing.T) {
	combined := regexp.MustCompile(InferenceRouteSuffixRegex())
	for _, tc := range inferenceRouteMatchCases {
		// The SQL side sees the silver path, which is query-free by construction —
		// mirror MatchInferenceRoute's defensive strip for the one query-tail case.
		path := tc.path
		if q := indexByteString(path, '?'); q >= 0 {
			path = path[:q]
		}
		if got := combined.MatchString(path); got != tc.wantOK {
			t.Errorf("generated regex verdict for %q = %v, want %v (must agree with MatchInferenceRoute)", path, got, tc.wantOK)
		}
	}
	if got := InferenceRouteSuffixRegex(); indexByteString(got, '\'') >= 0 {
		t.Errorf("generated regex %q must not contain single quotes (it embeds in a SQL string literal)", got)
	}
}

// TestIsInferenceShapedPath covers the WIDER SUB-8240 predicate: the adapter-backed
// registry UNION the adapter-less routes. It must agree with MatchInferenceRoute
// everywhere the registry has an opinion, and additionally admit the deferred OpenAI
// shape — which the models gate needs (it is real inference, so a model-less call on
// it must still mint an unnamed model row) even though the Conversations fold must
// keep refusing it (no adapter, so it would fold a promptless card).
func TestIsInferenceShapedPath(t *testing.T) {
	// Everything the registry admits, the wider predicate admits too.
	for _, tc := range inferenceRouteMatchCases {
		if tc.wantOK && !IsInferenceShapedPath(tc.path) {
			t.Errorf("IsInferenceShapedPath(%q) = false, want true (registry routes are a subset)", tc.path)
		}
	}

	// The adapter-less additions: inference for the models gate, still NOT registry routes.
	for _, path := range []string{"/v1/completions", "/api/v1/completions"} {
		if !IsInferenceShapedPath(path) {
			t.Errorf("IsInferenceShapedPath(%q) = false, want true (adapter-less inference wire)", path)
		}
		if _, ok := MatchInferenceRoute(path); ok {
			t.Errorf("MatchInferenceRoute(%q) = true — the adapter-backed registry must stay unchanged", path)
		}
	}

	// Non-inference paths stay out of BOTH. The Claude Code housekeeping corpus from
	// SUB-8240 plus the Ollama/embeddings wires that are owned by the path_class
	// taxonomy (egressclass.IsModelPathClass), which this predicate deliberately does
	// not restate.
	for _, path := range []string{
		"/api/event_logging/v2/batch",
		"/mcp-registry/v0/servers",
		"/api/claude_cli/bootstrap",
		"/api/claude_code/metrics",
		"/api/hello",
		"/api/eval/foo",
		"/metrics",
		"/healthz",
		"",
		"/api/generate",
		"/api/chat",
		"/v1/embeddings",
	} {
		if IsInferenceShapedPath(path) {
			t.Errorf("IsInferenceShapedPath(%q) = true, want false", path)
		}
	}
}

func indexByteString(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
