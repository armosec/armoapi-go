package aiproviders

import "testing"

// TestClassifyHost is the single-source-of-truth host→provider table. It is the
// SUPERSET of both legacy classifiers (the platform ai_llm_detector and the
// AI-Sandbox aggregator), so every host either classifier ever recognised must
// resolve here, to the EXACT canonical provider string.
func TestClassifyHost(t *testing.T) {
	tests := []struct {
		host     string
		provider string
	}{
		// --- OpenAI (detector: subdomain matching) ---
		{"openai.com", ProviderOpenAI},
		{"api.openai.com", ProviderOpenAI},
		{"chat.openai.com", ProviderOpenAI},
		{"v1.api.openai.com", ProviderOpenAI},
		// chatgpt.com is its own registrable domain — the openai.com rule above cannot
		// reach it. codex 0.149.x's model traffic lives here.
		{"chatgpt.com", ProviderOpenAI},
		{"api.chatgpt.com", ProviderOpenAI},
		// --- Anthropic ---
		{"anthropic.com", ProviderAnthropic},
		{"api.anthropic.com", ProviderAnthropic},
		{"console.anthropic.com", ProviderAnthropic},
		// --- AWS Bedrock (regional, control + data plane) ---
		{"bedrock.us-east-1.amazonaws.com", ProviderAWSBedrock},
		{"bedrock-runtime.us-east-1.amazonaws.com", ProviderAWSBedrock},
		{"bedrock-runtime.eu-north-1.amazonaws.com", ProviderAWSBedrock},
		{"bedrock-runtime.ap-southeast-1.amazonaws.com", ProviderAWSBedrock},
		{"bedrock-mantle.us-east-1.api.aws", ProviderAWSBedrock},
		{"bedrock-agent.us-west-2.amazonaws.com", ProviderAWSBedrock},
		{"bedrock-agent-runtime.eu-west-1.amazonaws.com", ProviderAWSBedrock},
		// --- AWS SageMaker (previously aggregator-only — now resolves in BOTH) ---
		{"api.sagemaker.ap-southeast-2.amazonaws.com", ProviderAWSSageMaker},
		// --- Azure OpenAI (must win over the generic openai.com branch) ---
		{"my-resource.openai.azure.com", ProviderAzureOpenAI},
		{"my-resource.cognitiveservices.azure.com", ProviderAzureOpenAI},
		// Azure AI Foundry / Azure AI Services — same Azure AI plane, maps to Azure OpenAI.
		{"mf-x.services.ai.azure.com", ProviderAzureOpenAI},
		// --- Google (previously aggregator-only — now resolves in BOTH) ---
		{"generativelanguage.googleapis.com", ProviderGoogle},
		{"aiplatform.googleapis.com", ProviderGoogle},
		{"us-central1-aiplatform.googleapis.com", ProviderGoogle},
		// --- Cohere (previously detector-only — now resolves in BOTH) ---
		{"cohere.ai", ProviderCohere},
		{"cohere.com", ProviderCohere},
		{"api.cohere.ai", ProviderCohere},
		{"dashboard.cohere.com", ProviderCohere},
		// --- Hugging Face ---
		{"huggingface.co", ProviderHuggingFace},
		{"api-inference.huggingface.co", ProviderHuggingFace},
		{"models.huggingface.co", ProviderHuggingFace},
		{"abc123.endpoints.huggingface.cloud", ProviderHuggingFace},
		// --- Replicate ---
		{"replicate.com", ProviderReplicate},
		{"api.replicate.com", ProviderReplicate},
		// --- Together AI ---
		{"together.ai", ProviderTogetherAI},
		{"together.xyz", ProviderTogetherAI},
		{"api.together.ai", ProviderTogetherAI},
		{"api.together.xyz", ProviderTogetherAI},
		// --- Perplexity ---
		{"perplexity.ai", ProviderPerplexity},
		{"api.perplexity.ai", ProviderPerplexity},
		// --- Mistral AI ---
		{"mistral.ai", ProviderMistralAI},
		{"api.mistral.ai", ProviderMistralAI},
		// --- AI21 Labs ---
		{"ai21.com", ProviderAI21Labs},
		{"api.ai21.com", ProviderAI21Labs},
		// --- xAI ---
		{"x.ai", ProviderXAI},
		{"api.x.ai", ProviderXAI},
		// --- DeepSeek ---
		{"deepseek.com", ProviderDeepSeek},
		{"api.deepseek.com", ProviderDeepSeek},
		// --- Groq ---
		{"groq.com", ProviderGroq},
		{"api.groq.com", ProviderGroq},
		// --- Fireworks AI ---
		{"fireworks.ai", ProviderFireworks},
		{"api.fireworks.ai", ProviderFireworks},
		// --- OpenRouter ---
		{"openrouter.ai", ProviderOpenRouter},
		{"api.openrouter.ai", ProviderOpenRouter},
		// --- Self-hosted / gateway AI infra (hyphen-component match) ---
		{"litellm.ai-platform.svc.cluster.local", ProviderLiteLLM},
		{"10-59-47-73.litellm.ai-platform.svc.cluster.local", ProviderLiteLLM},
		{"litellm.acme.cloud", ProviderLiteLLM},
		{"vllm.serving.svc.cluster.local", ProviderVLLM},
		{"ollama.ml.svc.cluster.local", ProviderOllama},
		// Compound service/pod names — the dominant real-world shape (a whole-label
		// rule missed every one of these).
		{"litellm-proxy.litellm-proxy.svc.cluster.local", ProviderLiteLLM},
		{"10-20-101-28.litellm-proxy.litellm-proxy.svc.cluster.local", ProviderLiteLLM},
		{"vllm-server.serving.svc.cluster.local", ProviderVLLM},
		{"ollama-webui.ml.svc.cluster.local", ProviderOllama},
		{"kserve-controller.kserve.svc.cluster.local", ProviderKServe},
		{"seldon-model.ml.svc.cluster.local", ProviderSeldon},
		// --- Self-hosted Bedrock proxy (NOT the AWS endpoint) ---
		{"bedrock-access-gateway.armo-platform.svc.cluster.local", ProviderAWSBedrock},
		// --- SageMaker data-plane resources (Studio / artifact buckets) ---
		{"sagemaker-studio-170156656001-azzhsnd4k3.s3.us-east-2.amazonaws.com", ProviderAWSSageMaker},
		{"aws-data-bonial-sagemaker.s3.eu-central-1.amazonaws.com", ProviderAWSSageMaker},
		// --- LLMOps / agent-framework tooling (all observed in prod traffic) ---
		{"beacon.langchain.com", ProviderLangChain},
		{"cloud.langfuse.com", ProviderLangfuse},
		{"langsmith-postgres-cp-aurora.cluster-cectmlk7oucz.us-east-2.rds.amazonaws.com", ProviderLangSmith},
		{"langsmith-redis.9fvfrb.ng.0001.use2.cache.amazonaws.com", ProviderLangSmith},
		{"langsmith-prod-097607883991.s3.eu-west-3.amazonaws.com", ProviderLangSmith},
		{"langsmith.infra.zeenea.app", ProviderLangSmith},
		{"vgw-sec-dev-ai-langfuse-1590420409.eu-west-1.elb.amazonaws.com", ProviderLangfuse},
		{"api.mlflow-telemetry.io", ProviderMLflow},
		{"zeenea-infra-mlflow-artifact-repo.s3.eu-west-3.amazonaws.com", ProviderMLflow},
		{"api.wandb.ai", ProviderWandB},
		// LangGraph Platform (SUB-8294) — the namespace label carries the token, so
		// per-agent hashed services AND the platform's own Redis/queue hosts match.
		{"agent-embedding-generator-e-90da2e56ee875afd82a20ef7775b65f6.langgraph-dataplane.svc.cluster.local", ProviderLangGraph},
		{"agent-embedding-generator-e-90da2e56ee875afd82a20ef7775b6-redis.langgraph-dataplane.svc.cluster.local", ProviderLangGraph},
		{"10-40-102-224.some-agent.langgraph-dataplane-dev.svc.cluster.local", ProviderLangGraph},
		// --- normalization: case, whitespace, trailing dot ---
		{"API.OPENAI.COM", ProviderOpenAI},
		{" Api.OpenAI.com ", ProviderOpenAI},
		{"api.anthropic.com.", ProviderAnthropic},
		{"bedrock-runtime.eu-north-1.amazonaws.com.", ProviderAWSBedrock},
		{"  BEDROCK-RUNTIME.US-EAST-1.AMAZONAWS.COM  ", ProviderAWSBedrock},
	}
	for _, tc := range tests {
		t.Run(tc.host, func(t *testing.T) {
			provider, ok := ClassifyHost(tc.host)
			if !ok {
				t.Fatalf("host %q must classify, got ok=false", tc.host)
			}
			if provider != tc.provider {
				t.Fatalf("host %q: provider = %q, want %q", tc.host, provider, tc.provider)
			}
		})
	}
}

