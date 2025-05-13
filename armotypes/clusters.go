package armotypes

import "time"

type ClusterInfo struct {
	Cluster        string       `json:"cluster"`
	NodeCount      int          `json:"nodeCount"`
	CPUSum         int          `json:"cpuSum"`
	CloudProvider  string       `json:"cloudProvider"`
	HelmVersion    string       `json:"helmVersion"`
	ClusterVersion string       `json:"clusterVersion"`
	LastReportTime *time.Time   `json:"lastReportTime,omitempty"`
	LastKeepAlive  *time.Time   `json:"lastKeepAlive,omitempty"`
	CreatedAt      *time.Time   `json:"createAt,omitempty"`
	IsConnected    bool         `json:"isConnected"`
	Capabilities   []Capability `json:"capabilities,omitempty"`
	Status         string       `json:"status,omitempty"`
	FailedFeatures []string     `json:"failedFeatures,omitempty"`
}

type Capability struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
