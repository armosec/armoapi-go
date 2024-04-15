package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

type BaseExceptionPolicy struct {
	PortalBase `json:",inline" bson:"inline"`
	PolicyType string `json:"policyType,omitempty" bson:"policyType,omitempty"`

	// IDs of the policies (SecurityRiskID, ControlID, etc.)
	PolicyIDs      []string                       `json:"policyIDs,omitempty" bson:"policyIDs,omitempty"`
	CreationTime   string                         `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason         string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate *time.Time                     `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy      string                         `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Resources      []identifiers.PortalDesignator `json:"resources" bson:"resources,omitempty"`
}