// TestClassifyHost_Rejects pins non-AI / malformed hosts to ok=false.
func TestClassifyHost_Rejects(t *testing.T) {
	for _, h := range []string{
		"", "   ",
		"example.com",
		"api.stripe.com",
		"service.default.svc.cluster.local",
		"redis.cache.svc.cluster.local",
		"s3.us-east-1.amazonaws.com",         // AWS but not AI
		"ec2.us-east-1.amazonaws.com",        // AWS but not AI
		"bedrock-runtime.us-east-1.evil.com", // not amazonaws.com
		"notopenai.com",                      // bare domain, not a real openai host
		"fakeopenai.com",
		"openai.com.evil.com",
		"notchatgpt.com",
		"chatgpt.com.evil.com",
		"claude.ai",              // claude.ai is NOT an AI/LLM API endpoint
		"services.ai.azure.com",  // bare apex, not a real Azure AI resource host
		"api.x.ai.evil.com",      // x.ai must be the domain suffix, not a label
		"notgroq.com",            // not a subdomain of groq.com
		"fireworks.ai.phish.com", // fireworks.ai must be the suffix
		// Gateway tokens match a whole hyphen-delimited COMPONENT, never a raw
		// substring — these have no "-" so the whole component must equal the token.
		"mylitellmproxy.com",
		"mylitellmthing.com",
		// The bare "ai" label is deliberately NOT a token. These are company
		// domains on the .ai ccTLD, not AI endpoints — adding "ai" as a token
		// would flood false positives across every such customer.
		"ensights.ai",
		"wint.ai",
		"db.ai-scraper.ai",
		// LLMOps/self-hosted tokens are components too, never raw substrings.
		"mymlflowthing.com",
		"notlangsmith.com",
		// "bedrock" is scoped to in-cluster hosts, so lookalike domains and
		// ambiguous corporate hosts must NOT match.
		"bedrock-access-gateway.evil.com",
		"vault.bedrock.i.example.se", // a Vault server on a host named "bedrock"
		// "sagemaker" is scoped to AWS domains.
		"sagemaker-studio-123.evil.com",
		"notollama.internal", // "notollama" != the "ollama" component
		"myvllmhost.com",     // "myvllmhost" != the "vllm" component
	} {
		if _, ok := ClassifyHost(h); ok {
			t.Errorf("host %q must be rejected (ok=false)", h)
		}
	}
}

