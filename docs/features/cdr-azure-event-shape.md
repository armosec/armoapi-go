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

## Read identity through the accessors

```go
caller := ev.EffectiveCaller()          // NOT ev.Caller
claims := ev.EffectiveClaims()          // NOT ev.Claims
auth   := ev.EffectiveAuthorization()   // NOT ev.Authorization
```

Each prefers the nested `identity` and falls back to the flat field, so callers
work on either shape. `EffectiveCaller` resolves the Event Hub case from the
`upn` / `name` claims (including their WS-Fed URI forms), since that path has no
`caller`. `EffectiveClaims` returns an empty map rather than nil, so indexing a
missing claim is always safe.

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

## `properties` may be an object or a string

Some operations deliver `properties` as a JSON-**encoded string** rather than an
object. `AzureProperties.UnmarshalJSON` accepts:

| Delivered | Result |
|---|---|
| `{"statusCode":"OK"}` | the map |
| `"{\"statusCode\":\"OK\"}"` | parsed into the map |
| `"something happened"` | `{"message": "something happened"}` |
| `null` / absent | nil |
| anything else | `{"value": …}` |

It never fails the surrounding decode. That matters because the failure was not
local: a failed record decode made the collector's alert builder return an error,
and the detection built from that record was discarded — a dropped alert caused by
a field the detection never referenced.

## Fields deliberately not modelled

The Event Hub wrapper fields (`RoleLocation`, `Stamp`, `ReleaseVersion`, `VmSku`)
carry no detection or attribution value. Unknown fields are ignored on decode, so
adding them is unnecessary; detection rules are evaluated against the **raw**
record JSON, not this struct, so a rule can still reference anything the export
sends.
