package armotypes

type NetworkPolicyStatus int

type MissingRuntimeInfoReason int

const (
	MissingRuntimeInfo    NetworkPolicyStatus = 1
	NetworkPolicyRequired NetworkPolicyStatus = 2
	NetworkPolicyApplied  NetworkPolicyStatus = 3
)

// MissingRuntimeInfoReason is used to store the reason why the runtime information is missing
const (
	UnknownReason            MissingRuntimeInfoReason = 0
	RestartRequired          MissingRuntimeInfoReason = 1
	UnscheduledNodeAgentPods MissingRuntimeInfoReason = 2
	IncompatibleKernel       MissingRuntimeInfoReason = 3
	RuncNotFound             MissingRuntimeInfoReason = 4
)

// NetworkPoliciesWorkload is used store information about workloads
// in the customer's clusters related to the NetworkPolicies feature
type NetworkPoliciesWorkload struct {
	ResourceHash               string                   `json:"resourceHash"`
	Name                       string                   `json:"name"`
	Kind                       string                   `json:"kind"`
	CustomerGUID               string                   `json:"customerGUID"`
	Namespace                  string                   `json:"namespace"`
	ClusterName                string                   `json:"cluster"`
	ClusterShortName           string                   `json:"clusterShortName"`
	AppliedNetworkPolicyType   string                   `json:"appliedNetworkPolicyType"`
	NetworkPolicyStatus        NetworkPolicyStatus      `json:"networkPolicyStatus"`
	NetworkPolicyStatusMessage string                   `json:"networkPolicyStatusMessage"`
	MissingRuntimeInfoReason   MissingRuntimeInfoReason `json:"missingRuntimeInfoReason"`
}
