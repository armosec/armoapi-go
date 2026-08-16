package armotypes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NetworkStreamEventProtocol string

const (
	NetworkStreamEventProtocolTCP NetworkStreamEventProtocol = "TCP"
	NetworkStreamEventProtocolUDP NetworkStreamEventProtocol = "UDP"
	NetworkStreamEventProtocolDNS NetworkStreamEventProtocol = "DNS"
)

type NetworkStreamEntityKind string

const (
	NetworkStreamEntityKindContainer NetworkStreamEntityKind = "container" // container
	NetworkStreamEntityKindHost      NetworkStreamEntityKind = "host"      // host
	// more types can be added here, but not one per platform: ECS entities are
	// container-kind and carry NetworkStreamEntityECS instead. See
	// docs/features/network-stream-workload-identity.md before adding a kind.
)

type EndpointKind string

const (
	EndpointKindPod     EndpointKind = "pod"
	EndpointKindService EndpointKind = "svc"
	EndpointKindRaw     EndpointKind = "raw"
)

// Process attribution on the network stream. Design, wire-format rationale and
// consumer guidance: docs/features/network-stream-process-attribution.md

// NetworkStreamProcessAttributionV1 is the current process-attribution revision.
// Bump it if the meaning of Processes or NetworkStreamEvent.ProcessRef changes.
const NetworkStreamProcessAttributionV1 = 1

const processRefSep = "/"

// ProcessRef identifies one process on the node, and is both the key of
// NetworkStream.Processes and the value of NetworkStreamEvent.ProcessRef.
// Resolve with NetworkStream.ProcessTreeFor.
//
// PID is a host-namespace pid. StartTimeNs is process CREATION time in
// nanoseconds since boot — not wall-clock, and not an exec timestamp (an
// exec-derived value changes on re-exec and breaks identity mid-life). Zero
// means unknown, leaving pid-only identity that is unsafe to compare across
// messages.
//
// The unit is nanoseconds; the resolution is not. procfs-derived values come
// from /proc/<pid>/stat field 22 in 10 ms clock ticks, so they are exact
// multiples of 10^7 ns. Fine for pid-reuse disambiguation, not a precise
// timestamp.
type ProcessRef struct {
	PID         uint32 `json:"pid,omitempty" bson:"pid,omitempty"`
	StartTimeNs uint64 `json:"startTimeNs,omitempty" bson:"startTimeNs,omitempty"`
}

// MarshalText renders a ProcessRef as "<pid>/<startTimeNs>". It never errors.
//
// encoding/json uses this for map keys and values alike; mongo-driver uses it
// for map keys only, so in BSON an event's ProcessRef is a subdocument instead.
func (p ProcessRef) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// String returns the same text as MarshalText, so a logged ref greps against
// the payload it came from.
func (p ProcessRef) String() string {
	return strconv.FormatUint(uint64(p.PID), 10) + processRefSep + strconv.FormatUint(p.StartTimeNs, 10)
}

// UnmarshalText parses "<pid>/<startTimeNs>", ignoring any further components
// and accepting leading zeros.
//
// Do not make this stricter. encoding/json aborts the whole enclosing decode
// when a TextUnmarshaler map key fails, and both stream consumers drop or nack
// such a message — so one bad key costs a node's entire interval. Appending
// components is the designed way to extend the format.
//
// The receiver is untouched on error, but only for direct callers: encoding/json
// leaves the allocated pointer installed, so ignoring its error yields a
// non-nil zero ref.
func (p *ProcessRef) UnmarshalText(text []byte) error {
	pidStr, rest, found := strings.Cut(string(text), processRefSep)
	if !found {
		return fmt.Errorf("invalid ProcessRef representation: %q", text)
	}
	startTimeStr, _, _ := strings.Cut(rest, processRefSep) // drop later-revision components

	pid, err := strconv.ParseUint(pidStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef pid %q: %w", pidStr, err)
	}
	startTimeNs, err := strconv.ParseUint(startTimeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef startTimeNs %q: %w", startTimeStr, err)
	}
	p.PID = uint32(pid) // range-checked by ParseUint above
	p.StartTimeNs = startTimeNs
	return nil
}

