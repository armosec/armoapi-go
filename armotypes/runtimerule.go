package armotypes

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/yaml.v3"
)

// copied from kubescape/node-agent/pkg/ruleengine/v1/rule.go
const (
	RuleSeverityNone        = 0
	RuleSeverityLow         = 1
	RuleSeverityMed         = 5
	RuleSeverityHigh        = 8
	RuleSeverityCritical    = 10
	RuleSeveritySystemIssue = 1000
)

func RuleSeverityFromString(severity string) int {
	switch severity {
	case "None":
		return RuleSeverityNone
	case "Low":
		return RuleSeverityLow
	case "Medium":
		return RuleSeverityMed
	case "High":
		return RuleSeverityHigh
	case "Critical":
		return RuleSeverityCritical
	case "System Issue":
		return RuleSeveritySystemIssue
	default:
		return RuleSeverityNone
	}
}

func RuleSeverityToString(severity int) string {
	switch severity {
	case RuleSeverityNone:
		return "None"
	case RuleSeverityLow:
		return "Low"
	case RuleSeverityMed:
		return "Medium"
	case RuleSeverityHigh:
		return "High"
	case RuleSeverityCritical:
		return "Critical"
	case RuleSeveritySystemIssue:
		return "System Issue"
	default:
		if severity < RuleSeverityMed {
			return "Low"
		} else if severity < RuleSeverityHigh {
			return "Medium"
		} else if severity < RuleSeverityCritical {
			return "High"
		}
		return "Unknown"
	}
}

type RuntimeRule struct {
	Enabled     bool            `json:"enabled" yaml:"enabled" bson:"enabled"`
	ID          string          `json:"id" yaml:"id" bson:"id"`
	Name        string          `json:"name" yaml:"name" bson:"name"`
	Description string          `json:"description" yaml:"description" bson:"description"`
	Expressions RuleExpressions `json:"expressions" yaml:"expressions" bson:"expressions"`
	// StateWrites declares what this rule remembers across events. Empty for
	// the overwhelming majority of rules, which are stateless predicates.
	StateWrites             []StateWrite         `json:"stateWrites,omitempty" yaml:"stateWrites,omitempty" bson:"stateWrites,omitempty"`
	ProfileDependency       ProfileDependency    `json:"profileDependency" yaml:"profileDependency" bson:"profileDependency"`
	ProfileDataRequired     *ProfileDataRequired `json:"profileDataRequired,omitempty" yaml:"profileDataRequired,omitempty" bson:"profileDataRequired,omitempty"`
	Severity                int                  `json:"severity" bson:"severity"`
	SeverityString          string               `json:"severityString" bson:"severityString"`
	SupportPolicy           bool                 `json:"supportPolicy" yaml:"supportPolicy" bson:"supportPolicy"`
	Tags                    []string             `json:"tags" yaml:"tags" bson:"tags"`
	State                   map[string]any       `json:"state,omitempty" yaml:"state,omitempty" bson:"state,omitempty"`
	AgentVersionRequirement string               `json:"agentVersionRequirement" yaml:"agentVersionRequirement" bson:"agentVersionRequirement"`
	IsTriggerAlert          bool                 `json:"isTriggerAlert" yaml:"isTriggerAlert" bson:"isTriggerAlert"`
	MitreTactic             string               `json:"mitreTactic" bson:"mitreTactic"`
	MitreTechnique          string               `json:"mitreTechnique" bson:"mitreTechnique"`
	Category                string               `json:"category" bson:"category"`
	IncidentTypeId          string               `json:"incidentTypeId" bson:"incidentTypeId"`
	// Documentation is the rendered markdown documentation for this rule,
	// sourced from the rule's README.md in the rulelibrary repo. Empty for
	// customer-authored rules.
	Documentation string `json:"documentation,omitempty" yaml:"documentation,omitempty" bson:"documentation,omitempty"`
	// Variants holds version-gated alternative bodies for this rule, so one rule ID/name can
	// serve different CEL bodies across agent-capability generations instead of forking
	// identity. The top-level Expressions field above always mirrors the lowest-MinAgentVersion
	// variant, so any reader that doesn't resolve Variants (console, K8s connector, policy
	// validation, incident-enrichment) still gets a working base body. Agent-facing resolver
	// responses omit this field entirely and project the resolved variant into Expressions
	// instead - see ResolveVariantExpressions.
	Variants []RuleVariant `json:"variants,omitempty" yaml:"variants,omitempty" bson:"variants,omitempty"`
}

