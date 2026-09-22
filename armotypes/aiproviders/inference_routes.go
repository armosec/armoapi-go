package aiproviders

// Inference-route registry (SUB-8184) — the SINGLE SOURCE OF TRUTH for "which HTTP
// request paths are LLM *inference* routes", consumed by the AI-Sandbox Conversations
// positive-inference predicate (design:
// specs/2026-08-21-ai-sandbox-conversations-positive-inference-predicate.md):
//
//   - the sessionizer's Go predicate calls MatchInferenceRoute per silver row — the
//     (E1) "positive: route" arm;
//   - the silver-read SQL prefilter embeds InferenceRouteSuffixRegex() as a Trino
//     `regexp_like(path, …)` OR-term (lakequery/silver_transactions.sql, clause 1) so
//     an unclassified-host inference call with NO derived signals still reaches the Go
//     predicate and its drift counter (default-deny WITH visibility) — a unit test in
//     lakequery pins the SQL literal to this generated regex, so the two can never
//     drift.
//
// Matching is PATH-SUFFIX and FAMILY-AGNOSTIC: the host is never consulted. A
// classified host may speak another family's wire shape — the in-cluster
// bedrock-access-gateway classifies as AWS Bedrock but serves the OpenAI-shaped
// `/api/v1/chat/completions` — so requiring host↔family agreement would drop real
// inference. Host normalization is not this registry's job either: the predicate's
// (N) arm calls ClassifyEndpoint, which port-strips + lowercases internally, and the
// path match below is case-insensitive.
//
// v1 = ADAPTER-BACKED ROUTES ONLY (design review finding #5): registering a route the
// fold has no normalizer/reassembler/session-extractor for would re-create the
// promptless-card class through the front door — an (E1) match without an adapter
// folds an unnamed card. `/v1/completions` (legacy text completions) is therefore
// deliberately NOT registered — it has no adapter — and joins the registry only when
// one exists. (`/v1/responses`, the OpenAI Responses API, WAS in that deferred set but
// GRADUATED into inferenceRoutes below once its adapter landed — see the deferred table
// at the bottom of this file.) The adapter-backed-only rule is pinned in BOTH directions
// by unit tests here and by the derive-helper parity test in ingesters/ai_sandbox_ingester.
//
// Route knowledge previously lived piecemeal in the derive-side helpers
// (bedrockModelFromPath, isOpenAIChatPath, isGeminiGeneratePath —
// ingesters/ai_sandbox_ingester); those migrate onto this registry in a named
// follow-up (design §Changes-3) to end the double maintenance.

import (
	"regexp"
	"strings"
)

// Inference-route families — one per provider wire family the fold has a full
// adapter (normalizer + stream reassembler + session extractor) for.
const (
	// RouteFamilyAnthropicMessages is the Anthropic Messages API (`/v1/messages`).
	RouteFamilyAnthropicMessages = "anthropic_messages"
	// RouteFamilyBedrockRuntime is the AWS Bedrock runtime data plane — the four
	// invocation routes `/model/{id}/invoke`, `/model/{id}/invoke-with-response-stream`,
	// `/model/{id}/converse`, `/model/{id}/converse-stream`. `{id}` may be a
	// percent-DECODED inference-profile ARN containing slashes (silver stores
	// req.URL.Path, which is decoded), so the id segment is matched as `.+`, never a
	// single path segment.
	RouteFamilyBedrockRuntime = "bedrock_runtime"
	// RouteFamilyOpenAIChat is the OpenAI-compatible Chat Completions wire
	// (`…/chat/completions` suffix) — genuine OpenAI, Azure OpenAI, and every
	// OpenAI-compatible gateway/self-hosted server (LiteLLM, vLLM, Ollama,
	// bedrock-access-gateway's `/api/v1/chat/completions`, …).
	RouteFamilyOpenAIChat = "openai_chat_completions"
	// RouteFamilyOpenAIResponses is the OpenAI Responses API wire — `/v1/responses`
	// over HTTP, and `/backend-api/codex/responses` over a WebSocket upgrade, which is
	// what codex 0.149.x speaks. One family: the streaming event set is the same, only
	// the transport differs, and the agent's wsdeflate has already reassembled the
	// WebSocket frames into JSON by the time the fold sees them.
	RouteFamilyOpenAIResponses = "openai_responses"
	// RouteFamilyGeminiGenerate is the Gemini / Vertex AI generate wire — the
	// `:generateContent` and `:streamGenerateContent` method suffixes.
	// `:embedContent` / `:countTokens` are different shapes and deliberately absent,
	// mirroring isGeminiGeneratePath.
	RouteFamilyGeminiGenerate = "gemini_generate"
)