// NetworkStream represents a collection of network traffic events for a specific pod/container
type NetworkStream struct {
	// <identifier> to <network stream entity>
	Entities map[string]NetworkStreamEntity `json:"entities,omitempty"`

	// CloudMetadata describes the machine that sent this message. Message scope
	// rather than entity scope: one sensor process reports for one host, so every
	// entity in a message shares it.
	//
	// Host alert policies are scoped by cloud provider, region and instance, and a
	// consumer that cannot place the sending host has to drop the alert. The
	// Kubernetes path already carries this on its own alert envelope; the network
	// stream had no slot for it, so alerts derived from a stream could not be
	// scoped at all.
	//
	// A pointer because encoding/json does not honour omitempty on a struct: a
	// value field would add an empty object to every payload that has never
	// carried one. Nil means the sender does not report it, which is not the same
	// as reporting that it knows nothing.
	CloudMetadata *CloudMetadata `json:"cloudMetadata,omitempty" bson:"cloudMetadata,omitempty"`

	// ProcessAttributionVersion is the revision the producing sensor implements.
	// Absent/zero means the sensor predates process attribution — not that it ran
	// and found nothing, which an empty Processes map legitimately means.
	ProcessAttributionVersion int `json:"processAttributionVersion,omitempty" bson:"processAttributionVersion,omitempty"`

	// Processes holds each contributing process tree once per message, keyed by
	// identity; connections reference it via NetworkStreamEvent.ProcessRef.
	// Resolve with ProcessTreeFor.
	//
	// Once per process rather than per connection because trees dominate the
	// payload (~2 KB median, unbounded cmdline) and the synchronizer envelope is
	// base64, costing x1.333 per byte. Message scope is safe because PIDs are
	// host-namespace and a message covers one node.
	Processes map[ProcessRef]*ProcessTree `json:"processes,omitempty" bson:"processes,omitempty"`
}

// ProcessTreeFor resolves the tree that opened the given connection, returning
// nil when attribution is absent, dangling or nil-valued. Prefer it to indexing
// Processes, which panics on a nil ref and cannot tell a missing entry from a
// present-but-nil one (reachable from a `"processes":{"1/2":null}` payload).
//
// It deliberately does not fall back to the per-event ProcessTree field.
func (n *NetworkStream) ProcessTreeFor(event *NetworkStreamEvent) *ProcessTree {
	if n == nil || event == nil || event.ProcessRef == nil {
		return nil
	}
	return n.Processes[*event.ProcessRef]
}

// NetworkStreamEntity represents an aggregation of network connections from/to a specific source
type NetworkStreamEntity struct {
	// entity kind
	Kind NetworkStreamEntityKind `json:"kind,omitempty"`
	// entity details
	NetworkStreamEntityContainer `json:",inline"`
	// Platform identity for hosts and ECS tasks, populated *alongside* the
	// Kubernetes-shaped fields above, never instead of them. Typed fields win
	// where present; the fields above are the compatibility bridge that keeps
	// today's consumers working unchanged. Both are all-omitempty, so a
	// Kubernetes entity serialises exactly as it does today.
	NetworkStreamEntityECS  `json:",inline" bson:",inline"`
	NetworkStreamEntityHost `json:",inline" bson:",inline"`
	// inbound network events
	Inbound map[string]NetworkStreamEvent `json:"inbound,omitempty"`
	// outbound network events
	Outbound map[string]NetworkStreamEvent `json:"outbound,omitempty"`
}