// TestClassifyEndpoint pins provider + kind + region for the AI-Sandbox
// aggregator. Providers without an explicit kind/region rule default to
// Inference / "" — that is the widening for the previously detector-only set.
func TestClassifyEndpoint(t *testing.T) {
	tests := []struct {
		host     string
		provider string
		kind     string
		region   string
	}{
		// GitHub Copilot — the whole githubcopilot.com family is inference. The
		// enterprise/business/individual API hosts and the sibling telemetry./exp.
		// hosts all sit on that domain; matching the registrable domain means a new
		// subdomain needs no catalog change.
		// chatgpt.com — codex 0.149.x's model plane. This is the arm the derive path
		// actually calls (fragment_derive.go -> ClassifyEndpoint), so pin kind here and
		// not only the provider: without KindInference the transaction classifies as a
		// non-inference call and never reaches the conversation fold.
		{"chatgpt.com", ProviderOpenAI, KindInference, ""},
		{"api.chatgpt.com", ProviderOpenAI, KindInference, ""},
		{"api.enterprise.githubcopilot.com", ProviderGitHubCopilot, KindInference, ""},
		{"api.githubcopilot.com", ProviderGitHubCopilot, KindInference, ""},
		{"api.business.githubcopilot.com", ProviderGitHubCopilot, KindInference, ""},
		{"telemetry.enterprise.githubcopilot.com", ProviderGitHubCopilot, KindInference, ""},
		{"githubcopilot.com", ProviderGitHubCopilot, KindInference, ""},
		// AWS Bedrock — runtime/agent-runtime inference; bedrock(.|-agent.) control-plane.
		{"bedrock-runtime.us-east-1.amazonaws.com", ProviderAWSBedrock, KindInference, "us-east-1"},
		{"bedrock-agent-runtime.eu-west-1.amazonaws.com", ProviderAWSBedrock, KindInference, "eu-west-1"},
		{"bedrock-agent.us-west-2.amazonaws.com", ProviderAWSBedrock, KindControlPlane, "us-west-2"},
		{"bedrock.us-east-1.amazonaws.com", ProviderAWSBedrock, KindControlPlane, "us-east-1"},
		// AWS SageMaker — api.* is the CONTROL plane (CreateEndpoint,
		// CreateTrainingJob); runtime.* is the data plane that hosts the model.
		{"api.sagemaker.ap-southeast-2.amazonaws.com", ProviderAWSSageMaker, KindControlPlane, "ap-southeast-2"},
		{"runtime.sagemaker.ap-southeast-2.amazonaws.com", ProviderAWSSageMaker, KindModelHost, "ap-southeast-2"},
		// AWS dual-stack .api.aws family (SUB-8288) — same prefixes, second AWS
		// domain; region attribution must survive. bedrock-mantle is Bedrock's
		// OpenAI-/Anthropic-compatible inference endpoint (data plane).
		{"bedrock-runtime.us-east-1.api.aws", ProviderAWSBedrock, KindInference, "us-east-1"},
		{"bedrock-mantle.us-east-1.api.aws", ProviderAWSBedrock, KindInference, "us-east-1"},
		{"bedrock-agent-runtime.eu-west-1.api.aws", ProviderAWSBedrock, KindInference, "eu-west-1"},
		{"bedrock.us-east-1.api.aws", ProviderAWSBedrock, KindControlPlane, "us-east-1"},
		{"api.sagemaker.eu-central-1.api.aws", ProviderAWSSageMaker, KindControlPlane, "eu-central-1"},
		{"runtime.sagemaker.eu-central-1.api.aws", ProviderAWSSageMaker, KindModelHost, "eu-central-1"},
		// Copilot's legacy completion proxy — exact host only (SUB-8288).
		{"copilot-proxy.githubusercontent.com", ProviderGitHubCopilot, KindInference, ""},
		// OpenAI.
		{"api.openai.com", ProviderOpenAI, KindInference, ""},
		{"foo.openai.com", ProviderOpenAI, KindInference, ""},
		// Azure OpenAI.
		{"my-resource.openai.azure.com", ProviderAzureOpenAI, KindInference, ""},
		{"my-resource.cognitiveservices.azure.com", ProviderAzureOpenAI, KindInference, ""},
		{"mf-x.services.ai.azure.com", ProviderAzureOpenAI, KindInference, ""},
		// Anthropic.
		{"api.anthropic.com", ProviderAnthropic, KindInference, ""},
		// Google.
		{"generativelanguage.googleapis.com", ProviderGoogle, KindInference, ""},
		{"aiplatform.googleapis.com", ProviderGoogle, KindInference, ""},
		{"us-central1-aiplatform.googleapis.com", ProviderGoogle, KindInference, ""},
		// IBM watsonx.ai — regional ML runtime (SUB-8080).
		{"us-south.ml.cloud.ibm.com", ProviderIBMWatsonx, KindInference, ""},
		{"eu-de.ml.cloud.ibm.com", ProviderIBMWatsonx, KindInference, ""},
		// LangGraph Platform (SUB-8294) — LLMOPS, deliberately NOT gateway: gateway
		// would put per-agent Redis/queue calls on the model plane and flood the
		// own-name AI-server badge with hashed deployment names.
		{"agent-x-4463b76f.langgraph-dataplane.svc.cluster.local", ProviderLangGraph, KindLLMOps, ""},
		// Datadog LLM Observability ingestion plane (SUB-8290) — telemetry-sink:
		// classifies (account-level signal) but never implies AI use for the caller.
		{"llmobs-intake.datadoghq.com", ProviderDatadogLLMObs, KindTelemetrySink, ""},
		{"llmobs-intake.us5.datadoghq.com", ProviderDatadogLLMObs, KindTelemetrySink, ""},
		{"llmobs-intake.datadoghq.eu", ProviderDatadogLLMObs, KindTelemetrySink, ""},
		// Widened (previously detector-only) — default to Inference, no region.
		{"api.cohere.ai", ProviderCohere, KindInference, ""},
		{"api-inference.huggingface.co", ProviderHuggingFace, KindInference, ""},
		{"api.replicate.com", ProviderReplicate, KindInference, ""},
		{"api.together.ai", ProviderTogetherAI, KindInference, ""},
		{"api.perplexity.ai", ProviderPerplexity, KindInference, ""},
		{"api.mistral.ai", ProviderMistralAI, KindInference, ""},
		{"api.ai21.com", ProviderAI21Labs, KindInference, ""},
		{"api.x.ai", ProviderXAI, KindInference, ""},
		{"api.deepseek.com", ProviderDeepSeek, KindInference, ""},
		{"api.groq.com", ProviderGroq, KindInference, ""},
		{"api.fireworks.ai", ProviderFireworks, KindInference, ""},
		{"openrouter.ai", ProviderOpenRouter, KindInference, ""},
		// Self-hosted / gateway AI infra — KindGateway, no region.
		{"litellm.ai-platform.svc.cluster.local", ProviderLiteLLM, KindGateway, ""},
		{"10-59-47-73.litellm.ai-platform.svc.cluster.local", ProviderLiteLLM, KindGateway, ""},
		{"litellm.acme.cloud", ProviderLiteLLM, KindGateway, ""},
		{"vllm.serving.svc.cluster.local", ProviderVLLM, KindGateway, ""},
		{"ollama.ml.svc.cluster.local", ProviderOllama, KindGateway, ""},
		// Compound (hyphenated) service/pod names still resolve to KindGateway.
		{"litellm-proxy.litellm-proxy.svc.cluster.local", ProviderLiteLLM, KindGateway, ""},
		{"10-20-101-28.litellm-proxy.litellm-proxy.svc.cluster.local", ProviderLiteLLM, KindGateway, ""},
		{"vllm-server.serving.svc.cluster.local", ProviderVLLM, KindGateway, ""},
		{"ollama-webui.ml.svc.cluster.local", ProviderOllama, KindGateway, ""},
		// normalization.
		{"  BEDROCK-RUNTIME.US-EAST-1.AMAZONAWS.COM  ", ProviderAWSBedrock, KindInference, "us-east-1"},
	}
	for _, tc := range tests {
		t.Run(tc.host, func(t *testing.T) {
			provider, kind, region, ok := ClassifyEndpoint(tc.host)
			if !ok {
				t.Fatalf("host %q must classify, got ok=false", tc.host)
			}
			if provider != tc.provider {
				t.Fatalf("host %q: provider = %q, want %q", tc.host, provider, tc.provider)
			}
			if kind != tc.kind {
				t.Fatalf("host %q: kind = %q, want %q", tc.host, kind, tc.kind)
			}
			if region != tc.region {
				t.Fatalf("host %q: region = %q, want %q", tc.host, region, tc.region)
			}
		})
	}
}

