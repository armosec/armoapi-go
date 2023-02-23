package armotypes

import "time"

type PostureExceptionPolicyActions string

const AlertOnly PostureExceptionPolicyActions = "alertOnly"
const Disable PostureExceptionPolicyActions = "disable"

type PostureExceptionPolicy struct {
	PortalBase      `json:",inline" bson:"inline"`
	PolicyType      string                          `json:"policyType,omitempty" bson:"policyType,omitempty"`
	CreationTime    string                          `json:"creationTime" bson:"creationTime"`
	Actions         []PostureExceptionPolicyActions `json:"actions" bson:"actions"`
	Resources       []PortalDesignator              `json:"resources" bson:"resources"`
	PosturePolicies []PosturePolicy                 `json:"posturePolicies" bson:"posturePolicies"`
	Reason          string                          `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate  *time.Time                      `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy       string                          `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}

type PosturePolicy struct {
	FrameworkName string `json:"frameworkName" bson:"frameworkName"`
	ControlName   string `json:"controlName,omitempty" bson:"controlName,omitempty"`
	ControlID     string `json:"controlID,omitempty" bson:"controlID,omitempty"`
	RuleName      string `json:"ruleName,omitempty" bson:"ruleName,omitempty"`
}

func (exceptionPolicy *PostureExceptionPolicy) IsAlertOnly() bool {
	if exceptionPolicy.IsDisable() {
		return false
	}

	for i := range exceptionPolicy.Actions {
		if exceptionPolicy.Actions[i] == AlertOnly {
			return true
		}
	}
	return false
}
func (exceptionPolicy *PostureExceptionPolicy) IsDisable() bool {
	for i := range exceptionPolicy.Actions {
		if exceptionPolicy.Actions[i] == Disable {
			return true
		}
	}
	return false
}