// NetworkStreamEntityContainer represents a container generating network events.
//
// On host and ECS entities these same fields carry bridged identity (e.g.
// WorkloadKind "Host"/"ECSService") so that consumers keying off them need no
// change. Read the typed structs below in preference to these.
type NetworkStreamEntityContainer struct {
	// ContainerName is the name of the container generating these network events
	ContainerName string `json:"containerName,omitempty"`
	// ContainerID is the unique identifier for the container
	ContainerID string `json:"containerID,omitempty"`
	// namespace is the namespace where the pod is deployed
	PodNamespace string `json:"podNamespace,omitempty"`
	// PodName is the name of the pod involved in the network traffic
	PodName string `json:"podName,omitempty"`
	// WorkloadName is the name of the parent workload (e.g., Deployment, StatefulSet)
	WorkloadName string `json:"workloadName,omitempty"`
	// WorkloadKind is the type of the parent workload (e.g., Deployment, StatefulSet)
	WorkloadKind string `json:"workloadKind,omitempty"`
}

// Host and ECS workload identity on the network stream. Design, the bridge to the
// Kubernetes-shaped fields above and consumer guidance:
// docs/features/network-stream-workload-identity.md

// NetworkStreamEntityECS identifies the ECS task a container belongs to.
//
// Entity Kind stays NetworkStreamEntityKindContainer. An ECS workload *is* a
// container, and the traffic-view writer drops every entity whose kind is not
// "container" (postgres-connector `services/networktraffic/service.go`,
// createNetworkTrafficEvent), so minting an "ecs" kind would silently delete
// these rows from the traffic view.
//
// Field names and JSON keys mirror RuntimeAlertECSDetails so the stream and the
// alert describe one workload the same way, with two additions the alert model has
// no counterpart for: TaskRevision, needed because the bridged WorkloadName is
// "<service>-<taskFamily>-<revision>", and ECSTaskID.
//
// Deliberately absent: ECSContainerName and ECSContainerID, whose values already
// ride the bridge ContainerName field and the Entities map key respectively —
// typed duplicates would give one value two sources of truth.
type NetworkStreamEntityECS struct {
	ECSClusterName string `json:"ecsClusterName,omitempty" bson:"ecsClusterName,omitempty"`
	ClusterARN     string `json:"clusterArn,omitempty" bson:"clusterArn,omitempty"`
	// ServiceName is empty on a standalone task — one launched with no ECS
	// service — which is also what makes the bridged WorkloadKind "ECSTask".
	ServiceName string `json:"serviceName,omitempty" bson:"serviceName,omitempty"`
	TaskFamily  string `json:"taskFamily,omitempty" bson:"taskFamily,omitempty"`
	// TaskRevision is a string because that is what the ECS task metadata
	// endpoint returns and what the sensor already holds; making it numeric here
	// would only move a parse onto the producer.
	TaskRevision string `json:"taskRevision,omitempty" bson:"taskRevision,omitempty"`
	TaskARN      string `json:"taskArn,omitempty" bson:"taskArn,omitempty"`
	// ECSTaskID is the stable per-task identifier — the short form of the same
	// identity TaskARN carries region-qualified. Prefer it as a grouping key; keep
	// TaskARN for anything that has to address the task in an AWS API.
	//
	// Named with the ECS prefix, unlike its TaskFamily/TaskRevision/TaskARN
	// neighbours, because these structs are inlined: every key here shares one flat
	// namespace with the container and host structs, so a bare "taskID" is the kind
	// of generic key a future non-ECS concept would collide with. See
	// TestNetworkStreamEntityIdentityKeysAreUnique.
	ECSTaskID string `json:"ecsTaskID,omitempty" bson:"ecsTaskID,omitempty"`
	// LaunchType is an unvalidated passthrough of what ECS reports; "EC2" and
	// "FARGATE" are the values observed, and ECS Anywhere would add "EXTERNAL".
	// Compare case-sensitively against a known value rather than assuming the set
	// is closed — RuntimeAlert.GetAlertSourcePlatform has no third branch: on an
	// ECS-shaped alert, anything not exactly "FARGATE" (including "EXTERNAL", or a
	// lowercased "fargate") returns AlertSourcePlatformECSAgent, its ECS-on-EC2 branch.
	//
	// It is also what keeps Fargate a value rather than a second schema: a Fargate
	// sensor that later learns to stream needs no change here.
	LaunchType string `json:"launchType,omitempty" bson:"launchType,omitempty"`
}

