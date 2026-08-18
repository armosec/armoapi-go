---
type: reality
status: active
owner: alonliwsky
scope: repo
related_code:
  - armotypes/networkstream.go
  - armotypes/networkstream_test.go
---

# Network Stream Cloud Metadata

Lets the backend place the host that sent a network stream — cloud provider,
region, account and instance. This library change defines the wire contract only;
populating it lives in private-node-agent and reading it in event-ingester-service.

Third sibling of
[network-stream-workload-identity.md](network-stream-workload-identity.md) (*which
workload*) and
[network-stream-process-attribution.md](network-stream-process-attribution.md)
(*which process*). This one covers *which machine*.

## The problem

Host reputation alerts never open an incident. Every one of them is acknowledged
and discarded.

The runtime incident ingester scopes host alert policies by cloud provider, region
and instance. When it cannot place the sending host it refuses the alert, in
`GetAlertPolicies`:

```go
if alert.AlertSourcePlatform == armotypes.AlertSourcePlatformHostAgent && cloudMetadata == nil {
    return false, nil, false, utils.NewAckError(...)
}
```

The failure is an ack error, so nothing is retried and no customer ever sees it.

The Kubernetes-shaped path does not hit this. The node agent's HTTP exporter puts
cloud metadata on the alert envelope it builds, so alerts sent straight from an
agent arrive placed. Reputation alerts take a different route: they are derived
from a network stream by the reputation ingester, which builds a **fresh** alert
envelope from the stream. The stream had no slot for cloud metadata, so that
envelope could not carry any — the data stopped at the agent.

The agent already resolves the whole struct. `cmd/host/main.go` reads it from IMDS,
falling back to machine-id and hostname, and already sends the host's identity on
the stream as `NetworkStreamEntityHost{HostID, HostName}`, where `HostID` resolves
to the instance ID. It simply never sent the rest of what it had.

## What this adds

One field on `NetworkStream`:

```go
CloudMetadata *CloudMetadata `json:"cloudMetadata,omitempty" bson:"cloudMetadata,omitempty"`
```

Nothing existing changed. No field was renamed or removed, and no entity struct was
touched.

`CloudMetadata` is the type this package already defines in `runtimeincidents.go` —
the same one the consumer's gate takes. Note it is **not** `armotypes/cdr`'s
`CloudMetadata`, a different two-field type for cloud audit events. Reusing the
consumer's own type means there is no mapping step and no second schema to drift.

## Why top level, not per entity

Cloud metadata describes the machine that sent the message, not any one workload on
it. One sensor process reports for one host, so every entity in a message shares
it. Per-entity would repeat identical data once per entity and invite the question
of what it means when two entities disagree — a state no producer can create.

This differs from workload identity, which is genuinely per entity, and that is the
line to hold on future additions: what varies per workload goes on the entity, what
describes the sender goes here.

## Why a pointer

`encoding/json` does not honour `omitempty` on a struct. A value field would add
`"cloudMetadata":{}` to every Kubernetes payload that has never carried one,
changing bytes on a path this change has no business touching.

Nil therefore means "the sender does not report this", which is not the same as a
sender reporting that it knows nothing.
`TestNetworkStreamCloudMetadataOmittedWhenAbsent` pins it, and the pre-existing
`TestNetworkStreamProcessAttributionOmitted` — which asserts the top-level key set
is exactly `["entities"]` — fails if a later change makes the field unconditional.

## Two alternatives, rejected

**Relax the gate.** Let host alerts through without cloud metadata. Rejected: the
gate enforces policy scoping shared by every host producer, and weakening it for
one producer degrades scoping for all of them.

**Carry the fields as Pulsar message properties.** About ten untyped string
properties on the message. Rejected: it diverges from the Kubernetes path, which
carries this in the payload, and the properties would be meaningful to exactly one
consumer.

## What this does not fix

Transmitting the field is necessary but not sufficient on a host with no cloud
provider.

The consumer does not read cloud metadata off the envelope directly. It converts it
to designator attributes and reconstructs it, and that round trip is lossy in two
ways.

It drops empty values, because `addCloudMetadata` writes only non-empty attributes
and `ExtractCloudMetadataFromDesignators` returns nil when the cloud-provider
attribute is absent. So metadata whose `Provider` is empty still fails the gate.

It also drops seven fields unconditionally — `MachineID`, `ClusterName`, `OrgID`,
`PrivateIP`, `PrivateIPs`, `PublicIPs` and `Services` — because no designator
attribute carries them. They reach the consumer in the payload and are readable
there, but they do not survive into the reconstructed struct the gate sees. Sending
the whole type is still right: the alternative freezes today's ten-attribute
consumer behaviour into the wire contract.

The agent leaves `Provider` empty when IMDS fails or is absent. That fallback is
not otherwise empty — it sets the hostname, the machine ID, and defaults `HostType`
to `other`, and it overwrites `InstanceID` with the machine ID whenever the provider
is empty or `other`. It is specifically the provider that stays unset, and the
provider is the one field the reconstruction requires. On such a host, reputation
alerts keep being acked.

This is **pre-existing and not specific to the stream path**: alerts sent straight
from a host agent take the identical designator round trip and are acked the same
way. This change reaches parity with the direct alert path; closing the gap for
providerless hosts is a separate decision, because it changes how every host alert
is scoped rather than how reputation alerts are carried.

Practical consequence: verify on a host where IMDS resolves. On one where it does
not, a correct implementation still produces no incident.

## bson tags

The field carries a `bson` tag per `AGENTS.md` ("Both must be maintained
together"), pinned by `TestNetworkStreamCloudMetadataBSONRoundTrip`. As with the two
sibling changes, **no BSON path consumes the network stream today** —
`postgres-connector` writes via gorm, the stream ingester has no mongo reference,
and the reputation path copies fields into `RuntimeAlert`.

## Why additive and not a kind bump

Unchanged from the process-attribution analysis: both consumers decode with plain
`json.Unmarshal` and reject a message only on a wrong CRD `kind` string or malformed
JSON, so an `omitempty` additive field is invisible to un-upgraded consumers.

No version marker is added. Absence is not ambiguous here: a sensor that predates
the field sends nothing, and a sensor that has it always has something to send. There
is no "ran and found nothing" state to distinguish, so this follows workload
identity rather than process attribution.

## Note for the sender

The sender splits an oversized payload when the receiver returns 413, and rebuilds
each half. A top-level field must be copied on **both** the build and the rebuild
paths. Copying it in only one place loses cloud metadata on split payloads — an
intermittent failure that appears only above a size threshold.
