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

// NetworkStreamProcessAttributionV1 is the current revision of the
// process-attribution extension to NetworkStream: NetworkStream.Processes plus
// the per-connection NetworkStreamEvent.Process reference. See
// NetworkStream.ProcessAttributionVersion for why an explicit marker exists at
// all, and bump this constant if the meaning of either field changes.
const NetworkStreamProcessAttributionV1 = 1

// processRefSep separates the two integer components of a ProcessRef's textual
// form.
//
// It is deliberately "/" and not the U+241F (␟) that CommPID in linuxobjects.go
// uses. CommPID needs an exotic separator because its first component (Comm) is
// a free-form string that may itself contain "/" or ":". Both components of a
// ProcessRef are decimal integers, so no escaping hazard exists, and "/" keeps
// this file internally consistent: NetworkStreamEntity.Inbound/Outbound are
// already keyed by the composite "<ip>/<port>/<protocol>". It is also 1 byte
// rather than 3 (UTF-8), which matters because this text appears once per
// connection — see NetworkStream.Processes on the wire budget.
const processRefSep = "/"

// ProcessRef identifies one OS process on the node that produced a
// NetworkStream message. It is both the key of NetworkStream.Processes and the
// value of NetworkStreamEvent.Process, so a consumer joins a connection to its
// process tree with `stream.Processes[*event.Process]` — no string formatting
// at the call site and no way for the two sides to disagree on key format.
//
// StartTime is CLOCK_BOOTTIME in clock ticks since boot (USER_HZ = 100, i.e.
// 10 ms per tick), NOT a wall-clock time.Time. This is deliberate:
//
//   - The only start-time source with full process coverage is
//     /proc/<pid>/stat field 22 (starttime), which is denominated in clock
//     ticks since boot.
//   - Rendering those ticks to wall-clock requires the boot timestamp from
//     /proc/stat `btime`, which has whole-second resolution only. A time.Time
//     on the wire would therefore be lossy by up to a second — for a value
//     whose entire job is to disambiguate identity, that is the wrong trade.
//   - Ticks are exact, need no boot-time lookup, and 10 ms granularity is
//     safe for pid-reuse disambiguation by 2-5 orders of magnitude (a kernel
//     needs to cycle the whole pid space to reuse a pid).
//
// A zero StartTime means the sensor could not read the process's start time
// (typically a process that exited before /proc could be sampled). Such a ref
// carries pid-only identity: it is still a valid key within the message, but it
// must not be treated as equal to a ref for the same pid with a known start
// time, and it is not safe to compare across messages.
//
// PID is a host-namespace pid, matching what the sensor's eBPF network gadget
// reports (the socket enricher's pid_tgid >> 32, i.e. task->tgid; there is no
// pid-namespace translation on that path). Host pids are node-unique, which is
// what makes a single message-scoped Processes map correct.
type ProcessRef struct {
	PID       uint32 `json:"pid,omitempty" bson:"pid,omitempty"`
	StartTime uint64 `json:"startTime,omitempty" bson:"startTime,omitempty"`
}

