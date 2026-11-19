// Package monitoredunits holds the wire types of the monitored units message:
// a sensor reports which units (a pod, later a VM or another kind) it monitors,
// and from when to when. The consumer builds monitoring coverage from it.
//
// The message is JSON. Wire contract and rules: docs/features/monitored-units.md
package monitoredunits

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

// SchemaVersion is the current value of Message.SchemaVersion. It is raised
// only by a change that breaks an existing field. A new optional field does
// not raise it.
const SchemaVersion = 1

// EntryType says whether an entry starts or stops the monitoring of a unit.
type EntryType string

const (
	EntryTypeStart EntryType = "start"
	EntryTypeStop  EntryType = "stop"
)

// Reason says why monitoring of a unit started or stopped.
//
// It is a plain string on purpose: a consumer keeps a value it does not know,
// and does not reject the entry.
type Reason string

const (
	// Start reasons, set by the sensor.

	// ReasonUnitStarted: the sensor saw the unit start.
	ReasonUnitStarted Reason = "unit-started"
	// ReasonAgentStarted: the unit was already running when the sensor started
	// or restarted.
	ReasonAgentStarted Reason = "agent-started"

	// Stop reasons, set by the sensor.

	// ReasonUnitEnded: the unit ended while it was monitored. The normal end.
	ReasonUnitEnded Reason = "unit-ended"
	// ReasonAgentStopping: the sensor shuts down while the unit still runs.
	ReasonAgentStopping Reason = "agent-stopping"

	// Stop reason, set by the consumer only. A sensor never sends it.

	// ReasonTimeout: the consumer closed the period, because the unit left the
	// periodic lists or the lists stopped arriving.
	ReasonTimeout Reason = "timeout"
)

// UnitKind names the kind of a monitored unit. It selects the shape of
// Unit.Attributes. A plain string, for the same reason as Reason.
type UnitKind string

const (
	// UnitKindK8sPod is a Kubernetes pod. Unit.ID is the pod UID, and
	// Unit.Attributes is a K8sPodAttributes.
	UnitKindK8sPod UnitKind = "k8s-pod"
)

// Message is one monitored units message. A sensor sends one per machine: when
// monitoring of a unit starts or stops, and as a periodic list.
type Message struct {
	SchemaVersion int `json:"schemaVersion"`
	// AgentVersion is the version of the sensor.
	AgentVersion string `json:"agentVersion"`
	// AgentStartedAt is when this sensor process started. It identifies the
	// sensor run: a change tells the consumer that the sensor restarted.
	AgentStartedAt time.Time `json:"agentStartedAt"`
	// ClusterName is the Kubernetes cluster of the sensor. A sensor in a cluster
	// must set it: it is the only cluster identity that the consumer uses, and a
	// k8s-pod entry of a message without it is dropped. Empty when the sensor
	// does not run in a cluster.
	ClusterName string `json:"clusterName,omitempty"`
	// Node is the node name in a cluster, or the host name elsewhere.
	Node string `json:"node"`
	// CloudMetadata describes the machine that the sensor runs on (provider,
	// account, region, instance). Nil when the sensor does not report it.
	CloudMetadata *armotypes.CloudMetadata `json:"cloudMetadata,omitempty"`
	// SentAt is when the message was sent.
	SentAt time.Time `json:"sentAt"`
	// IntervalSeconds is the interval of the periodic list, as configured on
	// this sensor.
	IntervalSeconds int `json:"intervalSeconds"`
	// Entries holds start and stop entries. It may be empty.
	Entries []Entry `json:"entries"`
	// MonitoredUnitIDs holds the IDs of all units that this sensor run monitors
	// now. It is set only in a periodic list. Nil means that the message is not
	// a periodic list. An empty list means that the sensor is alive and monitors
	// nothing. Use IsPeriodicList and SetMonitoredUnitIDs.
	MonitoredUnitIDs *[]string `json:"monitoredUnitIds,omitempty"`
}

