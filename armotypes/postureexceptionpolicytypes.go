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

/*
swagger:route POST /api/v1/postureExceptionPolicy
Create new exception

### Query params
`customerGUID` - required
`recommendationGUID` - mark this recommendation as solved

### Request body object
Exception object as described above

### Response body object
Exception object as described above
*/

/*
swagger:route GET /api/v1/postureExceptionPolicy
Returns exception object or list of objects

### Query params
`customerGUID` - required (Returns list of all policies of the given customer)
`policyName`
`policyGUID`
`list` - return just list of names of exceptions of current customer

#### Get all exceptions of the customer
/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7

#### Get filtered exceptions of the customer
/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&scope.cluster=cluster-name&scope.namespace&posturePolicies.frameworkName=NSA
/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7?name=exceptionPolicyName

#### Get specific policy
/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&policyName=reg-policy1
/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&policyGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7

### Response body object
Exception object as described above
*/

/*
swagger:route PUT /api/v1/postureExceptionPolicy
Updating existing exception

### Query params
`customerGUID` - required

### Request body object
Exception object as described above

### Response body object
Exception object as described above
*/

/*
swagger:route DELETE /api/v1/postureExceptionPolicy
Deleting existing exception

### Query params
`customerGUID` - required
`policyName` - Can appear multiple times to delete multiple policies
`policyGUID` - Can appear multiple times to delete multiple policies

### DELETE multiple exceptions
`/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&policyName=reg-policy1&policyName=bas-policy2`
`/api/v1/postureExceptionPolicy?customerGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&policyGUID=31fb54a9-6e8f-4289-8506-f4e875ac19f7&policyGUID=c62a6397-4777-4410-949e-d99c2efb1f79`
*/
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
