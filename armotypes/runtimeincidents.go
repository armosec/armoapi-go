package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
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
	IncidentCategory  IncidentCategory `json:"incidentCategory" bson:"incidentCategory" `
	CreationTimestamp time.Time        `json:"creationTimestamp" bson:"creationTimestamp"`
	Title             string           `json:"title" bson:"title"`
	Severity          string           `json:"incidentSeverity" bson:"incidentSeverity"`
	SeverityScore     int              `json:"severityScore" bson:"severityScore"`
	Mitigation        string           `json:"mitigation" bson:"mitigation"`
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

type AlertType int

const (
	AlertTypeRule AlertType = iota
	AlertTypeMalware
)

type BaseRuntimeAlert struct {
	// AlertName is either RuleName or MalwareName
	AlertName string `json:"alertName,omitempty" bson:"name,omitempty"`
	// Arguments of specific alerts (e.g. for unexpected files: open file flags; for unexpected process: return code)
	Arguments map[string]interface{} `json:"arguments,omitempty" bson:"arguments,omitempty"`
	// Command line of the process
	CommandLine string `json:"commandLine,omitempty" bson:"commandLine,omitempty"`
	// Fix suggestions
	FixSuggestions string `json:"fixSuggestions,omitempty" bson:"fixSuggestions,omitempty"`
	// Is part of the image
	IsPartOfImage bool `json:"isPartOfImage,omitempty" bson:"isPartOfImage,omitempty"`
	// Parent Process ID
	PPID uint32 `json:"ppid,omitempty" bson:"ppid,omitempty"`
	// PPIDComm - Parent Process Name
	PPIDComm string `json:"ppidComm,omitempty" bson:"ppidComm,omitempty"`
	// MD5 hash of the file that was infected
	MD5Hash string `json:"md5Hash,omitempty" bson:"md5Hash,omitempty"`
	// SHA1 hash of the file that was infected
	SHA1Hash string `json:"sha1Hash,omitempty" bson:"sha1Hash,omitempty"`
	// SHA256 hash of the file that was infected
	SHA256Hash string `json:"sha256Hash,omitempty" bson:"sha256Hash,omitempty"`
	// Severity of the alert
	Severity int `json:"severity,omitempty" bson:"severity,omitempty"`
	// Size of the file that was infected
	Size string `json:"size,omitempty" bson:"size,omitempty"`
	// Command line
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type RuleAlert struct {
	// Rule ID
	RuleID string `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
}

type MalwareAlert struct {
	MalwareDescription string `json:"malwareDescription,omitempty" bson:"malwareDescription,omitempty"`
}

type RuntimeAlertProcessDetails struct {
	// Process Name
	Comm string `json:"comm,omitempty" bson:"comm,omitempty"`
	// GID of the process
	GID uint32 `json:"gid,omitempty" bson:"gid,omitempty"`
	// Group name of the process
	GroupName string `json:"groupName,omitempty" bson:"groupName,omitempty"`
	// Path to the file that was infected
	Path string `json:"path,omitempty" bson:"path,omitempty"`
	// Process ID
	PID uint32 `json:"pid,omitempty" bson:"pid,omitempty"`
	// UID of the process
	UID uint32 `json:"uid,omitempty" bson:"uid,omitempty"`
	// User name of the process
	UserName string `json:"userName,omitempty" bson:"userName,omitempty"`
}

type RuntimeAlertK8sDetails struct {
	ClusterName       string `json:"clusterName" bson:"clusterName"`
	ContainerName     string `json:"containerName,omitempty" bson:"containerName,omitempty"`
	HostNetwork       bool   `json:"hostNetwork,omitempty" bson:"hostNetwork,omitempty"`
	Image             string `json:"image,omitempty" bson:"image,omitempty"`
	ImageDigest       string `json:"imageDigest,omitempty" bson:"imageDigest,omitempty"`
	Namespace         string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	NodeName          string `json:"nodeName,omitempty" bson:"nodeName,omitempty"`
	ContainerID       string `json:"containerID,omitempty" bson:"containerID,omitempty"`
	PodName           string `json:"podName,omitempty" bson:"podName,omitempty"`
	PodNamespace      string `json:"podNamespace,omitempty" bson:"podNamespace,omitempty"`
	WorkloadName      string `json:"workloadName" bson:"workloadName"`
	WorkloadNamespace string `json:"workloadNamespace,omitempty" bson:"namespace,omitempty"`
	WorkloadKind      string `json:"workloadKind" bson:"workloadKind"`
}

type RuntimeAlert struct {
	BaseRuntimeAlert           `json:",inline"`
	RuleAlert                  `json:",inline"`
	MalwareAlert               `json:",inline"`
	RuntimeAlertProcessDetails `json:",inline"`
	RuntimeAlertK8sDetails     `json:",inline"`
	AlertType                  AlertType `json:"alertType" bson:"alertType"`
	// Hostname is the name of the node agent pod
	HostName string `json:"hostName" bson:"hostName"`
	Message  string `json:"message" bson:"message"`
}
