// Package aiproviders is the SINGLE SOURCE OF TRUTH for the AI/LLM-provider
// vocabulary and the host→provider matcher, shared across services that need to
// recognize AI/LLM egress: event-ingester-service (the platform "AI client"
// badge + the AI-Sandbox aggregator/reassembler) and config-service (the
// HTTP-capture write-time host-verification warning, SUB-8712). It lives in
// armoapi-go — alongside armotypes/httpcapture, the CaptureConfig/CaptureRule
// contract these services already share — so no service owns the catalog and no
// consumer hand-mirrors it.
//
// Two classifiers historically diverged on both the provider set and the
// matching rules:
//   - the PLATFORM detector (event-ingester-service
//     ingesters/workloads_statuses_ingester) — ~11 providers,
//     DNS-substring/subdomain matching; writes
//     workload_statuses.ai_client_providers.
//   - the AI-SANDBOX aggregator (event-ingester-service
//     scheduled_tasks/ai_sandbox_aggregator) — 6 providers PLUS endpoint kind +
//     AWS region.
//
// This package merges both into one catalog (the superset, keeping the EXACT
// existing display strings — these ARE the ai_client_providers vocab; never
// rename them) and one matcher, so the platform badge and the AI-Sandbox
// rollups can never disagree on provider naming or coverage.
//
//   - ClassifyHost(host) (provider, ok)              — provider only.
//   - ClassifyEndpoint(host) (provider, kind, region, ok) — adds the aggregator's
//     endpoint kind and AWS region (providers without an explicit rule default
//     to KindInference / empty region).
//
// Both functions are pure and unit-tested.
package aiproviders

import (
	"net"
	"regexp"
	"strings"
)

// Canonical provider display strings — the SUPERSET of both legacy classifiers.
// These are the workload_statuses.ai_client_providers vocabulary; do NOT rename.
const (
	ProviderOpenAI       = "OpenAI"
	ProviderAnthropic    = "Anthropic"
	ProviderAWSBedrock   = "AWS Bedrock"
	ProviderAWSSageMaker = "AWS SageMaker"
	ProviderAzureOpenAI  = "Azure OpenAI"
	ProviderGoogle       = "Google"
	ProviderCohere       = "Cohere"
	ProviderHuggingFace  = "Hugging Face"
	ProviderReplicate    = "Replicate"
	ProviderTogetherAI   = "Together AI"
	ProviderPerplexity   = "Perplexity"
	ProviderMistralAI    = "Mistral AI"
	ProviderAI21Labs     = "AI21 Labs"
	ProviderXAI          = "xAI"
	ProviderDeepSeek     = "DeepSeek"
	// ProviderGitHubCopilot is GitHub Copilot's inference data-plane. It has NO
	// OTel gen_ai.system enum member, so providerGenAISystem deliberately omits it:
	// GenAISystemForProvider returns ok=false and the caller leaves gen_ai_system
	// empty, exactly as it does for any other host the OTel enum cannot name. The
	// provider column still populates, which is what Conversations and the Models
	// rollup key on.
	ProviderGitHubCopilot = "GitHub Copilot"
	ProviderGroq          = "Groq"
	ProviderFireworks     = "Fireworks AI"
	ProviderOpenRouter    = "OpenRouter"
	// Managed inference / embedding SaaS that the AI-Sandbox service catalog
	// (scheduled_tasks/ai_sandbox_aggregator/seed/service_catalog.json) already
	// named but the matcher did not recognise. The catalog's `provider` OVERRIDES
	// the host verdict, so a catalog label with no matching constant was a second
	// vocabulary — these constants (plus the subdomainProviders rules below) close
	// that gap in both directions.
	ProviderAnyscale  = "Anyscale"
	ProviderVoyageAI  = "Voyage AI"
	ProviderDeepInfra = "DeepInfra"
	// Self-hosted / in-cluster AI infrastructure (proxy gateways + inference
	// servers the customer runs themselves) — not a public SaaS API.
	ProviderLiteLLM = "LiteLLM"
	ProviderVLLM    = "vLLM"
	ProviderOllama  = "Ollama"
	ProviderKServe  = "KServe"
	ProviderSeldon  = "Seldon"
	// LLMOps / agent-framework tooling — observability, tracing, experiment
	// tracking and orchestration around models rather than the model API itself.
	ProviderLangSmith = "LangSmith"
	ProviderLangfuse  = "Langfuse"
	ProviderLangChain = "LangChain"
	// ProviderLangGraph is LangChain's LangGraph Platform — a self-hosted AGENT
	// RUNTIME (per-agent deployments, queues and Redis under langgraph-dataplane
	// namespaces). A distinct label rather than ProviderLangChain, mirroring how
	// LangSmith (same vendor) has its own: "LangGraph" tells the inventory it is
	// the agent runtime, not the telemetry beacon.
	ProviderLangGraph = "LangGraph"
	ProviderMLflow    = "MLflow"
	ProviderWandB     = "Weights & Biases"
	// ProviderDatadogLLMObs is Datadog's LLM Observability product — the
	// llmobs-intake.* ingestion plane. A workload sending here proves the CUSTOMER
	// instruments LLM calls (a high-value account signal), but the CALLER is always
	// the datadog-agent DaemonSet forwarding other pods' telemetry — so unlike the
	// LLMOps entries above this classifies as KindTelemetrySink, which
	// KindImpliesAIUse deliberately never includes (SUB-8289/SUB-8290: badging the
	// forwarder would misattribute someone else's AI usage).
	ProviderDatadogLLMObs = "Datadog LLM Observability"
	// MODEL VENDORS — labels for who BUILT a model, used only where a model id is
	// the evidence (the ai_sandbox_aggregator model-registry seed). Deliberately
	// NOT matched by ClassifyHost/ClassifyEndpoint and NOT aliased into
	// workloads_statuses_ingester.AILLMProvider: there is no DNS host that
	// resolves to them, so they can never appear in ai_client_providers.
	//
	// A vendor label is only correct when the model is reachable through MORE THAN
	// ONE service — Llama runs on Bedrock, Together, Groq and self-hosted vLLM, so
	// no single service label describes it. A vendor-exclusive model must use its
	// service label instead (Titan/Nova are Bedrock-only → ProviderAWSBedrock), so
	// the registry-fill and the host classifier agree on the same row.
	ProviderMeta = "Meta"
	// Matched by the OpenTelemetry gen_ai.system vocabulary (see
	// genAISystemProviders) AND by a DNS rule for the watsonx.ai runtime
	// (<region>.ml.cloud.ibm.com) in ClassifyEndpoint — so it reaches a row from
	// either an agent-reported semconv value or the host.
	ProviderIBMWatsonx = "IBM watsonx.ai"
)

