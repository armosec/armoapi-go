package armotypes

type WorkloadStatus struct {
	ResourceHash     string `json:"resourceHash,omitempty"`
	CustomerGUID     string `json:"customerGUID,omitempty"`
	IsInternetFacing *bool  `json:"isInternetFacing,omitempty"`
}