func TestClassifyEndpoint_Rejects(t *testing.T) {
	for _, h := range []string{
		"", "   ", "example.com", "s3.us-east-1.amazonaws.com", "claude.ai",
		// Plain GitHub is NOT Copilot inference. Copilot CLI calls api.github.com for
		// auth/PAT validation, and attributing that as AI usage would count every `gh`
		// invocation in the cluster — the false-positive class the bedrock
		// control-plane split exists to prevent. Also guards against a naive
		// "contains github" match.
		"github.com", "api.github.com", "raw.githubusercontent.com",
		"notgithubcopilot.com.evil.test",
		// copilot-proxy is an EXACT-host rule: other githubusercontent hosts and
		// suffix lookalikes stay unclassified.
		"avatars.githubusercontent.com",
		"copilot-proxy.githubusercontent.com.evil.test",
		// Datadog: only the llmobs-intake plane is an AI signal — ordinary Datadog
		// telemetry hosts and suffix lookalikes stay unclassified (SUB-8290).
		"agent-intake.datadoghq.com",
		"datadoghq.com",
		"llmobs-intake.datadoghq.com.evil.example",
		// .api.aws must be the host's SUFFIX — embedding it mid-host is a lookalike.
		"bedrock-runtime.us-east-1.api.aws.evil.example",
	} {
		if _, _, _, ok := ClassifyEndpoint(h); ok {
			t.Errorf("host %q must be rejected (ok=false)", h)
		}
	}
}