// genAISystemProviders maps an OpenTelemetry `gen_ai.system` value (renamed
// `gen_ai.provider.name` in current semconv, same values bar xai→x_ai) to this
// catalog's canonical label.
//
// The agent reports this attribute VERBATIM from the instrumented SDK, so its
// values are the spec's lowercase-dotted identifiers — "openai", "aws.bedrock",
// "gcp.gemini". Those are correct as wire values but are NOT display labels: the
// provider-fill path writes into the same serving column the host classifier and
// the model registry write, so passing them through unmapped puts "openai" beside
// "OpenAI" and "aws.bedrock" beside "AWS Bedrock" in one UI column — the exact
// split the catalog exists to prevent.
//
// Several spec values collapse onto one label ON PURPOSE, matching what
// ClassifyEndpoint already returns for the same vendor's hosts, so a row keyed by
// semconv and a row keyed by DNS agree: the three Google endpoints all classify
// as ProviderGoogle, and both Azure faces as ProviderAzureOpenAI.
//
// Source: https://opentelemetry.io/docs/specs/semconv/registry/attributes/gen-ai/
var genAISystemProviders = map[string]string{
	"anthropic":          ProviderAnthropic,
	"aws.bedrock":        ProviderAWSBedrock,
	"azure.ai.inference": ProviderAzureOpenAI,
	"azure.ai.openai":    ProviderAzureOpenAI,
	"cohere":             ProviderCohere,
	"deepseek":           ProviderDeepSeek,
	"gcp.gemini":         ProviderGoogle, // generativelanguage.googleapis.com
	"gcp.gen_ai":         ProviderGoogle, // any Google generative-AI endpoint
	"gcp.vertex_ai":      ProviderGoogle, // aiplatform.googleapis.com
	"groq":               ProviderGroq,
	"ibm.watsonx.ai":     ProviderIBMWatsonx,
	"mistral_ai":         ProviderMistralAI,
	"openai":             ProviderOpenAI,
	"perplexity":         ProviderPerplexity,
	"xai":                ProviderXAI, // gen_ai.system spelling
	"x_ai":               ProviderXAI, // gen_ai.provider.name spelling
	// Derived tokens (NOT official semconv): the AI-Sandbox derive path stamps
	// these for OpenAI-compatible providers the spec assigns no gen_ai.system, so a
	// stamped value maps back to the same label. Keep in sync with providerGenAISystem.
	"together_ai":  ProviderTogetherAI,
	"openrouter":   ProviderOpenRouter,
	"deepinfra":    ProviderDeepInfra,
	"anyscale":     ProviderAnyscale,
	"fireworks_ai": ProviderFireworks,
	"vllm":         ProviderVLLM,
	"ollama":       ProviderOllama,
}

