---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/azure.go
  - armotypes/cdr/azure_shapes_test.go
  - armotypes/cdr/azure_eventhub_capture_test.go
---

# Azure Activity Log event shape (two delivery paths)

`AzureActivityLogEvent` models the **same** Azure Activity Log record as delivered
by two different paths, which do not agree on where the caller identity lives.
Reading the wrong fields yields silent emptiness, not an error.

## The two shapes

| | Event Hub / diagnostic-settings export | Azure Monitor REST / management API |
|---|---|---|
| Who consumes it | the in-account CDR collector | `az monitor activity-log list`, portal, SDK queries |
| `caller` | **almost always absent** (9/339 — `Microsoft.Sql/servers/*` only) | present (UPN) |
| `authorization` | nested under `identity` | top level |
| `claims` | nested under `identity` | top level |
| `operationName`, `category` | plain **strings** (upper-cased) | `{value, localizedValue}` objects |
| Extras | `RoleLocation`, `Stamp`, `ReleaseVersion`, `VmSku` | — |

Both are captured as fixtures in `azure_shapes_test.go`. A sanitized **real**
Event Hub record lives in `azure_eventhub_capture_test.go`; prefer it when
reasoning about production behaviour.

## Measured against real data (SUB-7951)

Microsoft's docs contradict each other on identity placement, so it was settled
by capture: **339 records** off a live subscription-level Event Hub diagnostic
setting, 2026-08-18.

| Field | Present |
|---|---|
| top-level `authorization` | **0/339** |
| top-level `claims` | **0/339** |
| top-level `caller` | 9/339 (`Microsoft.Sql/servers/*`) |
| `identity.authorization` | 330/339 |
| `identity.claims` | 330/339 |
| `subscriptionId` | 9/339 — derive from `resourceId` (parseable 339/339) |
| `location` | **0/339** |
| `channels` | **0/339** |
| `resourceGroupName` | 0/339 (9 records send `ResourceGroup`, value `"DummyValue"`) |

All 339 decode through `AzureActivityLogEvent` with **zero errors**.

Two consequences worth knowing:

- **`upn` is never present** (0/339). `EffectiveCaller` lists it first, but real
  resolution happens via the `name` claim. Harmless, just not what fires.
- **8/339 records resolve to no actor at all** — service-principal calls whose
  claims carry `appid` but neither `name` nor `upn`. Those alerts ship with an
  empty `UserIdentity`, which is a source limitation, not a mapping bug.

## There is no region on an Activity Log record

`Location` is modelled but Azure never sends it — 0/339 on the Event Hub shape and
absent from the REST representation too. This is why Azure CDR incidents show an
empty region, and it is a **source limitation, not a pipeline bug**: nothing is
discarded in transit. AWS CloudTrail carries `awsRegion`, so the asymmetry between
providers is inherent to the sources.

Do **not** substitute `RoleLocation`. It is the Azure infrastructure stamp that
processed the API request, not where the resource lives — a resource in one region
is routinely served by a stamp in another, so mapping it to region would populate
the field with confidently wrong data. Real region would require an ARM lookup on
`resourceId` (an API call per resource, plus a permissions expansion).

**The struct models the REST shape's identity layout only.** A verbatim REST
response does not decode through it — `operationName` and `category` arrive as
`{value, localizedValue}` objects and the string fields reject them. That is
deliberate: nothing in the codebase decodes a REST response through this struct,
and the collector's records always carry the string form, so normalizing the
object form would be speculative. If a REST reader is ever added, the `restShape`
fixture is where the object form has to appear first, and the failure will be a
loud unmarshal error rather than a silent one.

## Read identity through the accessors

```go
caller := ev.EffectiveCaller()          // NOT ev.Caller
claims := ev.EffectiveClaims()          // NOT ev.Claims
auth   := ev.EffectiveAuthorization()   // NOT ev.Authorization
```

They do not all resolve the same way:

| Accessor | Precedence |
|---|---|
| `EffectiveClaims` | nested `identity.claims` → flat `claims` → non-nil empty map |
| `EffectiveAuthorization` | nested `identity.authorization` → flat `authorization` → nil |
| `EffectiveCaller` | flat `caller` **first** → then `upn`, `name`, and their WS-Fed URI forms out of `EffectiveClaims` → `""` |

`EffectiveCaller` is the odd one out: it reads the flat field first because
`caller`, when present, is the identity Azure itself resolved, and the Event Hub
shape has no `caller` at all — so the claim fallback only ever runs on the shape
that needs it.

