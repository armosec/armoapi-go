package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes/cdr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authentication/user"
)

type Process struct {
	PID        uint32    `json:"pid,omitempty" bson:"pid,omitempty"`
	Cmdline    string    `json:"cmdline,omitempty" bson:"cmdline,omitempty"`
	Comm       string    `json:"comm,omitempty" bson:"comm,omitempty"`
	PPID       uint32    `json:"ppid,omitempty" bson:"ppid,omitempty"`
	Pcomm      string    `json:"pcomm,omitempty" bson:"pcomm,omitempty"`
	Hardlink   string    `json:"hardlink,omitempty" bson:"hardlink,omitempty"`
	Uid        *uint32   `json:"uid,omitempty" bson:"uid,omitempty"`
	Gid        *uint32   `json:"gid,omitempty" bson:"gid,omitempty"`
	UserName   string    `json:"userName,omitempty" bson:"userName,omitempty"`
	GroupName  string    `json:"groupName,omitempty" bson:"groupName,omitempty"`
	StartTime  time.Time `json:"startTime,omitempty" bson:"startTime,omitempty"`
	UpperLayer *bool     `json:"upperLayer,omitempty" bson:"upperLayer,omitempty"`
	Cwd        string    `json:"cwd,omitempty" bson:"cwd,omitempty"`
	Path       string    `json:"path,omitempty" bson:"path,omitempty"`
	Children   []Process `json:"children,omitempty" bson:"children,omitempty"`
}

type CloudMetadata struct {
	// Provider is the cloud provider name (e.g. aws, gcp, azure).
	Provider     string   `json:"provider,omitempty" bson:"provider,omitempty"`
	InstanceID   string   `json:"instance_id,omitempty" bson:"instance_id,omitempty"`
	InstanceType string   `json:"instance_type,omitempty" bson:"instance_type,omitempty"`
	Region       string   `json:"region,omitempty" bson:"region,omitempty"`
	Zone         string   `json:"zone,omitempty" bson:"zone,omitempty"`
	PrivateIP    string   `json:"private_ip,omitempty" bson:"private_ip,omitempty"`
	PublicIP     string   `json:"public_ip,omitempty" bson:"public_ip,omitempty"`
	Hostname     string   `json:"hostname,omitempty" bson:"hostname,omitempty"`
	AccountID    string   `json:"account_id,omitempty" bson:"account_id,omitempty"`
	Services     []string `json:"services,omitempty" bson:"services,omitempty"`
}

type AlertType int

const (
	AlertTypeRule AlertType = iota
	AlertTypeMalware
	AlertTypeAdmission
	AlertTypeCdr
)

type StackFrame struct {
	// Frame ID
	FrameID string `json:"frameId,omitempty" bson:"frameId,omitempty"`
	// Function name
	Function string `json:"function,omitempty" bson:"function,omitempty"`
	// File name
	File string `json:"file,omitempty" bson:"file,omitempty"`
	// Line number
	Line *int `json:"line,omitempty" bson:"line,omitempty"`
	// Address
	Address string `json:"address,omitempty" bson:"address,omitempty"`
	// Arguments
	Arguments []string `json:"arguments,omitempty" bson:"arguments,omitempty"`
	// User/Kernel space
	UserSpace bool `json:"userSpace" bson:"userSpace"`
	// Native/Source code
	NativeCode *bool `json:"nativeCode,omitempty" bson:"nativeCode,omitempty"`
}

type Trace struct {
	// Trace ID
	TraceID string `json:"traceId,omitempty" bson:"traceId,omitempty"`
	// Stack trace
	Stack []StackFrame `json:"stack,omitempty" bson:"stack,omitempty"`
	// Package name
	Package string `json:"package,omitempty" bson:"package,omitempty"`
	// Language
	Language string `json:"language,omitempty" bson:"language,omitempty"`
}

type BaseRuntimeAlert struct {
	// AlertName is either RuleName or MalwareName
	AlertName string `json:"alertName,omitempty" bson:"name,omitempty"`
	// Arguments of specific alerts (e.g. for unexpected files: open file flags; for unexpected process: return code)
	Arguments map[string]interface{} `json:"arguments,omitempty" bson:"arguments,omitempty"`
	// Infected process id
	InfectedPID uint32 `json:"infectedPID,omitempty" bson:"infectedPID,omitempty"`
	// Process tree unique id
	ProcessTreeUniqueID uint32 `json:"processTreeUniqueID,omitempty" bson:"processTreeUniqueID,omitempty"`
	// Fix suggestions
	FixSuggestions string `json:"fixSuggestions,omitempty" bson:"fixSuggestions,omitempty"`
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
	// Timestamp of the alert
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	// Nanoseconds of the alert
	Nanoseconds uint64 `json:"nanoseconds,omitempty" bson:"nanoseconds,omitempty"`
	// Trace of the alert
	Trace Trace `json:"trace,omitempty" bson:"trace,omitempty"`
}

