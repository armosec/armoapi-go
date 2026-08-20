---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/cdr.go
---

# CDR heartbeat contract (`CdrAlertBatch.IsHeartbeat` / `ConnectionLevel` / `LogsSeen`)

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

type CdrAlertBatch struct {
    // ... existing fields ...
    ConnectionLevel ConnectionLevel `json:"connectionLevel,omitempty"`
}
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

## Logs seen (verification)

```go
type CdrAlertBatch struct {
    // ... existing fields ...
    LogsSeen *uint64 `json:"logsSeen,omitempty"`
}

func (b *CdrAlertBatch) LogsSeenValue() (count uint64, reported bool)
func NewLogsSeen(count uint64) *uint64
```

`LogsSeen` is how many raw audit records the collector has received through its log pipe
since start. It exists because **a heartbeat proves only that the collector process is
alive — it never traverses the log pipe**, so liveness alone cannot distinguish "everything
working" from "collector healthy, log routing silently dead."

That is not a theoretical concern. The GCP POC produced exactly that state: a Log Router
sink that was enabled, correctly filtered, and whose writer identity held
`roles/pubsub.publisher` on the destination topic, yet routed **zero** logs for 6.5
minutes — while raising no error anywhere. GCP's own telemetry was no safety net either
(`exports/error_count` stayed empty through ~11 silently-dropped entries). A
liveness-driven `Connected` would have reported a fully-blind account as healthy.

Reporting the count lets the backend flip `Pending → Connected` only once a log has
provably traversed sink → transport → collector.

### Absent and zero mean different things

The field is a **pointer**, and that is load-bearing. `omitempty` on a pointer tests nil,
not the pointee, so an explicit `0` is still serialized and stays distinguishable from
"not reported":

| Wire | Meaning | Consumer behavior |
|---|---|---|
| field absent | this producer does not report the signal (AWS and Azure today) | fall back to flipping on liveness alone |
| `"logsSeen": 0` | producer reports it; no log seen **by this instance** yet | do not advance to `Connected`; never regress from it |
| `"logsSeen": N` (N > 0) | a log provably traversed the pipe | flip `Pending → Connected` |

Conflating the first two rows breaks the gate in **both** directions: treating absent as
zero would hold every AWS/Azure connection `Pending` forever, and treating zero as absent
would report a silently-blind collector as `Connected` — the exact failure the field
exists to catch. Use `LogsSeenValue()`, which returns `(count, reported)`, rather than
dereferencing or defaulting the pointer.

### What it counts

Every audit record that arrived through the log pipe and decoded — **including** records
a detection gate skips and records no rule matches. It is evidence that the *pipe* is
live, not that a detection fired, so filtering it by either would defeat its purpose.

**It is cumulative per collector instance, and resets to zero when that instance is
replaced.** It is not a durable total and is not comparable across restarts. Instances
*are* replaced in normal operation — a deploy, a crash, or the platform recycling a
container (which Cloud Run does even at `min-instances=1`) — so a healthy,
already-`Connected` account will legitimately emit `logsSeen: 0` again for as long as it
takes the fresh instance to see its first log. See the latching rule below; this is the
single most likely way to misimplement the consumer side.

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
- `LogsSeen` — set it if the collector can count records off its log pipe. Report a
  genuine `0` (via `NewLogsSeen(0)`) rather than omitting the field: omitting it means
  "this producer does not report the signal", which is a different statement.

Cadence is the producer's responsibility: send one on startup (acts as the connect
signal) and then periodically at an interval **shorter than the disconnect threshold**.

> ⚠️ **The cadence must stay unconditional.** Do not withhold heartbeats until
> `LogsSeen > 0`. Heartbeat freshness is what keeps a quiet-but-healthy account
> `Connected`, so gating the beat itself would false-disconnect a dormant-but-correctly-
> wired account. Only the `Pending → Connected` *flip* is conditioned on the count.

## Consumer expectations

- Treat `IsHeartbeat == true` as a liveness signal: refresh connected / `LastKeepAlive`
  and return before the detection path.
- Route keep-alive on `ConnectionLevel` when set (`account` → key on `CloudAccountID`,
  `organization` → key on `OrgID`); when absent, fall back to `OrgID`-presence inference.
  An unrecognized value must fall back to inference, never error.
- Read `LogsSeen` via `LogsSeenValue()`. When `reported` is false, keep the existing
  liveness behavior (flip on first heartbeat). When it is true, gate `Pending →
  Connected` on `count > 0`; at zero, **leave the connection in whatever state it is
  already in** — do not advance it. Either way, keep refreshing `LastKeepAlive` on every
  message: the count gates the *connect* transition only, never the disconnect one, or a
  quiet account would drift to `Disconnected` while healthy.
- **The `Pending → Connected` flip latches — it is one-way.** Once an account is
  `Connected`, a later `logsSeen: 0` must never regress it to `Pending`. The count is
  per-instance and resets on restart (above), so reading zero as "not connected" rather
  than "no new evidence" would flap the connection on every container recycle. Treat a
  zero as *absence of new evidence*, never as evidence of absence.
- The legacy AWS `StackReady` rule-name path stays as-is (deploy-time connect) — a
  heartbeat does not replace it, and the two are dispatched independently.

## Non-goals

This type only declares the fields. The ingester dispatch — including the `Pending →
Connected` gate on `LogsSeen` — and the collector emitter are wired in their respective
repos (`event-ingester-service`, `cdr-agents`).
