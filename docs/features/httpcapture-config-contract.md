---
type: feature
status: done
owner: jonathang
scope: repo
---

# HTTP-capture Config Contract (`armotypes/httpcapture`)

## Purpose

`armotypes/httpcapture` is the **single, shared schema** for the AI-Sandbox HTTP(S)-capture
feature (SUB-7696) — the one place the whole platform agrees on the wire shape, so
config-service, the dashboard BE, careportsreceiver, and the node-agent don't each redefine
it and drift. A schema change is one edit here + a version bump that propagates to every
consumer, and the Go compiler enforces that they agree.

## `CaptureConfig` — the backend-driven capture config

The config the node-agent pulls (via careportsreceiver) and applies, fail-closed:

- **Scalars:** `protocolVersion`, `enabled`, size caps (`maxFragmentBytes`,
  `maxTransactionBytes`), volume caps (`maxTransactionsPerHour`,
  `maxTransactionsPerSandboxPerHour`), and credential masking
  (`maskKnownCredentialHeaders` *(pointer; absent ⇒ mask)*, `extraCredentialHeaders`).
  `maxFragmentBytes` absent/0 ⇒ `DefaultMaxFragmentBytes` = **256 KiB**: one fragment is one
  base64-encoded JSON record on the fragment→Iceberg Firehose hop, so the raw cap must leave
  room for the ~33% expansion under the ~1 MB Firehose per-record limit (256 KiB → ~350 KB;
  1 MiB → ~1.4 MB would exceed it and jam the stream).
- **Capture policy** (added for SUB-7696 — the migration the type's doc comment promised):
  - `defaultLevel` (string) — level when no rule matches; absent ⇒ `"none"` (fail-closed).
  - `rules` (`[]CaptureRule`) — ordered, **first-match-wins** upload policy.

## `CaptureRule`

Migrated from the node-agent's `uploadpolicy.Rule`. A transaction is matched by
`host` / `hostPrefix` / `hostSuffix` / `pathContains` / `protocol` (every set field is
AND-ed; empty = any) and assigned `level` — the wire string `"none" | "headers" | "full"`.
The agent maps `Level` to its internal `uploadpolicy.Level` on parse.

## `SandboxConfigResponse` — the served merged response

What the node-agent actually pulls is **one** response combining the capture policy with
the TLS offset-map section (TLS Offset-Map Distribution, SUB-7970):

```go
type SandboxConfigResponse struct {
    CaptureConfig `json:",inline"`                            // capture policy, top-level
    TLSOffsets json.RawMessage `json:"tlsOffsets,omitempty"`  // offsets, opaque raw JSON
}
```

Defined here so producer (cadashboardbe) and consumer (node-agent) share one envelope and
can't drift on the `tlsOffsets` key. `TLSOffsets` is raw (a `map[identity]tlsoffsets.Record`
once decoded) so a malformed offsets section can never fail the fail-closed capture decode;
`SetTLSOffsets` is the producer helper. **`CaptureConfig` is not widened** — it stays the
capture-only contract config-service stores; offsets keep their own type
([`armotypes/tlsoffsets`](tlsoffsets-record-contract.md)) and their own collection. The two
only ride together in this response. careportsreceiver forwards it as opaque bytes.

## Consumers

- **config-service** — stores it (`v1_http_capture_config`), global (`customers:[""]`) +
  per-customer docs.
- **careportsreceiver** — serves the resolved config to the agent at
  `/cloud/v1/aiSandbox/config`.
- **node-agent** — `pkg/httpcaptureconfig` decodes it and applies it fail-closed.

Additive and backward-compatible: adding the policy fields does not affect existing
`CaptureConfig` users.
