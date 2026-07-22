---
type: architecture
status: approved
owner: "@armosec/backend"
scope: repo
---
<!-- docs-system: scaffold v1 -->

# armoapi-go — Architecture

## Code layout

```
armoapi-go/
├── apis/                    # Command dispatch & websocket communication
│   ├── datastructures.go    # Command, Commands, SessionChain, JobTracking
│   ├── interfaces.go        # Connector, ImageScanCommand interfaces
│   ├── websocket*.go        # WebSocket scan command types & methods
│   ├── login*.go            # LoginObject, authentication data structures
│   └── backendconnector*.go # Backend connector implementation & methods
├── armotypes/               # Core platform domain types (~75 files)
│   ├── portaltypes.go       # PortalBase, PortalCluster, InstallationData
│   ├── configtypes.go       # CustomerConfig, Settings
│   ├── exceptionpolicy.go   # VulnerabilityExceptionPolicy
│   ├── attackchainstypes.go # AttackChain, SecurityRisk
│   ├── networkpolicies.go   # NetworkPoliciesWorkload
│   ├── registrytypes.go     # RegistryInfo, RegistryJobParams
│   ├── cdr/                 # Cloud Detection & Response types
│   │   ├── cdr.go           # CdrAlert, CloudMetadata, CdrAlertBatch
│   │   ├── aws.go           # CloudTrailEvent, AWS-specific types
│   │   └── azure.go         # AzureActivityLogEvent, Azure-specific types
│   └── common/              # Shared runtime sub-types (ProcessEntity, FileEntity)
├── identifiers/             # Resource designator system
│   ├── designators.go       # PortalDesignator, DesignatorType, constants
│   ├── armocontext.go       # ArmoContext helpers
│   └── utils.go             # WLID/SID parsing utilities
├── containerscan/           # Vulnerability scan report contracts
│   ├── interfaces.go            # Interfaces: ScanReport, VulnerabilityResult
│   ├── commondatastructures.go  # CommonContainerVulnerabilityResult, scan structs
│   ├── commonadapters.go        # Severity mapping, adapter logic
│   ├── consts.go                # Severity levels, scan status constants
│   └── v1/                      # Concrete implementation of interfaces
├── notifications/           # Alert & collaboration configuration
│   ├── alertchannelapi.go   # AlertChannelAPI, AlertConfig
│   ├── collaborationconfig.go # CollaborationConfig (Jira, Linear, etc.)
│   ├── integrationhealth.go # Integration token-health attributes + helpers
│   └── runtimeincidents.go  # Runtime incident notification types
├── broadcastevents/         # Platform event types for analytics pipeline
│   ├── events.go            # EventBase, LoginEvent, HelmInstalledEvent
│   └── factory.go           # Event factory/constructor helpers
├── package_versions/        # Multi-format version comparison
│   ├── version.go           # Version interface, NewVersion(), NewVersionFromPkgType()
│   ├── semantic_version.go  # Semver implementation
│   ├── apk_version.go       # Alpine APK version comparison
│   ├── deb_version.go       # Debian version comparison
│   ├── rpm_version.go       # RPM version comparison
│   ├── maven_version.go     # Maven version comparison
│   ├── pep440_version.go    # Python PEP440 version comparison
│   └── golang_version.go    # Go module version comparison
└── scanfailure/             # Scan failure reporting
    └── types.go             # ScanFailureReport, reason codes, ReasonFriendlyText()
```

## Build & test

**Build (library — no binary):**

```bash
go build ./...
```

**Test:**

```bash
go test ./...            # unit tests
go test -cover ./...     # with coverage
go test -race ./...      # race detection
```

No external dependencies needed for tests.

**Release:**

Tags follow `v0.0.XXX` pattern. Push a tag to trigger `.github/workflows/release.yaml`. Consumers pin versions via `go.mod`.

## Versioning & compatibility

- Module path: `github.com/armosec/armoapi-go`
- All consumers import via Go modules. Breaking changes to exported types require coordinated updates across consumer services.
- Interfaces (`containerscan.ScanReport`, `apis.Connector`) are stability contracts. Adding methods is a breaking change.
- Struct fields with `json`/`bson` tags are serialization contracts. Renaming or removing tags breaks wire compatibility.

## Conventions

- **Struct tags**: Every exported struct field that crosses a service boundary must have both `json` and `bson` tags. The `json` tag is the wire format name.
- **PortalBase embedding**: Domain entities stored in MongoDB embed `PortalBase` for GUID, name, cluster, and attribute metadata.
- **Interface-based contracts**: `containerscan/` defines interfaces; `containerscan/v1/` provides the implementation. Consumers depend on the interface package only.
- **Version comparison**: All version comparison must go through `package_versions.NewVersion()` or `NewVersionFromPkgType()`. Raw string comparison is incorrect for non-semver formats.
- **No runtime dependencies**: This library must not import HTTP servers, database drivers, or message brokers. It is a pure type + utility library.
- **Gojay for hot paths**: Types that are marshalled at high throughput (scan results, events) implement `gojay.MarshalerJSONObject` for zero-allocation encoding.
- **Test fixtures**: Test data lives in `fixtures/` or `testdata/` subdirectories within each package.
