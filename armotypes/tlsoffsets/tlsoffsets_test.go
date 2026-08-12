package tlsoffsets

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRecord_JSONRoundTrip_PreservesOpaquePayload(t *testing.T) {
	in := Record{
		Target:          "claude",
		Platform:        "linux",
		Arch:            "x86_64",
		Payload:         json.RawMessage(`{"buildID":"abc","machine":62}`),
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
	if out.Target != "claude" || out.Platform != "linux" || out.Arch != "x86_64" {
		t.Fatalf("envelope fields lost: %+v", out)
	}
	// Payload must survive verbatim (opaque to this package).
	var got map[string]any
	if err := json.Unmarshal(out.Payload, &got); err != nil {
		t.Fatalf("payload not preserved: %v", err)
	}
	if got["buildID"] != "abc" {
		t.Fatalf("payload content lost: %v", got)
	}
}

func TestRecord_JSONTagsAreCamelCase(t *testing.T) {
	b, err := json.Marshal(Record{Target: "claude", Payload: json.RawMessage(`null`)})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"target":"claude"`) || !strings.Contains(s, `"payload":`) {
		t.Fatalf("unexpected tags: %s", s)
	}
}
