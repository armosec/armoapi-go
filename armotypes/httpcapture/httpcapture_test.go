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
	// Every input must independently affect the id — change exactly one from the
	// baseline at a time, so the test fails if any part is dropped from the id.
	for _, tc := range []struct {
		name string
		id   string
	}{
		{"nonce", NewTransactionID("inst-1", 42, 100, 0xdeadbeef, 2)},
		{"instanceID", NewTransactionID("inst-2", 42, 100, 0xdeadbeef, 1)},
		{"pid", NewTransactionID("inst-1", 43, 100, 0xdeadbeef, 1)},
		{"procStartNS", NewTransactionID("inst-1", 42, 101, 0xdeadbeef, 1)},
		{"sslPtr", NewTransactionID("inst-1", 42, 100, 0xdeadbee0, 1)},
	} {
		if tc.id == a {
			t.Errorf("changing %s must change the transaction id, but both were %q", tc.name, a)
		}
	}
}

// A Fragment round-trips through JSON with the contract field names, and Data
// (bytes) survives — this is the binary-safe carrier for raw captured payload.
func TestFragmentJSONRoundTrip(t *testing.T) {
	in := Fragment{
		ProtocolVersion: CurrentProtocolVersion,
		TransactionID:   "inst-1:2a:64:deadbeef:1",
		Direction:       DirectionResponse,
		SequenceNumber:  3,
		EndOfStream:     true,
		Data:            []byte{0x00, 0x01, 0xff, 'h', 'i'}, // non-UTF8: proves bytes carrier
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	// Pin the exact wire shape (field names + order + base64 Data) so a JSON-tag
	// regression is caught — round-tripping the same Go type alone would not.
	const wantJSON = `{"protocolVersion":"1.0","transactionId":"inst-1:2a:64:deadbeef:1","direction":"response","sequenceNumber":3,"endOfStream":true,"data":"AAH/aGk="}`
	if string(b) != wantJSON {
		t.Fatalf("wire JSON =\n  %s\nwant\n  %s", b, wantJSON)
	}
	var out Fragment
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.ProtocolVersion != in.ProtocolVersion || out.TransactionID != in.TransactionID ||
		out.Direction != in.Direction || out.SequenceNumber != in.SequenceNumber ||
		out.EndOfStream != in.EndOfStream || string(out.Data) != string(in.Data) {
		t.Errorf("round-trip mismatch:\n in=%+v\nout=%+v", in, out)
	}
}
