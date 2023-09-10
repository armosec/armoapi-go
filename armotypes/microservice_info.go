package armotypes

import (
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CAContainerMetrics holds data of single container which runs in multiple pods
type CAContainerMetrics struct {
	core.Container    `json:",inline"`
	CAIntegrityStatus int `json:"caIntegrityStatus"`
}

// MicroserviceExtraDetails represent an overview of microservice, services, container data
// and cloud data
type MicroserviceExtraDetails struct {
	CAMicroserviceOverviewMetadata `json:",inline"`
	NumOfContainers                int                `json:"NumOfContainers"`
	Labels                         map[string]string  `json:"labels,omitempty"`
	Annotations                    map[string]string  `json:"annotations,omitempty"`
	ContainersSummary              []ContainerSummary `json:"containers"`
	ExternalFacing                 bool               `json:"isExternalFacingMS"`
}

// ContainerSummary - a must have summarized info of containers
type ContainerSummary struct {
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	IsPrivileged bool    `json:"root"`
	Probes       []Probe `json:"probes,omitempty"`
	Limitations  `json:"limitations,omitempty"`
}

// Probe - represent the various container probes
type Probe struct {
	Type string `json:"type"` // e,g liveness/readiness/<w.e>
	Data string `json:"data"` // actual probe data/settings
}

// Limitations - container defined limitations
type Limitations struct {
	CPU    int64 `json:"cpu,omitempty"`
	Memory int64 `json:"memory,omitempty"`
	Disk   int64 `json:"disk,omitempty"`
}

// CAMicroserviceOverview represnets it's name
type CAMicroserviceOverview struct {
	CAMicroserviceOverviewMetadata `json:",inline"`
}

// CAMicroserviceOverviewMetadata represnets it's name
type CAMicroserviceOverviewMetadata struct {
	CAK8SMeta     `json:",inline"`
	WLID          string   `json:"wlid"`
	Datacenter    string   `json:"datacenter,omitempty"`
	OVNamespace   string   `json:"namespace,omitempty"`
	Project       string   `json:"project,omitempty"`
	Orchestrator  string   `json:"orchestrator"`
	Kind          string   `json:"kind"`
	OperationType string   `json:"operationType"`
	OVName        string   `json:"name"`
	Categories    []string `json:"categories"`
	DisplayName   string   `json:"displayName,omitempty"`
	CloudProvider string   `json:"cloudProvider"`
}

// CAK8SMeta holds common metadata about k8s objects
type CAK8SMeta struct {
	CustomerGUID   string    `json:"customerGUID"`
	CAClusterName  string    `json:"caClusterName,omitempty"`
	LastUpdateTime time.Time `json:"caLastUpdate"`
	IsActive       bool      `json:"isActive"`
}

// MicroserviceMetadataView represent the model to return in metadata request
type MicroserviceMetadataView struct {
	CAMicroserviceOverviewMetadata
	metav1.ObjectMeta `json:"metadata"`
	Ancestor          K8SAncestor       `json:"uptreeOwner,omitempty"`
	UsageType         string            `json:"usageType,omitempty"`
	Categories        map[string]bool   `json:"categories"`
	CALabels          map[string]string `json:"caLabels"`
}

// MicroserviceInfo single microservice with CA metrics
type MicroserviceInfo struct {
	MicroserviceMetadataView `json:",inline"`
	PodSpecID                int64 `json:"podSpecId"` // will be sent from the cluster-agent to reconize this spec
	core.PodSpec             `json:"spec"`
	core.PodStatus           `json:"status" yaml:"status"`
	Containers               []CAContainerMetrics `json:"containers,omitempty"`
	K8SPodObjects            []K8SPodObject       `json:"k8sPodObjects,omitempty"`
	CAStartTime              time.Time            `json:"caStartTime"`
}

// K8SAncestor represents the kind of the microservice inside the k8s cluster
type K8SAncestor struct {
	Name           string      `json:"name"`
	Kind           string      `json:"kind"`
	FullDeclaraion interface{} `json:"ownerData,omitempty"`
}

// K8SPodObject represents actuall pod which run on particular node of the cluster
type K8SPodObject struct {
	CAK8SMeta         `json:",inline"`
	Name              string      `json:"podName"`
	CreatedAt         time.Time   `json:"startedAt,omitempty"`
	TerminatedAt      *time.Time  `json:"terminatedAt,omitempty"`
	PodIP             string      `json:"podIP"`
	NodeName          string      `json:"nodeName"`
	Namespace         string      `json:"namespace"`
	NominatedNodeName string      `json:"nominatedNodeName"`
	Ancestor          K8SAncestor `json:"uptreeOwner,omitempty"`
	PodSpecID         int64       `json:"podSpecId"`
	PodStatus         string      `json:"podStatus"`
	// ContainerStatuses  map[string]string `json:"containerStatuses"`
}

// ContainersStatusData holds the status of containers in runtime. This including the docker image tag + image hash
type ContainersStatusData map[string]map[string]string

// K8SNamespace represents single k8s namespace in cluster
type K8SNamespace struct {
	CAK8SMeta      `json:",inline"`
	Name           string `json:"name"`
	core.Namespace `json:",inline"`
}
