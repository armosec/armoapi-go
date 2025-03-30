package armotypes

import (
	"time"
)

type NetworkTrafficEventProtocol string

const (
	NetworkTrafficEventProtocolTCP NetworkTrafficEventProtocol = "TCP"
	NetworkTrafficEventProtocolUDP NetworkTrafficEventProtocol = "UDP"
)

type EndpointKind string

const (
	EndpointKindPod     EndpointKind = "pod"
	EndpointKindService EndpointKind = "svc"
	EndpointKindRaw     EndpointKind = "raw"
)

// NetworkTrafficEvents represents a collection of network traffic events for a specific pod/container
type NetworkTrafficEvents struct {
	// ContainerID to NetworkTrafficEventContainer
	Containers map[string]NetworkTrafficEventContainer `json:"containers,omitempty"`
}

// NetworkTrafficEventContainer represents an aggregation of network connections from/to a specific source
type NetworkTrafficEventContainer struct {
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
	// inbound network traffic events
	Inbound map[string]NetworkTrafficEvent `json:"inbound,omitempty"`
	// outbound network traffic events
	Outbound map[string]NetworkTrafficEvent `json:"outbound,omitempty"`
}

// NetworkTrafficEvent represents an aggregation of network connections from/to a specific source
type NetworkTrafficEvent struct {
	Timestamp time.Time                   `json:"timestamp,omitempty"`
	IPAddress string                      `json:"ipAddress,omitempty"`
	DNSName   string                      `json:"dnsName,omitempty"`
	Port      int32                       `json:"port,omitempty"`
	Protocol  NetworkTrafficEventProtocol `json:"protocol,omitempty"`
	// endpoint kind (pod, service, raw)
	Kind EndpointKind `json:"kind,omitempty"`
	// endpoint details in case of pod
	NetworkTrafficEventEndpointPodDetails `json:",inline"`
	// endpoint details in case of service
	NetworkTrafficEventEndpointServiceDetails `json:",inline"`
}

func (e *NetworkTrafficEvent) String() string {
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

type NetworkTrafficEventEndpointPodDetails struct {
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

type NetworkTrafficEventEndpointServiceDetails struct {
	ServiceName      string `json:"serviceName,omitempty"`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
}
