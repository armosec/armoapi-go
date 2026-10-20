---
type: feature
status: active
owner: ai-sandbox
scope: repo
---

# HTTP-capture capture source (`armotypes/httpcapture`)

## Purpose

The backend could only guess whether a captured HTTP transaction came from a TLS hook or a
plain socket by looking at the destination port (443). That guess breaks for memory-BIO
captures, where the destination tuple is lost (`server_port=0`), and it can never tell an
OpenSSL hook from a Go `crypto/tls` hook. The sensor knows which tap produced each fragment,
so it states it explicitly.

## `AttrCaptureSource` / `CaptureSource`

- Attribute key `http.capture.source` (`AttrCaptureSource`), carried on each fragment
  LogRecord next to the other `http.capture.*` keys.
- Values (`CaptureSource`):
  - `ssl` — OpenSSL/BoringSSL uprobes.
  - `gotls` — Go `crypto/tls` uprobes.
  - `cleartext` — plain (non-TLS) socket tap.
- `NormalizeCaptureSource(raw)` returns the typed value for a known string and `""` for
  anything else, so consumers keep the column bounded and never store free text.
  `CaptureSource.Valid()` reports whether a value is one of the known ones.

Absent or unknown means "not reported" (older sensors). Consumers must never infer the
source from the port or host instead.

## Consumers

- event-ingester writes it to the `http_fragment_bronze.capture_source` column, and the
  reassembler carries it into `http_transaction_silver` (armosec/backend, SUB-9362).
