// Package aisandboxconfig is the merged AI-Sandbox config the node-agent pulls: the
// HTTP-capture policy plus the TLS offset-map section. cadashboardbe produces it
// (stitches capture config + offsets); the node-agent consumes it. It is defined once,
// here, so producer and consumer cannot drift on the envelope shape or the "tlsOffsets"
// key — the same single-source-of-truth reason armotypes/httpcapture exists.
//
// The two sections stay independent types (httpcapture.CaptureConfig and
// tlsoffsets.Record); this package only defines how they ride the wire together. They
// are NOT combined in storage — config-service keeps them in separate collections.
package aisandboxconfig

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/armotypes/httpcapture"
	"github.com/armosec/armoapi-go/armotypes/tlsoffsets"
)

// Response is the wire shape served to the agent (via cadashboardbe → careportsreceiver).
// CaptureConfig is inlined at the top level (the unchanged capture contract); TLSOffsets
// is a separate section kept as RAW JSON so a malformed/oversized offsets section can
// never fail the (fail-closed) capture decode — the agent decodes it leniently, per
// record. The map it encodes is keyed by record identity (bare build-id or sha256:<hex>).
type Response struct {
	httpcapture.CaptureConfig `json:",inline"`

	TLSOffsets json.RawMessage `json:"tlsOffsets,omitempty"`
}

// SetTLSOffsets marshals the offset records into the raw TLSOffsets section. This is the
// producer-side (cadashboardbe) helper, so the offsets encoding lives in one place. An
// empty/nil map clears the section (omitted from the wire).
func (r *Response) SetTLSOffsets(offsets map[string]tlsoffsets.Record) error {
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
