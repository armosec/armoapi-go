---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/azure.go
  - armotypes/cdr/azure_shapes_test.go
---

# Azure Activity Log event shape (two delivery paths)

`AzureActivityLogEvent` models the **same** Azure Activity Log record as delivered
by two different paths, which do not agree on where the caller identity lives.
Reading the wrong fields yields silent emptiness, not an error.

## The two shapes

| | Event Hub / diagnostic-settings export | Azure Monitor REST / management API |
|---|---|---|
| Who consumes it | the in-account CDR collector | `az monitor activity-log list`, portal, SDK queries |
| `caller` | **absent** | present (UPN) |
| `authorization` | nested under `identity` | top level |
| `claims` | nested under `identity` | top level |
| `operationName`, `category` | plain **strings** (upper-cased) | `{value, localizedValue}` objects |
| Extras | `RoleLocation`, `Stamp`, `ReleaseVersion`, `VmSku` | — |

Both are captured as fixtures in `azure_shapes_test.go`.

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
carry no detection or attribution value. Unknown fields are ignored on decode, so
adding them is unnecessary; detection rules are evaluated against the **raw**
record JSON, not this struct, so a rule can still reference anything the export
sends.
