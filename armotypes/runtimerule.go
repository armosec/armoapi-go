package armotypes

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
	ProfileDependency       ProfileDependency `json:"profileDependency" yaml:"profileDependency" bson:"profileDependency"`
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
	UniqueID       string           `json:"uniqueID" yaml:"uniqueID" bson:"uniqueID"`
	RuleExpression []RuleExpression `json:"ruleExpression" yaml:"ruleExpression" bson:"ruleExpression"`
}

type RuleExpression struct {
	EventType  EventType `json:"eventType" yaml:"eventType" bson:"eventType"`
	Expression string    `json:"expression" yaml:"expression" bson:"expression"`
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
