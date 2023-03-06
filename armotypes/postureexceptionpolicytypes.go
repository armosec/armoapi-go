package armotypes

import (
	"fmt"
	"strings"
	"time"
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
	Resources       []PortalDesignator              `json:"resources" bson:"resources"`
	PosturePolicies []PosturePolicy                 `json:"posturePolicies,omitempty" bson:"posturePolicies,omitempty"`
	Reason          *string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate  *Date                           `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy       string                          `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}


type Date time.Time

const dtLayout = "02-01-2006"

func (dt *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, _ := time.Parse(dtLayout, s)
	*dt = Date(nt)
	return nil
}
func (dt Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(dt).Format(dtLayout))
	return []byte(stamp), nil
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
