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
// the per-connection NetworkStreamEvent.ProcessRef reference. See
// NetworkStream.ProcessAttributionVersion for why an explicit marker exists at
// all, and bump this constant if the meaning of either field changes.
const NetworkStreamProcessAttributionV1 = 1

// processRefSep separates the components of a ProcessRef's textual form.
//
// It is deliberately "/" and not the U+241F (␟) that CommPID in linuxobjects.go
// uses. CommPID needs an exotic separator because its first component (Comm) is
// a free-form string that may itself contain "/" or ":". Every component of a
// ProcessRef is a decimal integer, so no escaping hazard exists, and "/" keeps
// this file internally consistent: NetworkStreamEntity.Inbound/Outbound are
// already keyed by the composite "<ip>/<port>/<protocol>". It is also 1 byte
// rather than 3 (UTF-8), which matters because this text appears once per
// connection — see NetworkStream.Processes on the wire budget.
const processRefSep = "/"

// ProcessRef identifies one OS process on the node that produced a
// NetworkStream message. It is both the key of NetworkStream.Processes and the
// value of NetworkStreamEvent.ProcessRef, so a consumer resolves a connection's
// process tree with NetworkStream.ProcessTreeFor — no string formatting at the
// call site and no way for the two sides to disagree on key format.
//
// StartTimeNs is the process's creation time as nanoseconds since boot
// (CLOCK_BOOTTIME), NOT a wall-clock time. The name carries the unit because
// armotypes.Process.StartTime — the same concept on the tree this ref points
// at — is a wall-clock time.Time, and confusing the two would be silent. The
// unit costs nothing on the wire: ProcessRef is text-marshalled, so these field
// names never appear in JSON at all.
//
// The start time is in the identity purely to survive pid reuse — that is the
// stated mitigation in the network-reputation GA design (shared-designs-and-docs
// projects/epics/network-reputation-ga/network-reputation-consolidation-design.md
// §13: "join reliability across process restarts / pid reuse (mitigated by
// startTime in the key)"). Nothing renders it. Anything that needs a
// human-facing process start time should read Process.StartTime off the tree
// this ref resolves to, which is a wall-clock time.Time.
//
// Boot-relative rather than a wall-clock time.Time, deliberately: converting to
// wall-clock requires the boot timestamp from /proc/stat `btime`, which has
// whole-second resolution only. A time.Time on the wire would therefore be lossy
// by up to a second — for a value whose entire job is to disambiguate identity,
// that is the wrong trade. A boot-relative integer needs no boot-time lookup and
// is exact.
//
// NANOSECONDS ARE THE UNIT, NOT THE RESOLUTION. Do not read this as a precise
// timestamp. The only start-time source with full process coverage is
// /proc/<pid>/stat field 22, which is denominated in clock ticks (USER_HZ = 100,
// i.e. 10 ms), so procfs-derived values are exact multiples of 10,000,000 ns and
// carry far less precision than the unit implies. Nanoseconds are chosen anyway
// because conversion only loses information in the other direction — ticks to ns
// is an exact x10^7, ns to ticks is a lossy division — so no source has to
// discard anything, and there is one unit for one concept either side of the
// serialisation boundary. 10 ms granularity is safe for pid-reuse
// disambiguation by 2-5 orders of magnitude (a kernel must cycle the whole pid
// space to reuse a pid).
//
// The value MUST be process creation time. Never derive it from an exec event:
// exec does not create a process, so an exec-derived value changes when a
// process re-execs and the identity silently breaks mid-life.
//
// A zero StartTimeNs means the sensor could not read the process's start time
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
	PID         uint32 `json:"pid,omitempty" bson:"pid,omitempty"`
	StartTimeNs uint64 `json:"startTimeNs,omitempty" bson:"startTimeNs,omitempty"`
}

// MarshalText implements encoding.TextMarshaler, rendering a ProcessRef as the
// single string "<pid>/<startTimeNs>". It never returns an error.
//
// encoding/json uses this for both map keys and plain values, so JSON carries
// the text form everywhere. mongo-driver honours encoding.TextMarshaler for map
// KEYS ONLY, so in BSON the Processes key is this same string but a
// NetworkStreamEvent.ProcessRef is a subdocument built from the struct tags
// above. Both round-trip; only the query syntax differs.
func (p ProcessRef) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// String returns the same "<pid>/<startTimeNs>" text as MarshalText, so a
// logged ref is greppable against the payload it came from.
func (p ProcessRef) String() string {
	return strconv.FormatUint(uint64(p.PID), 10) + processRefSep + strconv.FormatUint(p.StartTimeNs, 10)
}

