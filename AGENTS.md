---
type: agents
status: approved
owner: "@armosec/backend"
scope: repo
---
# AGENTS.md — armoapi-go

> AI agents: this file is your high-signal entry point. Engineers: see [Architecture](docs/repo/architecture.md).

## What this repo is

Canonical shared type library for the Armosec platform. Defines domain types, command structures, and data contracts consumed by all backend services, in-cluster agents, and the portal.

## Services produced by this repo

(library — no services produced)

## Where docs live

- Repo-level (code structure, build, repo ADRs): `docs/repo/`
- Cross-cutting features: `docs/features/`
- Cross-repo projects: `shared-designs-and-docs/projects/`

## Key packages

| Package | Purpose |
|---------|---------|
| `apis` | Command dispatch, websocket scan commands, connector interfaces, pagination |
| `armotypes` | Core domain types: clusters, configs, policies, incidents, risk factors, registries |
| `armotypes/cdr` | Cloud Detection & Response alert types |
| `armotypes/ecs`, `armotypes/ecs/operator` | ECS wire types + AWS SDK converters shared with ecs-operator |
| `armotypes/aisandbox` | AI-Sandbox `identity_key` format + constructors shared with private-node-agent |
| `identifiers` | Workload/resource designator system (WLID, SID, PortalDesignator) |
| `containerscan` | Container vulnerability scan report interfaces and v1 implementation |
| `scanfailure` | Scan failure report types and reason codes |

## Consumers

event-ingester-service, cadashboardbe, config-service, kubevuln, kubescape, node-agent, notification-service, gateway, registry-scanner.

## Conventions

- This repo is the contract with agents and other outside consumers. A type that only
  backend services use belongs in `armosec/backend` (`internal/armosec-infra/apitypes`), not
  here. `notifications`, `broadcastevents`, `package_versions` and `armotypes/aiproviders`
  moved there in SUB-9864 and were removed from this repo in SUB-9866.

- All types use `json` struct tags for REST/Elasticsearch and `bson` tags for MongoDB serialization. Both must be maintained together.
- Interfaces in `containerscan/` decouple consuming services from concrete implementations (v1). Do not add methods to interfaces without checking all consumers.
- The `identifiers` package defines `PortalDesignator` — the universal resource-addressing type. Attribute constants (`AttributeCluster`, `AttributeNamespace`) must stay backward-compatible.
- New types in `armotypes` must include a `PortalBase` embed for consistent metadata (GUID, name, attributes).
- No database connections or runtime dependencies — this is a pure type/utility library.
- Releases are automated via `.github/workflows/release.yaml` on tag push (`v0.0.XXX`).
- High-performance JSON paths use `gojay` marshalling; standard `encoding/json` for everything else.

## When you change behavior

Refresh the matching doc. Pivots become ADRs. Cross-repo work has a project doc in `shared-designs-and-docs`. See `armosec-ai-shared-rules/docs-system/when-to-document.md`.
