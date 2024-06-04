package armotypes

type SeccompStatus int

const (
	SeccompStatusMissingRuntimeInfo SeccompStatus = 1
	SeccompStatusMissing            SeccompStatus = 2
	SeccompStatusOverlyPermissive   SeccompStatus = 3
	SeccompStatusOptimized          SeccompStatus = 4
)

type SeccompWorkload struct {
	Name                string        `json:"name"`
	Kind                string        `json:"kind"`
	Namespace           string        `json:"namespace"`
	ClusterName         string        `json:"clusterName"`
	ProfileStatus       SeccompStatus `json:"profileStatus"`
	SyscallsUsedCount   int           `json:"syscallsUsedCount"`
	SyscallsUnusedCount int           `json:"syscallsUnusedCount"`
	SyscallsUsed        []string      `json:"syscallsUsed"`
	SyscallUnused       []string      `json:"syscallsUnused"`
}
