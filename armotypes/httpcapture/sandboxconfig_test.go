package httpcapture

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/armosec/armoapi-go/armotypes/tlsoffsets"
)

func TestSandboxConfigResponse_CaptureInline_OffsetsSeparateSection(t *testing.T) {
	var r SandboxConfigResponse
	r.CaptureConfig = CaptureConfig{ProtocolVersion: CurrentProtocolVersion, Enabled: true}
	if err := r.SetTLSOffsets(map[string]tlsoffsets.Record{
		"bid123": {Target: "claude", Platform: "linux", Arch: "x86_64", Payload: json.RawMessage(`{"buildID":"bid123"}`)},
	}); err != nil {
		t.Fatalf("SetTLSOffsets: %v", err)
	}

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	// Capture fields flatten to the top level; offsets ride a sibling "tlsOffsets" key.
	if !strings.Contains(s, `"enabled":true`) || !strings.Contains(s, `"tlsOffsets":`) {
		t.Fatalf("unexpected wire shape: %s", s)
	}

	var out SandboxConfigResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.Enabled {
		t.Fatalf("capture section lost: %+v", out.CaptureConfig)
	}
	// TLSOffsets is opaque here; decoding it is the consumer's job.
	var offsets map[string]tlsoffsets.Record
	if err := json.Unmarshal(out.TLSOffsets, &offsets); err != nil {
		t.Fatalf("offsets section not preserved: %v", err)
	}
	if offsets["bid123"].Target != "claude" {
		t.Fatalf("offset record lost: %+v", offsets)
	}
}

func TestSandboxConfigResponse_NoOffsets_OmitsSection(t *testing.T) {
	var r SandboxConfigResponse
	r.CaptureConfig = CaptureConfig{ProtocolVersion: CurrentProtocolVersion}
	if err := r.SetTLSOffsets(nil); err != nil {
		t.Fatalf("SetTLSOffsets(nil): %v", err)
	}
	b, _ := json.Marshal(r)
	if strings.Contains(string(b), "tlsOffsets") {
		t.Fatalf("empty offsets must be omitted: %s", b)
	}
}

// SetTLSOffsets surfaces an invalid record payload: json.Marshal validates the embedded
// json.RawMessage and returns an error, which SetTLSOffsets propagates without setting the
// section (so a broken payload can never produce a malformed wire response).
func TestSandboxConfigResponse_SetTLSOffsets_RejectsInvalidPayload(t *testing.T) {
	var r SandboxConfigResponse
	err := r.SetTLSOffsets(map[string]tlsoffsets.Record{
		"bad": {Target: "claude", Payload: json.RawMessage("{not valid json")},
	})
	if err == nil {
		t.Fatal("expected an error for an invalid record payload")
	}
	if r.TLSOffsets != nil {
		t.Fatalf("TLSOffsets must remain unset when SetTLSOffsets errors, got %s", r.TLSOffsets)
	}
}
