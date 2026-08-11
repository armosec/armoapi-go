---
type: reality
status: active
owner: alonliwsky
scope: repo
related_code:
  - armotypes/networkstream.go
  - armotypes/networkstream_test.go
---

# Network Stream Workload Identity — host and ECS

Lets the backend tell one host or ECS service apart from another on the network
stream. This library change defines the wire contract only; populating it lives in
private-node-agent and reading it in event-ingester-service.

Sibling of [network-stream-process-attribution.md](network-stream-process-attribution.md),
which covers *which process* opened a connection. This covers *which workload*.

## The problem

`NetworkStreamEntity` described an entity as a `Kind` plus one inlined
`NetworkStreamEntityContainer` — container name/ID, pod namespace, pod name,
workload name/kind. Every field is Kubernetes-shaped, because only Kubernetes
workloads emitted the stream. `NetworkStreamEntityKindHost` existed as a `Kind`
constant with **no fields behind it**: a half-built extension point, so a host
entity carried no typed identity at all.

The reputation consumer builds an incident's dedup key as a hash of
`customerGUID | cluster | workloadKey | endpoint | ruleID`, where `workloadKey`
falls back to `pod/<namespace>/<podName>` when the workload fields are empty. A
host entity therefore hashes on the constant `pod//`, and because the host route's
envelope carries an **empty cluster** by design, every host in an account collapses
onto **one incident per endpoint**.

Scope honestly: standalone host and ECS agents do not export the stream over HTTP
at all yet (node-agent gates the send on `KubernetesMode`), so for them this is a
bug being designed out ahead of the sender work rather than one firing in
production. The live instance of it is the node's own host entity inside a
Kubernetes cluster, which is created with a `Kind` and nothing else — so today all
*nodes in a cluster* already share one incident per endpoint.

### Why not just read the `Entities` map key?

The key already holds something host-shaped (node name for hosts, container ID for
containers), so it looks like identity is free. It is not usable:

- **Both consumers discard it.** They iterate values only — `for _, entity := range
  ...Entities` in the reputation ingester and in postgres-connector's traffic
  writer. Using the key means changing consumer code, which is what the bridge
  below exists to avoid.
- **It is an undocumented producer-internal storage key**, typed as
  `map[string]NetworkStreamEntity` with no contract on its contents, and it is
  per-entity-kind polymorphic. Promoting it to identity would freeze an
  implementation detail.
- **It does not carry what is needed.** A node name is not the stable host ID, and
  ECS needs cluster, task family, revision and service — none of which fit one key.

## What this adds

Two structs, inlined into `NetworkStreamEntity` beside the container struct, the
way `RuntimeAlertECSDetails` sits beside `RuntimeAlertK8sDetails` in the alert
model. (Not character-for-character: the alert model writes `bson:"inline"` without
the leading comma. Both spellings do inline, but the form used here is the correct
one — don't copy the alert model's.)

| Struct | Fields |
|---|---|
| `NetworkStreamEntityECS` | `ecsClusterName`, `clusterArn`, `serviceName`, `taskFamily`, `taskRevision`, `taskArn`, `ecsTaskID`, `launchType` |
| `NetworkStreamEntityHost` | `hostID`, `hostName` |

Nothing existing changed. No field was renamed or removed, no new `Kind` value was
added, and `NetworkStreamEvent` and the process-identity fields were not touched.

## The bridge — and why nothing waits on a consumer change

New fields change nothing until a consumer reads them, and the consumer
(`event-ingester-service`) is owned elsewhere. That would make this a blocked
change — if the new fields were the *only* carrier of identity. They are not.

> **Senders populate both.** They emit the typed structs **and** keep filling
> `WorkloadKind` / `WorkloadName` / `PodNamespace` on the container struct.

The reputation consumer's key builder already branches on
`workloadKind != "" || workloadName != ""` before falling back to `pod//`, so
filling those fields yields correct per-workload keys **on day one, with zero
consumer changes**. The bridged values follow the conventions the shipped host and
ECS workload resolvers already use, so streams and alerts describe one workload the
same way:

| Slot | Kubernetes | Host | ECS |
|---|---|---|---|
| entity `Kind` | `container` | `host` | **`container`** |
| `WorkloadKind` | `Deployment`, … | `Host` | `ECSService`, or `ECSTask` with no service |
| `WorkloadName` | workload name | host ID | `<service>-<taskFamily>-<revision>`, or `<taskFamily>-<revision>` |
| `PodNamespace` | pod namespace | `"host"` | service name, else ECS cluster name |
| `PodName` | pod name | *empty* ⚠️ | *empty* ⚠️ |