// ProviderFromGenAISystem maps an OpenTelemetry gen_ai.system /
// gen_ai.provider.name value to its canonical catalog label. Matching is
// case-insensitive and trims surrounding whitespace; an empty or unrecognised
// value returns ok=false.
//
// Callers MUST NOT substitute "" on a miss — an unrecognised value is upstream
// drift (a provider the spec added after this map was written), and the raw
// string is more useful to an operator than a blank column. Keep it verbatim and
// add a mapping here.
func ProviderFromGenAISystem(system string) (string, bool) {
	provider, ok := genAISystemProviders[strings.ToLower(strings.TrimSpace(system))]
	return provider, ok
}

// providerGenAISystem is the DERIVE-path reverse of genAISystemProviders: the
// canonical gen_ai.system value the ai_sandbox_ingester stamps for a provider whose
// request/response body it can PARSE — the native adapters (Anthropic, AWS Bedrock,
// OpenAI) plus every provider served over the OpenAI-compatible /chat/completions
// wire (Azure OpenAI, Groq, Mistral, DeepSeek, Perplexity, xAI, Together, OpenRouter,
// DeepInfra, Anyscale, Fireworks, and self-hosted vLLM/Ollama).
//
// Values are the OpenTelemetry gen_ai.system identifiers where the spec defines one
// and a lowercase spec-style token where it does not; every value also appears in
// genAISystemProviders, so a stamped value maps back to the same label.
//
// Google / Gemini maps to the UMBRELLA "gcp.gen_ai" value rather than gcp.gemini vs
// gcp.vertex_ai: the derive path stamps gen_ai.system from the HOST-derived provider
// (ProviderGoogle), which cannot tell the Gemini-API sub-endpoint from Vertex, so it
// emits the umbrella value that maps back to Google either way.
//
// Providers with NO body parser are deliberately ABSENT (Cohere native, SageMaker,
// Hugging Face, Replicate, Voyage, AI21, and the gateway/LLMOps labels): they keep
// gen_ai.system="" so an empty system keeps its meaning — "provider identified, content
// not parsed".
var providerGenAISystem = map[string]string{
	ProviderAnthropic:   "anthropic",
	ProviderAWSBedrock:  "aws.bedrock",
	ProviderOpenAI:      "openai",
	ProviderGoogle:      "gcp.gen_ai",
	ProviderAzureOpenAI: "azure.ai.openai",
	ProviderGroq:        "groq",
	ProviderMistralAI:   "mistral_ai",
	ProviderDeepSeek:    "deepseek",
	ProviderPerplexity:  "perplexity",
	ProviderXAI:         "xai",
	ProviderTogetherAI:  "together_ai",
	ProviderOpenRouter:  "openrouter",
	ProviderDeepInfra:   "deepinfra",
	ProviderAnyscale:    "anyscale",
	ProviderFireworks:   "fireworks_ai",
	ProviderVLLM:        "vllm",
	ProviderOllama:      "ollama",
}

// GenAISystemForProvider returns the gen_ai.system value the derive path stamps for a
// catalog provider whose content the ingester parses. ok=false means the provider has
// no body parser — the caller MUST keep gen_ai.system empty rather than invent a
// value, so an empty system keeps meaning "identified but not parsed".
func GenAISystemForProvider(provider string) (string, bool) {
	system, ok := providerGenAISystem[provider]
	return system, ok
}

// Endpoint kinds — the face of a provider an endpoint exposes. The string values
// mirror armosec-infra/aisandbox.AiSandboxEndpointKind* (kept decoupled to avoid
// an import cycle / external dep in this leaf package; values must stay equal).
const (
	KindInference    = "inference"     // model data-plane (bedrock-runtime, api.openai.com)
	KindControlPlane = "control-plane" // management API (bedrock.<region>)
	KindModelHost    = "model-host"    // model hosting (sagemaker)
	KindGateway      = "gateway"       // self-hosted proxy/inference in the customer's own cluster/domain (litellm/vllm/ollama/kserve/seldon)
	KindLLMOps       = "llmops"        // tooling AROUND models — tracing/observability/experiment tracking (langsmith/langfuse/mlflow/wandb) or agent framework telemetry (langchain)
	// KindTelemetrySink is a vendor INGESTION endpoint that receives LLM telemetry
	// FORWARDED by an aggregation agent (llmobs-intake.*.datadoghq.com ← the
	// datadog-agent DaemonSet). The signal is account-level ("this customer
	// instruments LLM calls"), never caller-level: the calling workload is a
	// forwarder, not an AI client, so this kind is deliberately absent from
	// KindImpliesAIUse's allowlist (see the SUB-8289 forwarder false-positive
	// class) — surfaces can show it; the badge must never derive from it.
	KindTelemetrySink = "telemetry-sink"
)