// RuleVariant is a version-gated alternative rule body. MinAgentVersion is a bare semver (not a
// constraint string like AgentVersionRequirement), so "pick the highest satisfied variant" has a
// well-defined total order - see ValidateVariants and ResolveVariantExpressions.
type RuleVariant struct {
	MinAgentVersion string          `json:"minAgentVersion" yaml:"minAgentVersion" bson:"minAgentVersion"`
	Expressions     RuleExpressions `json:"expressions" yaml:"expressions" bson:"expressions"`
}

type RuleExpressions struct {
	Message        string           `json:"message" yaml:"message" bson:"message"`
	UniqueID       string           `json:"uniqueId" yaml:"uniqueId" bson:"uniqueId"`
	RuleExpression []RuleExpression `json:"ruleExpression" yaml:"ruleExpression" bson:"ruleExpression"`
}

type RuleExpression struct {
	EventType  EventType `json:"eventType" yaml:"eventType" bson:"eventType"`
	Expression string    `json:"expression" yaml:"expression" bson:"expression"`
}

// StateScope names the bucket a state entry belongs to. The scope determines
// which entries a read can see, and the engine — never the rule — resolves the
// concrete scope ID from the event, so cross-tenant leakage is not expressible.
//
// container/pod/node are node-agent scopes; identity is the operator's only
// scope (admission chains are defined by the caller who performed them).
type StateScope string

const (
	StateScopeContainer StateScope = "container"
	StateScopePod       StateScope = "pod"
	StateScopeNode      StateScope = "node"
	StateScopeIdentity  StateScope = "identity"
)

// StateWrite is one declarative write clause: on events of EventType, when the
// When guard holds, remember a marker under Name (optionally discriminated by
// Key) inside Scope for TTL.
//
// Writes are declared rather than expressed as a CEL setter so that remembering
// is decoupled from alerting: a rule may carry a write clause for an event type
// it never alerts on, and the guard is evaluated exactly once, unaffected by
// boolean short-circuiting or CEL's static optimiser.
type StateWrite struct {
	// EventType is the event stream that drives this write. It need not match
	// any eventType in Expressions.RuleExpression — that is what allows a rule
	// to remember on exec and alert on network.
	EventType EventType `json:"eventType" yaml:"eventType" bson:"eventType"`

	// When is a CEL boolean guard. Empty means "always write".
	When string `json:"when,omitempty" yaml:"when,omitempty" bson:"when,omitempty"`

	Scope StateScope `json:"scope" yaml:"scope" bson:"scope"`

	// Name is what kind of fact this is — a string literal, never an
	// expression, so it stays statically analysable and safe as a metric label.
	Name string `json:"name" yaml:"name" bson:"name"`

	// Key is who the fact is about: a CEL string expression, evaluated per
	// event. Optional — omit it for a fact that is true of the whole scope
	// rather than one subject, which yields exactly one entry per scope.
	Key string `json:"key,omitempty" yaml:"key,omitempty" bson:"key,omitempty"`

	// Value carries optional author-supplied extras as CEL expressions keyed by
	// name. Keys must not begin with "_", which is reserved for engine-stamped
	// provenance. Enforced by the consuming engine at rule load.
	Value map[string]any `json:"value,omitempty" yaml:"value,omitempty" bson:"value,omitempty"`

	// TTL is a Go duration string ("10m", "30m"). The consuming engine clamps
	// it to its configured maxTTL at load, so no rule can pin memory
	// indefinitely.
	TTL string `json:"ttl" yaml:"ttl" bson:"ttl"`
}

// ProfileDataPattern declares a single match pattern for a profile-data
// surface. Exactly one of the four fields must be non-empty.
type ProfileDataPattern struct {
	Exact    string `json:"exact,omitempty"    yaml:"exact,omitempty"    bson:"exact,omitempty"`
	Prefix   string `json:"prefix,omitempty"   yaml:"prefix,omitempty"   bson:"prefix,omitempty"`
	Suffix   string `json:"suffix,omitempty"   yaml:"suffix,omitempty"   bson:"suffix,omitempty"`
	Contains string `json:"contains,omitempty" yaml:"contains,omitempty" bson:"contains,omitempty"`
}

type EventType string

const (
	EventTypeExec         EventType = "exec"
	EventTypeOpen         EventType = "open"
	EventTypeCapabilities EventType = "capabilities"
	EventTypeDNS          EventType = "dns"
	EventTypeNetwork      EventType = "network"
	EventTypeSyscall      EventType = "syscall"
	EventTypeSymlink      EventType = "symlink"
	EventTypeHardlink     EventType = "hardlink"
	EventTypeSSH          EventType = "ssh"
	EventTypeHTTP         EventType = "http"
	EventTypeK8sAdmission EventType = "k8s-admission"

	// Node-telemetry event types emitted by node-agent. Added so a stateWrites
	// clause can name any node-agent event stream, not just the subset that
	// happened to be referenced by admission rules.
	EventTypePtrace  EventType = "ptrace"
	EventTypeBPF     EventType = "bpf"
	EventTypeKmod    EventType = "kmod"
	EventTypeUnshare EventType = "unshare"
	EventTypeIoUring EventType = "iouring"
	EventTypeRandomX EventType = "randomx"
	EventTypeProcfs  EventType = "procfs"
	EventTypeFork    EventType = "fork"
	EventTypeExit    EventType = "exit"

	// EventTypeAll is a wildcard used by rule bindings, not a real event stream.
	EventTypeAll EventType = "all"
)

