package armotypes

import "time"

type Inventory struct {
	WorkloadName       string     `json:"workloadName"`
	Kind               string     `json:"kind"` // will be deprecated in the future after type is introduced
	Type               string     `json:"type"`
	Cluster            string     `json:"cluster"`
	AccountID          string     `json:"accountId"`
	Region             string     `json:"region"`
	Provider           string     `json:"provider"`
	Namespace          string     `json:"namespace"`
	CreationTimestamp  *time.Time `json:"creationTimestamp,omitempty"`
	CompletionStatus   string     `json:"completionStatus,omitempty"`
	Status             string     `json:"status,omitempty"`
	LearningPeriod     string     `json:"learningPeriod,omitempty"`
	RiskFactors        []string   `json:"riskFactors,omitempty"`
	LearningPercentage *int       `json:"learningPercentage,omitempty"`
	HostName           string     `json:"hostName,omitempty"`
	AccountName        string     `json:"accountName,omitempty"`
	InstanceID         string     `json:"instanceId,omitempty"`
	HostIPs            []string   `json:"hostIPs,omitempty"`
	Tags               []string   `json:"tags,omitempty"`
	LaunchType         string     `json:"launchType,omitempty"`
	HostType           string     `json:"hostType,omitempty"`
	ResourceHash       string     `json:"resourceHash,omitempty"`
	// IsAgentic is the binary agentic verdict for the inventory/discovery badge.
	// It is derived from workload_statuses providers (ai_client_providers /
	// ai_server_providers) via armotypes.IsAgentic and is sandbox-table-independent
	// — it does NOT depend on the ai_sandboxes tables, so the discovery badge keeps
	// working regardless of the ai_sandboxes agentic-verdict migration state.
	IsAgentic bool `json:"isAgentic,omitempty"`
}
