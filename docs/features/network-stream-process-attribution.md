---
type: reality
status: active
owner: alonliwsky
scope: repo
related_code:
  - armotypes/networkstream.go
  - armotypes/networkstream_test.go
---

# Network Stream Process Attribution — shared model

Lets the backend answer "which process opened this connection?" for events on the
network stream. This library change defines the wire contract only; populating it
lives in node-agent and consuming it in event-ingester-service /
postgres-connector.

## The flow this sits on

node-agent builds `armotypes.GenericCRD[armotypes.NetworkStream]` and POSTs it.
It reaches the backend on Pulsar topic
`persistent://armo/internal/network-stream-v1`, wrapped in a synchronizer
envelope whose `patch` field is **base64-encoded JSON**. Two consumers read it:

| Consumer | Path | Sink |
|---|---|---|
| visibility | `event-ingester-service/ingesters/network_stream_ingester/` | Postgres |
| network reputation | `event-ingester-service/ingesters/network_reputation_ingester/` | `RuntimeAlert` (field-by-field copy) |
| upsert | `postgres-connector/services/networktraffic/` | Postgres (gorm) |

## What this adds

1. **`ProcessRef`** — `{PID uint32, StartTime uint64}`, text-marshalled to the
   single string `"<pid>/<startTime>"`.
2. **`NetworkStream.Processes map[ProcessRef]*ProcessTree`** — message-scoped,
   one tree per process.
3. **`NetworkStreamEvent.Process *ProcessRef`** — per-connection pointer into
   that map.
4. **`NetworkStream.ProcessAttributionVersion int`** + const
   `NetworkStreamProcessAttributionV1 = 1`.

The consumer join is `stream.Processes[*event.Process]`. The same `ProcessRef`
type is both the map key and the per-connection value, so the two sides cannot
disagree on key format.

## Why the tree is shipped once per process

Process trees dominate the payload cost, not the number of map entries. They are
chain-shaped and measure ~2.0 KB serialised at the median, ~5.1 KB at p90, of
which `cmdline` alone is ~40% and is **unbounded**. Because the synchronizer
`patch` field is base64, every byte added costs **×1.333** on the wire.

Inlining a tree per connection is the pre-existing
`NetworkStreamEvent.ProcessTree` shape, and it is exactly why the sensor nils
that field out on every event before sending. Deduplicating per process is what
makes attribution affordable: a worst-case observed message (109 container
entities, 283 connections, 142 KB on the wire) lands near 970 KB with full
attribution, against a 5 MiB broker limit — roughly 5.7× headroom in connection
count.

`NetworkStreamEvent.ProcessTree` is **retained unchanged**. An in-process
notification channel in node-agent feeds a separate host-network sensor from the
pre-nil struct, so removing or repurposing it would break that consumer.

## Why the map is scoped per message, not per container

`ProcessRef.PID` is a **host-namespace** pid. The sensor's eBPF network gadget
populates it from the socket enricher's `pid_tgid >> 32`
(`bpf_get_current_pid_tgid()` / `task->tgid`); there is no pid-namespace
translation on that path and no container-namespace pid is exposed at all. Host
pids are node-unique and a stream message covers exactly one node, so keys are
already unique message-wide. Per-container scoping would duplicate every tree
shared between containers and buy nothing.

## Why StartTime is ticks, not a `time.Time`

`StartTime` is CLOCK_BOOTTIME in clock ticks since boot (USER_HZ = 100, 10 ms
per tick).

- The only start-time source with full process coverage is `/proc/<pid>/stat`
  field 22, which is denominated in clock ticks.
- Converting to wall-clock needs `/proc/stat` `btime`, which has **whole-second
  resolution only**. A `time.Time` on the wire would therefore be lossy by up to
  a second — for a value whose entire job is to disambiguate identity, that is
  the wrong trade.
- Ticks are exact and need no boot-time lookup. 10 ms granularity is safe for
  pid-reuse disambiguation by 2–5 orders of magnitude.

`StartTime` is 64-bit because ticks outgrow `uint32` after ~497 days of uptime at
USER_HZ = 100.

