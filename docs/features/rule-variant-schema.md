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

- `ValidateVariants(variants) error` — rejects any `MinAgentVersion` that doesn't match
  `bareSemverPattern` (`^v?\d+\.\d+\.\d+$`: exactly three numeric segments, optional `v` prefix,
  **no** prerelease/build metadata) and rejects two variants sharing the same `MinAgentVersion`
  (after normalizing away the `v` prefix). Both checks exist for the same reason: without them,
  `LowestVariant` and `ResolveVariantExpressions` can disagree about which variant a given
  version means — a duplicate floor breaks their tie differently (stable-sort-keeps-first vs.
  loop-overwrites-on-tie), and prerelease/build metadata is exactly the input where this
  package's ordering (`hashiVer`, prerelease sorts below release) and rulelibrary's Python-side
  twin (`_parse_bare_semver` in `yaml_to_mongodb.py`, which strips prerelease/build before
  comparing) can silently rank two variants differently. Run this before persisting a rule with
  variants — every current writer does.
- `LowestVariant(variants) (RuleVariant, bool)` — the variant every writer derives the top-level
  `Expressions` mirror from.
- `ResolveVariantExpressions(variants, agentVersion) (RuleExpressions, bool)` — the agent-facing
  resolver's core function: picks the highest-`MinAgentVersion` variant satisfied by
  `agentVersion`. Falls back to the **lowest** variant (never the newest) whenever `agentVersion`
  is empty, unparsable, or below every variant's floor. That fallback prevents the specific
  failure this design exists to avoid — an old agent silently jumping to a body using a
  capability it doesn't have — but it is **not** by itself a guarantee that the returned body
  compiles on the requesting agent: if every variant's floor is above the agent's real version,
  the lowest variant is still returned, and nothing in this package enforces that a rule's
  lowest variant actually covers the oldest agent version the fleet runs. That floor-coverage
  discipline is a write-time authoring convention this function cannot see or enforce; callers
  that need the stronger guarantee must enforce it themselves. Returns `false` if `variants` is
  empty, in which case the caller should keep using the rule's plain top-level `Expressions`.

  **Asymmetric leniency, deliberately noted, not enforced.** Variant floors (`MinAgentVersion`)
  must match `ValidateVariants`' strict grammar — no prerelease/build metadata. `agentVersion` (the
  requester's own version) is not under this package's control and is parsed with
  `hashiVer.NewVersion`'s normal, lenient rules, which *does* accept prerelease metadata and
  orders it below the corresponding release. A requester reporting `1.5.0-rc1` therefore does not
  satisfy a `1.5.0` floor — conservative, not unsafe, but silent. As of this writing no known
  caller reports a prerelease-tagged agent version, so this is a documented assumption, not an
  observed problem.

## Rollout order — `scheduled-db-tasks` must consume this before any `variants:` rule ships

`rulelibrary`'s `yaml_to_mongodb.py` produces `variants:` in its release JSON independently of
this package (it's Python, no Go dependency on `armoapi-go`) — so it can, and does, ship ahead of
this schema landing anywhere else. But the managed-rule sync job
(`scheduled-db-tasks`' rule-library enricher) unmarshals that release asset into
`kdr.RuntimeRule` / `armotypes.RuntimeRule`, and its diff (`cmp.Diff`, no field allowlist) only
notices fields the struct it's compiled against actually has.

If that job's `armoapi-go` pin predates this `Variants` field, a `variants:`-bearing managed rule
is **silently accepted and stored as a plain single-body rule** — the field is dropped on
decode, and since `yaml_to_mongodb.py` already derives the top-level `expressions` from the
lowest variant, the stored document looks completely normal: readers get a working base body,
config-service's mirror backstop passes (there's no `Variants` present, so nothing to check), the
console shows an ordinary rule. Nothing signals that version gating never actually took effect —
newer agents get the same base body as old ones until the enricher is bumped and re-syncs.

**Do not author a `variants:`-bearing managed rule until `scheduled-db-tasks` has bumped past
this schema.** There is no dedicated PR for that bump in this rollout — it's a routine dependency
update once this PR (or whichever PR/tag supersedes it) is available, but it has to actually
happen before the rollout is complete, not just be theoretically possible.

## Compatibility

Purely additive — `omitempty`, no existing field renamed/retyped/retagged. Rules with no variants
serialize byte-identically to before (`TestRuntimeRule_Variants_YAMLRoundTrip` and the rest of
`runtimerule_variants_test.go` assert the new type directly; the pre-existing round-trip tests in
this package are unaffected since they never set `Variants`).

`Variants` is deliberately **omitted from agent-facing resolver responses** — the resolver
projects the resolved variant's `Expressions` into the response's top-level `Expressions` field
instead, so an agent's decode path needs no change at all.