// Every well-known OpenTelemetry gen_ai.system value maps to a canonical label.
// The values are the spec's, not ours — see the URL on genAISystemProviders.
func TestProviderFromGenAISystem(t *testing.T) {
	cases := map[string]string{
		"anthropic":          ProviderAnthropic,
		"aws.bedrock":        ProviderAWSBedrock,
		"azure.ai.inference": ProviderAzureOpenAI,
		"azure.ai.openai":    ProviderAzureOpenAI,
		"cohere":             ProviderCohere,
		"deepseek":           ProviderDeepSeek,
		"gcp.gemini":         ProviderGoogle,
		"gcp.gen_ai":         ProviderGoogle,
		"gcp.vertex_ai":      ProviderGoogle,
		"groq":               ProviderGroq,
		"ibm.watsonx.ai":     ProviderIBMWatsonx,
		"mistral_ai":         ProviderMistralAI,
		"openai":             ProviderOpenAI,
		"perplexity":         ProviderPerplexity,
		"xai":                ProviderXAI,
		"x_ai":               ProviderXAI,
		// Case and surrounding whitespace must not defeat the lookup.
		"  AWS.Bedrock  ": ProviderAWSBedrock,
		"OpenAI":          ProviderOpenAI,
	}
	for system, want := range cases {
		got, ok := ProviderFromGenAISystem(system)
		if !ok {
			t.Errorf("gen_ai.system %q: ok = false, want a canonical label", system)
			continue
		}
		if got != want {
			t.Errorf("gen_ai.system %q = %q, want %q", system, got, want)
		}
	}
}

