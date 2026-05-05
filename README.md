[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Farmosec%2Farmoapi-go.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Farmosec%2Farmoapi-go?ref=badge_shield)

# armoapi-go

Canonical shared type library for the Armosec platform. Defines domain types, command structures, and data contracts consumed by all backend services, in-cluster agents, and the portal.

## Packages

| Package | Description |
|---------|-------------|
| `apis` | Command dispatch, websocket scan commands, connector interfaces |
| `armotypes` | Core platform domain types (clusters, configs, policies, incidents, risks) |
| `identifiers` | Workload/resource designator system (WLID, SID, PortalDesignator) |
| `containerscan` | Container vulnerability scan report interfaces and implementations |
| `notifications` | Alert channel and collaboration configuration types |
| `broadcastevents` | Platform analytics/audit event types |
| `package_versions` | Multi-format package version parsing and comparison |
| `scanfailure` | Scan failure report types and reason codes |

## Usage

```go
import "github.com/armosec/armoapi-go/armotypes"
import "github.com/armosec/armoapi-go/apis"
import "github.com/armosec/armoapi-go/identifiers"
```

## Documentation

| Document | Description |
|----------|-------------|
| [AGENTS.md](AGENTS.md) | AI agent entry point — repo purpose, conventions, key packages |
| [Architecture](docs/repo/architecture.md) | Package structure, public API surface, versioning |

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Farmosec%2Farmoapi-go.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Farmosec%2Farmoapi-go?ref=badge_large)