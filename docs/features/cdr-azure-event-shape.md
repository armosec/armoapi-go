---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/azure.go
  - armotypes/cdr/azure_shapes_test.go
  - armotypes/cdr/azure_properties_bson_test.go
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

## `properties` is held raw — and crosses BSON, not just JSON

Some operations deliver `properties` as a JSON-**encoded string** rather than an
object, so the field is held raw:

```go
Properties AzureProperties `json:"properties,omitempty" bson:"properties,omitempty"`
```

Any well-formed JSON value decodes, and the bytes are forwarded into the alert
payload unchanged. Read it by unmarshalling into whatever you expect:

```go
var props map[string]any
_ = json.Unmarshal(ev.Properties, &props)   // handle the error; may be a string
```

**Why raw matters more than it looks.** A typed `map[string]any` rejects the
string form, and the failure is not local: the whole record decode fails, the
collector's alert builder returns an error, and the detection built from that
record is discarded — a dropped alert caused by a field the detection never
referenced. A raw field cannot fail that way.

**Why raw rather than a tolerant map.** No Go code reads this field; rules are
evaluated against the raw record JSON, not this struct. Coercing the string form
into a map would mean inventing keys (`message`, `value`) that no consumer reads,
and rewriting an audit record before it reaches the backend. Raw keeps it
byte-for-byte faithful to what Azure sent.

### Why it needs its own type, and not a bare `json.RawMessage`

**This struct is not only a wire shape.** config-service persists runtime
incidents with `AzureActivityLogEvent` embedded, so the bag also round-trips
through **BSON**. A bare `json.RawMessage` is a `[]byte` there, and the driver's
byte-slice codec accepts only BSON binary or string — never an embedded document,
which is what every incident stored before this type existed holds.

That shipped once. Typing the field `json.RawMessage` was reasoned about purely as
a JSON concern; it reached config-service two weeks later as an incidental
`go.mod` bump, and every Azure incident started returning 500:

```
error decoding key cdrevent.eventdata.azureactivitylog.properties:
cannot decode document into json.RawMessage
```

Writing was as quietly wrong in the other direction: `[]byte` marshals **to** BSON
binary, so incidents stored during that window hold an opaque blob that no Mongo
filter or index can reach.

`AzureProperties` therefore carries its own BSON codec, accepting every shape the
collection now holds and storing the bag in its natural BSON form:

| Stored as | Written by | Read back as |
|---|---|---|
| embedded document | the original `map[string]any` field, and any object bag | the object, as JSON |
| binary | the bare `json.RawMessage` field | the JSON payload |
| string | Azure's JSON-in-a-string operations, and any non-JSON bag | the JSON it holds, else the quoted string |
| array | an array bag | the array, as JSON |
| double / boolean / int | a top-level scalar bag | the scalar, as JSON |
| null / absent | an event with no bag | `nil` |

**Writing does not flatten the bag into a document.** An object becomes an
embedded document and an array a BSON array — the point being that a structured
bag stays queryable and indexable instead of becoming the opaque blob the binary
form produced. A top-level scalar is stored as the matching BSON scalar.

**Reading and writing are symmetric, and that is a property to preserve.** Every
type the writer can emit, the reader reads back. The first version of this codec
was not symmetric: it wrote top-level scalars as BSON doubles and booleans but
only read documents, binary and strings, so a scalar bag round-tripped to `nil`
with no error — silent data loss.
`TestAzureProperties_WriterAndReaderAgreeOnEveryShape` is the guard.

Decoding **never fails**: a bag of a foreign BSON type (a date, an ObjectID —
nothing writes those here) yields `nil` rather than sinking the incident it
belongs to, the same bargain the JSON side makes.

Documents reach JSON through `jsonifyBSON`, which rewrites the driver's ordered
`primitive.D` into a JSON object. BSON types with no JSON scalar form (dates,
binary) come back as their driver representation rather than their original text —
exact in practice for an Activity Log bag, which is plain JSON, but not a
byte-exact round trip for arbitrary BSON.

`azure_properties_bson_test.go` pins all four stored shapes and both write forms.
**Any raw-held field on a type that config-service persists needs the same
treatment** — a JSON round-trip test alone will not catch this class of bug.

## Fields deliberately not modelled

The Event Hub wrapper fields (`RoleLocation`, `Stamp`, `ReleaseVersion`, `VmSku`)
carry no detection or attribution value. Unknown fields are ignored on decode, so
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
