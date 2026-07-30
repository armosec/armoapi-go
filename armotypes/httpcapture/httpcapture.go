// Package httpcapture is the sensor↔backend wire contract for HTTP(S) capture
// (cloud-agent-sandbox spec §5). A sensor (the eBPF node-agent in Phase 1, an
// in-line proxy later) uploads captured HTTP transactions to the backend as a
// stream of message Fragments; the backend reassembles them. Both the sensor and
// the backend import THIS package so the wire shape is defined exactly once — the
// sensor never derives it from a backend/ingester package, and vice-versa.
//
// The sensor forwards raw bytes and does no content interpretation; everything
// semantic (parsing, redaction, dedup, token/usage derivation, reassembly) is the
// backend's. See spec §5.2 (wire unit) and §5.1 (capture fidelity).
package httpcapture

import (
	"fmt"

	"github.com/armosec/armoapi-go/armotypes"
)

// Direction identifies which half of a transaction a fragment belongs to
// (spec §5.2). It mirrors the capture-side direction: a request is what the
// workload sent (SSL_write), a response is what it received (SSL_read).
type Direction string

const (
	DirectionRequest  Direction = "request"
	DirectionResponse Direction = "response"
)

// Part marks which section of a direction a fragment carries: the header block or
// the body. It lets the backend split headers from body without re-parsing the
// reassembled stream, and — with EndOfStream — tells it whether to expect a body:
// a headers fragment with EndOfStream=true means the body was intentionally omitted
// (e.g. a headers-only capture level), not lost. Headers are typically one fragment
// (bounded ~4–16 KiB, captured in one shot); the body is streamed as many.
type Part string

const (
	PartHeaders Part = "headers"
	PartBody    Part = "body"
)

// CaptureFidelity marks how much of a record the sensor actually saw (spec §5.1),
// so the backend and UI never mistake a capture gap for an absence of activity.
type CaptureFidelity string

const (
	FidelityComplete     CaptureFidelity = "complete"      // whole transaction captured
	FidelityTruncated    CaptureFidelity = "truncated"     // body cut at a capture/size bound
	FidelityTupleMissing CaptureFidelity = "tuple-missing" // destination 5-tuple unrecoverable
	FidelityPartial      CaptureFidelity = "partial"       // a fragment was lost / reassembly timed out
)

// DefaultMaxFragmentBytes is the default per-fragment payload cap. It sits well
// under the ~4 MB OTLP/gRPC message ceiling to leave headroom for the record
// envelope and for the ~33% base64 expansion a bytes value incurs when a LogRecord
// is JSON-encoded; 1 MiB of raw payload encodes to well under 4 MB (spec §5.10).
const DefaultMaxFragmentBytes = 1 << 20 // 1 MiB

// CurrentProtocolVersion is the wire-contract version stamped on every record and
// every config response "from day one" (spec §5.10): with two sensor types and a
// fidelity enum evolving independently, versioning the contract now is cheap
// insurance against a breaking change later. Bump on a breaking wire change.
const CurrentProtocolVersion = "1.0"

// Attribute keys carried on each fragment LogRecord (spec §5.2). Defined here so
// the sensor that writes them and the backend that reads them share one source of
// truth for the key strings.
//
// These are the fragment-specific keys owned by this contract. Each Fragment is
// emitted as one OTLP LogRecord whose Body is the raw Data bytes; the identity fields
// below (K8s/Cloud/SandboxID/Server*/ProcessTree/…) are encoded with the shared
// sandbox-event-schema keys, NOT http.capture.* keys, so capture fragments carry one
// identity scheme with every other event family. The full LogRecord encoding
// (resource attrs + record attrs + Body + a worked example) is documented in
// private-node-agent docs/features/http-capture.md §3.1 "OTEL LogRecord encoding".
const (
	AttrTransactionID   = "http.capture.transaction_id"
	AttrSequenceNumber  = "http.capture.sequence_number"
	AttrDirection       = "http.capture.direction"
	AttrEndOfStream     = "http.capture.end_of_stream"
	AttrCaptureFidelity = "http.capture.fidelity"
	AttrFidelityReason  = "http.capture.fidelity_reason"
	AttrProtocolVersion = "http.capture.protocol_version"
	AttrPart            = "http.capture.part"
	AttrCapturedAt      = "http.capture.captured_at"
)

