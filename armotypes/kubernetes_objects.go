package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

// KubernetesObject represents a single Kubernetes object, either native or kubescape CRD
type KubernetesObject struct {
	Designators       identifiers.PortalDesignator `json:"designators"`
	ResourceHash      string                       `json:"resourceHash"`
	ResourceObjectRef string                       `json:"resourceObjectRef"`
	ResourceVersion   string                       `json:"resourceVersion"`
	Checksum          string                       `json:"checksum"`
	CreationTimestamp time.Time                    `json:"creationTimestamp"`

	OwnerReferenceName string `json:"ownerReferenceName"`
	OwnerReferenceKind string `json:"ownerReferenceKind"`

	// related only to kubescape DRDs.
	RelatedName            string `json:"relatedName"`
	RelatedKind            string `json:"relatedKind"`
	RelatedAPIGroup        string `json:"relatedAPIGroup"`
	RelatedNamespace       string `json:"relatedNamespace"`
	RelatedAPIVersion      string `json:"relatedAPIVersion"`
	RelatedResourceVersion string `json:"relatedResourceVersion"`
	Status                 string `json:"status"`
	CompletionStatus       string `json:"completionStatus"`

	NetworkPolicyStatus NetworkPolicyStatus `json:"networkPolicyStatus"`

	Labels map[string]string `json:"labels"`

	// pod selector labels of network policies
	NetworkPolicyPodSelectorLabels map[string]string `json:"podSelectorLabels,omitempty"`

	// pod spec labels of workloads
	PodSpecLabels map[string]string `json:"podSpecLabels,omitempty"`

	// pod selector labels of services
	ServicePodSelectorLabels map[string]string `json:"servicePodSelectorLabels,omitempty"`

	// roleRef of RoleBinding
	RoleBindingRoleRef *RoleBindingRoleRef `json:"roleRef,omitempty"`

	// subjects of RoleBinding
	RoleBindingSubjects []RoleBindingSubject `json:"subjects,omitempty"`
}

type Resource struct {
	K8sResourceHash  string `json:"k8sResourceHash,omitempty" bson:"k8sResourceHash,omitempty"`
	Cluster          string `json:"cluster,omitempty" bson:"cluster,omitempty"`
	ClusterShortName string `json:"clusterShortName"`
	Namespace        string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	Kind             string `json:"kind,omitempty" bson:"kind,omitempty"`
	Name             string `json:"name,omitempty" bson:"name,omitempty"`
}

type RoleBindingSubject struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
}

type RoleBindingRoleRef struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
}