// inferenceRoute is one registered inference wire route: a provider family plus a
// case-insensitive, END-ANCHORED regex fragment over the query-free request path.
// Fragments use only the syntax subset shared by Go's RE2 and Trino's regexp_like
// (alternation, `(?:…)` groups, `.`/`+` — no lookaround, no backreferences), because
// the SAME fragments render both the compiled Go matchers and the SQL prefilter
// regex (InferenceRouteSuffixRegex) — one table, two consumers, zero drift.
type inferenceRoute struct {
	family  string
	pattern string // end-anchored via the shared wrapper below; leading '/' or ':' guards the segment boundary
	re      *regexp.Regexp
}

// inferenceRoutes is the v1 registry table. Order is irrelevant (families are
// disjoint on real wires; MatchInferenceRoute returns the first match
// deterministically in table order).
var inferenceRoutes = []inferenceRoute{
	{family: RouteFamilyAnthropicMessages, pattern: `/v1/messages`},
	{family: RouteFamilyBedrockRuntime, pattern: `/model/.+/(?:invoke|invoke-with-response-stream|converse|converse-stream)`},
	{family: RouteFamilyOpenAIChat, pattern: `/chat/completions`},
	{family: RouteFamilyGeminiGenerate, pattern: `:(?:stream)?generatecontent`},
	// Both spellings are enumerated rather than matching a bare `/responses` suffix.
	// Matching never consults the host, so a bare suffix would classify a survey or
	// forms API's /api/survey/responses as LLM inference and fold cards from it. Same
	// reason the Bedrock entry above enumerates its four routes.
	{family: RouteFamilyOpenAIResponses, pattern: `(?:/v1|/backend-api/codex)/responses`},
}

// wrapRoutePattern turns a registry fragment into the full case-insensitive
// suffix-match regex — the ONE place the anchoring/flag convention lives, shared by
// the per-route Go matchers and the rendered SQL regex.
func wrapRoutePattern(pattern string) string {
	return `(?i)(?:` + pattern + `)$`
}

func init() {
	for i := range inferenceRoutes {
		inferenceRoutes[i].re = regexp.MustCompile(wrapRoutePattern(inferenceRoutes[i].pattern))
	}
	for i := range unadaptedInferenceRoutes {
		unadaptedInferenceRoutes[i].re = regexp.MustCompile(wrapRoutePattern(unadaptedInferenceRoutes[i].pattern))
	}
}

// MatchInferenceRoute reports whether a request path is a registered LLM inference
// route, and which family it belongs to. Matching is path-suffix, case-insensitive,
// and family-agnostic (the host is never consulted — see the package comment above).
//
// The silver `path` column is req.URL.Path — already query-free and percent-decoded —
// so no normalization is needed there. A `?query` tail is still stripped defensively
// because the derive-side callers this registry will absorb (isGeminiGeneratePath)
// receive raw capture paths that can carry one (`:streamGenerateContent?alt=sse`);
// on an already-clean path the strip is a no-op.
func MatchInferenceRoute(path string) (family string, ok bool) {
	if q := strings.IndexByte(path, '?'); q >= 0 {
		path = path[:q]
	}
	if path == "" {
		return "", false
	}
	for _, r := range inferenceRoutes {
		if r.re.MatchString(path) {
			return r.family, true
		}
	}
	return "", false
}

