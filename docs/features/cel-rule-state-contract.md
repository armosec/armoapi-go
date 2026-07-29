---
type: feature
status: in-progress
owner: ben@armosec.io
scope: repo
related_code:
  - armotypes/runtimerule.go
  - armotypes/runtimeincidents.go
---

# CEL rule state contract (cross-event correlation)

Shared types that let a CEL runtime rule remember a fact from one event and read it
back when a later event arrives, and that carry the remembered event into the alert.

Consumed by **node-agent** (node telemetry) and the **operator** (`k8s-admission`).
Both watch the same `Rules` CRD and both build `armotypes.RuntimeAlert`, so these
types are defined here once rather than per component.

Design spec: `shared-designs-and-docs/projects/2026-07-28-cel-rule-state-store/spec.md`.

## Why it exists

A CEL rule is otherwise a stateless predicate over a single event: it returns a
boolean and remembers nothing. That makes any detection whose signal is *two or more
events sharing a key within a time window* inexpressible — a process executed from a
mount that then connects out, a web server spawning a shell whose descendant makes an
outbound connection, or `create pod` → `exec pod` → `delete pod` by one subject.

## `RuntimeRule.StateWrites`

```go
StateWrites []StateWrite `json:"stateWrites,omitempty" yaml:"stateWrites,omitempty" bson:"stateWrites,omitempty"`
```

A rule declares what it remembers. Writes are **declared, not expressed as a CEL
setter** — a setter inside a boolean can be skipped by short-circuiting, may be
reordered by CEL's static optimiser, and would fuse "did this rule fire?" with "did
this rule remember something?", making write-without-alerting impossible. That last
property is exactly what the first leg of every correlation rule needs.

| Field | Meaning |
|---|---|
| `eventType` | The stream driving this write. Need not appear in `expressions.ruleExpression` — that is what allows remembering on `exec` and alerting on `network`. |
| `when` | CEL boolean guard. Empty means always write. A guard may itself read state, which is how multi-step sequences are expressed without a new construct. |
| `scope` | `container` / `pod` / `node` (node-agent) or `identity` (operator). |
| `name` | *What kind of fact.* A string literal, never an expression, so it stays statically analysable and safe as a metric label. |
| `key` | *Who the fact is about.* A CEL string expression. **Optional** — omit for a fact true of the whole scope, which yields one entry per scope. |
| `value` | Optional author extras as CEL expressions. Keys must not begin with `_`, which is reserved for engine-stamped provenance. |
| `ttl` | Go duration string. The consuming engine clamps it to its configured maximum at load, so no rule can pin memory indefinitely. |

### `name` vs `key`

These are the two most-confused fields. `name` is *what kind of fact*; `key` is *who
it is about*. One container may have 50 processes that exec'd from a mount: that is
**one** `name` with **50** `key`s. `name` is the variable, `key` indexes into it.

They are kept separate because the engine must be able to hold the fact type fixed
while varying the subject — that is what makes ancestry matching ("does any ancestor
of this PID carry a `mount_exec` marker?") mechanically possible. Fused into one
opaque string, the engine could not tell which substring is the PID.

### `StateScope`

```go
StateScopeContainer StateScope = "container"
StateScopePod       StateScope = "pod"
StateScopeNode      StateScope = "node"
StateScopeIdentity  StateScope = "identity"
```

The scope determines which entries a read can see. The **engine** resolves the
concrete scope ID from the event — never the rule — so cross-container and
cross-tenant leakage is not expressible regardless of what an author writes in `key`.

`identity` is the operator's only scope: an admission chain is defined by *who*
performed the sequence, so memory is bounded per caller.

## `CorrelationAlert` on `RuntimeAlert`

```go
type CorrelationAlert struct {
    Correlations []CorrelationEvidence
}
```

Inlined into `RuntimeAlert` alongside `RuleAlert`, `MalwareAlert`, `AdmissionAlert`
and the rest, so `correlations` appears at the **top level** of the alert payload.
Because both node-agent and the operator build `RuntimeAlert`, this single type
carries correlation evidence for both.

`CorrelationEvidence` describes the *remembered* end of a chain — the part the
triggering event cannot express. Without it, an alert would say only "a process made
an outbound connection" and drop the exec that makes it interesting.

`Process` and `Admission` are a **tagged union**: exactly one is set, determined by
the engine that wrote the entry. An admission request has no PID, and its useful
identity (who did what to which object) has no home in `Process`.

| Field | Meaning |
|---|---|
| `name` | The state variable the entry was stored under |
| `eventType` | Stream the remembered event came from |
| `timestamp` | Kernel/request time of the remembered event, **not** observation time — ordering comparisons rely on this |
| `scope`, `key` | Where it was stored; `key` empty for scope-wide markers |
| `process` | node-agent entries (reuses `armotypes.Process`) |
| `admission` | operator entries (`AdmissionEvidence`) |
| `values` | Author extras from the rule's `StateWrite.Value` |

`InfectedPID` and the alert's process tree continue to describe the **triggering**
event. Correlation enriches an incident; it does not re-key it, so backend incident
grouping is unchanged.

## `EventType` and `IsKnownEventType`

`EventType` previously defined only the 11 streams admission rules referenced.
node-agent has 20, so `ptrace`, `bpf`, `kmod`, `unshare`, `iouring`, `randomx`,
`procfs`, `fork`, `exit` and the `all` wildcard were added — a `stateWrites` clause
may name any of them.

```go
func IsKnownEventType(e EventType) bool
```

Rule loaders use this to reject a clause naming an event stream that does not exist.
`EventType` is a string type, so an unknown value would otherwise deserialize happily
and produce a rule that silently never matches — the failure mode this contract works
hardest to avoid.

## Compatibility

Purely additive. No existing field was renamed, retyped, or retagged. All new fields
are `omitempty` on `json`, `yaml` and `bson`, so rules and alerts that do not use
correlation serialize byte-identically to before. `runtimerule_statewrites_test.go`
and `correlation_test.go` assert this directly.
