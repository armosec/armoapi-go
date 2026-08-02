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
	// more types can be added here
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
	// inbound network events
	Inbound map[string]NetworkStreamEvent `json:"inbound,omitempty"`
	// outbound network events
	Outbound map[string]NetworkStreamEvent `json:"outbound,omitempty"`
}

// NetworkStreamEntityContainer represents a container generating network events
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