// Fragment is one unit on the wire (spec §5.2): a slice of one direction of one
// logical HTTP transaction. It is smaller than a whole request/response — a sensor
// only ever sees limited slices (one SSL_write unit, one proxy read buffer), and
// long/streamed bodies exceed the OTLP message ceiling — so a transaction is a
// STREAM of fragments, not one record.
//
// The backend reassembles per (TransactionID, Direction) ordered by SequenceNumber
// until EndOfStream. Completeness requires the terminal marker AND contiguous
// SequenceNumber coverage 0..N; a gap resolves to FidelityPartial.
type Fragment struct {
	// ProtocolVersion is the wire-contract version this record was produced under
	// (spec §5.10) — carried on every record from day one so the two sensor types
	// and the fidelity enum can evolve without a silent breaking change.
	ProtocolVersion string `json:"protocolVersion"`
	// TransactionID correlates all fragments of one logical HTTP transaction —
	// BOTH directions (request + response) share it. It MUST be globally unique
	// across every sensor instance and tenant (see NewTransactionID); the backend
	// treats it as opaque and only requires uniqueness.
	TransactionID string `json:"transactionId"`
	// Direction is request or response.
	Direction Direction `json:"direction"`
	// Part is which section of the direction this fragment carries (headers | body).
	// Headers fragment(s) come first (SequenceNumber 0…), body fragments after.
	Part Part `json:"part"`
	// SequenceNumber orders fragments within a direction, contiguous from 0
	// (across both parts). The sensor produces it; the backend re-sorts by it.
	SequenceNumber uint32 `json:"sequenceNumber"`
	// EndOfStream is set only on the final fragment of a direction, so the backend
	// knows the direction is done without a known total size. The sensor decides it
	// from transport framing (Content-Length satisfied / chunked terminator / close).
	EndOfStream bool `json:"endOfStream"`
	// CapturedAt is when this fragment was captured (unix nanoseconds). Per-fragment,
	// because a streamed transaction spans time.
	CapturedAt int64 `json:"capturedAt"`
	// Data is the raw captured payload slice (headers-block bytes or a body chunk),
	// verbatim — never parsed, masked, or reassembled by the sensor. Bounded by
	// DefaultMaxFragmentBytes.
	Data []byte `json:"data"`

	// --- Lightweight identity, on EVERY fragment (placement "B") ---
	// So the backend can attribute any fragment and query "in a sandbox" without a
	// join. These REUSE the existing armotypes identity structs — the agent does not
	// define its own process/cloud attribution. (The heavy process TREE is NOT here;
	// it rides once, on the first fragment — placement "A", see below.)
	//
	// This identity may also be propagated to OTLP resource attributes to avoid
	// repeating it per fragment; that is a later, easy-to-add optimization.
	CustomerGUID string                            `json:"customerGuid,omitempty"` // OTEL: customer_guid (resource attr)
	SandboxID    string                            `json:"sandboxId,omitempty"`    // OTEL: sandbox_id / instance_ref
	K8s          *armotypes.RuntimeAlertK8sDetails `json:"k8s,omitempty"`          // OTEL: wlid / pod_name / namespace / container.id / container_name
	Cloud        *armotypes.CloudMetadata          `json:"cloud,omitempty"`        // OTEL: cloud.account_id / cloud.region / cloud.provider

	// --- Per-transaction attribution, on the FIRST fragment only (placement A) ---
	// Carried once — on the request-headers fragment (SequenceNumber 0) — not a
	// separate record and not repeated per fragment (the process tree can be large).
	// nil/zero on every other fragment. Reuses armotypes.ProcessTree (RuntimeAlert
	// tree) — the agent defines no process attribution of its own.
	ProcessTree      *armotypes.ProcessTree `json:"processTree,omitempty"`      // OTEL: process_lineage (JSON, first fragment)
	ServerAddress    string                 `json:"serverAddress,omitempty"`    // OTEL: server.address (every fragment)
	ServerPort       uint16                 `json:"serverPort,omitempty"`       // OTEL: server.port (every fragment)
	SensorInstanceID string                 `json:"sensorInstanceId,omitempty"` // OTEL: sensor_instance_id (first fragment)

	// Fidelity marks a capture gap on this fragment (§5.1); empty ⇒ complete. On a
	// per-transaction size-cap truncation the terminal fragment carries
	// EndOfStream=true AND Fidelity=FidelityTruncated.
	Fidelity       CaptureFidelity `json:"fidelity,omitempty"`
	FidelityReason string          `json:"fidelityReason,omitempty"`
}

// CaptureConfig is the backend-supplied capture configuration the sensor polls
// (spec §5.3/§5.4). It carries enablement, the caps, and the masking controls. The
// sensor uses built-in defaults until a config is fetched (a 0 cap ⇒ agent default),
// and fails closed (captures nothing) with no config. The per-target upload-policy
// rules (host/path → full|headers|none) are the agent's uploadpolicy types today and
// migrate into this package next.
type CaptureConfig struct {
	ProtocolVersion string `json:"protocolVersion"`
	// Enabled is the master switch; false ⇒ the sensor captures nothing.
	Enabled bool `json:"enabled"`

	// Size caps in bytes (0 ⇒ use the agent default). MaxFragmentBytes bounds one
	// fragment; MaxTransactionBytes bounds a whole transaction — on hit the body is
	// truncated (headers are still sent) and the terminal fragment is marked
	// EndOfStream + FidelityTruncated.
	MaxFragmentBytes    int64 `json:"maxFragmentBytes,omitempty"`
	MaxTransactionBytes int64 `json:"maxTransactionBytes,omitempty"`

	// MaxTransactionsPerHour is the per-agent volume cap (0 ⇒ agent default). When it
	// is exhausted the agent uploads NOTHING for further transactions this window and
	// emits a "limit exhausted" signal to the backend.
	MaxTransactionsPerHour int64 `json:"maxTransactionsPerHour,omitempty"`

	// Credential-header masking (§5.6). MaskKnownCredentialHeaders is a pointer so
	// absent ⇒ default true (fail-safe: mask). ExtraCredentialHeaders adds names to
	// the masked set. Masking is a one-way (irreversible) fingerprint.
	MaskKnownCredentialHeaders *bool    `json:"maskKnownCredentialHeaders,omitempty"`
	ExtraCredentialHeaders     []string `json:"extraCredentialHeaders,omitempty"`
}

// NewTransactionID composes a globally-unique transaction id from readily-available
// capture-side parts (spec §5.2 "Global uniqueness (required)"). A raw ssl_ptr is a
// process-local pointer reused after SSL_free, so it collides within an instance
// over time and across instances immediately; namespacing it by the sensor-instance
// identity + process start time and disambiguating with a per-connection nonce makes
// a collision vanishingly unlikely. The backend does not parse the result.
func NewTransactionID(instanceID string, pid uint32, procStartNS, sslPtr, nonce uint64) string {
	return fmt.Sprintf("%s:%x:%x:%x:%x", instanceID, pid, procStartNS, sslPtr, nonce)
}
