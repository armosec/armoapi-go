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

type AlertType int

const (
	AlertTypeRule AlertType = iota
	AlertTypeMalware
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

type RuleAlert struct {
	// Rule Name
	RuleName string `json:"ruleName,omitempty" bson:"ruleName,omitempty"`
	// Rule ID
	RuleID string `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
	// Severity of the alert
	Severity int `json:"severity,omitempty" bson:"severity,omitempty"`
	// Process name
	ProcessName string `json:"processName,omitempty" bson:"processName,omitempty"`
	// Command line
	CommandLine string `json:"commandLine,omitempty" bson:"commandLine,omitempty"`
	// Fix suggestions
	FixSuggestions string `json:"fixSuggestions,omitempty" bson:"fixSuggestions,omitempty"`
	// Parent Process ID
	PPID uint32 `json:"ppid,omitempty" bson:"ppid,omitempty"`
	// PPIDComm - Parent Process Name
	PPIDComm string `json:"ppidComm,omitempty" bson:"ppidComm,omitempty"`
	// Is part of the image
	IsPartOfImage bool `json:"isPartOfImage,omitempty" bson:"isPartOfImage,omitempty"`
	// MD5 hash of the file that was infected
	MD5Hash string `json:"md5Hash,omitempty" bson:"md5Hash,omitempty"`
	// SHA1 hash of the file that was infected
	SHA1Hash string `json:"sha1Hash,omitempty" bson:"sha1Hash,omitempty"`
	// SHA256 hash of the file that was infected
	SHA256Hash string `json:"sha256Hash,omitempty" bson:"sha256Hash,omitempty"`
}

type MalwareAlert struct {
	// Name of the malware
	MalwareName string `json:"malwareName,omitempty" bson:"malwareName,omitempty"`
	// Description of the malware
	MalwareDescription string `json:"malwareDescription,omitempty" bson:"malwareDescription,omitempty"`
	// Path to the file that was infected
	Path string `json:"path,omitempty" bson:"path,omitempty"`
	// MD5 hash of the file that was infected
	MD5Hash string `json:"md5Hash,omitempty" bson:"md5Hash,omitempty"`
	// SHA1 hash of the file that was infected
	SHA1Hash string `json:"sha1Hash,omitempty" bson:"sha1Hash,omitempty"`
	// SHA256 hash of the file that was infected
	SHA256Hash string `json:"sha256Hash,omitempty" bson:"sha256Hash,omitempty"`
	// Size of the file that was infected
	Size string `json:"size,omitempty" bson:"size,omitempty"`
	// Is part of the image
	IsPartOfImage bool `json:"isPartOfImage,omitempty" bson:"isPartOfImage,omitempty"`
	// Severity of the malware
	Severity int `json:"severity,omitempty" bson:"severity,omitempty"`
	// Parent Process ID of the process that was infected
	PPID uint32 `json:"ppid,omitempty" bson:"ppid,omitempty"`
	// Parent Process Name of the process that was infected
	PPIDComm string `json:"ppidComm,omitempty" bson:"ppidComm,omitempty"`
	// Command Line of the process that was infected
	CommandLine string `json:"commandLine,omitempty" bson:"commandLine,omitempty"`
	// Fix suggestions
	FixSuggestions string `json:"fixSuggestions,omitempty" bson:"fixSuggestions,omitempty"`
}

type RuntimeAlert struct {
	RuleAlert     `json:",inline"`
	MalwareAlert  `json:",inline"`
	AlertType     AlertType `json:"alertType" bson:"alertType"`
	Message       string    `json:"message" bson:"message"`
	ContainerID   string    `json:"containerID,omitempty" bson:"containerID,omitempty"`
	ContainerName string    `json:"containerName,omitempty" bson:"containerName,omitempty"`
	Namespace     string    `json:"namespace,omitempty" bson:"namespace,omitempty"`
	PodName       string    `json:"podName,omitempty" bson:"podName,omitempty"`
	HostNetwork   bool      `json:"hostNetwork,omitempty" bson:"hostNetwork,omitempty"`
	Image         string    `json:"image,omitempty" bson:"image,omitempty"`
	ImageDigest   string    `json:"imageDigest,omitempty" bson:"imageDigest,omitempty"`
	HostName      string    `json:"hostName" bson:"hostName"`
	NodeName      string    `json:"nodeName" bson:"nodeName"`
	WorkloadName  string    `json:"workloadName" bson:"workloadName"`
	WorkloadKind  string    `json:"workloadKind" bson:"workloadKind"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
}
