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
	// SandboxStatus is the AI-sandbox "add to sandbox" lifecycle status of the
	// workload: "pending" or "in_sandbox" (empty = not-in-sandbox). It is DERIVED
	// in the inventory query, not stored: "in_sandbox" from the live pod-template
	// label (kubescape.io/sandbox on the synced workload), else "pending" from the
	// control-plane ai_sandbox_statuses request timestamp (LEFT JOIN on
	// (customer_guid, resource_hash)). Independent of the observed ai_sandboxes
	// telemetry tables.
	//
	// There is intentionally NO "failed" state (product decision, SUB-7442): the
	// lifecycle is not-in-sandbox -> pending -> in_sandbox. Deterministic enable
	// errors are handled in the write-path (unsupported kinds / not-found are
	// logged, not surfaced) rather than exposed as a status.
	//
	// Host coverage tracks resource_hash availability on the inventory row - k8s
	// and ECS carry one; cloud-host rows do not yet project one, so they read
	// empty until the platform surfaces a host resource_hash.
	SandboxStatus string `json:"sandboxStatus,omitempty"`
}
