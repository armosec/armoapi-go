---
type: feature
status: active
owner: alon@armosec.io
scope: repo
related_code:
  - armotypes/runtimeincidents.go
  - armotypes/runtimeincidents_trace_test.go
---

# Call-stack trace status, module identity and join metadata

`Trace` and `StackFrame` used to describe only what a call stack *was*: a list of
named frames with a package and a language. They said nothing about how well the
capture went. A stack with two frames and a stack cut short at two frames looked
identical, and an alert with no stack at all looked the same as an alert whose
stack was dropped by a rate limiter.

This change adds that missing half. Every field is additive and carries
`omitempty` on both tags, so nothing changes for a consumer until it bumps.

## What was added

| Where | Field | Carries |
|---|---|---|
| `StackFrame` | `Status` (`FrameStatus`) | how this frame was recovered and named |
| | `Runtime` (`FrameRuntime`) | which runtime lane the frame belongs to |
| | `ModuleIdx` | index into `Trace.Modules` |
| | `FileOffset` | offset inside the module file, hex |
| | `Source` | what named the frame: `symtab`, `cache`, `deferred` |
| `Trace` | `Status` (`StackStatus`) | whether the stack is complete, cut, or absent |
| | `StatusReason` | why, when `Status` is `not-captured` |
| | `TruncatedAt` | frames recovered before the walk stopped |
| | `Modules` (`[]TraceModule`) | raw identity of each executable mapping |
| | `Hook` | which hook the event came from |
| | `JoinKey` | the per-hook key that matched the stack to the event |
| | `Tid`, `*BootNs` | thread and boot-clock timestamps |
| | `MechanismVersion` | the capture mechanism's wire and vocabulary versions |

`FrameStatus`, `StackStatus` and `FrameRuntime` are string-underlying typed
enums, following `EventType` and `StateScope` in `runtimerule.go`. The wire
format is a plain string, so adding the type costs nothing on the wire and makes
a producer convert explicitly instead of passing an arbitrary value.

## Three decisions worth knowing

### Status is decomposed, not one string

`Status` is a closed vocabulary. `StatusReason` and `TruncatedAt` hold the
parametric detail. Mongo can index and group by a closed enum; it cannot group by
`"truncated at 12 frames"` without making every distinct depth its own status
value. Keep `Status` small so a dashboard facet stays useful, and put the detail
in the other two fields.

### Module identity lives on the Trace, not the frame

A stack of 64 frames typically spans 2-4 distinct executable mappings. Putting
path, build-id, inode and device on each frame would repeat the same four values
about twenty times per alert. `Trace.Modules` carries each mapping once and
`StackFrame.ModuleIdx` points at it.

Keeping the *raw* identity, not only the resolved symbol, is the point. A stack
captured before its symbols were available can be re-symbolized later from the
build-id, which a resolved-name-only record cannot.

### Addresses stay hex strings — the BSON uint64 trap

BSON has no unsigned integer type. The driver encodes `uint64` as a signed 64-bit
integer and returns `value out of range` for anything above `math.MaxInt64`. A
virtual address routinely exceeds that bound, so it can never be stored as a
number:

- `Address` and `FileOffset` are hex strings.
- `TraceModule.Inode` and `Device` are `uint64` **only because they stay below
  `MaxInt64` by construction**.

`TestTraceModule_InodeStaysBelowMaxInt64` pins both halves: `MaxInt64` marshals,
`MaxUint64` must error. BSON encoding is per-document, so one out-of-range field
fails the whole incident, not just that field. See
[CDR types are storage shapes, not just wire shapes](cdr-types-cross-bson.md) for
the two outages this rule comes from.

## Tags

Every new field names its `bson` key in camelCase, matching its `json` name and
matching the `frameId` / `traceId` tags already on these two structs. That is the
opposite of the rule for the `cdr` types, where naming a key would *rename* a
stored field. These fields are new, so there is no stored value to orphan, and
the surrounding struct already uses the camelCase convention.

## Tests

`armotypes/runtimeincidents_trace_test.go`:

- **JSON round trip** of a fully populated alert — catches a wrong or missing
  `json` tag.
- **BSON round trip** of the same — reaches the storage path a JSON test cannot,
  and surfaces duplicated inline keys, which the driver reports as an error where
  `encoding/json` silently drops both fields.
- **Stored-key pins** for every new field, so a later tag rename fails loudly.
- **Backward compatibility** — an alert that sets none of the new fields marshals
  byte-identically to a golden captured from `main`. Do not regenerate that golden
  to make the test pass; a diff against it means an existing consumer's payload
  changed.
- **Forward compatibility** — a new-shape trace decodes into a copy of the old
  `Trace` shape without erroring.

## Consumers

`private-node-agent` produces these fields. `event-ingester-service`,
`config-service` and `cadashboardbe` carry and store them. Each bumps
independently; until it does, it sees exactly the payload it saw before.
