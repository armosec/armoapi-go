// Package tlsoffsets defines the wire+store envelope for one TLS-capture offset
// record distributed to the node-agent (TLS Offset-Map Distribution). The envelope
// is shared by config-service (store), cadashboardbe (serve/merge), and the
// node-agent (consume) so the queryable fields have ONE source of truth; Payload is
// opaque here and understood only by the named Target's decoder in the agent.
package tlsoffsets

import "encoding/json"

// Record is one offset record. Target/Platform/Arch are the queryable, top-level
// envelope fields config-service filters on (InnerFilters) — no consumer parses
// Payload to filter. Payload is the target-specific offset data, opaque to store
// and transport (agent claude/opencode ⇒ internal/verify OffsetRecord JSON).
//
// Deliberately does NOT embed armotypes.PortalBase, mirroring httpcapture.CaptureConfig:
// this is an agent-facing wire+store envelope, so it stays minimal and free of portal
// metadata. The GUID/name/attributes are added one layer up by config-service's
// types.TLSOffsetRecord (PortalBase + Record), exactly as types.HTTPCaptureConfig wraps
// httpcapture.CaptureConfig. A record's identity is the config-service doc GUID (the
// served-map key), by design not a field here.
type Record struct {
	Target   string          `json:"target" bson:"target"`
	Platform string          `json:"platform" bson:"platform"`
	Arch     string          `json:"arch" bson:"arch"`
	Payload  json.RawMessage `json:"payload" bson:"payload"`
}
