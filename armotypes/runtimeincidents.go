package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes/cdr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authentication/user"
)

type AlertType int

const (
	AlertTypeRule AlertType = iota
	AlertTypeMalware
	AlertTypeAdmission
	AlertTypeCdr
	AlertTypeHttpRule
	AlertTypeNetworkScan
)

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
	UserSpace bool `json:"userSpace,omitempty" bson:"userSpace,omitempty"`
	// Native/Source code
	NativeCode *bool `json:"nativeCode,omitempty" bson:"nativeCode,omitempty"`
	// Anomaly flag
	Anomaly bool `json:"anomaly,omitempty" bson:"anomaly,omitempty"`
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
	MalwareFile        File             `json:"malwareFile,omitempty" bson:"malwareFile,omitempty"`
	Action             string           `json:"action,omitempty" bson:"action,omitempty"`
	DetectionMethod    string           `json:"detectionMethod,omitempty" bson:"detectionMethod,omitempty"`
	ProcessTree        ProcessTree      `json:"processTree,omitempty" bson:"processTree,omitempty"`
	Signature          MalwareSignature `json:"signature,omitempty" bson:"signature,omitempty"`
	MalwareDescription string           `json:"malwareDescription,omitempty" bson:"malwareDescription,omitempty"`
}

type HttpRuleAlert struct {
	Request struct {
		Method  string            `json:"method,omitempty" bson:"method,omitempty"`   // e.g., "GET"
		URL     string            `json:"url,omitempty" bson:"url,omitempty"`         // e.g., "/index.html"
		Header  map[string]string `json:"header,omitempty" bson:"header,omitempty"`   // e.g., "Content-Type" -> ["application/json"]
		Body    string            `json:"body,omitempty" bson:"body,omitempty"`       // e.g., "<html>...</html>"
		Proto   string            `json:"proto,omitempty" bson:"proto,omitempty"`     // e.g., "HTTP/1.1"
		Payload string            `json:"payload,omitempty" bson:"payload,omitempty"` // e.g., "'OR 1=1"'"
	} `json:"request,omitempty" bson:"request,omitempty"`

	Response struct {
		StatusCode   int               `json:"statusCode,omitempty" bson:"statusCode,omitempty"`     // e.g., 200
		Header       map[string]string `json:"header,omitempty" bson:"header,omitempty"`             // e.g., "Content-Type" -> ["application/json"]
		Body         string            `json:"body,omitempty" bson:"body,omitempty"`                 // e.g., "<html>...</html>"
		Proto        string            `json:"proto,omitempty" bson:"proto,omitempty"`               // e.g., "HTTP/1.1"
		FullResponse string            `json:"fullResponse,omitempty" bson:"fullResponse,omitempty"` // e.g., "{...}"
	} `json:"response,omitempty" bson:"response,omitempty"`

	SourcePodInfo RuntimeAlertK8sDetails `json:"sourcePodInfo,omitempty" bson:"podInfo,omitempty"`
	AttackerIp    string                 `json:"attackerIp,omitempty" bson:"attackerIp,omitempty"`
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
	ClusterName       string            `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
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
	WorkloadName      string            `json:"workloadName,omitempty" bson:"workloadName,omitempty"`
	WorkloadNamespace string            `json:"workloadNamespace,omitempty" bson:"workloadNamespace,omitempty"`
	WorkloadKind      string            `json:"workloadKind,omitempty" bson:"workloadKind,omitempty"`
}

type NetworkScanAlert struct {
	Domain    string   `json:"domain,omitempty" bson:"domain,omitempty"`
	Addresses []string `json:"addresses,omitempty" bson:"addresses,omitempty"`
}

type RuntimeAlert struct {
	BaseRuntimeAlert       `json:",inline" bson:"inline"`
	RuleAlert              `json:",inline" bson:"inline"`
	MalwareAlert           `json:",inline" bson:"inline"`
	AdmissionAlert         `json:",inline" bson:"inline"`
	RuntimeAlertK8sDetails `json:",inline" bson:"inline"`
	cdr.CdrAlert           `json:"cdrevent,omitempty" bson:"cdrevent"`
	HttpRuleAlert          `json:",inline" bson:"inline"`
	NetworkScanAlert       `json:"networkscan,inline" bson:"networkscan"`
	AlertType              AlertType `json:"alertType" bson:"alertType"`
	// Rule ID
	RuleID string `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
	// Hostname is the name of the node agent pod
	HostName string `json:"hostName" bson:"hostName"`
	Message  string `json:"message" bson:"message"`
}

type ProcessTree struct {
	ProcessTree Process `json:"processTree" bson:"processTree"`
	UniqueID    uint32  `json:"uniqueID,omitempty" bson:"uniqueID,omitempty"`
	ContainerID string  `json:"containerID,omitempty" bson:"containerID,omitempty"`
}

type KDRMonitoredEntitiesCounters struct {
	ClustersCount   int `json:"clustersCount"`
	NodesCount      int `json:"nodesCount"`
	NamespacesCount int `json:"namespacesCount"`
	PodsCount       int `json:"podsCount"`
	ContainersCount int `json:"containersCount"`
}

type KDRMonitoredClusters struct {
	MonitoredClusters    []string `json:"monitored,omitempty"`
	NotMonitoredClusters []string `json:"notMonitored,omitempty"`
}

type RuntimeIncidentExceptionPolicy struct {
	BaseExceptionPolicy `json:",inline"`
	Name                string `json:"name"`
	IncidentTypeId      string `json:"incidentTypeId"`
	Severity            string `json:"severity"`
	SeverityScore       int    `json:"severityScore"`
}

// FindProcessByPID searches for a process by PID in the process tree
func (pt *ProcessTree) FindProcessByPID(pid uint32) *Process {
	return findProcessRecursive(&pt.ProcessTree, pid)
}

// findProcessRecursive recursively searches for a process by PID
func findProcessRecursive(proc *Process, pid uint32) *Process {
	if proc.PID == pid {
		return proc
	}

	for _, child := range proc.Children {
		if found := findProcessRecursive(&child, pid); found != nil {
			return found
		}
	}
	return nil
}