// UnmarshalText implements encoding.TextUnmarshaler.
//
// Components beyond the first two are IGNORED rather than rejected. This is the
// forward-compatibility guarantee, and it is load-bearing: encoding/json aborts
// the entire enclosing decode when a TextUnmarshaler map key fails — unlike an
// int-keyed map, which skips the offending entry and carries on — and both
// stream consumers nack or drop a message whose decode returns an error. So
// without this tolerance, a future revision that appends a component to the key
// format would make every un-upgraded consumer reject whole messages, losing a
// node's entire interval of network visibility. Appending components is the
// designed way to extend this format; do not make the parser stricter.
//
// Leading zeros are likewise accepted, so non-canonical keys for the same
// process collapse to one map entry (last wins) instead of failing the message.
//
// The receiver is left untouched when an error is returned. Note that guarantee
// only reaches direct callers: encoding/json allocates the pointer target before
// calling this method and leaves it installed on failure, so a caller that
// ignores an Unmarshal error sees a non-nil, zero-valued ref.
func (p *ProcessRef) UnmarshalText(text []byte) error {
	pidStr, rest, found := strings.Cut(string(text), processRefSep)
	if !found {
		return fmt.Errorf("invalid ProcessRef representation: %q", text)
	}
	// Anything after the second component belongs to a later revision of this
	// format; drop it rather than failing.
	startTimeStr, _, _ := strings.Cut(rest, processRefSep)

	pid, err := strconv.ParseUint(pidStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef pid %q: %w", pidStr, err)
	}
	// Nanoseconds since boot need 64 bits: they exceed uint32 after ~4.3
	// seconds of uptime, and reach the uint64 ceiling after ~584 years.
	startTimeNs, err := strconv.ParseUint(startTimeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ProcessRef startTimeNs %q: %w", startTimeStr, err)
	}
	p.PID = uint32(pid) // range-checked by ParseUint's bitSize 32 above
	p.StartTimeNs = startTimeNs
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
	// NetworkStreamEvent.ProcessRef; resolve with ProcessTreeFor rather than
	// indexing the map directly.
	//
	// The map is scoped once per MESSAGE, not per container, and the tree is
	// shipped once per PROCESS, not once per connection. Both choices are about
	// size:
	//
	//   - The synchronizer envelope carries this payload base64-encoded in its
	//     `patch` field, so every byte added on this type costs x1.333 on the
	//     wire.
	//   - Process trees dominate that cost, not the number of map entries.
	//     They are chain-shaped, and measurement during design put them at
	//     ~2.0 KB serialised at the median and ~5.1 KB at p90, of which
	//     `cmdline` alone is ~40% and is unbounded. Inlining a tree per
	//     connection (the pre-existing NetworkStreamEvent.ProcessTree shape) is
	//     what forced the sensor to nil trees out entirely before sending;
	//     deduplicating per process is what makes attribution affordable at
	//     all.
	//   - Headroom against the 5 MiB broker limit is comfortable but NOT
	//     unbounded, and it is set by the tail rather than the median. Applying
	//     the p90 tree size to the largest message observed during design (109
	//     container entities, 283 connections, 142 KB on the wire) gives
	//     roughly 2 MB; the median tree size gives roughly 900 KB. Since
	//     `cmdline` is unbounded, a sensor-side cap on tree or cmdline size is
	//     the right lever if this ever approaches the limit — not dropping
	//     attribution.
	//   - Message scope is correct rather than merely cheaper: ProcessRef.PID
	//     is a host-namespace pid and a stream message covers exactly one node,
	//     so the keys are already unique across the message. Per-container
	//     scoping would duplicate every tree shared between containers and buy
	//     nothing.
	Processes map[ProcessRef]*ProcessTree `json:"processes,omitempty" bson:"processes,omitempty"`
}

// ProcessTreeFor resolves the process tree that opened the given connection,
// returning nil when attribution is absent, dangling or nil-valued.
//
// Use this instead of indexing Processes directly. A bare
// `n.Processes[*event.ProcessRef]` has three sharp edges: it panics when
// event.ProcessRef is nil (every event from a sensor that predates attribution,
// and any connection the sensor could not attribute), it cannot distinguish a
// missing entry from a present-but-nil one, and a `"processes":{"1/2":null}`
// payload decodes to exactly that present-but-nil case — so the map index
// returns ok == true with a nil tree and the next field access panics.
//
// This deliberately does not fall back to the deprecated per-event ProcessTree
// field; the two have different lifecycles (see NetworkStreamEvent.ProcessRef).
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
	// ProcessRef attributes this connection to an entry in the message-scoped
	// NetworkStream.Processes map. Resolve it with
	// NetworkStream.ProcessTreeFor(event), not by indexing the map.
	//
	// This is the field to use for anything travelling over the wire. The
	// ProcessTree field above inlines a whole tree per connection and is nilled
	// out by the sensor before the message is sent (node-agent
	// removeProcessTreeFromEvents); it survives because private-node-agent
	// wires an in-process notification channel that its host network sensor
	// reads pre-nil (pkg/hostnetworksensor/v1 reads outbound.ProcessTree
	// directly). Searching only the OSS node-agent makes it look dead — that
	// main passes nil for the channel.
	//
	// So: do not remove it as part of adding attribution. It does become
	// removable once HostNetworkSensor is retired, which the network-reputation
	// GA design (shared-designs-and-docs
	// projects/epics/network-reputation-ga/network-reputation-consolidation-design.md
	// §5) plans across private-node-agent's three mains and pkg/hostnetworksensor.
	// Drop it in that change, not this one.
	//
	// It is a pointer so `omitempty` actually drops it — encoding/json does not
	// omit zero structs, and an always-present "processRef":"0/0" on every
	// unattributed connection is exactly the per-connection waste this design
	// exists to avoid. A single text-marshalled field also beats two flat
	// scalars (pid + startTimeNs) on the wire by 6 bytes per connection,
	// because it spends one JSON key instead of two.
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
