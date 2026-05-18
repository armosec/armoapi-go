---
type: catalog
status: approved
owner: "@armosec/backend"
scope: repo
service: armoapi-go
---

# armoapi-go Library Documentation

## Purpose

`armoapi-go` is a shared Go types library that provides domain types, API definitions, and data contracts for the Armosec platform. It serves as the canonical source of truth for data structures exchanged between backend services, in-cluster agents, and the portal.

This is a **pure type/utility library** with no runtime services (no HTTP servers, no message brokers, no database connections). It does depend on serialization libraries (e.g., MongoDB BSON for struct tags) but does not establish any connections at import time.

## Package Inventory

| Package | Description |
|---------|-------------|
| `apis` | Command dispatch types, websocket scan commands, backend connector interfaces, and pagination helpers |
| `armotypes` | Core platform domain types: clusters, configurations, policies, incidents, security risks, registries, vulnerability exceptions |
| `armotypes/cdr` | Cloud Detection & Response alert types (CdrAlert, CloudTrailEvent, AWS-specific structures) |
| `armotypes/common` | Shared runtime sub-types (ProcessEntity, FileEntity) used across detection features |
| `identifiers` | Workload and resource designator system (WLID, SID, PortalDesignator) with parsing utilities |
| `containerscan` | Container vulnerability scan report interfaces defining the contract between scanner and consumers |
| `containerscan/v1` | Concrete implementation of the `containerscan` interfaces |
| `notifications` | Alert channel configuration, collaboration integrations (Jira, Linear, Slack, GitHub), SIEM, and incident notification types |
| `broadcastevents` | Platform analytics and audit event types with factory constructors |
| `package_versions` | Multi-format package version parsing and comparison (semver, APK, deb, RPM, Maven, PEP440, Go modules) |
| `scanfailure` | Scan failure report types, reason codes, and human-readable failure descriptions |

## Import Patterns

Import the library using standard Go modules:

```go
import (
    "github.com/armosec/armoapi-go/apis"
    "github.com/armosec/armoapi-go/armotypes"
    "github.com/armosec/armoapi-go/identifiers"
    "github.com/armosec/armoapi-go/containerscan"
    "github.com/armosec/armoapi-go/notifications"
    "github.com/armosec/armoapi-go/package_versions"
)
```

Add to your project:

```bash
go get github.com/armosec/armoapi-go@latest
```

Pin a specific version:

```bash
go get github.com/armosec/armoapi-go@v0.0.709
```

## Key Usage Patterns

### Resource Designators

The `identifiers.PortalDesignator` type is the universal resource-addressing mechanism:

```go
designator := identifiers.PortalDesignator{
    DesignatorType: identifiers.DesignatorAttributes,
    Attributes: map[string]string{
        identifiers.AttributeCluster:   "my-cluster",
        identifiers.AttributeNamespace: "default",
    },
}
```

### Version Comparison

Always use the approved entry points for version comparison:

```go
v, err := package_versions.NewVersion("1.2.3")
v2, err := package_versions.NewVersionFromPkgType("1.2.3-r0", "apk")
```

### Domain Types with PortalBase

Domain entities embed `PortalBase` for consistent metadata:

```go
type MyEntity struct {
    armotypes.PortalBase `json:",inline" bson:",inline"`
    // additional fields...
}
```

## Versioning Strategy

- **Module path**: `github.com/armosec/armoapi-go`
- **Tag format**: `v0.0.XXX` (auto-incremented on merge to main)
- **Release automation**: Tags are created automatically by `.github/workflows/release.yaml` on pushes to `main` that include non-Markdown changes
- **Consumer pinning**: Downstream services pin versions via `go.mod` and update as needed
- **Compatibility**: Exported struct fields with `json`/`bson` tags are serialization contracts. Renaming or removing tags is a breaking change.

## Build and Test

```bash
# Build (library only, no binary output)
go build ./...

# Run tests
go test ./...

# With coverage
go test -cover ./...

# Race detection
go test -race ./...
```

No external services or credentials are required for tests.

## Contribution Guidelines

1. All exported struct fields crossing service boundaries must have both `json` and `bson` struct tags
2. Interfaces in `containerscan/` are stability contracts; adding methods is a breaking change
3. New domain types should embed `PortalBase` for consistent metadata handling
4. Version comparison must use `package_versions.NewVersion()` or `NewVersionFromPkgType()`
5. High-throughput types should implement `gojay.MarshalerJSONObject` for performance
6. Test data belongs in `fixtures/` or `testdata/` subdirectories within each package