// IsPeriodicList reports whether the message carries the list of monitored
// unit IDs.
func (m *Message) IsPeriodicList() bool {
	return m.MonitoredUnitIDs != nil
}

// SetMonitoredUnitIDs makes the message a periodic list. A nil or empty ids is
// stored as an empty list, which is sent as [].
func (m *Message) SetMonitoredUnitIDs(ids []string) {
	if ids == nil {
		ids = []string{}
	}
	m.MonitoredUnitIDs = &ids
}

// Entry is the start or the stop of the monitoring of one unit.
type Entry struct {
	Type   EntryType `json:"type"`
	Reason Reason    `json:"reason"`
	// Unit is the monitored unit. A stop entry carries only Kind and ID.
	Unit Unit `json:"unit"`
	// MonitoredFrom is when the sensor started to monitor the unit. Set on a
	// start entry.
	MonitoredFrom *time.Time `json:"monitoredFrom,omitempty"`
	// MonitoredUntil is when the sensor stopped monitoring the unit. Set on a
	// stop entry.
	MonitoredUntil *time.Time `json:"monitoredUntil,omitempty"`
	// CI is the CI identity that the sensor reports about itself. Nil when the
	// sensor does not know it. See armotypes.CIIdentity.
	CI *armotypes.CIIdentity `json:"ci,omitempty"`
}

// Unit identifies a monitored unit.
type Unit struct {
	Kind UnitKind `json:"kind"`
	// ID is unique for the kind. For UnitKindK8sPod it is the pod UID.
	ID string `json:"id"`
	// StartedAt is when the unit started. Set on a start entry.
	StartedAt *time.Time `json:"startedAt,omitempty"`
	// Attributes describes the unit. Set on a start entry. Its shape depends on
	// Kind, so it is kept as raw JSON: a consumer can store and forward the
	// attributes of a kind it does not know. For UnitKindK8sPod use
	// K8sPodAttributes.
	Attributes json.RawMessage `json:"attributes,omitempty"`
}

// K8sPodAttributes is Unit.Attributes for UnitKindK8sPod.
type K8sPodAttributes struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
	// Labels holds the labels of the pod.
	Labels map[string]string `json:"labels,omitempty"`
	// WorkloadKind and WorkloadName name the pod's owner workload, as the sensor
	// resolves it. For a pod with no owner they name the pod itself.
	WorkloadKind string `json:"workloadKind,omitempty"`
	WorkloadName string `json:"workloadName,omitempty"`
	// WorkloadLabels holds the labels of the owner workload. Nil when the
	// workload is the pod itself. Non-nil whenever it is not, also when the
	// owner has no labels: an empty map is sent as {}. Same rule as
	// armotypes.RuntimeAlertK8sDetails.WorkloadLabels.
	WorkloadLabels map[string]string `json:"workloadLabels,omitzero"`
}

// NewK8sPodUnit builds the unit of a start entry for a pod.
func NewK8sPodUnit(podUID string, startedAt time.Time, attrs K8sPodAttributes) (Unit, error) {
	raw, err := json.Marshal(attrs)
	if err != nil {
		return Unit{}, fmt.Errorf("encode k8s-pod attributes: %w", err)
	}
	return Unit{Kind: UnitKindK8sPod, ID: podUID, StartedAt: &startedAt, Attributes: raw}, nil
}

// K8sPodAttributes decodes the attributes of a UnitKindK8sPod unit. It returns
// an error for another kind, and for a unit with no attributes (a stop entry).
func (u Unit) K8sPodAttributes() (*K8sPodAttributes, error) {
	if u.Kind != UnitKindK8sPod {
		return nil, fmt.Errorf("unit kind is %q, not %q", u.Kind, UnitKindK8sPod)
	}
	if len(u.Attributes) == 0 {
		return nil, fmt.Errorf("unit %q has no attributes", u.ID)
	}
	var attrs K8sPodAttributes
	if err := json.Unmarshal(u.Attributes, &attrs); err != nil {
		return nil, fmt.Errorf("decode k8s-pod attributes: %w", err)
	}
	return &attrs, nil
}
