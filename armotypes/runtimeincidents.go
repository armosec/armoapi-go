package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IncidentCategory string

const (
	RuntimeIncidentCategoryMalware IncidentCategory = "Malware"
	RuntimeIncidentCategoryAnomaly IncidentCategory = "Anomaly"
)

type RuntimeIncident struct {
	PortalBase `json:",inline" bson:"inline"`
	// details of the incident triggers
	RuntimeIncidentResource `json:",inline" bson:"inline"`
	RuntimeAlert            `json:",inline" bson:"inline"`
	// category of the incident
	IncidentCategory IncidentCategory `json:"incidentCategory" bson:"incidentCategory" `
	Timestamp        time.Time        `json:"timestamp" bson:"timestamp"`
	Title            string           `json:"title" bson:"title"`
	Severity         string           `json:"incidentSeverity" bson:"incidentSeverity"`
	SeverityScore    int              `json:"severityScore" bson:"severityScore"`
	Mitigation       string           `json:"mitigation" bson:"mitigation"`
	// alerts / events that are part of this incident
	RelatedAlerts []RuntimeAlert `json:"relatedAlerts" bson:"relatedAlerts"`
	// user gestures
	SeenAt                *time.Time `json:"seenAt" bson:"seenAt"`
	IsDismissed           bool       `json:"isDismissed" bson:"isDismissed"`
	MarkedAsFalsePositive bool       `json:"markedAsFalsePositive" bson:"markedAsFalsePositive"`
	// for future use
	RelatedResources []RuntimeIncidentResource `json:"relatedResources" bson:"relatedResources"`
}

type RuntimeIncidentResource struct {
	Spiffe      string                       `json:"spiffe" bson:"spiffe" `
	ResourceID  string                       `json:"resourceID" bson:"resourceID"` // hash of the resource to cross with various DBs
	Designators identifiers.PortalDesignator `json:"designators" bson:"designators"`
}

type RuleAlert struct {
	Severity       int    `json:"severity,omitempty" bson:"severity,omitempty"`       // PriorityToStatus(failedRule.Priority()),
	ProcessName    string `json:"processName,omitempty" bson:"processName,omitempty"` // failedRule.Event().Comm,
	FixSuggestions string `json:"fixSuggestions,omitempty" bson:"fixSuggestions,omitempty"`
	PID            uint32 `json:"pid,omitempty" bson:"pid,omitempty"`   // Process ID
	PPID           uint32 `json:"ppid,omitempty" bson:"ppid,omitempty"` //  Parent Process ID
	UID            uint32 `json:"uid,omitempty" bson:"uid,omitempty"`   // User ID of the process
	GID            uint32 `json:"gid,omitempty" bson:"gid,omitempty"`   // Group ID of the process
}

type MalwareAlert struct {
	MalwareName        string `json:"malwareName,omitempty" bson:"malwareName,omitempty"`
	MalwareDescription string `json:"malwareDescription,omitempty" bson:"malwareDescription,omitempty"`
	// Path to the file that was infected
	Path string `json:"path,omitempty" bson:"path,omitempty"`
	// Hash of the file that was infected
	Hash string `json:"hash,omitempty" bson:"hash,omitempty"`
	// Size of the file that was infected
	Size string `json:"size,omitempty" bson:"size,omitempty"`
	// Is part of the image
	IsPartOfImage bool `json:"isPartOfImage,omitempty" bson:"isPartOfImage,omitempty"`
	// K8s resource that was infected
	Resource schema.GroupVersionResource `json:"resource,omitempty" bson:"resource,omitempty"`
	// K8s container image that was infected
	ContainerImage string `json:"containerImage,omitempty" bson:"containerImage,omitempty"`
}

type RuntimeAlert struct {
	RuleAlert     `json:",inline"`
	MalwareAlert  `json:",inline"`
	RuleName      string    `json:"ruleName" bson:"ruleName"`
	RuleID        string    `json:"ruleID" bson:"ruleID"`
	Message       string    `json:"message" bson:"message"`
	ContainerID   string    `json:"containerID,omitempty" bson:"containerID,omitempty"`
	ContainerName string    `json:"containerName,omitempty" bson:"containerName,omitempty"`
	PodNamespace  string    `json:"podNamespace,omitempty" bson:"podNamespace,omitempty"`
	PodName       string    `json:"podName,omitempty" bson:"podName,omitempty"`
	HostName      string    `json:"hostName" bson:"hostName"`
	NodeName      string    `json:"nodeName" bson:"nodeName"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
}
