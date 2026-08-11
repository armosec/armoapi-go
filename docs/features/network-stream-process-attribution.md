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

For "which workload opened this connection?" on hosts and ECS tasks, see the
sibling doc [network-stream-workload-identity.md](network-stream-workload-identity.md).

## The flow this sits on

node-agent builds `armotypes.GenericCRD[armotypes.NetworkStream]` and POSTs it.
It reaches the backend on Pulsar topic
`persistent://armo/internal/network-stream-v1`, wrapped in a synchronizer
envelope whose `patch` field is **base64-encoded JSON**. Consumers:

| Consumer | Path | Sink | On decode error |
|---|---|---|---|
| visibility | `event-ingester-service/ingesters/network_stream_ingester/` | Postgres | **nack** → retry → DLQ |
| network reputation | `event-ingester-service/ingesters/network_reputation_ingester/` | `RuntimeAlert` (field-by-field copy) | log + ack (dropped) |
| upsert | `postgres-connector/services/networktraffic/` | Postgres (gorm) | — |

## What this adds

1. **`ProcessRef`** — `{PID uint32, StartTimeNs uint64}`, text-marshalled to
   the single string `"<pid>/<startTimeNs>"`, plus `String()`.
2. **`NetworkStream.Processes map[ProcessRef]*ProcessTree`** — message-scoped,
   one tree per process.
3. **`NetworkStreamEvent.ProcessRef *ProcessRef`** — per-connection pointer into
   that map.
4. **`NetworkStream.ProcessTreeFor(event)`** — the supported way to resolve a
   connection's tree.
5. **`NetworkStream.ProcessAttributionVersion int`** + const
   `NetworkStreamProcessAttributionV1 = 1`.

The same `ProcessRef` type is both the map key and the per-connection value, so
the two sides cannot disagree on key format.

## Resolve with `ProcessTreeFor`, never a bare map index

`stream.Processes[*event.ProcessRef]` has three sharp edges:

- it **panics** when `event.ProcessRef` is nil — every event from a pre-attribution
  sensor, and any connection the sensor could not attribute;
- it cannot distinguish a missing entry from a present-but-nil one;
- `{"processes":{"1/2":null}}` decodes to exactly that present-but-nil case, so
  the index returns `ok == true` with a nil tree and the next field access
  panics.

`ProcessTreeFor` collapses unattributed, dangling and nil-valued into a single
nil return. It does **not** fall back to the per-event `ProcessTree` field — the
two have different lifecycles.

## Why the tree is shipped once per process

Process trees dominate the payload cost, not the number of map entries. They are
chain-shaped, and measurement during design put them at ~2.0 KB serialised at the
median and ~5.1 KB at p90, of which `cmdline` alone is ~40% and is **unbounded**.
Because the synchronizer `patch` field is base64, every byte added costs
**×1.333** on the wire.

Inlining a tree per connection is the pre-existing
`NetworkStreamEvent.ProcessTree` shape, and it is exactly why the sensor nils
that field out on every event before sending. Deduplicating per process is what
makes attribution affordable at all.

**Headroom is comfortable but set by the tail, not the median.** Against the
5 MiB broker limit, the largest message observed during design (109 container
entities, 283 connections, 142 KB on the wire) comes to roughly 900 KB at the
median tree size, but roughly 2 MB at p90. Since `cmdline` is unbounded, the
right lever if this ever approaches the limit is a sensor-side cap on tree or
cmdline size — not dropping attribution. The tree-size distribution and the
observed message shape come from measurement during design and are not
reproducible from this repo.

`NetworkStreamEvent.ProcessTree` is **retained unchanged**. `private-node-agent`
wires a non-nil notification channel in all three of its mains
(`cmd/main.go`, `cmd/host/main.go`, `cmd/ecs/main.go`) and its host network
sensor reads `outbound.ProcessTree` directly off those events
(`pkg/hostnetworksensor/v1/hostnetworksensor.go`), exporting `.PID`, `.Comm`,
`.Uid` and friends. The OSS `node-agent/cmd/main.go` passes `nil` for that
channel, so searching only the OSS repo makes the field look dead. It is not.

## Why the map is scoped per message, not per container

`ProcessRef.PID` is a **host-namespace** pid. The sensor's eBPF network gadget
populates it from the socket enricher's `pid_tgid >> 32`
(`bpf_get_current_pid_tgid()` / `task->tgid`); there is no pid-namespace
translation on that path and no container-namespace pid is exposed at all. Host
pids are node-unique and a stream message covers exactly one node, so keys are
already unique message-wide. Per-container scoping would duplicate every tree
shared between containers and buy nothing.

## Why StartTimeNs is boot-relative nanoseconds, not a `time.Time`

`StartTimeNs` is the process's **creation** time as nanoseconds since boot
(CLOCK_BOOTTIME).