// KindImpliesAIUse reports whether an endpoint kind is evidence that the calling
// workload actually USES AI — the question the agentic/"AI client" badge asks.
// It is the ONE place that distinction is expressed, so every consumer of
// ClassifyEndpoint gates on the same rule.
//
// KindControlPlane is deliberately EXCLUDED. A control-plane endpoint
// (bedrock.<region>.amazonaws.com, bedrock-agent.<region>.amazonaws.com,
// api.sagemaker.<region>.amazonaws.com) is the MANAGEMENT API: ListFoundationModels,
// GetFoundationModel, ListAgents, DescribeEndpoint. Reaching it means the workload
// ENUMERATED or INSPECTED an AI service — it never means the workload sent a prompt
// or received a completion.
//
// PRODUCTION EVIDENCE (dev, verified in Postgres): a CSPM cloud-scanner workload was
// badged "AI" because its inventory sweep enumerates AWS AI services across 16-18
// regions. Its observed egress was EXCLUSIVELY control-plane hosts
// (bedrock.<region>, bedrock-agent.<region>, api.sagemaker.<region>) and contained
// ZERO data-plane hosts (bedrock-runtime.*, bedrock-agent-runtime.*,
// runtime.sagemaker.*) — the exact traffic shape of a scanner, badged as an AI
// consumer. Every security scanner in every customer's cluster has that shape.
//
// The switch is an explicit ALLOWLIST, not `kind != KindControlPlane`, on purpose: a
// kind added to this package in the future must OPT IN to implying AI use. A
// negative test would silently grant the badge to any new kind — including the next
// management-plane one — which is precisely the failure this function exists to fix.
//
// egressclass.Classify (pkg/egressclass/egressclass.go, the `modelHostPlane` term)
// already gates on kind the same way for the AI-Sandbox model rollup; routing the
// platform badge through this predicate makes the two consumers consistent instead
// of one of them quietly disagreeing.
func KindImpliesAIUse(kind string) bool {
	// KindTelemetrySink is likewise EXCLUDED on purpose: its caller is a telemetry
	// FORWARDER (the datadog-agent DaemonSet), and badging it as an AI client is
	// exactly the SUB-8289 misattribution class.
	switch kind {
	case KindInference, KindModelHost, KindGateway, KindLLMOps:
		return true
	default:
		return false
	}
}

// awsRegionFromHost pulls the AWS region token out of an `*.<region>.amazonaws.com`
// or `*.<region>.api.aws` host (e.g. bedrock-runtime.us-east-1.amazonaws.com or
// bedrock-mantle.us-east-1.api.aws → "us-east-1").
var awsRegionFromHost = regexp.MustCompile(`\.([a-z]{2,3}-[a-z]+-\d+)\.(?:amazonaws\.com|api\.aws)$`)

// onAWSDomain reports whether a host sits on one of AWS's OWN service domains:
// the classic `.amazonaws.com` family or the dual-stack `.api.aws` family
// (bedrock-runtime.<region>.api.aws, bedrock-mantle.<region>.api.aws, …).
// `.aws` is Amazon's brand TLD, so both suffixes are Amazon-operated end to end
// and the suffix check preserves the lookalike guard — a host merely CONTAINING
// the tokens (bedrock-runtime.us-east-1.api.aws.evil.example) matches neither.
func onAWSDomain(h string) bool {
	return strings.HasSuffix(h, ".amazonaws.com") || strings.HasSuffix(h, ".api.aws")
}

// subdomainProvider pairs a base domain with its provider for subdomain-style
// matching (host == base OR host ends with "."+base). This is the platform
// detector's matching style, carried into the shared matcher.
type subdomainProvider struct {
	base     string
	provider string
}