## Why additive fields and not a kind bump

Both consumers decode with plain `json.Unmarshal`, and `DisallowUnknownFields`
appears nowhere in either. They reject a message only on a wrong `kind` string
or malformed JSON. So `omitempty` additive fields are invisible to un-upgraded
consumers, and bumping the kind string to a v2 would have forced a coordinated
consumer rollout for no benefit.

## Why an explicit version marker

`ProcessAttributionVersion` absent/zero means *"this sensor predates process
attribution"* — **not** *"attribution ran and produced nothing"*. An upgraded
sensor can legitimately emit an empty `Processes` map when every connection in an
interval was unattributable, so `len(Processes) == 0` is not a usable proxy for
sensor capability.

This marker is new to the repo. Versioning here is otherwise done by CRD kind
string (`kubescape.io/v1/networkstreams`) plus topic name, and the established
pattern for additive changes is `omitempty` plus a documented consumer tolerance
(see `HostCVEOnFinishedMessage` in `armotypes/hot_cve.go`). That pattern alone is
insufficient here precisely because absence-of-data and absence-of-feature are
different states. The repo also auto-tags `v0.0.<github.run_number>` on push to
`main`, so the module version carries no semver signal a consumer could branch
on.

## Field shape and naming trade-offs

- **One text-marshalled field beats two flat scalars.** `,"process":"12345/91827364"`
  is 27 bytes against 33 for `,"pid":12345,"startTime":91827364` — one JSON key
  instead of two, ~6 bytes × 283 connections ≈ 2.3 KB base64 saved per
  worst-case message — *and* it is type-safe and identical to the map key.
- **`Process` is a pointer** so `omitempty` drops it. `encoding/json` does not
  omit zero structs; an always-present `"process":"0/0"` on every unattributed
  connection is the per-connection waste this design exists to avoid.
- **The JSON key is `process`, not `proc`.** Abbreviating saves ~1.1 KB base64
  per worst-case message (0.1%) at the cost of being the only abbreviated key in
  a file whose tags are all spelled-out camelCase (`ipAddress`, `podNamespace`,
  `workloadKind`).
- **The separator is `/`, not `CommPID`'s U+241F (`␟`).** `CommPID` needs an
  exotic separator because `Comm` is free-form and may contain `/` or `:`. Both
  `ProcessRef` components are decimal integers, so no escaping hazard exists;
  `/` matches the composite `"<ip>/<port>/<protocol>"` keys already used by
  `Inbound`/`Outbound` in the same file, and costs 1 byte rather than 3.

## bson tags

The new fields carry `bson` tags even though `armotypes/networkstream.go` is the
only file in `armotypes` with none, because:

- `AGENTS.md` is explicit: "Both must be maintained together."
- Empirically **no BSON path exists today**. `postgres-connector` writes via
  gorm; `network_stream_ingester` has no bson/mongo reference; the reputation
  path copies fields into `RuntimeAlert` rather than embedding a `NetworkStream*`
  type; `cadashboardbe` does not reference these types at all.
- Tags cost nothing at runtime and pre-empt a real trap: mongo-driver's default
  encoder would name the field `processattributionversion`, diverging from the
  JSON key.
- The identical shape already has bson tags in this repo —
  `Process.ChildrenMap map[CommPID]*Process` carries `bson:"childrenMap,omitempty"`.
  Verified that mongo-driver honours `encoding.TextMarshaler`/`TextUnmarshaler`
  for struct map keys, so `Processes` round-trips through bson with the same
  composite key as JSON (`TestProcessRefBSONRoundTrip`).

The remaining structs in this file still lack bson tags; that is pre-existing and
tracked separately.

## Known defect, not fixed here

`NetworkStreamEvent.String()` is broken: `string(e.Port)` converts an `int32` to
a rune rather than a decimal string (port 8080 becomes one garbage character),
and `string(e.IPAddress)` is a no-op on a `string`. No callers were found, but
the search was not exhaustive. Tracked as a separate change to keep this schema
addition reviewable.
