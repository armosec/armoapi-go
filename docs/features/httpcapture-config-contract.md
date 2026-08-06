---
type: feature
status: done
owner: jonathang
scope: service
service: armoapi-go
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
- **Capture policy** (added for SUB-7696 — the migration the type's doc comment promised):
  - `defaultLevel` (string) — level when no rule matches; absent ⇒ `"none"` (fail-closed).
  - `rules` (`[]CaptureRule`) — ordered, **first-match-wins** upload policy.

## `CaptureRule`

Migrated from the node-agent's `uploadpolicy.Rule`. A transaction is matched by
`host` / `hostPrefix` / `hostSuffix` / `pathContains` / `protocol` (every set field is
AND-ed; empty = any) and assigned `level` — the wire string `"none" | "headers" | "full"`.
The agent maps `Level` to its internal `uploadpolicy.Level` on parse.

## Consumers

- **config-service** — stores it (`v1_http_capture_config`), global (`customers:[""]`) +
  per-customer docs.
- **careportsreceiver** — serves the resolved config to the agent at
  `/cloud/v1/aiSandbox/config`.
- **node-agent** — `pkg/httpcaptureconfig` decodes it and applies it fail-closed.

Additive and backward-compatible: adding the policy fields does not affect existing
`CaptureConfig` users.
