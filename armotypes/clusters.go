package armotypes

import "time"

type ClusterInfo struct {
	Cluster        string       `json:"cluster"`
	NodeCount      int          `json:"nodeCount"`
	CPUSum         int          `json:"cpuSum"`
	CloudProvider  string       `json:"cloudProvider"`
	HelmVersion    string       `json:"helmVersion"`
	LastReportTime *time.Time   `json:"lastReportTime"`
	IsConnected    bool         `json:"isConnected"`
	Capabilities   []Capability `json:"capabilities,omitempty"`
}

type Capability struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