The start time is in the identity **purely to survive pid reuse**, which is the
mitigation the GA design itself names
(`projects/epics/network-reputation-ga/network-reputation-consolidation-design.md`
§13: "join reliability across process restarts / pid reuse (mitigated by
`startTime` in the key)"). Nothing renders this value. Anything needing a
human-facing process start time reads `Process.StartTime` off the resolved tree,
which is a wall-clock `time.Time`.

**Boot-relative, not wall-clock**, because converting to wall-clock needs
`/proc/stat btime`, which has whole-second resolution only. A `time.Time` on the
wire would be lossy by up to a second — wrong for a value whose entire job is
disambiguation. A boot-relative integer needs no boot-time lookup and is exact.

**Nanoseconds rather than clock ticks**, because that objection rules out a
*timestamp*, not a unit: both ticks and boot-ns are boot-relative integers, so it
says nothing about which. Between them:

- **Conversion only loses information in the ticks direction.** procfs ticks → ns
  is an exact ×10⁷; a boot-ns source → ticks is a lossy division. Choosing ns
  means no source ever has to discard anything.
- **One unit for one concept.** Two units either side of a serialisation boundary
  is a bug factory.
- **The cost is noise** — exactly 7 bytes per connection, ≈2.6 KB base64 at the
  largest observed message, against a 5 MiB limit. Exactly 7 because ns is
  ticks × 10⁷, which appends precisely 7 decimal digits at every magnitude.

> ⚠️ **Nanoseconds are the unit, not the resolution.** Do not read this as a
> precise timestamp. The only start-time source with full process coverage is
> `/proc/<pid>/stat` field 22, denominated in clock ticks (USER_HZ = 100, 10 ms),
> so procfs-derived values are exact multiples of 10,000,000 ns and carry far less
> precision than the unit implies. This does not affect correctness — 10 ms is
> safe for pid-reuse disambiguation by several orders of magnitude — but someone
> will eventually try to use it as a precise timestamp.

⚠️ **Never derive this from an exec event.** `exec` does not create a process, so
an exec-derived value changes when a process re-execs and the identity silently
breaks mid-life.

64-bit is required: boot-ns exceeds `uint32` after ~4.3 seconds of uptime, and
reaches the `uint64` ceiling after ~584 years. A **zero** value means the sensor
could not read the start time; such a ref carries pid-only identity and must not
be compared across messages.

The field name carries the unit because `armotypes.Process.StartTime` — the same
concept on the tree this ref points at — is a wall-clock `time.Time`. The unit is
free on the wire: `ProcessRef` is text-marshalled, so its field names never
appear in JSON at all.

### Producer note: node-agent has the field name, not yet a working value

node-agent already declares `StartTimeNs uint64` — *"Process start time in
nanoseconds for unique identification"* (`pkg/processtree/conversion/types.go:30`)
— so this schema aligns with an existing name and unit. It does **not** align with
a working value, and the sensor-side change cannot simply forward it:

