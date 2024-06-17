package armotypes

import (
	"time"
)

type ContainerType string

const (
	InitContainer      ContainerType = "initcontainer"
	Container          ContainerType = "container"
	EphemeralContainer ContainerType = "ephemeralcontainer"
)

type ContainerStatus struct {
	CustomerGUID string `json:"customerGUID"`
	ClusterName  string `json:"clusterName"`

	ResourceHash  string        `json:"resourceHash"`
	Name          string        `json:"name"`          // container name
	ContainerType ContainerType `json:"containerType"` // initcontainer, container, ephemeralcontainer

	Architectures []string `json:"architectures"` // architectures of the container
	WorkloadName  string   `json:"workloadName"`  // name of the workload
	Kind          string   `json:"kind"`          // kind of the workload
	Namespace     string   `json:"namespace"`     // namespace of the workload

	// seccomp related fields (coming from ApplicationProfile)
	// IsSeccompConfiguredWorkloadLevel  *bool    `json:"isSeccompConfiguredWorkloadLevel"` // if nil, seccomp is not configured
	IsSeccompConfiguredDefaultRuntime *bool    `json:"isSeccompConfiguredDefaultRuntime"` // if nil, seccomp is not configured
	SeccompConfiguredLocalhostProfile string   `json:"seccompConfiguredLocalhostProfile"`
	SeccompConfiguredSyscalls         []string `json:"seccompConfiguredSyscalls"`
	SeccompConfiguredArchitectures    []string `json:"seccompConfiguredArchitectures"`
	SyscallsUsed                      []string `json:"syscallsUsed"`

	ApplicationProfileLastUpdated  *time.Time `json:"applicationProfileLastUpdated"`  // last updated time of applicationProfile
	ApplicationProfileResourceHash string     `json:"applicationProfileResourceHash"` // resource hash of applicationProfile

}
