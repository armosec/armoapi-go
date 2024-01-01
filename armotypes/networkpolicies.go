package armotypes

// NetworkPoliciesWorkload is used store information about workloads
// in the customer's clusters related to the NetworkPolicies feature
type NetworkPoliciesWorkload struct {
	Name                       string `json:"name"`
	Kind                       string `json:"kind"`
	CustomerGUID               string `json:"customerGUID"`
	Namespace                  string `json:"namespace"`
	ClusterName                string `json:"cluster"`
	ClusterShortName           string `json:"clusterShortName"`
	NetworkPolicyStatus        int    `json:"networkPolicyStatus"`
	NetworkPolicyStatusMessage string `json:"networkPolicyStatusMessage"`
}

const (
	MissingRuntimeInfo    = 1
	NetworkPolicyRequired = 2
	NetworkPolicyApplied  = 3
)
