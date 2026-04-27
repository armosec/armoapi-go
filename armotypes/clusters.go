package armotypes

import "time"

type ClusterInfo struct {
	Cluster          string       `json:"cluster"`
	ClusterUID       string       `json:"clusterUID,omitempty"`
	NodeCount        int          `json:"nodeCount"`
	CPUSum           int          `json:"cpuSum"`
	CloudProvider    string       `json:"cloudProvider"`
	CloudRegion      string       `json:"cloudRegion,omitempty"`
	CloudAccountID   string       `json:"cloudAccountID,omitempty"`
	ResourceGroup    string       `json:"resourceGroup,omitempty"`
	HelmVersion      string       `json:"helmVersion"`
	ClusterVersion   string       `json:"clusterVersion"`
	LastReportTime   *time.Time   `json:"lastReportTime,omitempty"`
	LastKeepAlive    *time.Time   `json:"lastKeepAlive,omitempty"`
	CreatedAt        *time.Time   `json:"createdAt,omitempty"`
	IsConnected      bool         `json:"isConnected"`
	Capabilities     []Capability `json:"capabilities,omitempty"`
	Status           string       `json:"status,omitempty"`
	FailedFeatures   []string     `json:"failedFeatures,omitempty"`
	ConnectionTime   *time.Time   `json:"connectionTime,omitempty"`
	StatusChangeTime *time.Time   `json:"statusChangeTime,omitempty"`
}

type Capability struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
