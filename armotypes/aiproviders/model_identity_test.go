package aiproviders

import "testing"

// The load-bearing property (SUB-8068): every wire form of one logical model
// collapses to ONE identity, with the date/region lifted out — and a NEW model the
// registry has never seen still normalizes correctly.
func TestNormalizeModelIdentity(t *testing.T) {
	for _, tc := range []struct {
		name                           string
		raw, host                      string
		model, vendor, version, region string
	}{
		{
			name:  "bedrock full inference-profile id",
			raw:   "us.anthropic.claude-haiku-4-5-20251001-v1:0",
			host:  "bedrock-runtime.us-east-1.amazonaws.com",
			model: "anthropic.claude-haiku-4-5", vendor: "anthropic",
			version: "20251001", region: "us-east-1",
		},
		{
			name:  "bedrock stripped form (same model, no host)",
			raw:   "claude-haiku-4-5-20251001",
			model: "anthropic.claude-haiku-4-5", vendor: "anthropic", version: "20251001",
		},
		{
			name:  "public-API dateless alias UNIFIES with the dated snapshot identity",
			raw:   "claude-haiku-4-5",
			model: "anthropic.claude-haiku-4-5", vendor: "anthropic", version: "",
		},
		{
			name:  "provider-prefixed bedrock form",
			raw:   "anthropic.claude-3-5-sonnet-20241022-v2:0",
			host:  "bedrock-runtime.eu-west-1.amazonaws.com",
			model: "anthropic.claude-3-5-sonnet", vendor: "anthropic",
			version: "20241022", region: "eu-west-1",
		},
		{
			name:  "minor version is KEPT, not collapsed to family (opus 4 vs 4.1)",
			raw:   "claude-opus-4-1-20250805",
			model: "anthropic.claude-opus-4-1", vendor: "anthropic", version: "20250805",
		},
		{
			name:  "opus 4 stays distinct from opus 4.1",
			raw:   "claude-opus-4-20250514",
			model: "anthropic.claude-opus-4", vendor: "anthropic", version: "20250514",
		},
		{
			name:  "amazon nova via bedrock",
			raw:   "us.amazon.nova-lite-v1:0",
			model: "amazon.nova-lite", vendor: "amazon", version: "v1:0", region: "us",
		},
		{
			name:  "openai model, bare",
			raw:   "gpt-4o-mini",
			model: "openai.gpt-4o-mini", vendor: "openai",
		},
		{
			name:  "NEW model in a known family (unmapped) still normalizes",
			raw:   "claude-opus-99-20260101",
			model: "anthropic.claude-opus-99", vendor: "anthropic", version: "20260101",
		},
		{
			name:  "genuinely unknown vendor: bare identity, NEVER guessed",
			raw:   "acme-frobnicator-7",
			model: "acme-frobnicator-7", vendor: "", version: "",
		},
		{
			// The o-series is matched by pattern, not enumeration: o1/o3/o4 still
			// resolve, and a FUTURE o5/o6 resolves too (the fork the enumeration caused).
			name:  "openai o1 (bare) resolves",
			raw:   "o1",
			model: "openai.o1", vendor: "openai",
		},
		{
			name:  "openai o3-mini resolves",
			raw:   "o3-mini",
			model: "openai.o3-mini", vendor: "openai",
		},
		{
			name:  "openai o1-preview with snapshot date",
			raw:   "o1-preview-20250101",
			model: "openai.o1-preview", vendor: "openai", version: "20250101",
		},
		{
			name:  "NEW o-series generation (o5) resolves via the pattern, no fork",
			raw:   "o5-mini",
			model: "openai.o5-mini", vendor: "openai",
		},
		{
			name:  "o5 bedrock wire form collapses to the SAME identity as the bare alias",
			raw:   "openai.o5-mini",
			model: "openai.o5-mini", vendor: "openai",
		},
		{
			// Guard: an unrecognized family that merely STARTS with "o" must not be
			// mis-attributed to openai — the pattern is `o` + a DIGIT.
			name:  "non-openai model starting with o is NOT guessed (olmo)",
			raw:   "olmo-7b",
			model: "olmo-7b", vendor: "",
		},
		{
			// Bedrock GLOBAL cross-region inference profile: the `global.` prefix must be
			// recognized and stripped so it unifies with the plain identity (not fork).
			name:  "bedrock global cross-region profile unifies with plain identity",
			raw:   "global.anthropic.claude-sonnet-4-5-v1:0",
			model: "anthropic.claude-sonnet-4-5", vendor: "anthropic",
			version: "v1:0", region: "global",
		},
		// ── Family-boundary near-misses: a root that only PREFIXES a longer word must
		// NOT be attributed to that vendor (documented "unrecognized vendor is NEVER
		// guessed"). Each keeps a bare identity.
		{
			name:  "gptfoo is not openai (gpt + letter)",
			raw:   "gptfoo-1",
			model: "gptfoo-1", vendor: "",
		},
		{
			name:  "claudette is not anthropic (claude + letter)",
			raw:   "claudette-1",
			model: "claudette-1", vendor: "",
		},
		{
			name:  "commandant is not cohere (command + letter)",
			raw:   "commandant-1",
			model: "commandant-1", vendor: "",
		},
		{
			name:  "o7zip is not an openai o-series model (o7 + letter)",
			raw:   "o7zip",
			model: "o7zip", vendor: "",
		},
		// ── Boundary POSITIVES: a digit or separator after the root is a real family
		// member and MUST still resolve.
		{
			name:  "llama3 (digit boundary) still resolves to meta",
			raw:   "llama3",
			model: "meta.llama3", vendor: "meta",
		},
		{
			name:  "o5-mini (future o-series) still resolves to openai",
			raw:   "o5-mini",
			model: "openai.o5-mini", vendor: "openai",
		},
		{
			name:  "empty input yields empty identity",
			raw:   "",
			model: "", vendor: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeModelIdentity(tc.raw, tc.host)
			if got.Model != tc.model {
				t.Errorf("Model = %q, want %q", got.Model, tc.model)
			}
			if got.Vendor != tc.vendor {
				t.Errorf("Vendor = %q, want %q", got.Vendor, tc.vendor)
			}
			if got.Version != tc.version {
				t.Errorf("Version = %q, want %q", got.Version, tc.version)
			}
			if got.Region != tc.region {
				t.Errorf("Region = %q, want %q", got.Region, tc.region)
			}
		})
	}
}

// The whole point: the three wire forms of Haiku 4.5 that render as three separate
// rows today must share ONE identity after normalization.
func TestNormalizeModelIdentity_CollapsesWireForms(t *testing.T) {
	forms := []string{
		"us.anthropic.claude-haiku-4-5-20251001-v1:0",
		"claude-haiku-4-5-20251001",
		"claude-haiku-4-5",
	}
	const want = "anthropic.claude-haiku-4-5"
	for _, f := range forms {
		if got := NormalizeModelIdentity(f, "").Model; got != want {
			t.Errorf("%q → %q, want %q", f, got, want)
		}
	}
}