// knownEventTypes is the authoritative set of event types the rule contract
// accepts. Kept in sync with node-agent's utils.EventType.
var knownEventTypes = map[EventType]struct{}{
	EventTypeExec: {}, EventTypeOpen: {}, EventTypeCapabilities: {},
	EventTypeDNS: {}, EventTypeNetwork: {}, EventTypeSyscall: {},
	EventTypeSymlink: {}, EventTypeHardlink: {}, EventTypeSSH: {},
	EventTypeHTTP: {}, EventTypeK8sAdmission: {}, EventTypePtrace: {},
	EventTypeBPF: {}, EventTypeKmod: {}, EventTypeUnshare: {},
	EventTypeIoUring: {}, EventTypeRandomX: {}, EventTypeProcfs: {},
	EventTypeFork: {}, EventTypeExit: {}, EventTypeAll: {},
}

// IsKnownEventType reports whether e is a recognised event type, including the
// EventTypeAll binding wildcard. Rule loaders use this to reject an entry
// naming an event stream that does not exist, rather than silently never
// matching.
//
// For state-write clauses use IsValidStateWriteEventType instead: a write must
// be driven by a concrete stream, and EventTypeAll is not one.
func IsKnownEventType(e EventType) bool {
	_, ok := knownEventTypes[e]
	return ok
}

// IsValidStateWriteEventType reports whether e may drive a StateWrite.
//
// This is stricter than IsKnownEventType: EventTypeAll is a rule-binding
// wildcard rather than an event stream, so it cannot drive a write. Accepting it
// would let `stateWrites.eventType: all` pass validation and then never match a
// concrete event -- a rule that loads cleanly and silently never fires, which is
// the failure mode this contract is built to avoid.
func IsValidStateWriteEventType(e EventType) bool {
	return e != EventTypeAll && IsKnownEventType(e)
}

// ProfileDataField is a tagged union: either All == true (the rule needs every
// entry on this surface) or Patterns is non-empty (the rule needs entries
// matching any pattern). The YAML/JSON form is either the literal string "all"
// or a list of ProfileDataPattern objects.
type ProfileDataField struct {
	All      bool
	Patterns []ProfileDataPattern
}

const profileDataFieldAllSentinel = "all"

func (f ProfileDataField) MarshalYAML() (any, error) {
	if f.All {
		return profileDataFieldAllSentinel, nil
	}
	return f.Patterns, nil
}

func (f *ProfileDataField) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var s string
		if err := node.Decode(&s); err != nil {
			return fmt.Errorf("profileDataField: scalar must be string %q: %w", profileDataFieldAllSentinel, err)
		}
		if s != profileDataFieldAllSentinel {
			return fmt.Errorf("profileDataField: scalar value must be %q, got %q", profileDataFieldAllSentinel, s)
		}
		f.All = true
		f.Patterns = nil
		return nil
	case yaml.SequenceNode:
		var patterns []ProfileDataPattern
		if err := node.Decode(&patterns); err != nil {
			return err
		}
		f.All = false
		f.Patterns = patterns
		return nil
	default:
		return fmt.Errorf("profileDataField: must be string %q or list, got %v", profileDataFieldAllSentinel, node.Kind)
	}
}

func (f ProfileDataField) MarshalJSON() ([]byte, error) {
	if f.All {
		return json.Marshal(profileDataFieldAllSentinel)
	}
	return json.Marshal(f.Patterns)
}

func (f *ProfileDataField) UnmarshalJSON(data []byte) error {
	if string(bytes.TrimSpace(data)) == "null" {
		*f = ProfileDataField{}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s != profileDataFieldAllSentinel {
			return fmt.Errorf("profileDataField: scalar must be %q, got %q", profileDataFieldAllSentinel, s)
		}
		f.All = true
		f.Patterns = nil
		return nil
	}
	var patterns []ProfileDataPattern
	if err := json.Unmarshal(data, &patterns); err != nil {
		return fmt.Errorf("profileDataField: must be string %q or list: %w", profileDataFieldAllSentinel, err)
	}
	f.All = false
	f.Patterns = patterns
	return nil
}

