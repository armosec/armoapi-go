---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/cdr.go
---

# CDR heartbeat contract (`CdrAlertBatch.IsHeartbeat`)

A provider-neutral wire flag that lets an in-tenant CDR collector send a periodic
**liveness** message on the same channel as detection alerts, without it being
mistaken for a detection.

Produced by the cloud CDR collectors (Azure/GCP; AWS may adopt it later) and consumed
by **event-ingester** (`cdrprocessor`). Both sides share the `CdrAlertBatch` type, so
the flag is defined here once.

Design/analysis: `shared-designs-and-docs/cdr/cdr-backend-azure-readiness-audit.md`.

## Why it exists

CADR "connected" status is kept fresh by keep-alive: the ingester stamps the CADR
feature's `LastKeepAlive` on every message it processes, and a disconnect cron flips
an account to `Disconnected` once that stamp goes stale. On a quiet subscription no
real alerts fire, so the collector must send a periodic heartbeat or it would be
falsely disconnected.

AWS encodes its one-shot, deploy-time connect signal as a rule failure named
`"StackReady"` (fired by a CloudFormation custom resource). That is the wrong vehicle
for a periodic, cross-cloud heartbeat: Azure and GCP have no "stack" concept, and
`StackReady` already means something else (deploy-time, once). `IsHeartbeat` models the
periodic liveness signal as a first-class, provider-neutral field instead of
overloading a magic rule name.

## The field

```go
// IsHeartbeat marks a periodic liveness message rather than a detection batch.
IsHeartbeat bool `json:"isHeartbeat,omitempty"`
```

- **Absent / `false`** — a normal alert batch. Existing AWS producers are unaffected;
  the field is additive and back-compatible.
- **`true`** — a heartbeat. The batch carries **no `RuleFailures`**. The ingester uses
  it to refresh the CADR feature's connected / `LastKeepAlive` state and **skips the
  alert pipeline** (no incident is created).

## Connection level (routing)

```go
type ConnectionLevel string

const (
    ConnectionLevelAccount      ConnectionLevel = "account"
    ConnectionLevelOrganization ConnectionLevel = "organization"
)

// ConnectionLevel ConnectionLevel `json:"connectionLevel,omitempty"`
```

`ConnectionLevel` states whether the batch belongs to a single-account connection or an
organization/tenant-wide one, so the ingester routes keep-alive by **stated intent** rather
than inferring it from `OrgID` presence. It applies to **both** heartbeats and real alert
batches.

- **`account`** — keep-alive is keyed on `CloudAccountID` (the account/subscription), even
  when `OrgID` (the tenant) is also carried as data.
- **`organization`** — keep-alive is keyed on `OrgID`.
- **Absent** — legacy inference: `OrgID == ""` → account, else organization. Existing AWS
  producers set nothing, so their routing is unchanged.

This lets a producer always send both `CloudAccountID` and `OrgID` without `OrgID`'s presence
being misread as "this is an org connection." Design:
`shared-designs-and-docs/docs/superpowers/specs/2026-08-11-cdr-connection-level-dispatch-design.md`.

## Producer expectations

Because a heartbeat has no `RuleFailures`, the batch **must** set the batch-level
identity itself — the consumer cannot fall back to a rule's `CloudMetadata`:

- `Provider` — required (the consumer's provider resolution otherwise reads
  `RuleFailures[0].CloudMetadata.Provider`, which is absent here).
- `CustomerGUID`, and `CloudAccountID` (single account) and/or `OrgID` (org/tenant-wide).
- `ConnectionLevel` — recommended, so routing is explicit rather than inferred. Set
  `account` for a single-account connection (carry the tenant in `OrgID` as data if known)
  and `organization` for a tenant-wide one. Omitting it falls back to legacy `OrgID`-presence
  inference.

Cadence is the producer's responsibility: send one on startup (acts as the connect
signal) and then periodically at an interval **shorter than the disconnect threshold**.

## Consumer expectations

- Treat `IsHeartbeat == true` as a liveness signal: refresh connected / `LastKeepAlive`
  and return before the detection path.
- Route keep-alive on `ConnectionLevel` when set (`account` → key on `CloudAccountID`,
  `organization` → key on `OrgID`); when absent, fall back to `OrgID`-presence inference.
  An unrecognized value must fall back to inference, never error.
- The legacy AWS `StackReady` rule-name path stays as-is (deploy-time connect) — a
  heartbeat does not replace it, and the two are dispatched independently.

## Non-goals

This type only declares the flag. The ingester dispatch and the collector emitter are
wired in their respective repos (`event-ingester-service`, `cdr-agents`).