⚠️ **`PodName` is the one row that deliberately diverges from the shipped
resolvers.** The host/ECS *alert* path sets `PodName` to the workload name (the host
ID), whereas the stream leaves it empty — see the classifier warning below. Every
other row matches the shipped resolvers exactly.

**Precedence rule: typed fields win where present; the Kubernetes-shaped fields are
the compatibility bridge.** A consumer taught to read the typed structs should
prefer them and fall back to the bridge. **Do not remove the bridge** — removing it
regresses every consumer that has not been taught yet.

Teaching the consumer to read the typed structs is therefore an *improvement*, not
a dependency: it makes the keys semantically clean (`host/<hostID>` rather than the
`workload/Host/host/<hostID>` that the bridge row above produces) and lets ECS
alerts be labelled with the ECS platform instead of host-agent. It is scheduled
work, not a blocker.

Cost of populating both: a few short strings per entity, against the ~2 KB process
trees the payload already carries.

⚠️ **Never put a hostname in `PodName`.** The alert platform classifier reads an
empty `PodName` as "not Kubernetes"; filling it misclassifies the alert as a
Kubernetes alert.

## Why ECS entities keep `Kind: container`

This is the one place the design departs from "a new platform gets a new kind".

The traffic-view writer **drops every entity whose kind is not `container`**
(`postgres-connector` `services/networktraffic/service.go`,
`createNetworkTrafficEvent`, which returns early on
`entity.Kind != NetworkStreamEntityKindContainer`). An ECS workload *is* a
container, so minting an `ecs` kind would silently delete every ECS row from the
traffic view — a regression, in a repo that is out of scope to change.

So ECS entities are `container`-kind and carry `NetworkStreamEntityECS`; host
entities keep the existing `host` kind and carry `NetworkStreamEntityHost`. The
kind constant block carries a comment saying so, because that block is exactly
where someone will reach for a new kind.

Consequence to know: because ECS entities pass that filter, ECS rows **will start
being written** to `network_traffic_events` with ECS workload names once agents
stream. Whether they are then *visible* is a separate question — the read path
filters on an exact three-way match of `workload_name` + `workload_kind` +
`workload_namespace` and rejects empty values, so a caller has to query with
`workload_kind = "ECSService"` and the ECS service name as the namespace — or, for
a standalone task, `workload_kind = "ECSTask"` with the ECS **cluster** name as the
namespace, per the bridge table above. A wrong namespace returns zero rows rather
than an error, so missing the standalone-task shape looks like "no ECS traffic". Host
entities are not written at all — host streams feed reputation only, unless that
one-line kind check is widened. (DNS-protocol events are dropped for every
platform, kind notwithstanding.)

## Fargate needs no schema change

ECS-on-EC2 and ECS-on-Fargate share one struct and differ only by `launchType`
(here `"EC2"` and `"FARGATE"` — examples, not a closed set: the field is an
unvalidated passthrough and ECS Anywhere reports `"EXTERNAL"`), matching how
the alert model already distinguishes them. A Fargate sensor that later learns to
emit the stream needs **no schema negotiation** — only sensor-side capture work.
Nothing in `NetworkStreamEntityECS` is EC2-specific.

## Field-shape notes

- **Names and JSON keys mirror `RuntimeAlertECSDetails`** (`clusterArn`, `taskArn` —
  not `clusterARN`), so the stream and the alert cannot drift on spelling. The set
  is a deliberate **subset** — see the next two bullets for what is left out and
  why.
- **`taskRevision` and `ecsTaskID` have no alert-model counterpart.** The revision is
  needed because the bridged `WorkloadName` is `<service>-<taskFamily>-<revision>`.
  `ecsTaskID` is the stable per-task identifier: the short form of the identity
  `taskArn` carries region-qualified, so prefer it as a grouping key and keep
  `taskArn` for addressing the task in an AWS API.
- **`ecsContainerName` and `ecsContainerID` are omitted on purpose**, not for lack of
  data. Their values already reach the wire — the container name via the bridge
  `containerName` field, the container ID as the `Entities` map key. Typed
  duplicates would give one value two sources of truth.
- **`taskDefinitionArn`, `containerArn`, `containerInstance`, `availabilityZone`,
  `ecsImage` and `ecsImageDigest` are omitted** because nothing needs them to
  identify a workload. They can be added later under the same mirroring rule.
- **`ecsTaskID` carries the ECS prefix** where its `taskFamily`/`taskRevision`/
  `taskArn` neighbours do not. Inlining means every key here shares one flat
  namespace with the container and host structs, and a bare `taskID` is exactly the
  generic key a future non-ECS concept would collide with — a collision that costs
  both fields silently (see the hazards section).