`EffectiveClaims` always returns a **non-nil** map, so ranging over the result is
safe on an event with no claims. (Indexing a nil map is already safe in Go; the
guarantee here is about ranging and about handing the map to code that appends.)

**Reading `ev.Claims` or `ev.Caller` directly is a bug** for anything processing
collector traffic: those are always empty on the Event Hub shape. That defect
shipped once — every Azure CDR alert carried no actor, `UserIdentity` was never
populated in `common.Identifiers`, and risk-acceptance advanced-scope matching on
a user could not match. It produced no log line, because an unknown JSON object
is silently ignored.

## `operationName` / `category` are strings on the collector's path

CEL detection rules compare them with `.matches("(?i)…")`, which requires a
string. The `{value, localizedValue}` object form exists **only** on the REST
path — do not "fix" a rule to use `.value` based on a REST sample, or it will stop
matching real collector traffic. A regression test asserts the string decode.

Compare case-insensitively: this path upper-cases `operationName` and
`resourceId`, while `category` and `resultType` keep their casing.

## `properties` is held raw

Some operations deliver `properties` as a JSON-**encoded string** rather than an
object, so the field is typed `json.RawMessage`:

```go
Properties json.RawMessage `json:"properties,omitempty"`
```

Any well-formed JSON value decodes, and the bytes are forwarded into the alert
payload unchanged. Read it by unmarshalling into whatever you expect:

```go
var props map[string]any
_ = json.Unmarshal(ev.Properties, &props)   // handle the error; may be a string
```

**Why this matters more than it looks.** A typed `map[string]any` rejects the
string form, and the failure is not local: the whole record decode fails, the
collector's alert builder returns an error, and the detection built from that
record is discarded — a dropped alert caused by a field the detection never
referenced. `json.RawMessage` cannot fail that way.

**Why raw rather than a tolerant map.** No Go code reads this field; rules are
evaluated against the raw record JSON, not this struct. Coercing the string form
into a map would mean inventing keys (`message`, `value`) that no consumer reads,
and rewriting an audit record before it reaches the backend. Raw keeps it
byte-for-byte faithful to what Azure sent.

## Fields deliberately not modelled

The Event Hub wrapper fields (`RoleLocation`, `Stamp`, `ReleaseVersion`, `VmSku`)
carry no detection or attribution value. The capture also surfaced `durationMs`
(330/339, mostly `"0"`), `jobId`/`jobType` (58/339, Azure async-job bookkeeping
such as `ResourceConsistencyJob`), and `LogicalServerName`/`ResourceGroup` on the
9 SQL records — where `ResourceGroup` is the literal string `"DummyValue"`. None
are security-relevant, so none are modelled. Unknown fields are ignored on decode, so
adding them is unnecessary; detection rules are evaluated against the **raw**
record JSON, not this struct, so a rule can still reference anything the export
sends.

## Subscription-id casing and identity matching

The subscription id is derived from `resourceId`, which the Event Hub export
emits **upper-cased** (`/SUBSCRIPTIONS/<UPPER>/…`). Incidents therefore store an
upper-cased `cloudMetadata.account_id`, while account onboarding/storage keep the
customer's **original** casing (no normalization). So any match on the Azure
account id must be **case-insensitive** — normalizing to a single case is not
safe because no single layer owns the canonical form.

Two matchers rely on this:

- **Account lookup / keep-alive** — `event-ingester`'s
  `BuildGetAccountWithFeatureQuery` builds an anchored, escaped regex with the
  `ignorecase` option for Azure, so heartbeat (onboarded casing) and real-alert
  (upper-cased) batches resolve the same account.
- **Risk-acceptance / exception matching** — `GetRuntimeIncidentsRequestFilterFromExceptionPolicy`
  (`armotypes/exceptionpolicy.go`) emits the `cloudMetadata.account_id` filter the
  same way for Azure (`^<escaped>$|regex&ignorecase`). Without it an exact-equality
  match silently missed and the risk-acceptance never suppressed the incident.
  AWS/GCP keep exact equality (numeric / already-lowercase ids).

The recommended behavior is **(a) preserve the original casing everywhere and match
case-insensitively** at the comparison points above — not (b) rewrite ids into a
canonical case. Do not normalize casing at any layer: not the raw event (it must
stay byte-for-byte faithful to what Azure sent; rules are evaluated against it) and
not the derived identity fields (no single layer owns the canonical form, so a
partial rewrite just reintroduces mismatches). Case-insensitive matching is the
whole strategy.
