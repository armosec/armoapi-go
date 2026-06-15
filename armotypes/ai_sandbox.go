package armotypes

import (
	"encoding/json"
	"time"
)

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
	// resolution shows it can contact — e.g. bedrock-runtime, *.openai.com,
	// *.anthropic.com, generativelanguage.googleapis.com. This is DNS-derived
	// REACHABILITY (the workload resolved/can reach these hosts), NOT confirmed
	// usage: presence here does not imply the workload actually issued an AI call.
	// Stored on the subject row so it flows into AiSandboxView via the embed.
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

// AiSandboxEvent is the shared READ DTO for one row of the per-sandbox Activity
// & Events surface (UI Track D). It mirrors the lake-backed ai_sandbox_events
// serving table (a hot, short-TTL table). The json tags are the armoapi DTO
// (camelCase) convention shared with the sibling AiSandbox* types; the stored
// BRONZE/serving columns are the snake_case envelope keys (event_time,
// instance_ref, exe_path, …) the aggregator/SILVER writes and cadashboardbe
// reads, mapped to these fields. This is a clean read model: producers project
// rows into it, the UI binds to it; there is no separate write split because
// events are append-only and never upserted.
//
// Envelope keys identify the sandboxed subject and the running unit that emitted
// the event, joining back to ai_sandboxes / ai_sandbox_instances on
// (customer_guid, resource_hash) and instance_ref.
//
// The typed process fields (ExePath, Argv, Cwd, Pid, Ppid, ParentExePath) are
// populated NOW for R-PROC-1 (process-exec, EventType == EventTypeExec) events;
// other event families (network/dns/http/file/syscall) carry their type-specific
// payload in Detail until they get first-class fields. Verdict/Rule/PolicyRef are
// enforcement stubs reserved per the policies-descoped decision — always empty
// for now.
//
// NOTE: the contract keeps a generic `detail` jsonb column for not-yet-typed
// event families; it is surfaced here as a raw json.RawMessage so the read model
// stays lossless without committing to a schema for those families yet.
type AiSandboxEvent struct {
	CustomerGUID string `json:"customerGUID"`
	ResourceHash string `json:"resourceHash"`
	WLID         string `json:"wlid"`
	InstanceRef  string `json:"instanceRef"`

	EventTime time.Time `json:"eventTime"`

	// EventFamily is the node-agent's `event.family` envelope value — the R-*
	// family the agent emits (R-PROC-1, R-NET-1, R-NET-2/DNS, R-L7-1/HTTP,
	// R-FILE-1, R-SYSCALL-1). Kept as string to mirror the stored column; it is
	// NOT a new net|proc|file vocabulary.
	EventFamily string `json:"eventFamily"`
	// EventType is the specific event within the family, using the shared
	// armotypes.EventType vocabulary (EventTypeExec=="exec", EventTypeDNS=="dns",
	// EventTypeHTTP=="http", …) — the same string values the node-agent emits in
	// the `event_type` envelope key. Process-exec events carry EventTypeExec.
	EventType EventType `json:"eventType"`
	// Summary is a short human-readable description for the activity feed.
	Summary string `json:"summary,omitempty"`

	// R-PROC-1 process-exec typed fields (populated NOW for EventTypeExec events).
	// Field names mirror the agent's structured attrs / BRONZE columns:
	// exe_path, argv, cwd, pid, ppid, parent_exe_path.
	ExePath       string   `json:"exePath,omitempty"`
	Argv          []string `json:"argv,omitempty"`
	Cwd           string   `json:"cwd,omitempty"`
	Pid           int      `json:"pid,omitempty"`
	Ppid          int      `json:"ppid,omitempty"`
	ParentExePath string   `json:"parentExePath,omitempty"`

	// SessionID and Turn are stamped during SILVER sessionization (R-L7).
	SessionID string `json:"sessionID,omitempty"`
	Turn      int    `json:"turn,omitempty"`

	// Detail carries the type-specific payload for event layers that do not yet
	// have first-class fields (the contract's `detail` jsonb column).
	Detail json.RawMessage `json:"detail,omitempty"`

	// Enforcement stubs — descoped (policies). Always empty for now.
	Verdict   string `json:"verdict,omitempty"`
	Rule      string `json:"rule,omitempty"`
	PolicyRef string `json:"policyRef,omitempty"`
}