- **`taskRevision` is a string**, because that is what the ECS task metadata
  endpoint returns and what the sensor already holds. Making it numeric would only
  move a parse onto the producer.
- **`hostID`, not `hostId`.** Matches the closest precedent in this package —
  `HostPathInfo` pairs `hostID` with `hostName` — and the file-local `containerID`.
  Note `identifiers/designators.go` uses `hostId` for the *designator attribute*
  key; that is a different layer and does not constrain this one.
- **`hostID` is the identity, `hostName` is display.** Group by `hostID`; a hostname
  can change under a host. But `hostID` is deliberately specified as **opaque** —
  the shipped host agent resolves it to the cloud instance ID, falling back to
  `/etc/machine-id`, then the hostname, and it can end up a non-unique sentinel
  when all three fail. Consumers must not join it against any one of those sources.
  ⚠️ **Open cross-repo item:** the host alert path derives its host ID from
  `cloudMetadata.InstanceID` while the plan's mapping named `ResolveHostID()`
  (machine-id first). If the sender picks the second, stream incidents and host
  alerts will not correlate on host. The sender work package owns that choice; this
  contract accommodates either.
- **Value-embedded, all fields `omitempty`.** An empty identity struct contributes
  zero keys, so a Kubernetes entity serialises exactly as before — same keys, and
  same order too, since promoted fields emit in embed order and both new embeds sit
  after the container embed. `TestNetworkStreamEntityIdentityIsFlattened` pins the
  key *set* per platform; it does not assert byte order.

## Two hazards the tests exist to catch

**`json:",inline"` is not an `encoding/json` feature.** It is a Kubernetes
convention; `encoding/json` sees an empty tag name on an anonymous field and
promotes its fields, which is what actually flattens these structs. The tag is
carried for consistency with the file, but the flattening is promotion, and
`TestNetworkStreamEntityIdentityIsFlattened` pins the resulting key set per
platform. A nested identity object would be a silent break — consumers read these
keys at the entity's top level.

**A duplicate JSON key across two promoted structs makes `encoding/json` drop
both fields**, with no error. Three inlined structs at the same depth make that a
live risk on every future field addition, so
`TestNetworkStreamEntityIdentityKeysAreUnique` reflects over every field of all
three, assigns a unique marker, and asserts each reaches the wire under its own
key. Verified to fail (with the offending key named) when a colliding field is
added.

## bson tags

The new fields carry `bson` tags per `AGENTS.md` ("Both must be maintained
together"), and the embedded fields carry `bson:",inline"` so BSON agrees with JSON
on their keys.

The pre-existing `NetworkStreamEntityContainer` embed has no bson tags, so it lands
as a **nested, default-lowercased subdocument**: a Mongo query needs
`networkstreamentitycontainer.workloadname` where JSONB needs `workloadName`. This
change deliberately does not fix that — doing so means editing an existing field's
tags, which is exactly the backward-incompatible edit this change is scoped to
avoid. It is pre-existing debt shared with the rest of the file and is not tracked
under a ticket. Both halves are
pinned by `TestNetworkStreamEntityIdentityBSONRoundTrip` so neither changes by
accident. As with process attribution, **no BSON path consumes the network stream
today**: `postgres-connector` writes via gorm, the stream ingester has no mongo
reference, and the reputation path copies fields into `RuntimeAlert`.

## Why additive fields and not a kind bump

Unchanged from the process-attribution analysis: both consumers decode with plain
`json.Unmarshal` and reject a message only on a wrong CRD `kind` string or
malformed JSON, so `omitempty` additive fields are invisible to un-upgraded
consumers. Bumping `kubescape.io/v1/networkstreams` to a v2 would force a
coordinated consumer rollout for no benefit.

No version marker is added here. Unlike process attribution, absence of these
fields is not ambiguous: a sensor that predates them sends no identity, and a
sensor that has them always has *something* to send for a host or ECS task. There
is no "ran and found nothing" state to distinguish.

## Known degradations, accepted

- **No `WLID` on stream-derived incidents, for any platform.** The reputation
  ingester never builds one — `Wlid` and `Spiffe` do not appear in that package at
  all — so a WLID is not something this contract can supply or withhold. (The
  separate host/ECS *alert* path in private-node-agent does build one, using a
  `cluster-unknown` pseudo-cluster for hosts. That is a different producer and out
  of scope here.)
- **ECS reputation alerts classify as host-agent alerts** until the consumer is
  taught to populate `RuntimeAlertECSDetails` from these fields. Functional, just
  mislabelled.
- **Host traffic-view rows are absent** — see the kind-filter section above.
