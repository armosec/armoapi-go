package armotypes

import "time"

// AiSandboxInfo is the shared serving DTO for one sandboxed subject — a
// Kubernetes workload, an ECS service, or a host — in the AI-Sandbox
// (AI-SPM) feature. It crosses repos: postgres-connector maps it to/from
// the ai_sandboxes serving row and cadashboardbe serves it to the UI.
//
// Only STORED serving fields appear here. Counters (instance_count,
// models_count, connectors_count) are computed at serving time, policy
// state and row-lifecycle timestamps (created/updated/deleted) are not
// part of the serving contract, and live-join fields (risk_factors,
// is_internet_facing, entity_type, ...) are resolved at serving time and
// are intentionally absent.
type AiSandboxInfo struct {
	CustomerGUID string `json:"customerGUID"`
	ResourceHash string `json:"resourceHash"`

	// HostType is the kind of host running the subject:
	// kubernetes|ecs-ec2|ecs-fargate|ec2.
	HostType string `json:"hostType"`
	WLID     string `json:"wlid"`
	Name     string `json:"name"`

	Account       string `json:"account"`
	Region        string `json:"region"`
	CloudProvider string `json:"cloudProvider"`

	Cluster      string `json:"cluster"`
	Namespace    string `json:"namespace"`
	Kind         string `json:"kind"`
	WorkloadName string `json:"workloadName"`

	HostID string `json:"hostID"`

	ExposureFlags []string `json:"exposureFlags"`

	// EnablementState is observing|active.
	EnablementState string    `json:"enablementState"`
	FirstSeen       time.Time `json:"firstSeen"`
	LastSeen        time.Time `json:"lastSeen"`
}

// AiSandboxInstanceInfo is the shared serving DTO for one running unit of
// a sandboxed subject — a Kubernetes pod or an ECS task.
type AiSandboxInstanceInfo struct {
	CustomerGUID string `json:"customerGUID"`
	InstanceRef  string `json:"instanceRef"`

	WLID         string `json:"wlid"`
	ResourceHash string `json:"resourceHash"`

	Status  string `json:"status"`
	Node    string `json:"node"`
	PodName string `json:"podName"`

	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"lastSeen"`
}