// subdomainProviders are the non-AWS, non-Azure, non-Google providers matched by
// exact-or-subdomain rule. Order is irrelevant — bases are disjoint.
var subdomainProviders = []subdomainProvider{
	{"openai.com", ProviderOpenAI}, // Azure (.openai.azure.com) is handled earlier
	// chatgpt.com is a SEPARATE registrable domain from openai.com, so the rule above
	// does not reach it. It carries real inference traffic: codex 0.149.x moved its
	// model calls to `GET /backend-api/codex/responses` on chatgpt.com (WebSocket +
	// permessage-deflate), where earlier versions used api.openai.com.
	//
	// Unlike the github.com exclusion below, the consumer web UI on this domain is NOT
	// a false positive to guard against: chatgpt.com serves no non-AI product, so a
	// browser session there is genuine AI usage. Contrast github.com, which Copilot also
	// calls for auth and repos and which must stay unclassified.
	{"chatgpt.com", ProviderOpenAI},
	{"anthropic.com", ProviderAnthropic},
	{"cohere.ai", ProviderCohere},
	{"cohere.com", ProviderCohere},
	{"huggingface.co", ProviderHuggingFace},
	{"huggingface.cloud", ProviderHuggingFace}, // Inference Endpoints domain
	{"replicate.com", ProviderReplicate},
	{"together.ai", ProviderTogetherAI},
	{"together.xyz", ProviderTogetherAI},
	{"perplexity.ai", ProviderPerplexity},
	{"mistral.ai", ProviderMistralAI},
	{"ai21.com", ProviderAI21Labs},
	{"x.ai", ProviderXAI},              // xAI Grok API (api.x.ai)
	{"deepseek.com", ProviderDeepSeek}, // DeepSeek API (api.deepseek.com)
	{"groq.com", ProviderGroq},         // Groq API (api.groq.com)
	{"fireworks.ai", ProviderFireworks},
	{"openrouter.ai", ProviderOpenRouter},
	// The base covers api.anyscale.com AND api.endpoints.anyscale.com.
	{"anyscale.com", ProviderAnyscale},
	{"voyageai.com", ProviderVoyageAI},   // Voyage AI embeddings/rerank API
	{"deepinfra.com", ProviderDeepInfra}, // DeepInfra managed inference
}

// componentToken is the classification a self-hosted / LLMOps DNS token resolves to.
type componentToken struct {
	provider string
	kind     string
	// requireSuffixes, when non-empty, restricts the token to hosts ending in ANY of
	// these suffixes. Some tokens are only safe inside a known domain: "bedrock"
	// appears in AWS's own endpoint names, so an unrestricted token would match
	// lookalike domains such as bedrock-runtime.<region>.evil.com. Scoping keeps the
	// token useful for the self-hosted case while preserving the lookalike guard.
	// A slice (not a single suffix) because AWS resource names live on two domain
	// families: .amazonaws.com and the dual-stack .api.aws.
	requireSuffixes []string
}

// componentTokenProviders maps an unambiguous AI DNS token to its canonical provider
// and endpoint kind. Tokens are matched as a whole HYPHEN-DELIMITED COMPONENT of a
// DNS label (host split on "."; each label split on "-"; a component must EQUAL the
// token) — never as a raw substring.
//
// Component matching is what makes these tokens usable at all, because AI infra is
// almost never named as a bare label. It recognises:
//   - bare labels — litellm.<ns>.svc.cluster.local, litellm.<customer>.cloud
//   - compound Kubernetes service/pod names — litellm-proxy.litellm-proxy.svc.cluster.local,
//     bedrock-access-gateway.<ns>.svc.cluster.local, vllm-server.serving.svc.cluster.local
//   - self-hosted cloud resource names — langsmith-prod-<acct>.s3.<region>.amazonaws.com,
//     sagemaker-studio-<acct>-<id>.s3.<region>.amazonaws.com,
//     langsmith-postgres-cp-aurora.cluster-<id>.<region>.rds.amazonaws.com
//
// while still excluding raw substrings: mylitellmthing.com has no "-", so its only
// component is "mylitellmthing", which does not equal "litellm".
//
// This map is checked LAST, after every public-provider rule, so genuine provider
// domains (api.sagemaker.<region>.amazonaws.com, bedrock-runtime.<region>.amazonaws.com,
// api.openai.com, ...) always keep their existing, more specific classification.
//
// DELIBERATELY EXCLUDED — do not add:
//   - "ai" — it matches the .ai ccTLD (ensights.ai, wint.ai, db.ai-scraper.ai) and
//     would flood the badge with false positives.
//   - "triton", "tgi" — too ambiguous (NVIDIA Triton shares its name with unrelated
//     software; "tgi" is a three-letter token).
var componentTokenProviders = map[string]componentToken{
	// --- self-hosted proxy / inference serving ---
	"litellm": {ProviderLiteLLM, KindGateway, nil},
	"vllm":    {ProviderVLLM, KindGateway, nil},
	"ollama":  {ProviderOllama, KindGateway, nil},
	"kserve":  {ProviderKServe, KindGateway, nil},
	"seldon":  {ProviderSeldon, KindGateway, nil},
	// Self-hosted Bedrock proxy (bedrock-access-gateway), SCOPED to in-cluster hosts.
	// "bedrock" also appears in AWS's own endpoint names, so an unrestricted token
	// would match lookalike domains (bedrock-runtime.<region>.evil.com) and ambiguous
	// corporate hosts (a Vault server on a host named "bedrock"). A self-hosted
	// Bedrock proxy lives in-cluster, so requiring .svc.cluster.local keeps the real
	// case and drops both false-positive classes. Real AWS endpoints
	// (bedrock*.amazonaws.com) are matched earlier by their own rules.
	"bedrock": {ProviderAWSBedrock, KindGateway, []string{".svc.cluster.local"}},
	// SageMaker data-plane resources — Studio buckets, artifact stores, and the
	// runtime.sagemaker.<region> invoke plane — SCOPED to AWS domains (classic
	// .amazonaws.com plus dual-stack .api.aws): scoping prevents sagemaker-*.evil.com
	// lookalikes. The api.sagemaker.<region> control-plane rule matches earlier.
	"sagemaker": {ProviderAWSSageMaker, KindModelHost, []string{".amazonaws.com", ".api.aws"}},
	// --- LLMOps / agent-framework tooling ---
	// Unscoped: these appear on vendor domains (cloud.langfuse.com), self-hosted cloud
	// resource names (langsmith-prod-<acct>.s3.<region>.amazonaws.com) AND corporate
	// domains (langsmith.infra.<org>.app), so no single suffix would cover them. The
	// token names are distinctive enough that lookalike risk is negligible.
	"langsmith": {ProviderLangSmith, KindLLMOps, nil},
	"langfuse":  {ProviderLangfuse, KindLLMOps, nil},
	"langchain": {ProviderLangChain, KindLLMOps, nil},
	// LangGraph Platform (SUB-8294) — the token matches EVERY host whose label
	// carries the component, which under a langgraph-dataplane namespace includes
	// the platform's per-agent Redis and queue services, not only agent APIs. That
	// is exactly why the kind is LLMOPS and NOT gateway, on two consumer-visible
	// grounds: (1) gateway is in egressclass's model-plane set, so calls to a
	// per-agent Redis would mint MODEL rows — the agents' real model calls are
	// separately attributed at the LLM gateway they invoke; (2) gateway+unscoped
	// tokens feed the ClassifyWorkloadName AI-server badge, and a fleet of hashed
	// per-agent deployment names would flood it (observed: 142 workloads on one
	// production fleet). LLMOps keeps the caller-side AI badge (KindImpliesAIUse)
	// without either side effect.
	"langgraph": {ProviderLangGraph, KindLLMOps, nil},
	"mlflow":    {ProviderMLflow, KindLLMOps, nil},
	"wandb":     {ProviderWandB, KindLLMOps, nil},
}