// NetworkStreamEntityHost identifies the host itself, on entities of kind
// NetworkStreamEntityKindHost — which until now was a kind with no fields behind
// it. The only host identifier on the wire was the Entities map key, which both
// consumers discard, so hosts shared one identity downstream.
type NetworkStreamEntityHost struct {
	// HostID is an opaque, producer-chosen stable host identifier — treat it as
	// a grouping key, not as any particular system value. The shipped host agent
	// resolves it to the cloud instance ID, falling back to /etc/machine-id and
	// then the hostname, so consumers must not join it against any one of those.
	// It can also be a non-unique sentinel when every source fails.
	//
	// Prefer it to HostName, which is display-oriented and can change under a host.
	HostID string `json:"hostID,omitempty" bson:"hostID,omitempty"`
	// HostName is the host's reported hostname.
	HostName string `json:"hostName,omitempty" bson:"hostName,omitempty"`
}

// NetworkStreamEvent represents an aggregation of network connections from/to a specific source
type NetworkStreamEvent struct {
	Timestamp   time.Time                  `json:"timestamp,omitempty"`
	IPAddress   string                     `json:"ipAddress,omitempty"`
	DNSName     string                     `json:"dnsName,omitempty"`
	Port        int32                      `json:"port,omitempty"`
	Protocol    NetworkStreamEventProtocol `json:"protocol,omitempty"`
	ProcessTree *ProcessTree               `json:"processTree,omitempty"`
	// ProcessRef is the wire-side attribution: it points at an entry in the
	// message-scoped NetworkStream.Processes map. Resolve with ProcessTreeFor.
	// Use this rather than ProcessTree above, which the sensor nils out before
	// sending and which only survives for private-node-agent's in-process host
	// network sensor — do not remove it here.
	ProcessRef *ProcessRef `json:"processRef,omitempty" bson:"processRef,omitempty"`
	// endpoint kind (pod, service, raw)
	Kind EndpointKind `json:"kind,omitempty"`
	// endpoint details in case of pod
	NetworkStreamEventEndpointPodDetails `json:",inline"`
	// endpoint details in case of service
	NetworkStreamEventEndpointServiceDetails `json:",inline"`
}

func (e *NetworkStreamEvent) String() string {
	switch e.Kind {
	case EndpointKindPod:
		return "p/" + e.PodName + "/" + e.PodNamespace + "/" + string(e.IPAddress) + "/" + string(e.Port) + "/" + string(e.Protocol)
	case EndpointKindService:
		return "s/" + e.ServiceName + "/" + e.ServiceNamespace + "/" + string(e.IPAddress) + "/" + string(e.Port) + "/" + string(e.Protocol)
	case EndpointKindRaw:
		return "r/" + e.IPAddress + "/" + string(e.Port) + "/" + string(e.Protocol)
	default:
		return e.IPAddress + "/" + string(e.Port) + "/" + string(e.Protocol)
	}
}

type NetworkStreamEventEndpointPodDetails struct {
	// PodName is the name of the pod
	PodName string `json:"podName,omitempty"`
	// PodNamespace is the namespace of the pod
	PodNamespace string `json:"podNamespace,omitempty"`
	// WorkloadName is the name of the parent workload (e.g., Deployment, StatefulSet)
	WorkloadName string `json:"workloadName,omitempty"`
	// WorkloadNamespace is the namespace of the parent workload
	WorkloadNamespace string `json:"workloadNamespace,omitempty"`
	// WorkloadKind is the kind of the parent workload (e.g., Deployment, StatefulSet)
	WorkloadKind string `json:"workloadKind,omitempty"`
}

type NetworkStreamEventEndpointServiceDetails struct {
	ServiceName      string `json:"serviceName,omitempty"`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
}
