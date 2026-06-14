package armotypes

import "time"

// AiSandboxInfo is the WRITE model for one sandboxed subject — a Kubernetes
// workload, an ECS service, or a host — in the AI-Sandbox (AI-SPM) feature.
// Every field here is a STORED column: exactly what the aggregator upserts and
// postgres-connector persists in the ai_sandboxes serving table. Row-lifecycle
// timestamps (created/updated/deleted) live on the gorm model, not here.
//
// The UI does not consume AiSandboxInfo directly — it consumes AiSandboxView,
// which embeds this record and adds the fields produced at serving time
// (counters, risk_factors, is_internet_facing, entity_type, ai_*_providers).
// Splitting the two keeps the write contract (what we store) distinct from the
// read contract (what we serve).
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

// AiSandboxView is the READ model for the AI-Sandbox UI: one row in the sandbox
// list (and the source for the header). It composes the persisted AiSandboxInfo
// record with the fields that are NOT stored but produced at serving time, so
// the UI binds to one complete shape regardless of each field's provenance.
//
// Provenance of the added fields:
//   - Counters are computed at serving time by COUNTing the per-subject detail
//     tables — there are no stored counter columns, which avoids drift:
//     instanceCount   ← COUNT(ai_sandbox_instances)              [available now]
//     modelsCount     ← COUNT(distinct models)                   [R-L7; 0 until]
//     connectorsCount ← COUNT(mcp_servers)+COUNT(external_svcs)  [R-L7; 0 until]
//   - The remaining fields are joined live from workload_statuses on
//     (customer_guid, resource_hash); they are display-only header/badge inputs
//     the platform already owns.
//
// The aggregator never populates these; the serving query (cadashboardbe) does.
// On the write path they are simply absent.
type AiSandboxView struct {
	AiSandboxInfo `json:",inline"`

	// Computed at serving time (COUNT over the per-subject detail tables).
	InstanceCount   int `json:"instanceCount"`
	ModelsCount     int `json:"modelsCount"`
	ConnectorsCount int `json:"connectorsCount"`

	// Joined live from workload_statuses on (customer_guid, resource_hash).
	AIClientProviders []string `json:"aiClientProviders"`
	AIServerProviders []string `json:"aiServerProviders"`
	RiskFactors       []string `json:"riskFactors"`
	IsInternetFacing  bool     `json:"isInternetFacing"`
	EntityType        string   `json:"entityType"`
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
