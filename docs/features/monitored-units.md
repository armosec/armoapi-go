---
type: feature
status: active
owner: "@armosec/runtime"
scope: repo
related_code:
  - armotypes/monitoredunits/monitoredunits.go
  - armotypes/ciidentity.go
  - armotypes/runtimeincidents.go
---

# Monitored units message, `workloadLabels` and CI identity

Wire types that a sensor and the backend share, so that the backend can tell which units a
sensor monitored and when, and can link an alert to the CI job that ran in the unit.

Design and rules of the contract:
`shared-designs-and-docs/projects/epics/github-runners-protection/sensor-backend-contract.md`.
This doc describes the types. The contract doc is the source of truth for the behavior of both
sides.

| Type | Package | Sent by | Read by |
|------|---------|---------|---------|
| `Message`, `Entry`, `Unit`, `K8sPodAttributes` | `armotypes/monitoredunits` | private-node-agent | event-ingester-service |
| `RuntimeAlertK8sDetails.WorkloadLabels` | `armotypes` | node-agent | event-ingester-service |
| `CIIdentity`, `GitHubCIIdentity` | `armotypes` | no sensor yet | event-ingester-service |

## Monitored units message (`armotypes/monitoredunits`)

A sensor reports the units it monitors. A unit is a pod today (`UnitKindK8sPod`), and can be
another kind later. The sensor sends a `Message`:

- when monitoring of a unit starts or stops: an `Entry` of type `start` or `stop`;
- periodically, with the IDs of all units it monitors now (`MonitoredUnitIDs`).

The message is JSON. The package has no `bson` tags: how a consumer stores the data is its own
design.

### `Message`

| Field | JSON | Meaning |
|-------|------|---------|
| `SchemaVersion` | `schemaVersion` | `SchemaVersion` (1). Raised only by a change that breaks an existing field |
| `AgentVersion` | `agentVersion` | Version of the sensor |
| `AgentStartedAt` | `agentStartedAt` | When the sensor process started. It identifies the sensor run |
| `ClusterName` | `clusterName` | The Kubernetes cluster of the sensor. Required from a sensor in a cluster. Absent outside a cluster |
| `Node` | `node` | Node name in a cluster, host name elsewhere |
| `CloudMetadata` | `cloudMetadata` | The machine of the sensor: provider, account, region, instance. Optional. The existing `armotypes.CloudMetadata`, so its fields are snake_case as on other messages |
| `SentAt` | `sentAt` | When the message was sent |
| `IntervalSeconds` | `intervalSeconds` | Interval of the periodic list on this sensor |
| `Entries` | `entries` | Start and stop entries. May be empty |
| `MonitoredUnitIDs` | `monitoredUnitIds` | Only in a periodic list. See below |

The message carries the cluster and the machine itself. A consumer behind a queue does not
depend on the HTTP query parameters of the request that delivered the message. `ClusterName` in
the body is the only cluster identity that the consumer uses: a `k8s-pod` entry of a message
without it is dropped, and the receiver rejects a request whose `clusterName` query parameter
differs from it.

**Absent list and empty list are different.** `MonitoredUnitIDs` is a slice with the JSON tag
`omitzero`:

- a nil slice (the field is absent): the message is not a periodic list;
- a non-nil empty slice (`[]`): the sensor is alive and monitors nothing.

`Entries` is always sent as a list. `Message` has a `MarshalJSON` that encodes a nil `Entries`
as `[]`, so a periodic list with no entries does not send `null`.

Use `SetMonitoredUnitIDs` to build a periodic list and `IsPeriodicList` to test for one.

### `Entry` and `Unit`

| Field | In | Meaning |
|-------|----|---------|
| `Type` | both | `EntryTypeStart` or `EntryTypeStop` |
| `Reason` | both | See the reasons below |
| `Unit.Kind`, `Unit.ID` | both | The kind and the ID of the unit. For a pod the ID is the pod UID |
| `Unit.StartedAt` | start | When the unit started |
| `Unit.Attributes` | start | Raw JSON whose shape depends on the kind. For a pod: `K8sPodAttributes` |
| `MonitoredFrom` | start | When the sensor started to monitor the unit |
| `MonitoredUntil` | stop | When the sensor stopped monitoring the unit |
| `CI` | both, optional | See "CI identity" |

