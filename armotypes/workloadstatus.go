package armotypes

type WorkloadStatus struct {
	ResourceHash     string   `json:"resourceHash"`
	CustomerGUID     string   `json:"customerGUID"`
	ClusterName      string   `json:"clusterName"`
	IsInternetFacing *bool    `json:"isInternetFacing"`
	RiskFactors      []string `json:"riskFactors"`
}
