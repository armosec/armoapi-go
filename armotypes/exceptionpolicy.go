package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

type BaseExceptionPolicy struct {
	Designators    identifiers.PortalDesignator `json:"designators" bson:"designators,omitempty"`
	PolicyType     string                       `json:"policyType,omitempty" bson:"policyType,omitempty"`
	CreationTime   string                       `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason         string                       `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate *time.Time                   `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy      string                       `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}
