package cdr

import "time"

// CdrRuleOrigin distinguishes centrally-managed rules from customer-authored ones. Customer
// authoring is not yet supported; the field exists so the delivery contract can carry custom rules
// in future without a wire change. The server always sets it (managed or custom); a consumer should
// treat an absent value as managed.
type CdrRuleOrigin string

const (
	// CdrRuleOriginManaged is a rule maintained centrally by ARMO.
	CdrRuleOriginManaged CdrRuleOrigin = "managed"
	// CdrRuleOriginCustom is a customer-authored rule.
	CdrRuleOriginCustom CdrRuleOrigin = "custom"
)

// CdrRule is the sensor-facing projection of a CDR detection rule served to the in-account
// collector over the managed rule-delivery endpoint. It is deliberately slim and decoupled from the
// internal runtime-rule storage model so the wire contract stays stable as that storage evolves.
// Its fields align with the alert a match produces, so rule metadata passes straight through to the
// emitted CdrAlert.
type CdrRule struct {
	// RuleID is the unique identifier of the rule.
	RuleID string `json:"ruleID"`
	// Name is the human-readable rule name.
	Name string `json:"name"`
	// Description is the rule description.
	Description string `json:"description,omitempty"`
	// Service is the target log stream the rule evaluates (e.g. activitylogs, cloudtrail).
	Service CloudService `json:"service"`
	// Expression is the CEL rule text the collector compiles and evaluates.
	Expression string `json:"expression"`
	// Priority is the severity of the rule (matches CdrAlert.Priority for pass-through).
	Priority string `json:"priority,omitempty"`
	// MitreTactic is the MITRE ATT&CK tactic.
	MitreTactic string `json:"mitreTactic,omitempty"`
	// MitreTechnique is the MITRE ATT&CK technique.
	MitreTechnique string `json:"mitreTechnique,omitempty"`
	// Tags is the rule tags.
	Tags []string `json:"tags,omitempty"`
	// Message is the alert-message template emitted on a match.
	Message string `json:"message,omitempty"`
	// UniqueID is the templated dedup key for a match.
	UniqueID string `json:"uniqueID,omitempty"`
	// Origin is managed (default) or custom.
	Origin CdrRuleOrigin `json:"origin,omitempty"`
}

// CdrRuleBundle is the response body of the managed rule-delivery endpoint: the provider-scoped CEL
// rule set for a tenant plus an opaque content version for change detection. The collector
// downloads it, authenticated with its per-tenant access key, instead of shipping rules baked into
// the deployed image. Version is emitted as the response ETag; the collector echoes it via
// If-None-Match to get a cheap 304 when the set is unchanged.
type CdrRuleBundle struct {
	// Version is an opaque content hash of the served rule set; it changes if and only if the set
	// changes. Clients must treat it as opaque.
	Version string `json:"version"`
	// Provider is the cloud provider the bundle is scoped to.
	Provider CloudProvider `json:"provider"`
	// GeneratedAt is when the bundle was assembled (informational; not part of Version).
	GeneratedAt time.Time `json:"generatedAt"`
	// Rules is the enabled, provider-scoped CDR rule set for the tenant.
	Rules []CdrRule `json:"rules"`
}
