package armotypes

import (
	"encoding/json"
	"time"
)

// AiSandboxAIEndpoint observation source: how the endpoint was detected.
const (
	AiSandboxEndpointSourceDNS = "dns" // DNS resolution — reachability, NOT confirmed usage
	AiSandboxEndpointSourceL7  = "l7"  // R-L7 structured payload — confirmed usage
)

// AiSandboxAIEndpoint plane/kind: which face of the provider the endpoint is.
const (
	AiSandboxEndpointKindInference    = "inference"     // model data-plane (bedrock-runtime, api.openai.com)
	AiSandboxEndpointKindControlPlane = "control-plane" // management API (bedrock.<region>)
	AiSandboxEndpointKindModelHost    = "model-host"    // model hosting (sagemaker)
)

// AiSandboxAIEndpoint is one AI-provider endpoint a sandboxed workload is
// observed contacting. v1 entries are DNS-derived (Source="dns" → reachability,
// not confirmed usage); the R-L7 payload later adds confirmed-usage entries
// (Source="l7"). It is the lightweight precursor of the
// ai_sandbox_external_services detail model — same host/provider keying, ready
// to merge when typed L7 lands.
type AiSandboxAIEndpoint struct {
	Host      string     `json:"host"`                // e.g. bedrock-runtime.us-east-1.amazonaws.com
	Provider  string     `json:"provider,omitempty"`  // aligns to workload_statuses ai_client_providers vocab, e.g. "AWS Bedrock"
	Kind      string     `json:"kind,omitempty"`      // one of AiSandboxEndpointKind*
	Region    string     `json:"region,omitempty"`    // e.g. us-east-1 (when derivable from the host)
	Source    string     `json:"source"`              // one of AiSandboxEndpointSource*
	FirstSeen *time.Time `json:"firstSeen,omitempty"` // first observation in the rollup window
	LastSeen  *time.Time `json:"lastSeen,omitempty"`  // most recent observation
}

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

	// AIEndpoints is the set of AI-provider endpoints the workload is observed
	// contacting (e.g. bedrock-runtime.us-east-1, api.openai.com), each with its
	// provider/plane/region. v1 entries are DNS-derived (Source="dns" →
	// REACHABILITY, not confirmed usage); the R-L7 payload later adds confirmed
	// rows (Source="l7") with model/token detail. This is the lightweight
	// precursor of the ai_sandbox_external_services detail model. Stored on the
	// subject row so it flows into AiSandboxView via the embed.
	AIEndpoints []AiSandboxAIEndpoint `json:"aiEndpoints,omitempty"`

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

// AiSandboxBoundRole is one Role/ClusterRole bound to a K8s identity — a unit of
// its permissions. Kind is "ClusterRole" or "Role". Binding is the
// (Cluster)RoleBinding that grants it.
type AiSandboxBoundRole struct {
	Kind    string `json:"kind"` // "ClusterRole" | "Role"
	Name    string `json:"name"`
	Binding string `json:"binding,omitempty"` // the (Cluster)RoleBinding name
}

// AiSandboxIdentity is an identity a sandbox subject runs as, plus its
// permissions. For K8s: section="k8s", the ServiceAccount, and the Roles bound to
// it. Sourced read-time from the subject's pod spec (serviceAccountName, from the
// synced object) joined to the PG role-binding graph — NOT from compliance
// controls. IsAdmin is true when any bound role is cluster-admin/admin. Usage
// fields (in_use/last_used) are a later, agent-blocked overlay — absent in v1.
type AiSandboxIdentity struct {
	Section string               `json:"section"`         // "k8s" (cloud/api later)
	Name    string               `json:"name"`            // ServiceAccount name
	Type    string               `json:"type"`            // "ServiceAccount"
	Scope   string               `json:"scope,omitempty"` // namespace
	IsAdmin bool                 `json:"isAdmin"`         // bound to cluster-admin/admin
	Roles   []AiSandboxBoundRole `json:"roles,omitempty"` // bound Roles/ClusterRoles
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

// AiSandboxEvent.EventLayer — UI grouping of an activity event by protocol/kind.
// Derived at serving time from the agent's event_type (exec→proc, network→net,
// open→file, …); the UI groups/filters the Activity feed by these.
const (
	AiSandboxEventLayerNetwork = "net"
	AiSandboxEventLayerDNS     = "dns"
	AiSandboxEventLayerHTTP    = "http"
	AiSandboxEventLayerProcess = "proc"
	AiSandboxEventLayerFile    = "file"
	AiSandboxEventLayerSyscall = "syscall"
)

// AiSandboxEvent.Direction — for the network/L7 layers (egress vs ingress, audit R5).
const (
	AiSandboxEventDirectionIngress = "ingress"
	AiSandboxEventDirectionEgress  = "egress"
)

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
// Type-specific payload (R-PROC-1 exe_path/argv/cwd/pid/ppid/parent_exe_path now;
// host/method/model for L7 later) lives in the generic Detail jsonb — so the
// hot-table column set stays STABLE as network/dns/http/file/syscall layers come
// online, no column-per-type sprawl. Verdict/Rule/PolicyRef are enforcement stubs
// reserved per the policies-descoped decision — always empty for now.
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

	// EventLayer is the UI grouping of the event by protocol/kind — one of
	// AiSandboxEventLayer* (net/dns/http/proc/file/syscall). Derived at serving
	// time from EventType (exec→proc, network→net, open→file, …).
	EventLayer string `json:"eventLayer,omitempty"`
	// Direction applies to the network/L7 layers — AiSandboxEventDirection*
	// (ingress/egress). Empty for proc/file/syscall.
	Direction string `json:"direction,omitempty"`
	// EventType is the specific event, using the shared
	// armotypes.EventType vocabulary (EventTypeExec=="exec", EventTypeDNS=="dns",
	// EventTypeHTTP=="http", …) — the same string values the node-agent emits in
	// the `event_type` envelope key. Process-exec events carry EventTypeExec.
	EventType EventType `json:"eventType"`
	// Summary is a short human-readable description for the activity feed.
	Summary string `json:"summary,omitempty"`

	// LogRecordUID carries the OTel log.record.uid — the agent's globally-unique
	// per-event id — and is the canonical dedup key for the ai_sandbox_events
	// serving table (the lake already dedups on it; PG should too, instead of a
	// fragile composite key).
	LogRecordUID string `json:"logRecordUID,omitempty"`

	// SessionID and Turn are stamped during SILVER sessionization (R-L7).
	// Turn is a pointer so that nil means "not yet sessionized" and a stamped
	// first-turn value of 0 is preserved (plain int + omitempty would silently
	// drop a real zero, making turn 0 indistinguishable from "unsessionized").
	SessionID string `json:"sessionID,omitempty"`
	Turn      *int   `json:"turn,omitempty"`

	// Detail carries the type-specific payload for event layers that do not yet
	// have first-class fields (the contract's `detail` jsonb column).
	Detail json.RawMessage `json:"detail,omitempty"`

	// Enforcement stubs — descoped (policies). Always empty for now.
	Verdict   string `json:"verdict,omitempty"`
	Rule      string `json:"rule,omitempty"`
	PolicyRef string `json:"policyRef,omitempty"`
}