// A semconv value and the SAME vendor's DNS host must produce the identical
// label — the two fill paths write one serving column, so any disagreement is a
// vendor rendered two ways in the UI. This is the invariant, not the table above.
func TestGenAISystemAgreesWithHostClassifier(t *testing.T) {
	cases := []struct{ system, host string }{
		{"openai", "api.openai.com"},
		{"anthropic", "api.anthropic.com"},
		{"aws.bedrock", "bedrock-runtime.us-east-1.amazonaws.com"},
		{"gcp.gemini", "generativelanguage.googleapis.com"},
		{"gcp.vertex_ai", "us-central1-aiplatform.googleapis.com"},
		{"azure.ai.openai", "my-resource.openai.azure.com"},
		{"cohere", "api.cohere.ai"},
		{"deepseek", "api.deepseek.com"},
		{"groq", "api.groq.com"},
		{"mistral_ai", "api.mistral.ai"},
		{"perplexity", "api.perplexity.ai"},
		{"xai", "api.x.ai"},
		{"ibm.watsonx.ai", "us-south.ml.cloud.ibm.com"},
	}
	for _, tc := range cases {
		fromSystem, ok := ProviderFromGenAISystem(tc.system)
		if !ok {
			t.Errorf("gen_ai.system %q must map", tc.system)
			continue
		}
		fromHost, ok := ClassifyHost(tc.host)
		if !ok {
			t.Errorf("host %q must classify", tc.host)
			continue
		}
		if fromSystem != fromHost {
			t.Errorf("split vocabulary: gen_ai.system %q = %q but host %q = %q",
				tc.system, fromSystem, tc.host, fromHost)
		}
	}
}

// An unrecognised or empty value reports ok=false so the caller can keep the raw
// string. Blanking it would trade a debuggable value for an empty UI column.
func TestProviderFromGenAISystem_Unknown(t *testing.T) {
	for _, s := range []string{"", "   ", "some.new.provider", "not-a-system"} {
		if provider, ok := ProviderFromGenAISystem(s); ok {
			t.Errorf("value %q must not map (got %q)", s, provider)
		}
	}
}

