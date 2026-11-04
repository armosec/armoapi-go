---
type: feature
status: active
owner: matthiasb@armosec.io
scope: repo
related_code:
  - apis/interfaces.go
  - apis/websocketdatastructures.go
  - apis/websocketmethods.go
---

# Registry credentials type (`registry.AuthConfig` from `moby/moby/api`)

The scan commands in `apis` carry registry credentials as `registry.AuthConfig`. They
are used by `ImageScanCommand` and `RegistryScanCommand` (`GetCreds`, `GetCredentialsList`,
`SetCredentialsList`) and by the `WebsocketScanCommand.Credentials` and
`ImageScanParams.Credentialslist` fields. That type comes from
`github.com/moby/moby/api/types/registry`, not from `github.com/docker/docker/api/types/registry`.

## Why

`github.com/docker/docker` is frozen at v28.5.2. Moby now publishes its client and API
types as `github.com/moby/moby/api` and `github.com/moby/moby/client`, and the daemon as
`github.com/moby/moby/v2`. The Go vulnerability database entries GO-2026-4887
(CVE-2026-34040) and GO-2026-4883 (CVE-2026-33997) flag the **whole** `docker/docker`
module, with no package or symbol list and no fixed version. Because of that,
`govulncheck` fails for any binary that links even a plain type package from it. The
vulnerable code is daemon-side, but that doesn't help. Consumers that link `apis`
(for example kubevuln) could not pass `govulncheck` while `apis` imported `docker/docker`.

## Wire format

Unchanged. The struct has the same fields and JSON tags, except for the deprecated
`Email` field, which `moby/moby/api` dropped:

| Field | JSON |
|---|---|
| `Username` | `username` |
| `Password` | `password` |
| `Auth` | `auth` |
| `ServerAddress` | `serveraddress` |
| `IdentityToken` | `identitytoken` |
| `RegistryToken` | `registrytoken` |

A sender that still includes `"email"` is not rejected; the value is ignored on decode.

## Consumer migration

This is a Go type change, so consumers must switch in lockstep when they bump
armoapi-go. In each file that names `registry.AuthConfig` for values passed to or from `apis`:

```go
// before
import "github.com/docker/docker/api/types/registry"
// after
import "github.com/moby/moby/api/types/registry"
```

Remove any use of `AuthConfig.Email`. Then run `go mod tidy`. `docker/docker` may remain
as `// indirect` through other dependencies. `govulncheck -mode=binary` only considers
packages linked into the binary, so an indirect `go.mod` entry by itself does not fail the
check. A transitive dependency can still link a `docker/docker` package, though. To
confirm none is linked, run `go list -deps <main package> | grep docker/docker`, which
should print nothing. This change only removes the import from `apis`. Other armoapi-go
packages (for example `package_versions`, through grype) still reach `docker/docker`, so a
consumer that links those is not cleared by this change.
