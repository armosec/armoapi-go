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
- **Dynamic capability toggles** (added for [ADR 0009](https://github.com/armosec/shared-designs-and-docs/blob/main/adrs/0009-node-agent-capture-config-rename-and-otel-filter.md) — all three are `*bool` so an absent
  field is distinguishable from an explicit `false`, and each has its own polarity chosen to
  match how sensitive the capability it gates is):
  - `otelEventsEnabled` *(pointer; absent ⇒ **allow**, fail-safe)* — gates the node-agent's
    dynamic, per-event OTel raw-eBPF-event export filter. Raw eBPF telemetry is a detection
    signal, so losing it silently on a transient config-fetch issue is worse than a bounded
    over-collection window; hence fail-*safe*, not fail-closed like the other two.
  - `httpCaptureTapEnabled` *(Go field `HTTPCaptureTapEnabled`; pointer; absent ⇒ **deny**,
    fail-closed)* — the dedicated, live-toggleable switch for whether the node-agent's
    HTTP-capture tap is *constructed* at all. Distinct from `enabled` above, which is a lossless
    per-transaction policy decision on an already-existing tap; this field instead controls the
    tap's existence. Extends HTTP-capture's own on/off surface, so it inherits `enabled`'s
    fail-closed posture — captured content is potentially PII-bearing.
  - `workloadScanEnabled` *(pointer; absent ⇒ **deny**, fail-closed)* — gates the node-agent's
    NA-3 workload scanner (snapshots the overlay upperdir, secret/SA-token mount targets, and
    `/proc/<pid>/environ`; emits identity-material findings at least as sensitive as
    HTTP-capture content, hence fail-closed).
- **`mergeWithGlobal`** *(pointer; absent ⇒ **true**)* — controls whether a per-customer
  document is merged with the global default document (customer fields winning on conflict) or
  resolved as a standalone document (the global document is ignored entirely). Set explicitly to
  `false` on a customer document to opt that customer out of merging. **Only meaningful on a
  customer document** — `resolveHTTPCaptureConfig` never reads it off the global document
  (reading it there to decide whether to merge INTO the global document would be circular).
  Setting `mergeWithGlobal: false` on the global document itself has no effect.

**Merge granularity, part of this contract (not left to each consumer to decide):** when
`mergeWithGlobal` is `true` (the default) and both a customer and a global document exist,
`cadashboardbe`'s `resolveHTTPCaptureConfig` merges them **per field, not per document**, and the
merge is NOT uniform across every field:
- The four pointer-typed toggles above (`maskKnownCredentialHeaders`, `otelEventsEnabled`,
  `httpCaptureTapEnabled`, `workloadScanEnabled`) fall back to the global document's value
  **individually** when the customer document leaves them `nil` — `nil` is the one unambiguous
  "not set" signal these fields have. This matters *because* their fail-open/fail-closed
  polarities differ: a customer document that merely omits `otelEventsEnabled` correctly
  inherits global's value either way, rather than an omission on one field silently reading as
  "allow" and on another as "deny" depending on which document happened to be picked wholesale.
- Every other field (`protocolVersion`, `enabled`, the size/volume caps, `defaultLevel`, `rules`,
  `extraCredentialHeaders`) comes from the customer document **as a whole** when one exists —
  these are plain scalars/slices whose Go zero value (`false`/`0`/`""`/empty slice) can't be told
  apart from "the customer didn't set this," so merging them field-by-field would risk silently
  overwriting a customer's deliberate zero value (e.g. `enabled: false`) with global's.
- When `mergeWithGlobal` is explicitly `false` on the customer document, none of the above
  applies: the customer document is returned exactly as stored, global is never consulted for
  any field, pointer or scalar.

See `cadashboardbe`'s `docs/features/ai-sandbox-httpcapture-config.md` for the resolution
function itself.

## `CaptureRule`

Migrated from the node-agent's `uploadpolicy.Rule`. A transaction is matched by
`host` / `hostPrefix` / `hostSuffix` / `pathContains` / `protocol` (every set field is
AND-ed; empty = any) and assigned `level` — the wire string `"none" | "headers" | "full"`.
The agent maps `Level` to its internal `uploadpolicy.Level` on parse.

## `CaptureConfigResponse` — the served merged response

What the node-agent actually pulls is **one** response combining the capture policy with
the TLS offset-map section (TLS Offset-Map Distribution, SUB-7970):

```go
type CaptureConfigResponse struct {
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
