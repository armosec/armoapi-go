package armotypes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/armosec/armoapi-go/armotypes/cdr"
	"github.com/armosec/armoapi-go/armotypes/common"
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

type AlertSourcePlatform int

const (
	AlertSourcePlatformUnknown AlertSourcePlatform = iota
	AlertSourcePlatformK8s
	AlertSourcePlatformHost
	AlertSourcePlatformCloud
	AlertSourcePlatformECS
)

type ProfileType int

const (
	ApplicationProfile ProfileType = iota
	NetworkProfile
)

type ProfileDependency int

const (
	Required ProfileDependency = iota
	Optional
	NotRequired
)

type ProfileMetadata struct {
	Status            string            `json:"status,omitempty" bson:"status,omitempty"`
	Completion        string            `json:"completion,omitempty" bson:"completion,omitempty"`
	Name              string            `json:"name,omitempty" bson:"name,omitempty"`
	FailOnProfile     bool              `json:"failOnProfile" bson:"failOnProfile"`
	Type              ProfileType       `json:"type" bson:"type"`
	ProfileDependency ProfileDependency `json:"profileDependency,omitempty" bson:"profileDependency,omitempty"`
	Error             string            `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
}

type HostType string

const (
	HostTypeAci        HostType = "aci"
	HostTypeAks        HostType = "aks"
	HostTypeAutopilot  HostType = "autopilot"
	HostTypeAzureVm    HostType = "azurevm"
	HostTypeCloudRun   HostType = "cloudrun"
	HostTypeDoks       HostType = "doks"
	HostTypeDroplet    HostType = "droplet"
	HostTypeEc2        HostType = "ec2"
	HostTypeEcsEc2     HostType = "ecs-ec2"
	HostTypeEcsFargate HostType = "ecs-fargate"
	HostTypeEksEc2     HostType = "eks-ec2"
	HostTypeEksFargate HostType = "eks-fargate"
	HostTypeGce        HostType = "gce"
	HostTypeGke        HostType = "gke"
	HostTypeKubernetes HostType = "kubernetes"
	HostTypeOther      HostType = "other"
)

type Provider string

const (
	ProviderAlibaba      Provider = "alibaba"
	ProviderAws          Provider = "aws"
	ProviderAzure        Provider = "azure"
	ProviderDigitalOcean Provider = "digitalocean"
	ProviderEquinixMetal Provider = "equinixmetal" // formerly Packet
	ProviderExoscale     Provider = "exoscale"
	ProviderGcp          Provider = "gcp"
	ProviderHetzner      Provider = "hetzner"
	ProviderIBM          Provider = "ibm"
	ProviderLinode       Provider = "linode"
	ProviderOpenStack    Provider = "openstack"
	ProviderOracle       Provider = "oracle"
	ProviderOther        Provider = "other"
	ProviderScaleway     Provider = "scaleway"
	ProviderVMware       Provider = "vmware"
	ProviderVultr        Provider = "vultr"
)

type CloudMetadata struct {
	AccountID    string   `json:"account_id,omitempty" bson:"account_id,omitempty"`
	HostType     HostType `json:"host_type,omitempty" bson:"host_type,omitempty"`
	Hostname     string   `json:"hostname,omitempty" bson:"hostname,omitempty"`
	InstanceID   string   `json:"instance_id,omitempty" bson:"instance_id,omitempty"`
	InstanceType string   `json:"instance_type,omitempty" bson:"instance_type,omitempty"` // m5.large, ...
	OrgID        string   `json:"org_id,omitempty" bson:"org_id,omitempty"`
	PrivateIP    string   `json:"private_ip,omitempty" bson:"private_ip,omitempty"`
	PrivateIPs   []string `json:"private_ips,omitempty" bson:"private_ips,omitempty"`
	// Provider is the cloud provider name (e.g. aws, gcp, azure).
	Provider      Provider `json:"provider,omitempty" bson:"provider,omitempty"`
	PublicIP      string   `json:"public_ip,omitempty" bson:"public_ip,omitempty"`
	PublicIPs     []string `json:"public_ips,omitempty" bson:"public_ips,omitempty"`
	Region        string   `json:"region,omitempty" bson:"region,omitempty"`
	ResourceGroup string   `json:"resource_group,omitempty" bson:"resource_group,omitempty"` // Azure Resource Group
	Services      []string `json:"services,omitempty" bson:"services,omitempty"`
	Zone          string   `json:"zone,omitempty" bson:"zone,omitempty"`
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
	AlertName string `json:"alertName,omitempty" bson:"alertName,omitempty"`
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
	// Unique ID of the alert
	UniqueID string `json:"uniqueID,omitempty" bson:"uniqueID,omitempty"`
	// Profile metadata
	ProfileMetadata *ProfileMetadata `json:"profileMetadata,omitempty" bson:"profileMetadata,omitempty"`
	// Identifiers of the alert
	Identifiers *common.Identifiers `json:"identifiers,omitempty" bson:"identifiers,omitempty"`
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

	// Enrichment fields for the layer 7 alert
	PastActivityCount *int   `json:"pastActivityCount,omitempty" bson:"pastActivityCount,omitempty"`
	Country           string `json:"country,omitempty" bson:"country,omitempty"`
	City              string `json:"city,omitempty" bson:"city,omitempty"`
	Explain           string `json:"explain,omitempty" bson:"explain,omitempty"`
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
	ClusterUID        string            `json:"clusterUID,omitempty" bson:"clusterUID,omitempty"`
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
	PodUID            string            `json:"podUID,omitempty" bson:"podUID,omitempty"`
	PodLabels         map[string]string `json:"podLabels,omitempty" bson:"podLabels,omitempty"`
	WorkloadName      string            `json:"workloadName,omitempty" bson:"workloadName,omitempty"`
	WorkloadNamespace string            `json:"workloadNamespace,omitempty" bson:"workloadNamespace,omitempty"`
	WorkloadKind      string            `json:"workloadKind,omitempty" bson:"workloadKind,omitempty"`
	WorkloadUID       string            `json:"workloadUID,omitempty" bson:"workloadUID,omitempty"`
}

type RuntimeAlertECSDetails struct {
	ClusterARN        string `json:"clusterArn,omitempty" bson:"clusterArn,omitempty"`
	ECSClusterName    string `json:"ecsClusterName,omitempty" bson:"ecsClusterName,omitempty"`
	ServiceName       string `json:"serviceName,omitempty" bson:"serviceName,omitempty"`
	TaskARN           string `json:"taskArn,omitempty" bson:"taskArn,omitempty"`
	TaskFamily        string `json:"taskFamily,omitempty" bson:"taskFamily,omitempty"`
	TaskDefinitionARN string `json:"taskDefinitionArn,omitempty" bson:"taskDefinitionArn,omitempty"`
	ECSContainerName  string `json:"ecsContainerName,omitempty" bson:"ecsContainerName,omitempty"`
	ContainerARN      string `json:"containerArn,omitempty" bson:"containerArn,omitempty"`
	ECSContainerID    string `json:"ecsContainerID,omitempty" bson:"ecsContainerID,omitempty"`
	ContainerInstance string `json:"containerInstance,omitempty" bson:"containerInstance,omitempty"` // EC2 instance ID (EC2 launch type only)
	LaunchType        string `json:"launchType,omitempty" bson:"launchType,omitempty"`               // EC2 or FARGATE
	AvailabilityZone  string `json:"availabilityZone,omitempty" bson:"availabilityZone,omitempty"`
	ECSImage          string `json:"ecsImage,omitempty" bson:"ecsImage,omitempty"`
	ECSImageDigest    string `json:"ecsImageDigest,omitempty" bson:"ecsImageDigest,omitempty"`
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
	RuntimeAlertECSDetails `json:",inline" bson:"inline"`
	cdr.CdrAlert           `json:"cdrevent,omitempty" bson:"cdrevent"`
	HttpRuleAlert          `json:",inline" bson:"inline"`
	NetworkScanAlert       `json:"networkscan,inline" bson:"networkscan"`
	AlertType              AlertType           `json:"alertType" bson:"alertType"`
	AlertSourcePlatform    AlertSourcePlatform `json:"alertSourcePlatform" bson:"alertSourcePlatform"`
	// Rule ID
	RuleID string `json:"ruleID,omitempty" bson:"ruleID,omitempty"`
	// IsTriggerAlert indicates if this alert is a trigger alert
	IsTriggerAlert bool `json:"isTriggerAlert,omitempty" bson:"isTriggerAlert,omitempty"`
	// Hostname is the name of the node agent pod
	HostName string          `json:"hostName" bson:"hostName"`
	Message  string          `json:"message" bson:"message"`
	Fields   json.RawMessage `json:"fields,omitempty" bson:"fields,omitempty"`
}

func (ra *RuntimeAlert) GetAlertSourcePlatform() AlertSourcePlatform {
	if ra.AlertType == AlertTypeCdr {
		return AlertSourcePlatformCloud
	}

	if ra.PodName != "" {
		return AlertSourcePlatformK8s
	}

	if ra.TaskARN != "" || ra.ClusterARN != "" {
		return AlertSourcePlatformECS
	}

	return AlertSourcePlatformHost
}

func (ra *RuntimeAlert) Validate() error {
	if ra.RuleID == "" {
		return fmt.Errorf("ruleID is required")
	}

	switch ra.AlertSourcePlatform {
	case AlertSourcePlatformK8s:
		requiredFields := map[string]string{
			"WorkloadNamespace": ra.WorkloadNamespace,
			"WorkloadKind":      ra.WorkloadKind,
			"WorkloadName":      ra.WorkloadName,
			"PodNamespace":      ra.PodNamespace,
			"PodName":           ra.PodName,
			"ContainerName":     ra.RuntimeAlertK8sDetails.ContainerName,
		}
		for fieldName, fieldValue := range requiredFields {
			if fieldValue == "" {
				return fmt.Errorf("%s is required", fieldName)
			}
		}
	case AlertSourcePlatformHost, AlertSourcePlatformCloud, AlertSourcePlatformUnknown, AlertSourcePlatformECS:
		return nil
	}

	return nil
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

	proc.MigrateToMap()

	for _, child := range proc.ChildrenMap {
		if found := findProcessRecursive(child, pid); found != nil {
			return found
		}
	}

	return nil
}