// matchComponentToken splits an already-normalised host into DNS labels, then splits
// each label into hyphen-delimited components, and returns the classification of the
// FIRST component that is a known AI token. Host order makes the result deterministic.
func matchComponentToken(host string) (componentToken, bool) {
	for _, label := range strings.Split(host, ".") {
		for _, component := range strings.Split(label, "-") {
			if t, ok := componentTokenProviders[component]; ok {
				if len(t.requireSuffixes) > 0 && !hasAnySuffix(host, t.requireSuffixes) {
					continue
				}
				return t, true
			}
		}
	}
	return componentToken{}, false
}

// hasAnySuffix reports whether host ends with any of the given suffixes.
func hasAnySuffix(host string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(host, s) {
			return true
		}
	}
	return false
}

// normalizeHost lowercases, trims surrounding whitespace, and strips a trailing
// FQDN dot.
func normalizeHost(host string) string {
	h := strings.ToLower(strings.TrimSpace(host))
	return strings.TrimSuffix(h, ".")
}

func matchesSubdomain(host, base string) bool {
	return host == base || strings.HasSuffix(host, "."+base)
}

// ClassifyWorkloadName maps a Kubernetes workload's OWN name to the self-hosted
// AI provider it is named after — the "AI server" half of the badge: a workload
// NAMED vllm-inference-0 IS the AI infrastructure, regardless of its egress.
//
// It reuses the componentTokenProviders catalog with the same matching rule as
// hosts (split on "." then "-"; a component must EQUAL the token — never a raw
// substring, so llama-index-web does not match "ollama"), restricted to the
// UNSCOPED KindGateway tokens:
//   - kind == KindGateway: only self-hosted proxy/inference serving counts as
//     "this workload IS an AI server". LLMOps tokens (langfuse, mlflow, wandb, …)
//     are tooling AROUND models, and public-provider rules are DNS-only.
//   - no requireSuffixes: a suffix-scoped token is only safe inside its domain,
//     and a bare workload name has no domain to check — "bedrock" without its
//     .svc.cluster.local scope would badge a Minecraft Bedrock server.
//
// Pure; ok=false means the name is not recognisably AI infrastructure.
func ClassifyWorkloadName(name string) (provider string, ok bool) {
	n := normalizeHost(name)
	if n == "" {
		return "", false
	}
	for _, label := range strings.Split(n, ".") {
		for _, component := range strings.Split(label, "-") {
			if t, tok := componentTokenProviders[component]; tok && t.kind == KindGateway && len(t.requireSuffixes) == 0 {
				return t.provider, true
			}
		}
	}
	return "", false
}

