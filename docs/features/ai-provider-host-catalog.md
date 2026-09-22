---
type: feature
status: done
owner: eranm
scope: repo
---

# AI/LLM provider host catalog (`armotypes/aiproviders`)

## Purpose

`armotypes/aiproviders` is the **single source of truth** for the AI/LLM-provider
vocabulary and the host→provider matcher, shared by every service that recognizes AI/LLM
egress. It lives here — next to `armotypes/httpcapture` (the `CaptureConfig`/`CaptureRule`
contract those same services share) — so no service owns the catalog and no consumer
hand-mirrors it.

It was moved into armoapi-go from event-ingester-service's `pkg/aiproviders` (SUB-8712):
event-ingester-service still owns the *behavior* (badge derivation, the AI-Sandbox
reassembler's provider/model tagging) but now imports the catalog from here, and
config-service imports the same catalog for its HTTP-capture write-time host-verification
warning — so the two can never drift on which hosts are "recognized".

## Core matcher (moved verbatim)

- `ClassifyEndpoint(host) (provider, kind, region, ok)` — the provider PLUS endpoint kind
  (`inference` / `control-plane` / `model-host` / `gateway` / `llmops`) and AWS region.
  **The one to call.**
- `KindImpliesAIUse(kind) bool` — does this endpoint kind mean the caller actually USES
  AI? (Excludes `control-plane` management APIs.)
- `ClassifyHost(host) (provider, ok)` — provider only; ⚠️ discards the kind.
- `MatchInferenceRoute(path)` / `InferenceRouteSuffixRegex()` — the path-level inference
  registry (`inference_routes.go`).
- `ClassifyWorkloadName(name)` — classifies a workload's own name (self-hosted "AI server"
  badge).

The canonical provider display strings ARE the `ai_client_providers` vocabulary — never
rename them.

## Capture-recognition helpers (`capture_recognition.go`)

Added for SUB-8712 so config-service's HTTP-capture warning derives its "recognized" set
from the exact catalog the reassembler tags from:

- `IsCoarseGoogleGenAIHost(host) bool` — a **deliberately narrow** widening for Google:
  a curated set of Gen-AI naming tokens (`generativelanguage`, `generativeai`,
  `aiplatform`, `vertexai`, `cloudcode`) under `googleapis.com`, **not** a blanket
  `*.googleapis.com` match (that domain also serves Storage, Compute, BigQuery, …). It
  exists because `ClassifyEndpoint`'s Google rule is an exact whitelist, unlike the
  prefix/suffix rules for Bedrock/Azure/OpenAI — the gap the SUB-8712 motivating host
  (`daily-cloudcode-pa.googleapis.com`, Google's internal Code-Assist/Antigravity wire)
  fell into. The AI-Sandbox reassembler's host-only fallback (SUB-8712 item 4) consults
  this too, so tagging and warning agree.
- `IsRecognizedInferenceHost(host) bool` — `ClassifyEndpoint` ok && `KindImpliesAIUse`;
  the confident "this is really an AI inference host" anchor.
- `RecognizedForCapture(host) bool` — `IsRecognizedInferenceHost || IsCoarseGoogleGenAIHost`;
  the exact set the derivation code will attach a provider/gen_ai.system for.
- `HostDomainFamily(host) string` — last-two-labels registrable-domain key
  (`googleapis.com`, `amazonaws.com`, …) for grouping hosts by provider family. A
  documented simplification of the public-suffix list, correct for every catalog domain.

All are port-agnostic (strip an optional `:port`, mirroring `ClassifyEndpoint`) and pure.

## Consumers

- **event-ingester-service** — platform "AI client" badge, AI-Sandbox aggregator, and the
  reassembler's per-transaction provider/model derivation. Imports this package (was its
  own `pkg/aiproviders`).
- **config-service** — the `v1_http_capture_config` write-time warning: flags a
  full-capture rule for a host not `RecognizedForCapture` that shares a `HostDomainFamily`
  with a recognized AI host in the same config (SUB-8712 item 5).