A stop entry carries only `Unit.Kind` and `Unit.ID`.

Reasons:

| Constant | Value | Entry | Set by |
|----------|-------|-------|--------|
| `ReasonUnitStarted` | `unit-started` | start | sensor |
| `ReasonAgentStarted` | `agent-started` | start | sensor |
| `ReasonUnitEnded` | `unit-ended` | stop | sensor |
| `ReasonAgentStopping` | `agent-stopping` | stop | sensor |
| `ReasonTimeout` | `timeout` | stop | consumer only. A sensor never sends it |

`NewK8sPodUnit` builds the unit of a start entry. `Unit.K8sPodAttributes()` decodes the
attributes, and returns an error for another kind or for a unit with no attributes (absent, or a
JSON `null`).

### Compatibility rules

- `EntryType`, `Reason` and `UnitKind` are plain strings. Decoding does not reject a value that
  has no constant. A consumer keeps what it does not know.
- `Unit.Attributes` is raw JSON for the same reason: the attributes of an unknown kind pass
  through unchanged.
- A new optional field does not raise `SchemaVersion`.
- The package has no `Validate` method. Which entries a consumer accepts is its own decision.

The examples of the contract are the golden files in `armotypes/monitoredunits/testdata/`. The
tests decode and encode them and compare the JSON.

## `workloadLabels` on the alert

`RuntimeAlertK8sDetails.WorkloadLabels` holds the labels of the pod's owner workload, the one
named in `WorkloadKind` and `WorkloadName`. It is needed when a label that identifies the
workload is on the owner and not on the pod, as for a pod that a Job creates.

| Case | Value | JSON |
|------|-------|------|
| The workload is the pod itself | nil | field absent |
| The owner has labels | the labels | `"workloadLabels": {"k": "v"}` |
| The owner has no labels | empty map | `"workloadLabels": {}` |

The JSON tag is `omitzero`, not `omitempty`, so that the empty map is sent. A consumer uses
this to recognize an old sender with no version check: a workload kind that is not `Pod`, with
no `workloadLabels` field, means that the sender predates the field. A sender must therefore set
a non-nil map whenever the workload is not the pod.

Limits:

- `omitzero` needs Go 1.24 and `encoding/json`. A service that encodes the alert again with
  another JSON library must check that the library keeps the empty map.
- The `bson` tag keeps `omitempty`, so the difference between absent and empty is lost in
  storage. A consumer that needs it reads it from the JSON at ingest.

`K8sPodAttributes.WorkloadLabels` follows the same rule.

## CI identity

`CIIdentity` is the CI identity that a sensor reports about itself. It is optional enrichment:
`RuntimeAlert.CI` (JSON `ci`), next to the Kubernetes and ECS details, and `Entry.CI`.

```json
"ci": { "provider": "github", "github": { "runnerName": "vm-runner-7", "runId": 9912, "attempt": 2, "job": "build" } }
```

- **Each provider has its own block.** The keys of a run differ between CI systems, so the
  fields are not generalized. A new provider adds a `CIProvider` constant and a block next to
  `GitHub`. A consumer that does not know a provider still reads `provider` and ignores the block.
- **Every field of a block is optional.** A sensor sends what it knows.
- **The values are keys for a lookup, not facts.** The reporting workload may run untrusted
  code. A consumer confirms the values with the CI system before it shows them as verified.
- **No sensor sends it yet.** A sensor that reports pods leaves it nil, and the consumer derives
  the identity from the pod labels and `workloadLabels`. It is defined now so that a later
  sensor with no pod labels (for example on a VM runner) needs no change of the contract.

## Consumers

- private-node-agent builds and sends the monitored units message.
- node-agent fills `WorkloadLabels` on the alert.
- event-ingester-service (armosec/backend) reads all three.
