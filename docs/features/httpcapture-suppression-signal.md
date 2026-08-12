---
type: feature
status: done
owner: matthiasb@armosec.io
scope: repo
---

# HTTP-capture non-loss suppression signal (`armotypes/httpcapture`)

## Purpose

`Fragment.Fidelity`/`FidelityReason` are contractually **loss-only**: they mean "content
the sensor should have captured but didn't" — a capture gap the backend must not mistake
for an absence of activity. There was no field for the opposite case: a sensor
**deliberately** withholding content by policy, because it decided the content had no
analytical value (e.g. a periodic keepalive/heartbeat pattern). Reusing `Fidelity` for
that would misreport intentional noise reduction as a capture defect — exactly the
conflation this signal exists to prevent.

## `Fragment.SuppressedReason` / `Fragment.SuppressedCount`

Added to `Fragment`, orthogonal to `Fidelity`/`FidelityReason`:

- **`suppressedReason`** (string) — a short, sensor-defined tag identifying *what* was
  withheld, e.g. `"ws-keepalive"`.
- **`suppressedCount`** (uint32) — how many items were withheld under that reason,
  scoped to this fragment's transaction.

Both empty/zero on every existing fragment — additive, backward-compatible. A sensor
withholding more than one *kind* of thing on the same transaction uses more than one
non-suppressed fragment to carry them; this pair is deliberately not a map, to keep the
wire shape simple for the first (and likely only, for now) consumer.

OTLP attribute keys, alongside the existing `http.capture.*` keys:
`AttrSuppressedReason` (`http.capture.suppressed_reason`),
`AttrSuppressedCount` (`http.capture.suppressed_count`).

### Sequence-number contract (the one thing that changes how a consumer handles it)

**Suppression never creates a `SequenceNumber` gap.** A withheld item never reaches the
fragment stream in the first place, so `SequenceNumber` stays contiguous (0..N,
unbroken) across whatever fragments *are* emitted — there are simply fewer of them. A
consumer must not infer `FidelityPartial`, or any other loss signal, from
`SuppressedCount > 0`.

Verified against the reference sensor implementation (private-node-agent,
`pkg/wscapture`): a suppressed item's code path returns before the call that assigns a
`SequenceNumber`, so the counter backing it is simply never incremented for that item —
by construction, not by convention a future sensor could accidentally violate.

## Consumers

- **private-node-agent** (`pkg/wscapture`) — the first consumer. WS keepalive filtering
  (`docs/features/ws-keepalive-filtering.md` in that repo) currently uses a per-exchange,
  log-only interim trace pending this field reaching a released `armoapi-go` version and
  a `go.mod` bump there.
- Backend/ingester consumption of the new fields is a follow-up, tracked separately —
  this PR only adds the wire-schema capability.
