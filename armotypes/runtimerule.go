package armotypes

type RuntimeRules struct {
	// Identifiers for future rules
	CustomerGUID   string        `json:"customerGUID,omitempty" bson:"customerGUID,omitempty"`
	HostIdentifier string        `json:"hostIdentifier,omitempty" bson:"hostIdentifier,omitempty"`
	Rules          []RuntimeRule `json:"rules" bson:"rules"`
}

type RuntimeRule struct {
	Parameters map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
	RuleID     string                 `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
	RuleName   string                 `json:"ruleName,omitempty" bson:"ruleName,omitempty"`
	RuleTags   []string               `json:"ruleTags,omitempty" bson:"ruleTags,omitempty"`
	Severity   string                 `json:"severity,omitempty" bson:"severity,omitempty"`
}