// MarshalText implements encoding.TextMarshaler. Both encoding/json and
// mongo-driver's bson use it for map keys and for plain values, so a ProcessRef
// is always the single string "<pid>/<startTime>" on the wire.
func (p ProcessRef) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(p.PID), 10) + processRefSep + strconv.FormatUint(p.StartTime, 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. The receiver is left
// untouched on error so a caller that ignores the error cannot end up
// attributing a connection to a half-parsed process.
func (p *ProcessRef) UnmarshalText(text []byte) error {
	pidStr, startTimeStr, found := strings.Cut(string(text), processRefSep)
	if !found {
		return fmt.Errorf("invalid ProcessRef representation: %q", text)
	}
	pid, err := strconv.ParseUint(pidStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef pid %q: %w", pidStr, err)
	}
	// startTime is ticks since boot; it outgrows uint32 after ~497 days of
	// uptime at USER_HZ = 100, so it is parsed and stored as 64-bit.
	startTime, err := strconv.ParseUint(startTimeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef startTime %q: %w", startTimeStr, err)
	}
	p.PID = uint32(pid) // range-checked by ParseUint's bitSize 32 above
	p.StartTime = startTime
	return nil
}

// NetworkStream represents a collection of network traffic events for a specific pod/container
type NetworkStream struct {
	// <identifier> to <network stream entity>
	Entities map[string]NetworkStreamEntity `json:"entities,omitempty"`

	// ProcessAttributionVersion states which revision of process attribution
	// the producing sensor implements; see NetworkStreamProcessAttributionV1.
	//
	// Absent or zero means "this sensor predates process attribution", NOT
	// "attribution ran and produced nothing". The backend genuinely needs to
	// tell those apart: an upgraded sensor can legitimately emit an empty
	// Processes map (every connection unattributable in that interval), so
	// `len(Processes) == 0` is not a usable proxy for sensor capability.
	//
	// This marker is new to the repo — there is no existing schema-version-field
	// convention here. Versioning is otherwise done by the CRD kind string
	// (kubescape.io/v1/networkstreams) plus Pulsar topic name, and the
	// established pattern for additive changes is `omitempty` plus a documented
	// consumer tolerance (see HostCVEOnFinishedMessage in hot_cve.go). That
	// pattern alone is not sufficient here precisely because absence of data
	// and absence of the feature are different states. Bumping the kind string
	// to a v2 was rejected: both consumers decode with plain json.Unmarshal and
	// neither sets DisallowUnknownFields, so a purely additive change needs no
	// breaking version, and a kind bump would force a coordinated consumer
	// rollout. Note also that this repo auto-tags v0.0.<run_number> on push to
	// main, so the module version carries no semver signal a consumer could
	// branch on.
	ProcessAttributionVersion int `json:"processAttributionVersion,omitempty" bson:"processAttributionVersion,omitempty"`

	// Processes maps a process identity to the process tree that produced the
	// connections referencing it. Connections point here via
	// NetworkStreamEvent.Process.
	//
	// The map is scoped once per MESSAGE, not per container, and the tree is
	// shipped once per PROCESS, not once per connection. Both choices are about
	// size:
	//
	//   - The synchronizer envelope carries this payload base64-encoded in its
	//     `patch` field, so every byte added on this type costs x1.333 on the
	//     wire.
	//   - Process trees dominate that cost, not the number of map entries:
	//     they are chain-shaped and measure ~2.0 KB serialised at the median
	//     and ~5.1 KB at p90, of which `cmdline` alone is ~40% and is
	//     unbounded. Inlining a tree per connection (the pre-existing
	//     NetworkStreamEvent.ProcessTree shape) is what forced the sensor to
	//     nil trees out entirely before sending; deduplicating per process is
	//     what makes attribution affordable at all. A worst case observed
	//     message (109 container entities, 283 connections, 142 KB on the
	//     wire) lands near 970 KB with full attribution, against a 5 MiB
	//     broker limit.
	//   - Message scope is correct rather than merely cheaper: ProcessRef.PID
	//     is a host-namespace pid and a stream message covers exactly one node,
	//     so the keys are already unique across the message. Per-container
	//     scoping would duplicate every tree shared between containers and buy
	//     nothing.
	Processes map[ProcessRef]*ProcessTree `json:"processes,omitempty" bson:"processes,omitempty"`
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
	// Process attributes this connection to a process in the message-scoped
	// NetworkStream.Processes map: `stream.Processes[*event.Process]`.
	//
	// This is the field to use. The ProcessTree field above inlines a whole
	// tree per connection and is nilled out by the sensor before the message is
	// sent; it survives only because an in-process notification channel feeds a
	// separate host-network sensor from the pre-nil struct.
	//
	// It is a pointer so `omitempty` actually drops it — encoding/json does not
	// omit zero structs, and an always-present "process":"0/0" on every
	// unattributed connection is exactly the per-connection waste this design
	// exists to avoid. A single text-marshalled field also beats two flat
	// scalars (pid + startTime) on the wire by ~6 bytes per connection, because
	// it spends one JSON key instead of two.
	Process *ProcessRef `json:"process,omitempty" bson:"process,omitempty"`
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