// GenAISystemForProvider stamps a gen_ai.system for every provider whose body the
// ingester parses (native adapters plus the OpenAI-compatible wire) and each stamped
// value must round-trip back to the same label — a split here is a vendor rendered
// two ways in the serving column.
func TestGenAISystemForProvider(t *testing.T) {
	parseable := map[string]string{
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
	for provider, wantSystem := range parseable {
		system, ok := GenAISystemForProvider(provider)
		if !ok {
			t.Errorf("provider %q: ok = false, want gen_ai.system %q", provider, wantSystem)
			continue
		}
		if system != wantSystem {
			t.Errorf("provider %q = %q, want %q", provider, system, wantSystem)
		}
		if back, ok := ProviderFromGenAISystem(system); !ok || back != provider {
			t.Errorf("round-trip: provider %q -> system %q -> %q (ok=%v)", provider, system, back, ok)
		}
	}
}

// A provider the ingester cannot parse must NOT be stamped: an empty gen_ai.system is
// the signal "identified but content not parsed". Guards a future edit from stamping a
// system for a provider with no body adapter. Google moved OUT of this set when its
// :generateContent adapter landed (SUB-8081) — it is now in the parseable map above.
func TestGenAISystemForProvider_Unparseable(t *testing.T) {
	for _, provider := range []string{
		ProviderCohere, ProviderAWSSageMaker, ProviderHuggingFace,
		ProviderReplicate, ProviderVoyageAI, ProviderAI21Labs, ProviderLiteLLM,
		ProviderIBMWatsonx, "", "Not A Provider",
	} {
		if system, ok := GenAISystemForProvider(provider); ok {
			t.Errorf("provider %q must not be stamped (got %q)", provider, system)
		}
	}
}

// TestClassifyWorkloadName pins the own-name "AI server" matcher: a workload
// NAMED like self-hosted AI serving infra classifies to that provider, with the
// same whole-component rule as hosts (never raw substrings).
func TestClassifyWorkloadName(t *testing.T) {
	tests := []struct {
		name     string
		provider string
	}{
		// bare token names
		{"litellm", ProviderLiteLLM},
		{"vllm", ProviderVLLM},
		{"ollama", ProviderOllama},
		{"kserve", ProviderKServe},
		{"seldon", ProviderSeldon},
		// compound Kubernetes workload names — token as a "-"-split component
		{"vllm-inference-0", ProviderVLLM},
		{"litellm-proxy", ProviderLiteLLM},
		{"ollama-server", ProviderOllama},
		{"my-vllm-server", ProviderVLLM},
		{"kserve-controller-manager", ProviderKServe},
		{"seldon-controller-manager", ProviderSeldon},
		// normalisation: case + surrounding whitespace
		{"  VLLM-Inference  ", ProviderVLLM},
		// dotted names (DNS-1123 subdomain names are legal for some kinds)
		{"vllm.serving", ProviderVLLM},
	}
	for _, tc := range tests {
		provider, ok := ClassifyWorkloadName(tc.name)
		if !ok {
			t.Errorf("workload name %q must classify (want %q)", tc.name, tc.provider)
			continue
		}
		if provider != tc.provider {
			t.Errorf("workload name %q → %q, want %q", tc.name, provider, tc.provider)
		}
	}
}

// TestClassifyWorkloadName_Rejects pins ok=false for non-AI names, raw-substring
// traps, and tokens that are host-scoped or non-serving in the catalog.
func TestClassifyWorkloadName_Rejects(t *testing.T) {
	for _, n := range []string{
		"", "   ",
		"nginx", "frontend", "payments-api",
		// whole-component rule: substrings never match
		"llama-index-web", // "llama" != "ollama"
		"myvllmhost",      // no "-", single component "myvllmhost" != "vllm"
		"mylitellmproxy",  // same trap for litellm
		"ollamanager",     // "ollamanager" != "ollama"
		"vllmish-worker",  // "vllmish" != "vllm"
		// suffix-scoped catalog tokens have no domain to check in a bare name —
		// they must NOT match (a "bedrock-server" is likely Minecraft, and
		// "sagemaker" resource names are only meaningful under amazonaws.com).
		"bedrock-server",
		"bedrock-access-gateway",
		"sagemaker-studio",
		// LLMOps tokens are tooling AROUND models, not AI servers.
		"mlflow", "langfuse-web", "langsmith", "wandb-local", "langchain-app",
		// LangGraph is llmops, not gateway — hashed per-agent deployment names must
		// never flood the AI-server badge (SUB-8294).
		"langgraph-api", "agent-embedding-generator-e-90da2e56ee875afd82a20ef7775b65f6",
	} {
		if provider, ok := ClassifyWorkloadName(n); ok {
			t.Errorf("workload name %q must be rejected (got %q)", n, provider)
		}
	}
}

// TestKindImpliesAIUse pins the allowlist: only kinds that represent a workload
// TALKING TO a model (or to model-adjacent tooling) imply AI use. The
// control-plane kind must not — reaching bedrock.<region> / api.sagemaker.<region>
// is enumerating an AI service, not using one (the CSPM cloud-scanner false
// positive this predicate exists to fix).
func TestKindImpliesAIUse(t *testing.T) {
	for _, kind := range []string{KindInference, KindModelHost, KindGateway, KindLLMOps} {
		if !KindImpliesAIUse(kind) {
			t.Errorf("kind %q must imply AI use", kind)
		}
	}
	for _, kind := range []string{
		KindControlPlane,
		KindTelemetrySink, // caller is a telemetry FORWARDER (SUB-8290) — never an AI client
		"",                // an unclassified endpoint is not evidence of anything
		"future-kind",     // a NEW kind must OPT IN — the switch is an allowlist,
		"control-plane",   // and this literal guards the KindControlPlane value itself
	} {
		if KindImpliesAIUse(kind) {
			t.Errorf("kind %q must NOT imply AI use", kind)
		}
	}
}

// TestControlPlaneHostsClassifyButDoNotImplyAIUse is the end-to-end statement of
// the fix at the catalog layer: the control-plane hosts a scanner sweeps still
// CLASSIFY (they really are AWS Bedrock / SageMaker endpoints, and the AI-Sandbox
// service catalog wants them named), but their kind must not imply AI use — while
// the corresponding data-plane hosts must.
func TestControlPlaneHostsClassifyButDoNotImplyAIUse(t *testing.T) {
	// Several regions each: the scanner's traffic is the same host family repeated
	// across 16-18 regions, so a per-region rule regression must fail here.
	controlPlane := []string{
		"bedrock.us-east-1.amazonaws.com",
		"bedrock.eu-west-1.amazonaws.com",
		"bedrock.ap-southeast-2.amazonaws.com",
		"bedrock-agent.us-east-1.amazonaws.com",
		"bedrock-agent.eu-central-1.amazonaws.com",
		"api.sagemaker.us-east-1.amazonaws.com",
		"api.sagemaker.ap-southeast-2.amazonaws.com",
	}
	for _, host := range controlPlane {
		provider, kind, _, ok := ClassifyEndpoint(host)
		if !ok || provider == "" {
			t.Errorf("host %q must still classify to a provider (the catalog names it)", host)
		}
		if kind != KindControlPlane {
			t.Errorf("host %q: kind = %q, want %q", host, kind, KindControlPlane)
		}
		if KindImpliesAIUse(kind) {
			t.Errorf("host %q must NOT imply AI use (control plane = enumeration, not inference)", host)
		}
	}

	dataPlane := []string{
		"bedrock-runtime.us-east-1.amazonaws.com",
		"bedrock-runtime.eu-west-1.amazonaws.com",
		"bedrock-agent-runtime.us-east-1.amazonaws.com",
		"bedrock-agent-runtime.ap-southeast-2.amazonaws.com",
		"runtime.sagemaker.us-east-1.amazonaws.com",
		"runtime.sagemaker.ap-southeast-2.amazonaws.com",
	}
	for _, host := range dataPlane {
		_, kind, _, ok := ClassifyEndpoint(host)
		if !ok {
			t.Errorf("host %q must classify", host)
		}
		if !KindImpliesAIUse(kind) {
			t.Errorf("host %q (kind %q) MUST imply AI use — it is the data plane", host, kind)
		}
	}
}

// TestClassifyHost_PortStripped: a :port on the host must not defeat classification
// (regression — a captured Host header may carry an explicit port: an in-cluster gateway
// on :8080, or even the standard port spelled out as api.openai.com:443). host:port must
// classify == host.
func TestClassifyHost_PortStripped(t *testing.T) {
	cases := []struct {
		host, wantProvider string
	}{
		{"bedrock-access-gateway.ai-sandbox-sim.svc.cluster.local:8080", ProviderAWSBedrock},
		{"api.openai.com:443", ProviderOpenAI},
		{"bedrock-runtime.us-east-1.amazonaws.com:443", ProviderAWSBedrock},
	}
	for _, tc := range cases {
		got, ok := ClassifyHost(tc.host)
		if !ok || got != tc.wantProvider {
			t.Errorf("ClassifyHost(%q) = %q,%v; want %q,true", tc.host, got, ok, tc.wantProvider)
		}
	}
	// explicit host==host:port parity
	for _, h := range []string{"bedrock-access-gateway.x.svc.cluster.local", "api.openai.com"} {
		p1, ok1 := ClassifyHost(h)
		p2, ok2 := ClassifyHost(h + ":8080")
		if p1 != p2 || ok1 != ok2 {
			t.Errorf("port changed classification for %q: %q,%v vs %q,%v", h, p1, ok1, p2, ok2)
		}
	}
}
