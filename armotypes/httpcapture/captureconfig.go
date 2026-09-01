package httpcapture

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/armotypes/tlsoffsets"
)

// CaptureConfigResponse is the merged capture config the node-agent pulls: the
// HTTP-capture policy (CaptureConfig, inlined at the top level) plus the TLS offset-map
// section. cadashboardbe produces it (stitches capture + offsets); the node-agent
// consumes it. It is defined once, here, so producer and consumer cannot drift on the
// envelope shape or the "tlsOffsets" key.
//
// Renamed from SandboxConfigResponse (ADR 0009 §4): the type is wire-neutral —
// CaptureConfig is inlined via anonymous-field promotion, so the JSON shape on the wire
// is unchanged by this rename.
//
// It is a SERVE/wire type only. CaptureConfig remains the standalone capture contract —
// it is what config-service stores (v1_http_capture_config) and is NOT widened with
// offsets; the offsets live in their own type (tlsoffsets.Record) and their own
// collection. The two only ride together in this response.
type CaptureConfigResponse struct {
	CaptureConfig `json:",inline"`

	// TLSOffsets is the offset-map section, kept as RAW JSON (a map[identity]tlsoffsets.Record
	// once decoded) so a malformed/oversized offsets section can never fail the fail-closed
	// capture decode — the agent decodes it leniently, per record.
	TLSOffsets json.RawMessage `json:"tlsOffsets,omitempty"`
}

// SetTLSOffsets marshals the offset records into the raw TLSOffsets section. Producer-side
// (cadashboardbe) helper, so the offsets encoding lives in one place. An empty/nil map
// clears the section (omitted from the wire).
func (r *CaptureConfigResponse) SetTLSOffsets(offsets map[string]tlsoffsets.Record) error {
	if len(offsets) == 0 {
		r.TLSOffsets = nil
		return nil
	}
	b, err := json.Marshal(offsets)
	if err != nil {
		return err
	}
	r.TLSOffsets = b
	return nil
}
