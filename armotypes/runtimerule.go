package armotypes

type RuntimeRules struct {
	// Identifiers for future rules
	PortalBase     `json:",inline" bson:",inline"`
	CustomerGUID   string        `json:"customerGUID" bson:"customerGUID"`
	HostIdentifier string        `json:"hostIdentifier,omitempty" bson:"hostIdentifier,omitempty"`
	Rules          []RuntimeRule `json:"rules" bson:"rules"`
}

type RuntimeRule struct {
	Enabled           bool              `json:"enabled" yaml:"enabled" bson:"enabled"`
	ID                string            `json:"id" yaml:"id" bson:"id"`
	Name              string            `json:"name" yaml:"name" bson:"name"`
	Description       string            `json:"description" yaml:"description" bson:"description"`
	Expressions       RuleExpressions   `json:"expressions" yaml:"expressions" bson:"expressions"`
	ProfileDependency ProfileDependency `json:"profile_dependency" yaml:"profile_dependency" bson:"profile_dependency"`
	Severity          int               `json:"severity" yaml:"severity" bson:"severity"`
	SupportPolicy     bool              `json:"support_policy" yaml:"support_policy" bson:"support_policy"`
	Tags              []string          `json:"tags" yaml:"tags" bson:"tags"`
	State             map[string]any    `json:"state,omitempty" yaml:"state,omitempty" bson:"state,omitempty"`
}

type RuleExpressions struct {
	Message        string           `json:"message" yaml:"message" bson:"message"`
	UniqueID       string           `json:"unique_id" yaml:"unique_id" bson:"unique_id"`
	RuleExpression []RuleExpression `json:"rule_expression" yaml:"rule_expression" bson:"rule_expression"`
}

type RuleExpression struct {
	EventType  string `json:"event_type" yaml:"event_type" bson:"event_type"` // TODO: change to enum
	Expression string `json:"expression" yaml:"expression" bson:"expression"`
}