// ClassifyHost maps a DNS host to its canonical AI-provider name. ok=false means
// the host is not a recognised AI-provider endpoint. Pure.
//
// ⚠️ It DISCARDS the endpoint kind, so `ok` means only "this host belongs to a
// recognised AI provider" — NOT "this workload uses AI". A control-plane host
// (bedrock.<region>, api.sagemaker.<region>) returns ok=true here even though
// reaching it is an inventory/management call, never inference.
//
// Do NOT use this for the agentic / "AI client" badge or for any other
// does-this-workload-use-AI decision. Call ClassifyEndpoint and gate on
// KindImpliesAIUse(kind) instead — see that function for the production
// false-positive this warning exists to prevent.
func ClassifyHost(host string) (provider string, ok bool) {
	p, _, _, ok := ClassifyEndpoint(host)
	return p, ok
}

// ClassifyEndpoint maps a DNS host to its provider, endpoint kind, and AWS region
// (region empty when the host encodes none). Providers without an explicit
// kind/region rule resolve to KindInference and an empty region. ok=false means
// the host is not a recognised AI-provider endpoint.
func ClassifyEndpoint(host string) (provider, kind, region string, ok bool) {
	// A captured L7 host / Host header may carry a :port — an in-cluster gateway
	// on :8080, or api.openai.com:443. Strip it so classification is port-agnostic
	// (a host:port is the same host). IPv6-safe via SplitHostPort; a host with no
	// port errors and is left unchanged.
	if hostOnly, _, err := net.SplitHostPort(host); err == nil {
		host = hostOnly
	}
	h := normalizeHost(host)
	if h == "" {
		return "", "", "", false
	}

	awsRegion := ""
	if m := awsRegionFromHost.FindStringSubmatch(h); m != nil {
		awsRegion = m[1]
	}

	switch {
	// --- GitHub Copilot -----------------------------------------------------
	// Copilot CLI / IDE inference is served from the githubcopilot.com family:
	// api.enterprise.githubcopilot.com (Copilot Enterprise), api.githubcopilot.com
	// (individual/business), api.business.githubcopilot.com, plus the sibling
	// telemetry./exp. hosts on the same domain. Matched on the registrable domain
	// so a new subdomain does not need a catalog change.
	//
	// Deliberately NOT github.com or api.github.com: those are the ordinary GitHub
	// API (auth, repos, PATs), which Copilot also calls, and attributing them as
	// inference would count every `gh` invocation in the cluster as AI usage — the
	// same false-positive class the bedrock control-plane split exists to avoid.
	case h == "githubcopilot.com" || strings.HasSuffix(h, ".githubcopilot.com"):
		return ProviderGitHubCopilot, KindInference, "", true
	// Copilot's LEGACY completion endpoint — older IDE/CLI builds (and CI runners
	// on them) still call it. EXACT host only: githubusercontent.com serves all of
	// GitHub's user content (raw files, avatars, pages assets), so widening to the
	// domain would classify every raw.githubusercontent.com fetch as AI inference —
	// the same false-positive class the github.com exclusion above guards against.
	case h == "copilot-proxy.githubusercontent.com":
		return ProviderGitHubCopilot, KindInference, "", true

	// --- AWS Bedrock (control + data plane, any region) ---------------------
	// Each prefix is accepted on BOTH AWS domain families (.amazonaws.com and the
	// dual-stack .api.aws) — real prod traffic runs on both.
	case strings.HasPrefix(h, "bedrock-runtime.") && onAWSDomain(h):
		return ProviderAWSBedrock, KindInference, awsRegion, true
	case strings.HasPrefix(h, "bedrock-agent-runtime.") && onAWSDomain(h):
		return ProviderAWSBedrock, KindInference, awsRegion, true
	// bedrock-mantle.<region>.api.aws is Bedrock's OpenAI-/Anthropic-COMPATIBLE
	// inference endpoint (OpenAI Responses + Chat Completions, Anthropic Messages —
	// https://docs.aws.amazon.com/bedrock/latest/userguide/bedrock-mantle.html).
	// Pure data plane → KindInference. It ships only on .api.aws today; accepting
	// both families via onAWSDomain is harmless and future-proof.
	case strings.HasPrefix(h, "bedrock-mantle.") && onAWSDomain(h):
		return ProviderAWSBedrock, KindInference, awsRegion, true
	case strings.HasPrefix(h, "bedrock-agent.") && onAWSDomain(h):
		return ProviderAWSBedrock, KindControlPlane, awsRegion, true
	case strings.HasPrefix(h, "bedrock.") && onAWSDomain(h):
		return ProviderAWSBedrock, KindControlPlane, awsRegion, true

	// --- AWS SageMaker ------------------------------------------------------
	// api.sagemaker.<region> is the CONTROL PLANE (CreateEndpoint,
	// CreateTrainingJob, …); the data plane is runtime.sagemaker.<region>
	// (InvokeEndpoint), which falls through to the "sagemaker" component token
	// below and resolves to KindModelHost. This rule used to return KindModelHost
	// for both, which disagreed with the service catalog's (correct) control-plane
	// label for the same host — and with the componentTokenProviders comment that
	// already refers to this as "the control-plane rule".
	case strings.HasPrefix(h, "api.sagemaker.") && onAWSDomain(h):
		return ProviderAWSSageMaker, KindControlPlane, awsRegion, true

	// --- Hugging Face (hub vs inference plane) ------------------------------
	// The hub itself HOSTS model artifacts; the inference API
	// (api-inference.huggingface.co) and Inference Endpoints (huggingface.cloud)
	// are a data plane. The exact-or-subdomain rule below cannot tell them apart —
	// it returns KindInference for every match — so the hub host gets an explicit
	// rule ahead of it. Everything else under huggingface.co still falls through
	// to the subdomain rule and stays KindInference.
	case h == "huggingface.co", h == "www.huggingface.co":
		return ProviderHuggingFace, KindModelHost, "", true

	// --- Azure OpenAI (MUST precede the generic openai.com branch) ----------
	case strings.HasSuffix(h, ".openai.azure.com"):
		return ProviderAzureOpenAI, KindInference, "", true
	case strings.HasSuffix(h, ".cognitiveservices.azure.com"):
		return ProviderAzureOpenAI, KindInference, "", true
	// Azure AI Foundry / Azure AI Services — same Azure AI plane as Azure
	// OpenAI (unified *.services.ai.azure.com endpoint), so it maps to the
	// EXISTING ProviderAzureOpenAI vocab, not a new provider string.
	case strings.HasSuffix(h, ".services.ai.azure.com"):
		return ProviderAzureOpenAI, KindInference, "", true

	// --- Google (Gemini / Vertex AI) ----------------------------------------
	case h == "generativelanguage.googleapis.com",
		h == "aiplatform.googleapis.com",
		strings.HasSuffix(h, "-aiplatform.googleapis.com"):
		return ProviderGoogle, KindInference, "", true

	// --- Datadog LLM Observability (telemetry sink — never badges the caller) --
	// The llmobs-intake ingestion plane, generically across Datadog's regional
	// domains: llmobs-intake.datadoghq.com, llmobs-intake.us5.datadoghq.com,
	// llmobs-intake.datadoghq.eu. Prefix+registrable-suffix keeps the lookalike
	// guard (llmobs-intake.datadoghq.com.evil.example matches neither suffix) and
	// deliberately does NOT cover the rest of datadoghq.com — ordinary Datadog
	// telemetry is not an AI signal.
	case strings.HasPrefix(h, "llmobs-intake.") &&
		(strings.HasSuffix(h, ".datadoghq.com") || strings.HasSuffix(h, ".datadoghq.eu")):
		return ProviderDatadogLLMObs, KindTelemetrySink, "", true

	// --- IBM watsonx.ai -----------------------------------------------------
	// The watsonx.ai foundation-model runtime is <region>.ml.cloud.ibm.com
	// (us-south.ml.cloud.ibm.com, eu-de.ml.cloud.ibm.com, …) — the /ml/v1/text
	// inference plane. gen_ai.system already maps ibm.watsonx.ai; this is the
	// matching DNS rule so a watsonx call classifies from the host too, not only
	// from an agent-reported semconv value.
	case h == "ml.cloud.ibm.com", strings.HasSuffix(h, ".ml.cloud.ibm.com"):
		return ProviderIBMWatsonx, KindInference, "", true
	}

	// --- exact-or-subdomain providers (OpenAI, Anthropic, Cohere, ...) -------
	for _, sp := range subdomainProviders {
		if matchesSubdomain(h, sp.base) {
			return sp.provider, KindInference, "", true
		}
	}

	// --- component-token rules: self-hosted AI infra + LLMOps tooling --------
	// Checked LAST, AFTER every public-provider rule, so real provider domains always
	// win (api.sagemaker.<region>.amazonaws.com stays control/model-plane,
	// bedrock-runtime.<region>.amazonaws.com stays inference). See
	// componentTokenProviders for the matching rule and the excluded tokens.
	//
	// awsRegion is passed through: it is non-empty only when the host genuinely ends
	// in .<region>.amazonaws.com, which is exactly the case for the self-hosted cloud
	// resource names (S3 buckets etc.) this rule is meant to catch.
	if t, cok := matchComponentToken(h); cok {
		return t.provider, t.kind, awsRegion, true
	}

	return "", "", "", false
}
