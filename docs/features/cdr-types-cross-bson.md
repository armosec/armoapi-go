---
type: feature
status: active
owner: shanyl@armosec.io
scope: repo
related_code:
  - armotypes/cdr/gcp.go
  - armotypes/cdr/azure.go
  - armotypes/cdr/gcp_bson_test.go
  - armotypes/cdr/azure_properties_bson_test.go
---

# CDR types are storage shapes, not just wire shapes

The `cdr` types are not only decoded from provider JSON. **config-service persists
runtime incidents with them embedded**, inserting and decoding typed Go structs
(`db/utils.go:824`, `:145`). So every field also round-trips through **BSON**, and a
field's Go type decides both what BSON gets written and what BSON the decoder will
accept.

Two production outages came from the same mistake — a type chosen for JSON
correctness on a struct that is also persisted. Both were invisible to the tests,
because the tests round-trip JSON only.

## What went wrong, twice

| | `AzureActivityLogEvent.Properties` | `GcpAuditLogPayload.NumResponseItems` |
|---|---|---|
| Type chosen | `json.RawMessage` | `json.Number` |
| Sound JSON reason | the bag arrives as an object *or* a JSON-encoded string | protojson emits int64 as `"3"`, but `3` is also valid |
| BSON failure | `[]byte` cannot decode a stored embedded document | empty value parses as neither int nor float |
| Direction | **read** — existing incidents 500 | **write** — new incidents never persist |
| Blast radius | every Azure incident | every GCP incident |

Neither failure is local. BSON encoding and decoding are per-**document**, so one
bad field fails the whole incident — a dropped or unreadable alert over a field
nothing reads.

### The `omitempty` trap

`NumResponseItems` *had* `omitempty` — in its `json` tag:

```go
NumResponseItems json.Number `json:"numResponseItems,omitempty"`   // BROKEN
```

The driver only consults `json` tags when `useJSONStructTags` is enabled, and it is
**off by default**. With no `bson` tag the field gets no `omitempty` at all, so the
empty value was always encoded, and always failed. The fix is the tag:

```go
NumResponseItems json.Number `json:"numResponseItems,omitempty" bson:"numResponseItems,omitempty"`
```

**A `json` tag is not a `bson` tag.** On a persisted type, spell out both.

## Rules for these types

1. **Every field needs an explicit `bson` tag.** Name and `omitempty` do not carry
   over from `json`.
2. **A `[]byte`-backed type needs its own BSON codec** — it writes as binary, which
   no Mongo filter or index can reach, and it cannot read a stored document. See
   `AzureProperties` in `azure.go` for the pattern.
3. **A type that cannot represent "absent" needs `bson:",omitempty"`** — `json.Number`
   is the one in this package; its zero value is the empty string, which does not
   encode.
4. **Test the BSON path.** A JSON round-trip test passes on all of the above.
   `gcp_bson_test.go` and `azure_properties_bson_test.go` are the guards.

## Finding the next one

The failure mode is that a zero or absent field breaks the document, so a plain
zero-value test misses it — nested pointers are nil, and the broken field is never
reached. Populate every pointer first, then marshal:

```go
p := reflect.New(reflect.TypeOf(doc))
fillAllPointers(p.Elem())          // allocate nested structs
_, err := bson.Marshal(p.Elem().Interface())
```

A reflective sweep of all 20 types config-service persists, run at
`armoapi-go v0.0.757`, found `NumResponseItems` as the **only** remaining
document-breaking field. It also flagged, as accepted and not bugs:

- `RuntimeAlert.Fields` (`json.RawMessage`) — stores as binary, so it is
  unqueryable, but it has *always* been this type, so every stored document is
  binary and reads are self-consistent. It would break the way `Properties` did if
  a writer ever stored it as a document.
- ~428 fields with `json` `omitempty` and no `bson` tag — zero values are written
  rather than omitted. Storage bloat, not breakage.
