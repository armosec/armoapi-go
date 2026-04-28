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
	Enabled                 bool              `json:"enabled" yaml:"enabled" bson:"enabled"`
	ID                      string            `json:"id" yaml:"id" bson:"id"`
	Name                    string            `json:"name" yaml:"name" bson:"name"`
	Description             string            `json:"description" yaml:"description" bson:"description"`
	Expressions             RuleExpressions   `json:"expressions" yaml:"expressions" bson:"expressions"`
	ProfileDependency       ProfileDependency      `json:"profileDependency" yaml:"profileDependency" bson:"profileDependency"`
	ProfileDataRequired    *ProfileDataRequired   `json:"profileDataRequired,omitempty" yaml:"profileDataRequired,omitempty" bson:"profileDataRequired,omitempty"`
	Severity                int               `json:"severity" bson:"severity"`
	SeverityString          string            `json:"severityString" bson:"severityString"`
	SupportPolicy           bool              `json:"supportPolicy" yaml:"supportPolicy" bson:"supportPolicy"`
	Tags                    []string          `json:"tags" yaml:"tags" bson:"tags"`
	State                   map[string]any    `json:"state,omitempty" yaml:"state,omitempty" bson:"state,omitempty"`
	AgentVersionRequirement string            `json:"agentVersionRequirement" yaml:"agentVersionRequirement" bson:"agentVersionRequirement"`
	IsTriggerAlert          bool              `json:"isTriggerAlert" yaml:"isTriggerAlert" bson:"isTriggerAlert"`
	MitreTactic             string            `json:"mitreTactic" bson:"mitreTactic"`
	MitreTechnique          string            `json:"mitreTechnique" bson:"mitreTechnique"`
	Category                string            `json:"category" bson:"category"`
	IncidentTypeId          string            `json:"incidentTypeId" bson:"incidentTypeId"`
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
)

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
		return fmt.Errorf("profileDataField: must be string %q or list, got null", profileDataFieldAllSentinel)
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
