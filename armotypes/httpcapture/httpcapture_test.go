package httpcapture

import (
	"encoding/json"
	"testing"
)

// The wire string values are part of the contract — both sensor and backend match
// on them, so a rename is a breaking change and must be caught here.
func TestEnumWireValues(t *testing.T) {
	if DirectionRequest != "request" || DirectionResponse != "response" {
		t.Errorf("direction wire values changed: %q %q", DirectionRequest, DirectionResponse)
	}
	fid := map[CaptureFidelity]string{
		FidelityComplete: "complete", FidelityTruncated: "truncated",
		FidelityTupleMissing: "tuple-missing", FidelityPartial: "partial",
	}
	for got, want := range fid {
		if string(got) != want {
			t.Errorf("fidelity wire value = %q, want %q", got, want)
		}
	}
}

// NewTransactionID must be deterministic in its inputs (same parts → same id) and
// distinguish connections by the nonce even when every other part repeats — the
// ssl_ptr-reuse collision the id exists to defeat (spec §5.2).
func TestNewTransactionID(t *testing.T) {
	a := NewTransactionID("inst-1", 42, 100, 0xdeadbeef, 1)
	b := NewTransactionID("inst-1", 42, 100, 0xdeadbeef, 1)
	if a != b {
		t.Errorf("same inputs must yield same id: %q vs %q", a, b)
	}
	// identical except the nonce — must differ (reused ssl_ptr)
	if c := NewTransactionID("inst-1", 42, 100, 0xdeadbeef, 2); c == a {
		t.Errorf("nonce must disambiguate a reused ssl_ptr, got %q for both", a)
	}
	// different instance — must differ
	if d := NewTransactionID("inst-2", 42, 100, 0xdeadbeef, 1); d == a {
		t.Errorf("instance id must namespace the transaction id, got %q for both", a)
	}
}

// A Fragment round-trips through JSON with the contract field names, and Data
// (bytes) survives — this is the binary-safe carrier for raw captured payload.
func TestFragmentJSONRoundTrip(t *testing.T) {
	in := Fragment{
		TransactionID:  "inst-1:2a:64:deadbeef:1",
		Direction:      DirectionResponse,
		SequenceNumber: 3,
		EndOfStream:    true,
		Data:           []byte{0x00, 0x01, 0xff, 'h', 'i'}, // non-UTF8: proves bytes carrier
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out Fragment
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.TransactionID != in.TransactionID || out.Direction != in.Direction ||
		out.SequenceNumber != in.SequenceNumber || out.EndOfStream != in.EndOfStream ||
		string(out.Data) != string(in.Data) {
		t.Errorf("round-trip mismatch:\n in=%+v\nout=%+v", in, out)
	}
}
