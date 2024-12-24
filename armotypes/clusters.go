package armotypes

import "time"

type ClusterInfo struct {
	Cluster        string     `json:"cluster"`
	NodeCount      int        `json:"nodeCount"`
	CPUSum         int        `json:"cpuSum"`
	CloudProvider  string     `json:"cloudProvider"`
	HelmVersion    string     `json:"helmVersion"`
	LastReportTime *time.Time `json:"lastReportTime"`
	IsConnected    bool       `json:"isConnected"`
}
