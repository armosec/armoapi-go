package armotypes

import (
	"time"
)

type BaseExceptionPolicy struct {
	ExceptionPolicyHash string `json:"exceptionPolicyHash,omitempty" bson:"exceptionPolicyHash,omitempty"`
	PolicyType          string `json:"policyType,omitempty" bson:"policyType,omitempty"`

	// IDs of the policies (SecurityRiskID, ControlID, etc.)
	PolicyIDs                []string   `json:"policyIDs,omitempty" bson:"policyIDs,omitempty"`
	CreationTime             string     `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason                   string     `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationTime           *time.Time `json:"expirationTime,omitempty" bson:"expirationTime,omitempty"`
	CreatedBy                string     `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	ExceptionPolicyResources []Resource `json:"exceptionPolicyResources,omitempty" bson:"exceptionPolicyResources,omitempty"`
}