- The exec/fork/exit feeders set it from the **event timestamp**
  (`pkg/processtree/conversion/convert.go`, *"Use event timestamp for
  consistency"*). For an exec event that is the exec instant, not process
  creation — semantically wrong for identity, per the warning above.
- The procfs feeder passes through `procfsEvent.StartTimeNs`, which nothing
  upstream ever assigns, so it is **always zero for procfs-sourced processes** —
  the one source with full coverage. Verified: no `/proc/<pid>/stat` field-22
  parse exists anywhere in node-agent.

Populating it correctly is separately tracked sensor-side work.

## The parser is deliberately permissive

`UnmarshalText` **ignores components beyond the first two** and accepts leading
zeros. This is load-bearing, not sloppiness.

`encoding/json` aborts the **entire enclosing decode** when a `TextUnmarshaler`
map key fails — `decode.go` does `return err` for that case, where an int-keyed
map does `saveError; break` and carries on. Measured: one bad key in
`{"processes":{"BAD":…,"1/2":…}}` discards **both** entries and returns an error.
The visibility consumer nacks on that error, so every connection row for that
node's whole interval is lost.

So a future revision that appends a component to the key format must not be able
to make un-upgraded consumers reject whole messages. Appending is the designed
extension point; **do not make the parser stricter.** The trade-off accepted in
exchange: non-canonical keys for one process collapse to a single entry, last
wins (`TestProcessRefUnmarshalTextAcceptsNonCanonical`).

`UnmarshalText` leaves the receiver untouched on error, but that only reaches
direct callers — `encoding/json` allocates the pointer target before calling it
and leaves it installed on failure, so a caller that ignores an `Unmarshal` error
sees a non-nil, zero-valued ref (`TestNetworkStreamEventMalformedRefInstallsZeroPointer`).
Both consumers do check the error.

## Why additive fields and not a kind bump

Both consumers decode with plain `json.Unmarshal`. There is exactly one
`DisallowUnknownFields()` in all four consumer repos, in an unrelated
cadashboardbe webhook handler; the only non-stdlib JSON in use is `ujson`,
confined to SBOM parsing. Consumers reject a message only on a wrong `kind`
string or malformed JSON. So `omitempty` additive fields are invisible to
un-upgraded consumers, and bumping the kind string to a v2 would have forced a
coordinated consumer rollout for no benefit.

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

Note the marker cannot rescue a key-format break: decode dies before anything can
read it. That is what the permissive parser above is for.

## Field shape and naming trade-offs

- **One text-marshalled field beats two flat scalars.**
  `,"processRef":"12345/918273640000000"` (37 bytes) against
  `,"pid":12345,"startTimeNs":918273640000000` (42 bytes) — one JSON key instead
  of two saves **5 bytes per connection**, ≈1.9 KB base64 on a 283-connection
  message. The saving is exactly 5 regardless of the pid and start-time values,
  since it is purely the difference in key overhead; it does move if either field
  is renamed. And it is type-safe and identical to the map key.
- **The field is a pointer** so `omitempty` drops it. `encoding/json` does not
  omit zero structs; an always-present `"processRef":"0/0"` on every unattributed
  connection is the per-connection waste this design exists to avoid.
- **The field is named `ProcessRef`, not `Process`.** `armotypes.Process` is a
  different, heavily-used type in the same package, and its `StartTime` is a
  wall-clock `time.Time` where `ProcessRef.StartTimeNs` is boot-relative
  nanoseconds — `event.Process` would read as `*armotypes.Process` to anyone who
  had not read this file. Naming the field after its type also matches the
  adjacent `ProcessTree *ProcessTree`. Costs 3 bytes per connection over
  `process`, ≈1.1 KB base64 on a 283-connection message: the cheap side of the
  trade.
- **Keys are spelled out, not abbreviated.** `proc` would save 6 bytes per
  connection over `processRef`, ≈2.3 KB base64, at the cost of being the only
  abbreviated key in a file whose tags are all spelled-out camelCase
  (`ipAddress`, `podNamespace`, `workloadKind`).
- **The separator is `/`, not `CommPID`'s U+241F (`␟`).** `CommPID` needs an
  exotic separator because `Comm` is free-form and may contain `/` or `:`. Every
  `ProcessRef` component is a decimal integer, so no escaping hazard exists; `/`
  matches the composite `"<ip>/<port>/<protocol>"` keys already used by
  `Inbound`/`Outbound` in the same file, and costs 1 byte rather than 3.

## bson tags

The new fields carry `bson` tags even though `armotypes/networkstream.go` is the
only file in `armotypes` with none, because:

- `AGENTS.md` is explicit: "Both must be maintained together."
- The identical shape already has them —
  `Process.ChildrenMap map[CommPID]*Process` carries `bson:"childrenMap,omitempty"`.
- Tags cost nothing at runtime and pre-empt a real trap: mongo-driver's default
  encoder would name the field `processattributionversion`, diverging from JSON.

**JSON and BSON do not agree on the per-connection representation.**
mongo-driver honours `encoding.TextMarshaler` for map **keys only**, not for
values. So `Processes` is keyed by the same composite string in both, but the
event-level ref is a string in JSON and a **subdocument** in BSON:

```
JSON: "processRef":"4242/3141590000000"
BSON: "processRef":{"pid":NumberLong(4242),"startTimeNs":NumberLong(3141590000000)}
```

Both round-trip (`TestProcessRefBSONRoundTrip`), but a Mongo query needs
`processRef.pid` where JSONB needs the string. Consequence for the struct tags:
`ProcessRef`'s inner **`json`** tags are dead weight (MarshalText always wins);
its inner **`bson`** tags are live and are what name that subdocument.

Empirically **no BSON path exists today**: `postgres-connector` writes via gorm;
`network_stream_ingester` has no bson/mongo reference; the reputation path copies
fields into `RuntimeAlert` rather than embedding a `NetworkStream*` type;
`cadashboardbe` does not reference these types at all. The tags are therefore
speculative, and because the rest of the file is untagged, a hypothetical BSON
document would be mixed-case (`processRef` beside default-lowercased `ipaddress`,
`dnsname`, `processtree`, and embedded structs emitted nested rather than
inlined). Tagging the whole file is pre-existing work tracked separately.

## Known defect, not fixed here

`NetworkStreamEvent.String()` is broken: `string(e.Port)` converts an `int32` to
a rune rather than a decimal string (port 8080 becomes one garbage character),
and `string(e.IPAddress)` is a no-op on a `string`. `go vet` cannot catch it
because `rune` *is* `int32`, making the conversion indistinguishable from a
legitimate one. No callers were found, but the search was not exhaustive. Tracked
as a separate change to keep this schema addition reviewable.
