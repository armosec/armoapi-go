package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

const (
	// SecurityRiskPolicy - policy for security risks
	SecurityRiskExceptionPolicyType PolicyType = "securityRiskExceptionPolicy"

	// RuntimeIncidentPolicy - policy for runtime incidents
	RuntimeIncidentExceptionPolicyType PolicyType = "runtimeIncidentExceptionPolicy"

	// CSPM - policy for CSPM
	CSPMExceptionPolicyType PolicyType = "cspmExceptionPolicy"
)

type BaseExceptionPolicy struct {
	PortalBase `json:",inline" bson:"inline"`
	PolicyType PolicyType `json:"policyType,omitempty" bson:"policyType,omitempty"`

	// IDs of the policies (SecurityRiskID, ControlID, etc.)
	PolicyIDs      []string                       `json:"policyIDs,omitempty" bson:"policyIDs,omitempty"`
	CreationTime   string                         `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason         string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate *time.Time                     `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy      string                         `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Resources      []identifiers.PortalDesignator `json:"resources,omitempty" bson:"resources,omitempty"`
}
