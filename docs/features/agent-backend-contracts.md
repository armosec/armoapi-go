---
type: feature
status: active
owner: eranm
scope: repo
related_code:
  - armotypes/ecs/types.go
  - armotypes/ecs/operator/types.go
  - armotypes/ecs/operator/converters.go
  - armotypes/aisandbox/ai_sandbox_identity_key.go
---

# Agent↔backend contracts: ECS wire types and AI-Sandbox identity keys

## Purpose

Two customer-side agents share types with the backend that used to live in
`github.com/armosec/armosec-infra`:

| Package | Contract | Agent side | Backend side |
|---|---|---|---|
| `armotypes/ecs` | `EcsResourceType`, `TaskType` enums | ecs-operator | backend ECS ingestion |
| `armotypes/ecs/operator` | `EcsResource`, `EcsClusterResource`, the `New*Resource` converters from AWS SDK `ecs/types`, `Unmarshal*`, `ComputeChecksum`, `ClassifyTaskType` | ecs-operator (discovery + state change detection) | backend ECS ingestion |
| `armotypes/aisandbox` | `identity_key` format: `<tag>:<canonical-body>` tag constants and constructors (`K8sServiceAccountKey`, `AWSRoleKey`, `AWSSessionKey`, `BearerSessionKey`, `OpaqueIdentityKey`, `IdentityKeyTag`, …) | private-node-agent (`K8sServiceAccountKey`) | backend AI-Sandbox identity derivation |

`armosec-infra` is being consolidated into the `armosec/backend` monorepo under `internal/`,
and its repo is archived. Go can't import `internal/` packages from outside the monorepo, so
the agents could never get a newer version of these types. They move here because
armoapi-go is already the versioned agent↔backend types module, and both agents already
depend on it.

## Moved verbatim (SUB-8699)

The files are copied byte-for-byte from `armosec-infra` `main` (unchanged since v0.0.594).
The only change is the import of `armosec-infra/ecs` in `ecs/operator`, which now points
at `armoapi-go/armotypes/ecs`, so the JSON wire format is identical to the `armosec-infra`
copies. Compared with the agents' current pins (ecs-operator v0.0.499, private-node-agent
v0.0.592), the only differences are additions: `TaskType`/`ClassifyTaskType` and the
`bearer-session` identity kind. The rest of `armosec-infra/aisandbox` (identity records, display
labels, permissions, observed credentials) is backend-only and did **not** move.

One test in the source test file (`TestIdentityDisplayLabel_BearerSession`) exercises the
backend-only display label, so it stayed behind.

## Rollout

1. This package is released (armoapi-go tag).
2. ecs-operator and private-node-agent swap their `armosec-infra` imports for these packages.
3. After the monorepo cutover, the backend's `internal/armosec-infra/{ecs,aisandbox}` copies
   become type aliases and forwarding funcs to these packages, so there is one source of truth.

Until step 3 lands, a change to any of these types must be made in **both** places.

## Dependencies

`armotypes/ecs/operator` imports `github.com/aws/aws-sdk-go-v2/service/ecs/types`
(the converters take the SDK's `Cluster`, `Service`, `Task` and `TaskDefinition`). armoapi-go
requires `service/ecs` v1.71.0, the lowest version either agent uses, which raises the
module graph's `aws-sdk-go-v2` floor to v1.41.1 and `smithy-go` to v1.24.0. Consumers that
don't import `armotypes/ecs/operator` don't compile the SDK.