type RuleAlert struct {
	// Rule Description
	RuleDescription string `json:"ruleDescription,omitempty" bson:"ruleDescription,omitempty"`
}

type MalwareAlert struct {
	MalwareDescription string `json:"malwareDescription,omitempty" bson:"malwareDescription,omitempty"`
}

type AdmissionAlert struct {
	Kind             schema.GroupVersionKind     `json:"kind,omitempty" bson:"kind,omitempty"`
	RequestNamespace string                      `json:"requestNamespace,omitempty" bson:"requestNamespace,omitempty"`
	ObjectName       string                      `json:"objectName,omitempty" bson:"objectName,omitempty"`
	Resource         schema.GroupVersionResource `json:"resource,omitempty" bson:"resource,omitempty"`
	Subresource      string                      `json:"subresource,omitempty" bson:"subresource,omitempty"`
	Operation        admission.Operation         `json:"operation,omitempty" bson:"operation,omitempty"`
	Options          *unstructured.Unstructured  `json:"options,omitempty" bson:"options,omitempty"`
	DryRun           bool                        `json:"dryRun,omitempty" bson:"dryRun,omitempty"`
	Object           *unstructured.Unstructured  `json:"object,omitempty" bson:"object,omitempty"`
	OldObject        *unstructured.Unstructured  `json:"oldObject,omitempty" bson:"oldObject,omitempty"`
	UserInfo         *user.DefaultInfo           `json:"userInfo,omitempty" bson:"userInfo,omitempty"`
}

type RuntimeAlertK8sDetails struct {
	ClusterName       string            `json:"clusterName" bson:"clusterName"`
	ContainerName     string            `json:"containerName,omitempty" bson:"containerName,omitempty"`
	HostNetwork       *bool             `json:"hostNetwork,omitempty" bson:"hostNetwork,omitempty"`
	OldImage          string            `json:"oldImage,omitempty" bson:"oldImage,omitempty"`
	Image             string            `json:"image,omitempty" bson:"image,omitempty"`
	ImageDigest       string            `json:"imageDigest,omitempty" bson:"imageDigest,omitempty"`
	Namespace         string            `json:"namespace,omitempty" bson:"namespace,omitempty"`
	NodeName          string            `json:"nodeName,omitempty" bson:"nodeName,omitempty"`
	ContainerID       string            `json:"containerID,omitempty" bson:"containerID,omitempty"`
	PodName           string            `json:"podName,omitempty" bson:"podName,omitempty"`
	PodNamespace      string            `json:"podNamespace,omitempty" bson:"podNamespace,omitempty"`
	PodLabels         map[string]string `json:"podLabels,omitempty" bson:"podLabels,omitempty"`
	WorkloadName      string            `json:"workloadName" bson:"workloadName"`
	WorkloadNamespace string            `json:"workloadNamespace,omitempty" bson:"workloadNamespace,omitempty"`
	WorkloadKind      string            `json:"workloadKind" bson:"workloadKind"`
}

type RuntimeAlert struct {
	BaseRuntimeAlert       `json:",inline" bson:"inline"`
	RuleAlert              `json:",inline" bson:"inline"`
	MalwareAlert           `json:",inline" bson:"inline"`
	AdmissionAlert         `json:",inline" bson:"inline"`
	RuntimeAlertK8sDetails `json:",inline" bson:"inline"`
	cdr.CdrAlert           `json:"cdrevent" bson:"cdrevent"`
	AlertType              AlertType `json:"alertType" bson:"alertType"`
	// Rule ID
	RuleID string `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
	// Hostname is the name of the node agent pod
	HostName string `json:"hostName" bson:"hostName"`
	Message  string `json:"message" bson:"message"`
}

type ProcessTree struct {
	ProcessTree Process `json:"processTree" bson:"processTree"`
	UniqueID    uint32  `json:"uniqueID" bson:"uniqueID"`
	ContainerID string  `json:"containerID" bson:"containerID"`
}

type KDRMonitoredEntitiesCounters struct {
	ClustersCount   int `json:"clustersCount"`
	NodesCount      int `json:"nodesCount"`
	NamespacesCount int `json:"namespacesCount"`
	PodsCount       int `json:"podsCount"`
	ContainersCount int `json:"containersCount"`
}

type RuntimeIncidentExceptionPolicy struct {
	BaseExceptionPolicy `json:",inline"`
	Name                string `json:"name"`
	IncidentTypeId      string `json:"incidentTypeId"`
	Severity            string `json:"severity"`
	SeverityScore       int    `json:"severityScore"`
}
