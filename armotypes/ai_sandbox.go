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

	// HostType is the kind of host running the subject. Values are the
	// platform's canonical HostType constants (this package):
	// HostTypeKubernetes|HostTypeEcsEc2|HostTypeEcsFargate|HostTypeEc2 — never a
	// hand-written literal. Kept as string here to match the stored column type;
	// producers must assign string(armotypes.HostType*) values.
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

	// AIReachableEndpoints is the set of AI-provider hosts the workload's DNS
	// activity shows it can reach (e.g. bedrock-runtime.*.amazonaws.com,
	// *.openai.com, *.anthropic.com, generativelanguage.googleapis.com). This is
	// DNS-derived REACHABILITY — the workload resolved/contacted these hosts — and
	// is NOT confirmed model usage (no proof a request was made or that a call
	// succeeded). Confirmed usage comes from R-L7 model events, served elsewhere.
	AIReachableEndpoints []string `json:"aiReachableEndpoints,omitempty"`

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

// AiSandboxEvent is the READ DTO for one row on the per-sandbox Activity &
// Events surface (Track D). It mirrors the `ai_sandbox_events` serving table
// (serving-contract.md): the envelope keys that scope/join the event, the
// event-time + event-layer/event-type discriminators, and the typed fields the
// aggregator flattens out of the raw record for the first event types.
//
// This is a clean read model — the aggregator writes the underlying row, the
// serving query (cadashboardbe) returns this shape, and the UI binds to it. The
// table is a hot table (TTL ~7d), so this is recent-activity, not history.
//
// First event types served are R-PROC process-exec events; the process-typed
// fields below are populated for the `proc` layer. Other layers
// (net/dns/http/file/syscall) reuse the same envelope and carry their own
// payload in Detail until they are typed out.
//
// NOTE: enforcement is descoped — Verdict/Rule/PolicyRef are nullable stubs,
// empty until a policy engine exists. They are kept on the contract so the
// shape is stable when enforcement lands.
type AiSandboxEvent struct {
	// Envelope keys — scope the event and join back to the sandbox / instance.
	CustomerGUID string `json:"customerGUID"`
	ResourceHash string `json:"resourceHash"`
	WLID         string `json:"wlid"`
	InstanceRef  string `json:"instanceRef,omitempty"`

	EventTime time.Time `json:"eventTime"`

	// EventLayer is the family of the event:
	// net|dns|http|proc|file|syscall (incl. ingress direction). EventType is the
	// finer discriminator within the layer (e.g. process exec). Summary is a
	// human-readable one-line description for the activity feed.
	EventLayer string `json:"eventLayer"`
	EventType  string `json:"eventType"`
	Summary    string `json:"summary,omitempty"`

	// Detail is the per-type payload that is not (yet) flattened into typed
	// columns — the jsonb `detail` column passed through as-is.
	Detail map[string]interface{} `json:"detail,omitempty"`

	// Process-typed fields (proc layer / R-PROC process-exec). Populated for
	// process events; empty for other layers.
	ExePath   string   `json:"exePath,omitempty"`
	Argv      []string `json:"argv,omitempty"`
	Cwd       string   `json:"cwd,omitempty"`
	PID       int      `json:"pid,omitempty"`
	PPID      int      `json:"ppid,omitempty"`
	ParentExe string   `json:"parentExe,omitempty"`

	// Sessionization (R-L7, stamped in SILVER). Empty until R-L7 lands.
	SessionID string `json:"sessionID,omitempty"`
	Turn      int    `json:"turn,omitempty"`

	// Enforcement stubs (descoped). Nullable/empty until a policy engine exists.
	Verdict   string `json:"verdict,omitempty"`
	Rule      string `json:"rule,omitempty"`
	PolicyRef string `json:"policyRef,omitempty"`
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
