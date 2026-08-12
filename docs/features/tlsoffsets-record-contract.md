---
type: feature
status: active
owner: amirm
scope: repo
---

# TLS Offset-Map Record Contract (`armotypes/tlsoffsets`)

## Purpose

`armotypes/tlsoffsets` is the **single, shared envelope** for one TLS-capture offset record
distributed to the node-agent (TLS Offset-Map Distribution, SUB-7970). It is the one place
config-service, cadashboardbe, and the node-agent agree on the record's shape, so none of
them redefine it and drift — the same pattern as [`armotypes/httpcapture`](httpcapture-config-contract.md).

It is a **separate type from `httpcapture.CaptureConfig`**, and it is stored in a **separate
config-service collection (`v1_tls_offsets`)** — never mixed into the capture config's type or
its `v1_http_capture_config` collection. The two only meet in the *serve-time response body*
(cadashboardbe splices a `tlsOffsets` section alongside the resolved capture config); in
storage and in the shared type they stay independent so each can evolve on its own.

## `Record` — one offset record

- **Queryable envelope fields (top-level):** `target` (`claude` | `opencode`; `aead-awslc`
  reserved), `platform` (`linux`), `arch` (`x86_64` | `arm64`). config-service filters records
  by these via `InnerFilters` — no consumer parses the payload to filter.
- **`payload`** (`json.RawMessage`) — the target-specific offset data, **opaque** to store and
  transport. For `claude`/`opencode` it is the standard `encoding/json` of the agent's
  internal offset record (numeric offsets, base64 byte-signatures). config-service and
  cadashboardbe never open it.
- **`protocolVersion`** — forward-compat handle (`CurrentProtocolVersion`); same major ⇒
  additive fields tolerated.

**Identity** of a record is the config-service doc GUID (a bare hex build-id, or
`sha256:<hex>` for a content-digest fallback), used as the key in the served `tlsOffsets`
map — it is **not** a field on `Record`.

## Served envelope

A record reaches the agent inside the merged **`httpcapture.SandboxConfigResponse`** (capture
policy inline + a raw `tlsOffsets` section keyed by identity) — documented in
[httpcapture-config-contract.md](httpcapture-config-contract.md#sandboxconfigresponse--the-served-merged-response).
`Record` stays independent of that envelope and of the capture config.

## Consumers

- **config-service** — stores it in `v1_tls_offsets`, global (`customers:[""]`) docs; filters
  by `platform`/`arch`.
- **cadashboardbe** — queries the collection and merges a `tlsOffsets` section into the
  `httpCaptureConfig` response (careportsreceiver proxies it through).
- **node-agent** — decodes each record, routes the payload to its `target`'s validator, and
  swaps that target's in-memory offset table.
