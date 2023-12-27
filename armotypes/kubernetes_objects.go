package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

type NetworkPolicyStatus int

const (
	StatusNetworkPolicyApplied    NetworkPolicyStatus = 1
	StatusNetworkPolicyNotApplied NetworkPolicyStatus = 2
	StatusNetworkPolicyUknown     NetworkPolicyStatus = 3
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
	NetworkPolicyStatus    string `json:"networkPolicyStatus"` // DEPRECATED

	NetworkPolicyAppliedCustomer  bool `json:"networkPolicyAppliedCustomer"`
	NetworkPolicyAppliedKubescape bool `json:"networkPolicyAppliedKubescape"`
	NetworkPolicyStatusKnown      bool `json:"networkPolicyStatusKnown"`

	Labels map[string]string `json:"labels"`
}

func (ko *KubernetesObject) GetNetworkPolicyStatus() NetworkPolicyStatus {
	if !ko.NetworkPolicyStatusKnown {
		return StatusNetworkPolicyUknown
	}

	if ko.NetworkPolicyAppliedCustomer || ko.NetworkPolicyAppliedKubescape {
		return StatusNetworkPolicyApplied
	}

	return StatusNetworkPolicyNotApplied
}
