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

import "fmt"

// Direction identifies which half of a transaction a fragment belongs to
// (spec §5.2). It mirrors the capture-side direction: a request is what the
// workload sent (SSL_write), a response is what it received (SSL_read).
type Direction string

const (
	DirectionRequest  Direction = "request"
	DirectionResponse Direction = "response"
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
// is JSON-encoded; 1 MiB of raw payload encodes to well under 4 MB.
const DefaultMaxFragmentBytes = 1 << 20 // 1 MiB

// Attribute keys carried on each fragment LogRecord (spec §5.2). Defined here so
// the sensor that writes them and the backend that reads them share one source of
// truth for the key strings.
const (
	AttrTransactionID   = "http.capture.transaction_id"
	AttrSequenceNumber  = "http.capture.sequence_number"
	AttrDirection       = "http.capture.direction"
	AttrEndOfStream     = "http.capture.end_of_stream"
	AttrCaptureFidelity = "http.capture.fidelity"
	AttrFidelityReason  = "http.capture.fidelity_reason"
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
	// TransactionID correlates all fragments of one logical HTTP transaction —
	// BOTH directions (request + response) share it. It MUST be globally unique
	// across every sensor instance and tenant (see NewTransactionID); the backend
	// treats it as opaque and only requires uniqueness.
	TransactionID string `json:"transactionId"`
	// Direction is request or response.
	Direction Direction `json:"direction"`
	// SequenceNumber orders fragments within a direction, contiguous from 0.
	SequenceNumber uint32 `json:"sequenceNumber"`
	// EndOfStream is set only on the final fragment of a direction, so the backend
	// knows the direction is done without a known total size.
	EndOfStream bool `json:"endOfStream"`
	// Data is the raw captured payload slice (headers-block bytes or a body chunk),
	// verbatim — never parsed, masked, or reassembled by the sensor. Bounded by
	// DefaultMaxFragmentBytes.
	Data []byte `json:"data"`
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
