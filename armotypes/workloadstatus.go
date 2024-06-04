package armotypes

type WorkloadStatus struct {
	ResourceHash     string `json:"resourceHash"`
	CustomerGUID     string `json:"customerGUID"`
	IsInternetFacing *bool  `json:"isInternetFacing"`
}
