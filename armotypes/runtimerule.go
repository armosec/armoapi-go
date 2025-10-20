package armotypes

type RuntimeRules struct {
	// Identifiers for future rules
	PortalBase     `json:",inline" bson:",inline"`
	CustomerGUID   string        `json:"customerGUID" bson:"customerGUID"`
	HostIdentifier string        `json:"hostIdentifier,omitempty" bson:"hostIdentifier,omitempty"`
	Rules          []RuntimeRule `json:"rules" bson:"rules"`
}

type RuntimeRule struct {
	Enabled                 bool                  `json:"enabled" yaml:"enabled" bson:"enabled"`
	ID                      string                `json:"id" yaml:"id" bson:"id"`
	Name                    string                `json:"name" yaml:"name" bson:"name"`
	EventType               string                `json:"eventType" bson:"eventType"`
	Description             string                `json:"description" yaml:"description" bson:"description"`
	ViolationMessage        string                `json:"violationMessage" bson:"violationMessage"`
	Expressions             RuleExpressions       `json:"expressions" yaml:"expressions" bson:"expressions"`
	ProfileDependency       ProfileDependency     `json:"profileDependency" yaml:"profileDependency" bson:"profileDependency"`
	Severity                SecurityIssueSeverity `json:"severity" bson:"severity"`
	SeverityScore           int                   `json:"severityScore" bson:"severityScore"`
	SupportPolicy           bool                  `json:"supportPolicy" yaml:"supportPolicy" bson:"supportPolicy"`
	Tags                    []string              `json:"tags" yaml:"tags" bson:"tags"`
	State                   map[string]any        `json:"state,omitempty" yaml:"state,omitempty" bson:"state,omitempty"`
	AgentVersionRequirement string                `json:"agentVersionRequirement" yaml:"agentVersionRequirement" bson:"agentVersionRequirement"`
	IsTriggerAlert          bool                  `json:"isTriggerAlert" yaml:"isTriggerAlert" bson:"isTriggerAlert"`
	MitreTactic             string                `json:"mitreTactic" bson:"mitreTactic"`
	MitreTechnique          string                `json:"mitreTechnique" bson:"mitreTechnique"`
	Category                string                `json:"category" bson:"category"`
	IncidentTypeId          string                `json:"incidentTypeId" bson:"incidentTypeId"`
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