// unadaptedInferenceRoutes are wire routes that ARE genuine LLM inference but that
// the v1 registry above deliberately does NOT admit, because no fold adapter exists
// for them (see the adapter-backed-only rule in the package comment). They are kept
// in a SEPARATE table, and MatchInferenceRoute / InferenceRouteSuffixRegex never
// consult it, so the Conversations predicate and its pinned SQL prefilter see exactly
// the route set they saw before — registering these there would fold promptless
// cards, which is the failure the adapter rule exists to prevent.
//
// They matter to a DIFFERENT question. SUB-8240's models gate asks only "was this
// call inference?", to decide whether a model-id-less observation may mint an unnamed
// ("model not identified") row on the Models surface. For that question the answer is
// yes: a /v1/completions or /v1/responses call whose model could not be derived is
// real inference we failed to name, and dropping it would hide model usage — the
// exact defect SUB-8240 fixes — whereas folding a conversation card for it would
// invent content. Same fragments, same RE2/Trino-safe subset as the table above; a
// route graduates from here into inferenceRoutes when its adapter lands.
// (compiled by the same init as inferenceRoutes; family is left empty — these carry
// no family because nothing downstream may dispatch an adapter on them).
var unadaptedInferenceRoutes = []inferenceRoute{
	{pattern: `/v1/completions`}, // legacy OpenAI text completions
	// `/v1/responses` GRADUATED into inferenceRoutes above when its adapter landed
	// (openai_responses.go). It is intentionally absent here: leaving a graduated
	// route in both tables would make IsInferenceShapedPath's two arms overlap, and
	// the whole point of this table is the routes the fold cannot reconstruct.
}

// IsInferenceShapedPath reports whether a request path is an LLM INFERENCE wire of
// any kind — the adapter-backed registry (MatchInferenceRoute) UNION the
// adapter-less routes above. It is the predicate for consumers that must decide
// "was this inference?" without caring whether the fold can reconstruct it, and its
// only caller today is the SUB-8240 models minting gate
// (scheduled_tasks/ai_sandbox_aggregator/egress_model_inference_gate.go), which uses
// it to keep telemetry/registry/bootstrap/health traffic to a model-plane host from
// minting unnamed model rows.
//
// Callers that need the family, or that must agree with the SQL prefilter, want
// MatchInferenceRoute instead — this predicate is deliberately WIDER than the
// registry the Conversations surface is pinned to.
//
// Note it answers from the PATH only: Ollama's native routes (/api/generate,
// /api/chat, /api/embed) and the embeddings wires (/v1/embeddings, /api/embeddings)
// are absent here on purpose — they are already a first-class part of the projected
// path_class taxonomy (egressclass.IsModelPathClass), which is the group-exact signal
// the models gate checks first, so restating them here would be a second, driftable
// copy of a list that already has an owner.
// Normalization (the `?query` strip and the empty-path guard) happens ONCE, at the
// top, and both arms then match the same query-free string — so the two can never
// drift on what they consider a path (Copilot, PR #2216). MatchInferenceRoute strips
// defensively on its own too; on an already-clean path that is a no-op, and it is
// left untouched here because the Conversations predicate and the pinned Trino
// literal are both bound to its exact current behaviour.
func IsInferenceShapedPath(path string) bool {
	if q := strings.IndexByte(path, '?'); q >= 0 {
		path = path[:q]
	}
	if path == "" {
		return false
	}
	if _, ok := MatchInferenceRoute(path); ok {
		return true
	}
	for _, r := range unadaptedInferenceRoutes {
		if r.re.MatchString(path) {
			return true
		}
	}
	return false
}

// InferenceRouteSuffixRegex renders the registry as ONE case-insensitive,
// end-anchored regex string for the Trino `regexp_like(path, …)` prefilter in
// lakequery/silver_transactions.sql (design §Changes-1). Generated from the same
// table MatchInferenceRoute matches against, so SQL and Go can never disagree on the
// registered route set; the SQL literal is pinned to this exact string by a
// lakequery unit test. Contains no single quotes, so it embeds verbatim inside a
// SQL string literal.
func InferenceRouteSuffixRegex() string {
	fragments := make([]string, 0, len(inferenceRoutes))
	for _, r := range inferenceRoutes {
		fragments = append(fragments, r.pattern)
	}
	return `(?i)(?:` + strings.Join(fragments, "|") + `)$`
}
