package tlsoffsets

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRecord_JSONRoundTrip_PreservesAllFieldsAndOpaquePayload(t *testing.T) {
	in := Record{
		Target:          "claude",
		Platform:        "linux",
		Arch:            "x86_64",
		Payload:         json.RawMessage(`{"buildID":"abc","machine":62,"read":{"fileOffset":123}}`),
		ProtocolVersion: CurrentProtocolVersion,
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Record
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Every envelope field must round-trip — including ProtocolVersion.
	if out.Target != in.Target || out.Platform != in.Platform || out.Arch != in.Arch || out.ProtocolVersion != in.ProtocolVersion {
		t.Fatalf("envelope fields not preserved:\n got  %+v\n want %+v", out, in)
	}

	// The opaque payload must survive BYTE-FOR-BYTE (this package never interprets it).
	// Compact both sides so only insignificant whitespace is ignored.
	var inC, outC bytes.Buffer
	if err := json.Compact(&inC, in.Payload); err != nil {
		t.Fatalf("compact in: %v", err)
	}
	if err := json.Compact(&outC, out.Payload); err != nil {
		t.Fatalf("compact out: %v", err)
	}
	if !bytes.Equal(inC.Bytes(), outC.Bytes()) {
		t.Fatalf("opaque payload not preserved:\n got  %s\n want %s", outC.Bytes(), inC.Bytes())
	}
}

func TestRecord_JSONTagsAreCamelCase(t *testing.T) {
	b, err := json.Marshal(Record{
		Target:          "claude",
		Platform:        "linux",
		Arch:            "x86_64",
		Payload:         json.RawMessage(`{}`),
		ProtocolVersion: CurrentProtocolVersion,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	// Assert the COMPLETE camelCase key set (every field, so a wrong tag is caught).
	for _, key := range []string{`"target":`, `"platform":`, `"arch":`, `"payload":`, `"protocolVersion":`} {
		if !strings.Contains(s, key) {
			t.Fatalf("missing expected camelCase key %s in %s", key, s)
		}
	}
}
