package armotypes

type NetworkPolicyStatus int

const (
	MissingRuntimeInfo    NetworkPolicyStatus = 1
	NetworkPolicyRequired NetworkPolicyStatus = 2
	NetworkPolicyApplied  NetworkPolicyStatus = 3
)

// NetworkPoliciesWorkload is used store information about workloads
// in the customer's clusters related to the NetworkPolicies feature
type NetworkPoliciesWorkload struct {
	Name                       string              `json:"name"`
	Kind                       string              `json:"kind"`
	CustomerGUID               string              `json:"customerGUID"`
	Namespace                  string              `json:"namespace"`
	ClusterName                string              `json:"cluster"`
	ClusterShortName           string              `json:"clusterShortName"`
	NetworkPolicyStatus        NetworkPolicyStatus `json:"networkPolicyStatus"`
	NetworkPolicyStatusMessage string              `json:"networkPolicyStatusMessage"`
}