func (f ProfileDataField) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if f.All {
		return bson.MarshalValue(profileDataFieldAllSentinel)
	}
	return bson.MarshalValue(f.Patterns)
}

func (f *ProfileDataField) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bsontype.Null {
		*f = ProfileDataField{}
		return nil
	}
	raw := bson.RawValue{Type: t, Value: data}
	switch t {
	case bsontype.String:
		s, ok := raw.StringValueOK()
		if !ok {
			return fmt.Errorf("profileDataField: bson string decode failed")
		}
		if s != profileDataFieldAllSentinel {
			return fmt.Errorf("profileDataField: bson scalar must be %q, got %q", profileDataFieldAllSentinel, s)
		}
		f.All = true
		f.Patterns = nil
		return nil
	case bsontype.Array:
		var patterns []ProfileDataPattern
		if err := raw.Unmarshal(&patterns); err != nil {
			return err
		}
		f.All = false
		f.Patterns = patterns
		return nil
	default:
		return fmt.Errorf("profileDataField: bson type must be string or array, got %v", t)
	}
}

// ProfileDataRequired declares which subsets of the container profile this rule
// needs at runtime. Each surface field is optional; absence means "this rule
// does not query this surface".
type ProfileDataRequired struct {
	Opens            *ProfileDataField `json:"opens,omitempty"            yaml:"opens,omitempty"            bson:"opens,omitempty"`
	Execs            *ProfileDataField `json:"execs,omitempty"            yaml:"execs,omitempty"            bson:"execs,omitempty"`
	Capabilities     *ProfileDataField `json:"capabilities,omitempty"     yaml:"capabilities,omitempty"     bson:"capabilities,omitempty"`
	Syscalls         *ProfileDataField `json:"syscalls,omitempty"         yaml:"syscalls,omitempty"         bson:"syscalls,omitempty"`
	Endpoints        *ProfileDataField `json:"endpoints,omitempty"        yaml:"endpoints,omitempty"        bson:"endpoints,omitempty"`
	EgressDomains    *ProfileDataField `json:"egressDomains,omitempty"    yaml:"egressDomains,omitempty"    bson:"egressDomains,omitempty"`
	EgressAddresses  *ProfileDataField `json:"egressAddresses,omitempty"  yaml:"egressAddresses,omitempty"  bson:"egressAddresses,omitempty"`
	IngressDomains   *ProfileDataField `json:"ingressDomains,omitempty"   yaml:"ingressDomains,omitempty"   bson:"ingressDomains,omitempty"`
	IngressAddresses *ProfileDataField `json:"ingressAddresses,omitempty" yaml:"ingressAddresses,omitempty" bson:"ingressAddresses,omitempty"`
}

// IsEmpty reports whether every surface field is nil. Used to detect
// "profileDataRequired: {}" — structurally valid YAML that declares nothing.
func (p *ProfileDataRequired) IsEmpty() bool {
	if p == nil {
		return true
	}
	return p.Opens == nil && p.Execs == nil && p.Capabilities == nil &&
		p.Syscalls == nil && p.Endpoints == nil && p.EgressDomains == nil &&
		p.EgressAddresses == nil && p.IngressDomains == nil && p.IngressAddresses == nil
}

// Validate reports schema violations. Single source of truth for both the
// rulelibrary lint and node-agent's load-time check.
func (p *ProfileDataRequired) Validate() error {
	if p == nil {
		return nil
	}
	for name, field := range map[string]*ProfileDataField{
		"opens":            p.Opens,
		"execs":            p.Execs,
		"capabilities":     p.Capabilities,
		"syscalls":         p.Syscalls,
		"endpoints":        p.Endpoints,
		"egressDomains":    p.EgressDomains,
		"egressAddresses":  p.EgressAddresses,
		"ingressDomains":   p.IngressDomains,
		"ingressAddresses": p.IngressAddresses,
	} {
		if field == nil {
			continue
		}
		if err := validateProfileDataField(name, field); err != nil {
			return err
		}
	}
	return nil
}

func validateProfileDataField(name string, f *ProfileDataField) error {
	if f.All && len(f.Patterns) > 0 {
		return fmt.Errorf("profileDataRequired.%s: cannot have both 'all' and pattern list", name)
	}
	if !f.All && len(f.Patterns) == 0 {
		return fmt.Errorf("profileDataRequired.%s: must declare 'all' or at least one pattern", name)
	}
	for i, pat := range f.Patterns {
		filled := 0
		for _, v := range []string{pat.Exact, pat.Prefix, pat.Suffix, pat.Contains} {
			if v != "" {
				filled++
			}
		}
		if filled != 1 {
			return fmt.Errorf("profileDataRequired.%s[%d]: exactly one of {exact, prefix, suffix, contains} must be set, got %d", name, i, filled)
		}
	}
	return nil
}
