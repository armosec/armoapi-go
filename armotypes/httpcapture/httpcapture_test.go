package httpcapture

import (
	"encoding/json"
	"testing"

	"github.com/armosec/armoapi-go/armotypes"
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
		Part:            PartBody,
		SequenceNumber:  3,
		EndOfStream:     true,
		CapturedAt:      1700000000000000000,
		Data:            []byte{0x00, 0x01, 0xff, 'h', 'i'}, // non-UTF8: proves bytes carrier
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	// Pin the exact wire shape (field names + order + base64 Data) so a JSON-tag
	// regression is caught — round-tripping the same Go type alone would not.
	const wantJSON = `{"protocolVersion":"1.0","transactionId":"inst-1:2a:64:deadbeef:1","direction":"response","part":"body","sequenceNumber":3,"endOfStream":true,"capturedAt":1700000000000000000,"data":"AAH/aGk="}`
	if string(b) != wantJSON {
		t.Fatalf("wire JSON =\n  %s\nwant\n  %s", b, wantJSON)
	}
	var out Fragment
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.ProtocolVersion != in.ProtocolVersion || out.TransactionID != in.TransactionID ||
		out.Direction != in.Direction || out.Part != in.Part ||
		out.SequenceNumber != in.SequenceNumber || out.EndOfStream != in.EndOfStream ||
		out.CapturedAt != in.CapturedAt || string(out.Data) != string(in.Data) {
		t.Errorf("round-trip mismatch:\n in=%+v\nout=%+v", in, out)
	}
}

// The reused identity (per fragment, "B") and the reused process tree (per
// transaction, "A") round-trip — proving we compose the existing armotypes structs
// rather than inventing our own process/cloud attribution.
func TestReusedIdentityAndProcessTree(t *testing.T) {
	// Fragment carries the lightweight identity (B).
	f := Fragment{
		ProtocolVersion: CurrentProtocolVersion,
		TransactionID:   "t1", Direction: DirectionRequest, Part: PartHeaders,
		CustomerGUID: "cust-abc", SandboxID: "sbx-1",
		K8s:   &armotypes.RuntimeAlertK8sDetails{PodName: "orders-7c9", WorkloadName: "orders", NodeName: "node-7"},
		Cloud: &armotypes.CloudMetadata{AccountID: "1234", Region: "eu-central-1"},
	}
	var fo Fragment
	if b, err := json.Marshal(f); err != nil {
		t.Fatal(err)
	} else if err := json.Unmarshal(b, &fo); err != nil {
		t.Fatal(err)
	}
	if fo.K8s == nil || fo.K8s.PodName != "orders-7c9" || fo.Cloud == nil || fo.Cloud.Region != "eu-central-1" || fo.CustomerGUID != "cust-abc" {
		t.Errorf("fragment identity did not round-trip: %+v", fo)
	}

	// The FIRST fragment (request headers, seq 0) carries the full process lineage
	// (A): P1 spawned by P2. On other fragments ProcessTree is nil.
	first := Fragment{
		ProtocolVersion: CurrentProtocolVersion, TransactionID: "t1",
		Direction: DirectionRequest, Part: PartHeaders, SequenceNumber: 0,
		ProcessTree: &armotypes.ProcessTree{ProcessTree: armotypes.Process{
			PID: 200, Comm: "sandbox-agent",
			ChildrenMap: map[armotypes.CommPID]*armotypes.Process{
				{Comm: "python3", PID: 201}: {PID: 201, PPID: 200, Comm: "python3"},
			},
		}},
		ServerAddress: "api.example.com", ServerPort: 443, SensorInstanceID: "node-7",
	}
	var firstO Fragment
	if b, err := json.Marshal(first); err != nil {
		t.Fatal(err)
	} else if err := json.Unmarshal(b, &firstO); err != nil {
		t.Fatal(err)
	}
	if firstO.ProcessTree == nil || firstO.ProcessTree.ProcessTree.PID != 200 || firstO.ProcessTree.FindProcessByPID(201) == nil {
		t.Errorf("process tree/lineage did not round-trip on the first fragment: %+v", firstO.ProcessTree)
	}
	if fo.ProcessTree != nil {
		t.Error("non-first fragment must not carry a process tree")
	}
}

// CaptureConfig round-trips the caps + masking controls the backend delivers.
func TestCaptureConfigRoundTrip(t *testing.T) {
	mask := false
	in := CaptureConfig{
		ProtocolVersion: CurrentProtocolVersion, Enabled: true,
		MaxFragmentBytes: 1 << 20, MaxTransactionBytes: 8 << 20, MaxTransactionsPerHour: 5000,
		MaskKnownCredentialHeaders: &mask, ExtraCredentialHeaders: []string{"x-acme-signature"},
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out CaptureConfig
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if !out.Enabled || out.MaxTransactionsPerHour != 5000 || out.MaxTransactionBytes != 8<<20 ||
		out.MaskKnownCredentialHeaders == nil || *out.MaskKnownCredentialHeaders != false ||
		len(out.ExtraCredentialHeaders) != 1 {
		t.Errorf("config did not round-trip: %+v", out)
	}
}
