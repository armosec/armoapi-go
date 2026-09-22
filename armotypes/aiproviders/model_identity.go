package aiproviders

import (
	"net/url"
	"regexp"
	"strings"
)

// SUB-8068 — canonical model IDENTITY normalization.
//
// The Models surface (and the Conversations model facet/filter, which shares the
// same silver-derived model string) must show ONE row per logical model, keyed on
// a canonical identity, with the snapshot date and the region carried as ATTRIBUTES
// — the shape the UI spec pins (visibility-ui §331: `anthropic.claude-sonnet-4-6` |
// `AWS Bedrock, eu-west-1` | …; §455 detail panel = provider/path, region, versions
// seen). Today the stored model_id is the RAW wire string, so one model forks into
// several rows (`claude-haiku-4-5`, `claude-haiku-4-5-20251001`,
// `us.anthropic.claude-haiku-4-5-20251001-v1:0`).
//
// This normalizer is DETERMINISTIC and REGISTRY-INDEPENDENT on purpose:
//   - the model registry LAGS new models (there is a modelRegistryMisses drift
//     signal for exactly that), so identity must not require a registry hit;
//   - the registry `family` field OVER-MERGES (claude-opus-4 and claude-opus-4-1
//     both map to `claude-4-opus`), which would fuse two distinct models.
//
// So identity is derived from the WIRE FORM alone: keep the vendor + the marketing
// minor version, strip the region prefix, the Bedrock `-vN:M` tail and the
// `-YYYYMMDD` snapshot date (both lifted into Version). A NEW model in a known
// family (`claude-opus-99`) still normalizes correctly via the family→vendor
// prefix heuristic; a genuinely unrecognized vendor keeps a bare identity and is
// NEVER guessed (SUB-7677). The registry, when it does have the model, only
// enriches the display label / provider elsewhere — it is never required here.
var (
	// modelRegionPrefix matches a leading Bedrock cross-region inference-profile
	// segment ("us." in "us.anthropic.claude-…"). us-gov before us so the longer
	// match wins. Mirrors the aggregator's bedrockRegionPrefix.
	modelRegionPrefix = regexp.MustCompile(`^(us-gov|us|eu|apac)\.`)
	// modelVendorPrefix matches a leading Bedrock model-vendor namespace. Superset
	// of the aggregator's bedrockProviderPrefix (adds google/openai for the rare
	// wire forms that carry them). The captured token IS the vendor.
	modelVendorPrefix = regexp.MustCompile(`^(anthropic|amazon|meta|cohere|mistral|ai21|stability|deepseek|google|openai)\.`)
	// modelVersionSuffix matches the Bedrock version tail (`-v1:0`, `:0`). Mirrors
	// the aggregator's bedrockVersionSuffix.
	modelVersionSuffix = regexp.MustCompile(`(-v\d+)?:\d+$|-v\d+$`)
	// modelDateSuffix matches a trailing `-YYYYMMDD` snapshot date. Mirrors the
	// aggregator's modelDateSuffix.
	modelDateSuffix = regexp.MustCompile(`-\d{8}$`)
)

// vendorByFamilyPrefix derives the model VENDOR from a bare family name (a
// public-API alias like `claude-haiku-4-5` carries no wire vendor prefix). Keyed
// on the family ROOT so a NEW model in a known family (`claude-opus-99`,
// `gpt-6-mini`) still resolves. Empty for an unrecognized root — never guessed.
//
// OpenAI's reasoning family (`o1`, `o3-mini`, `o4-mini`, …) is deliberately NOT
// enumerated here — it is matched by openAIOSeriesFamily below, so a future
// generation (`o5`, `o6`) resolves without a code change. Enumerating it (the
// pre-2026-08 `{"o1","openai"},{"o3","openai"},{"o4","openai"}` entries) meant an
// unlisted `o5-mini` fell to vendor="" and FORKED from its `openai.o5-mini` wire
// form on the Models surface.
var vendorByFamilyPrefix = []struct{ prefix, vendor string }{
	{"claude", "anthropic"},
	{"chatgpt", "openai"}, {"gpt", "openai"},
	{"gemini", "google"},
	{"nova", "amazon"}, {"titan", "amazon"},
	{"llama", "meta"},
	{"codestral", "mistral"}, {"mixtral", "mistral"}, {"mistral", "mistral"},
	{"command", "cohere"},
	{"deepseek", "deepseek"},
	{"jamba", "ai21"},
}

