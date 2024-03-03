package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

type ExceptionPolicyResource struct {
	ExceptionPolicyHash string `json:"exceptionPolicyHash,omitempty" bson:"exceptionPolicyHash,omitempty"`
	K8sResourceHash     string `json:"k8sResourceHash,omitempty" bson:"k8sResourceHash,omitempty"`
	Cluster             string `json:"cluster,omitempty" bson:"cluster,omitempty"`
	Namespace           string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	Kind                string `json:"kind,omitempty" bson:"kind,omitempty"`
	Name                string `json:"name,omitempty" bson:"name,omitempty"`
}

type BaseExceptionPolicy struct {
	Designators              []identifiers.PortalDesignator `json:"designators" bson:"designators,omitempty"`
	PolicyType               string                         `json:"policyType,omitempty" bson:"policyType,omitempty"`
	Policies                 []string                       `json:"policies,omitempty" bson:"policies,omitempty"`
	CreationTime             string                         `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason                   string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationTime           *time.Time                     `json:"expirationTime,omitempty" bson:"expirationTime,omitempty"`
	CreatedBy                string                         `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	ExceptionPolicyResources []ExceptionPolicyResource      `json:"exceptionPolicyResources,omitempty" bson:"exceptionPolicyResources,omitempty"`
}
