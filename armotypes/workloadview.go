package armotypes

import "time"

type WorkloadViews struct {
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
	// IsAgentic is the binary agentic-classification badge for the inventory.
	// It is derived from workload_statuses (ai_client_providers /
	// ai_server_providers) via armotypes.IsAgentic — it does NOT depend on the
	// AI-Sandbox tables, so the dashboard badge keeps working regardless of the
	// ai_sandboxes agentic-verdict migration state.
	IsAgentic bool `json:"isAgentic,omitempty"`
}