// openAIOSeriesFamily matches OpenAI's reasoning-model family root — an `o`
// immediately followed by a version digit (`o1`, `o1-preview`, `o3-mini`,
// `o4-mini`, and any future `o5`/`o6`…). A PATTERN, not an enumerated list, so a
// new generation resolves to `openai` on arrival. `^o\d` is the faithful
// generalization of the old `o1`/`o3`/`o4` HasPrefix entries; it cannot match a
// non-OpenAI family that merely starts with "o" (e.g. `olmo-7b` → `o` then `l`).
var openAIOSeriesFamily = regexp.MustCompile(`^o\d`)

// ModelIdentity is the canonical, new-model-safe decomposition of a raw wire model
// id into the identity plus its attributes.
type ModelIdentity struct {
	// Model is the canonical identity/grouping key: `<vendor>.<family-minor>`,
	// dateless (e.g. `anthropic.claude-haiku-4-5`). Bare `<family-minor>` when the
	// vendor is unrecognized. "" only for an empty input.
	Model string
	// Vendor is the model vendor ("anthropic", "openai", …); "" when unrecognized.
	Vendor string
	// Version is the snapshot date (`20251001`), or the Bedrock version tail
	// (`v1:0`) when no date is present; "" for a floating alias.
	Version string
	// Region is the AWS region from the host (`us-east-1`) when available, else the
	// coarse inference-profile geo (`us`); "" when neither is present.
	Region string
}

// NormalizeModelIdentity decomposes a raw wire model id (optionally aided by the
// captured host, which carries the authoritative AWS region) into a canonical
// identity + attributes. Pure and deterministic; see the package doc above.
func NormalizeModelIdentity(rawModel, host string) ModelIdentity {
	id := ModelIdentity{}

	// Region from the host first — the authoritative AWS region
	// (bedrock-runtime.us-east-1.amazonaws.com → "us-east-1"), the §331 form.
	if host != "" {
		if _, _, region, ok := ClassifyEndpoint(host); ok && region != "" {
			id.Region = region
		}
	}

	m := strings.ToLower(strings.TrimSpace(rawModel))
	if m == "" {
		return id
	}
	if decoded, err := url.PathUnescape(m); err == nil {
		m = decoded
	}

	// Region prefix: the inference-profile geo ("us.") is a coarse fallback used
	// only when the host gave us nothing; either way it is stripped from identity.
	if loc := modelRegionPrefix.FindString(m); loc != "" {
		if id.Region == "" {
			id.Region = strings.TrimSuffix(loc, ".")
		}
		m = modelRegionPrefix.ReplaceAllString(m, "")
	}

	// Bedrock `-vN:M` tail — captured as a fallback version, stripped from identity.
	bedrockVer := ""
	if v := modelVersionSuffix.FindString(m); v != "" {
		bedrockVer = strings.TrimPrefix(v, "-")
		m = modelVersionSuffix.ReplaceAllString(m, "")
	}
	// `-YYYYMMDD` snapshot date — the primary version signal, stripped from identity.
	if d := modelDateSuffix.FindString(m); d != "" {
		id.Version = strings.TrimPrefix(d, "-")
		m = modelDateSuffix.ReplaceAllString(m, "")
	} else {
		id.Version = bedrockVer
	}

	// Vendor: keep the wire prefix when present, else derive from the family root.
	if vp := modelVendorPrefix.FindString(m); vp != "" {
		id.Vendor = strings.TrimSuffix(vp, ".")
		m = modelVendorPrefix.ReplaceAllString(m, "")
	} else {
		for _, e := range vendorByFamilyPrefix {
			if strings.HasPrefix(m, e.prefix) {
				id.Vendor = e.vendor
				break
			}
		}
		if id.Vendor == "" && openAIOSeriesFamily.MatchString(m) {
			id.Vendor = "openai"
		}
	}

	if id.Vendor != "" {
		id.Model = id.Vendor + "." + m
	} else {
		id.Model = m // unrecognized vendor — bare identity, never guessed
	}
	return id
}
