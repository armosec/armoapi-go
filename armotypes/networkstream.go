package armotypes

import (
	"time"
)

type NetworkStreamEventProtocol string

const (
	NetworkStreamEventProtocolTCP  NetworkStreamEventProtocol = "TCP"
	NetworkStreamEventProtocolUDP  NetworkStreamEventProtocol = "UDP"
	NetworkStreamEventProtocolICMP NetworkStreamEventProtocol = "DNS"
)

type NetworkStreamEntityKind string

const (
	NetworkStreamEntityKindContainer NetworkStreamEntityKind = "container" // container
	// more types can be added here
)

type EndpointKind string

const (
	EndpointKindPod     EndpointKind = "pod"
	EndpointKindService EndpointKind = "svc"
	EndpointKindRaw     EndpointKind = "raw"
)

// NetworkStream represents a collection of network traffic events for a specific pod/container
type NetworkStream struct {
	// <identifier> to <network stream entity>
	Entities map[string]NetworkStreamEntity `json:"entities,omitempty"`
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
	Timestamp time.Time                  `json:"timestamp,omitempty"`
	IPAddress string                     `json:"ipAddress,omitempty"`
	DNSName   string                     `json:"dnsName,omitempty"`
	Port      int32                      `json:"port,omitempty"`
	Protocol  NetworkStreamEventProtocol `json:"protocol,omitempty"`
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
