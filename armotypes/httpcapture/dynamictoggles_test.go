package httpcapture

import (
	"encoding/json"
	"strings"
	"testing"
)

// boolPtr is a test helper: CaptureConfig's dynamic toggles are all *bool, and Go has no
// address-of-literal syntax.
func boolPtr(b bool) *bool { return &b }

// TestDynamicToggles_AbsentWhenNil pins the fail-safe/fail-closed contract's foundation: with
// every dynamic toggle left nil, none of their wire keys appear at all. A future accidental
// switch from *bool to bool would make this test fail immediately (a plain bool always
// serializes, defeating "absent ⇒ allow/deny" for every consumer that branches on presence).
func TestDynamicToggles_AbsentWhenNil(t *testing.T) {
	cfg := CaptureConfig{ProtocolVersion: CurrentProtocolVersion, Enabled: true}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{
		"otelEventsEnabled", "httpCaptureTapEnabled", "workloadScanEnabled", "mergeWithGlobal",
	} {
		if strings.Contains(s, key) {
			t.Errorf("wire output must omit %q when nil, got: %s", key, s)
		}
	}
}

// TestDynamicToggles_ExplicitFalseSerializes is the assertion that stops a future
// "simplify *bool -> bool" refactor from silently erasing every explicit false in the fleet:
// omitempty elides only a NIL pointer, not a pointer-to-false, so an explicit false must
// still appear on the wire, distinguishable from "absent."
func TestDynamicToggles_ExplicitFalseSerializes(t *testing.T) {
	cfg := CaptureConfig{
		ProtocolVersion:       CurrentProtocolVersion,
		OTelEventsEnabled:     boolPtr(false),
		HTTPCaptureTapEnabled: boolPtr(false),
		WorkloadScanEnabled:   boolPtr(false),
		MergeWithGlobal:       boolPtr(false),
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"otelEventsEnabled":false`, `"httpCaptureTapEnabled":false`,
		`"workloadScanEnabled":false`, `"mergeWithGlobal":false`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in wire output, got: %s", want, s)
		}
	}
}

// TestDynamicToggles_RoundTripPreservesNilVsFalse is the round-trip half: a decoded absent
// field must come back nil (not a pointer-to-false), and a decoded explicit false must come
// back a pointer-to-false (not nil) -- the two states this whole *bool design exists to keep
// distinguishable, verified through actual JSON marshal/unmarshal rather than asserted on the
// Go struct alone.
func TestDynamicToggles_RoundTripPreservesNilVsFalse(t *testing.T) {
	cfg := CaptureConfig{
		ProtocolVersion: CurrentProtocolVersion,
		// OTelEventsEnabled left nil (absent).
		HTTPCaptureTapEnabled: boolPtr(false), // explicit false.
		// WorkloadScanEnabled left nil (absent).
		MergeWithGlobal: boolPtr(true), // explicit true.
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out CaptureConfig
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.OTelEventsEnabled != nil {
		t.Errorf("OTelEventsEnabled: absent field must round-trip to nil, got %v", *out.OTelEventsEnabled)
	}
	if out.WorkloadScanEnabled != nil {
		t.Errorf("WorkloadScanEnabled: absent field must round-trip to nil, got %v", *out.WorkloadScanEnabled)
	}
	if out.HTTPCaptureTapEnabled == nil {
		t.Fatal("HTTPCaptureTapEnabled: explicit false must round-trip to non-nil")
	} else if *out.HTTPCaptureTapEnabled {
		t.Error("HTTPCaptureTapEnabled: explicit false must round-trip to false, got true")
	}
	if out.MergeWithGlobal == nil {
		t.Fatal("MergeWithGlobal: explicit true must round-trip to non-nil")
	} else if !*out.MergeWithGlobal {
		t.Error("MergeWithGlobal: explicit true must round-trip to true, got false")
	}
}
