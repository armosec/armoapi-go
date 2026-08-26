---
type: feature
status: active
owner: "@armosec/backend"
scope: repo
related_code:
  - armotypes/runtimerule.go
  - armotypes/runtimerule_variants.go
---

# Rule variants (version-gated rule bodies)

`RuntimeRule.Variants` lets one rule ID/name serve different CEL bodies across agent-capability
generations, instead of forking rule identity every time a new CEL library function is added and
part of the fleet doesn't have it yet.

Design spec: `shared-designs-and-docs/projects/rule-improvement/cel-rule-version-resolution-design.md`
(§7, "Backend changes", is the section this schema implements — Phase 2 of that design).

## Why it exists

A rule expression calling a CEL function an agent hasn't registered fails to compile, and today
that makes the **entire rule silently and permanently stop firing** for that agent, not just the
part using the new capability. The backend wants to ship rules using new capabilities without
regressing agents that haven't upgraded yet, while keeping one rule ID/name across generations
(forking identity per capability generation was explicitly rejected — see the design's §11
"Alternative B").

## Shape

```go
type RuleVariant struct {
    MinAgentVersion string          `json:"minAgentVersion" yaml:"minAgentVersion" bson:"minAgentVersion"`
    Expressions     RuleExpressions `json:"expressions" yaml:"expressions" bson:"expressions"`
}

// on RuntimeRule:
Variants []RuleVariant `json:"variants,omitempty" yaml:"variants,omitempty" bson:"variants,omitempty"`
```

`MinAgentVersion` is a **bare semver** (`1.2.3`), not a constraint string like the sibling
`AgentVersionRequirement` field (`>=1.2.3`) — a bare version gives "pick the highest variant the
requesting agent satisfies" a well-defined total order across all variants on a rule.

## The mirror invariant

`RuntimeRule.Expressions` (the pre-existing top-level field) must always equal the
lowest-`MinAgentVersion` variant's `Expressions`. This is **not enforced by this package** — every
reader that doesn't know about `Variants` (console, K8s connector, policy validation,
incident-enrichment) reads only the top-level field, so it degrades to a working base body instead
of an empty/stale one. Each writer that constructs a `RuntimeRule` document is responsible for
deriving this itself using `LowestVariant`; see `cadashboardbe`'s custom-rule write path and
`rulelibrary`'s `yaml_to_mongodb.py` for the two current writers.

## API

- `ValidateVariants(variants) error` — rejects any `MinAgentVersion` that isn't a parsable bare
  semver. Run this before persisting a rule with variants.
- `LowestVariant(variants) (RuleVariant, bool)` — the variant every writer derives the top-level
  `Expressions` mirror from.
- `ResolveVariantExpressions(variants, agentVersion) (RuleExpressions, bool)` — the agent-facing
  resolver's core function: picks the highest-`MinAgentVersion` variant satisfied by
  `agentVersion`. Falls back to the **lowest** variant (never the newest) whenever `agentVersion`
  is empty, unparsable, or below every variant's floor — an agent whose capability can't be
  confirmed must never be handed a body it might not be able to compile. Returns `false` if
  `variants` is empty, in which case the caller should keep using the rule's plain top-level
  `Expressions`.

## Compatibility

Purely additive — `omitempty`, no existing field renamed/retyped/retagged. Rules with no variants
serialize byte-identically to before (`TestRuntimeRule_Variants_YAMLRoundTrip` and the rest of
`runtimerule_variants_test.go` assert the new type directly; the pre-existing round-trip tests in
this package are unaffected since they never set `Variants`).

`Variants` is deliberately **omitted from agent-facing resolver responses** — the resolver
projects the resolved variant's `Expressions` into the response's top-level `Expressions` field
instead, so an agent's decode path needs no change at all.
