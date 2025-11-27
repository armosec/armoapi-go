package armotypes

type WorkloadStatus struct {
	ResourceHash      string   `json:"resourceHash"`
	CustomerGUID      string   `json:"customerGUID"`
	ClusterName       string   `json:"clusterName"`
	IsInternetFacing  *bool    `json:"isInternetFacing"`
	AiClientProviders []string `json:"aiClientProviders"`
	AiServerProviders []string `json:"aiServerProviders"`
	RiskFactors       []string `json:"riskFactors"`
}
