package armotypes

type WorkloadStatus struct {
	WorkloadName        string   `json:"workloadName"`
	Namespace           string   `json:"namespace"`
	Kind                string   `json:"kind"`
	ClusterName         string   `json:"clusterName"`
	CustomerGUID        string   `json:"customerGUID"`
	ResourceHash        string   `json:"resourceHash"`
	ResourceHashNumeric float64  `json:"resourceHashNumeric"`
	workloadLabelsCount int      `json:"workloadLabelsCount"`
	riskFactorsCount    int      `json:"riskFactorsCount"`
	WorkloadLabels      []string `json:"workloadLabels"`
	RiskFactors         []string `json:"riskFactors"`
	IsInternetFacing    *bool    `json:"isInternetFacing"`
}
