package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

type PostureExceptionPolicyActions string

const AlertOnly PostureExceptionPolicyActions = "alertOnly"
const Disable PostureExceptionPolicyActions = "disable"

type PolicyType string

const PostureExceptionPolicyType PolicyType = "postureExceptionPolicy"
const VulnerabilityExceptionPolicyType PolicyType = "vulnerabilityExceptionPolicy"

type PostureExceptionPolicy struct {
	PortalBase      `json:",inline" bson:"inline"`
	PolicyType      string                          `json:"policyType,omitempty" bson:"policyType,omitempty"`
	CreationTime    string                          `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Actions         []PostureExceptionPolicyActions `json:"actions,omitempty" bson:"actions,omitempty"`
	Resources       []identifiers.PortalDesignator  `json:"resources" bson:"resources,omitempty"`
	PosturePolicies []PosturePolicy                 `json:"posturePolicies,omitempty" bson:"posturePolicies,omitempty"`
	Reason          *string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate  *time.Time                      `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy       string                          `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}

type PosturePolicy struct {
	FrameworkName string  `json:"frameworkName" bson:"frameworkName"`
	ControlName   string  `json:"controlName,omitempty" bson:"controlName,omitempty"`
	ControlID     string  `json:"controlID,omitempty" bson:"controlID,omitempty"`
	RuleName      string  `json:"ruleName,omitempty" bson:"ruleName,omitempty"`
	BaseScore     float32 `json:"baseScore,omitempty" bson:"baseScore,omitempty"`
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
