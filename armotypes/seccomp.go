package armotypes

type SeccompStatus int

const (
	SeccompStatusUnknown            SeccompStatus = 0
	SeccompStatusMissingRuntimeInfo SeccompStatus = 1
	SeccompStatusMissing            SeccompStatus = 2
	SeccompStatusOverlyPermissive   SeccompStatus = 3
	SeccompStatusOptimized          SeccompStatus = 4
	SeccompStatusMisconfigured      SeccompStatus = 5
)

var MandatorySeccompSyscalls = []string{"epoll_wait", "tgkill", "sched_yield"}

type SeccompWorkload struct {
	Name                     string                   `json:"name"`
	Kind                     string                   `json:"kind"`
	Namespace                string                   `json:"namespace"`
	ClusterName              string                   `json:"clusterName"`
	K8sResourceHash          string                   `json:"k8sResourceHash"`
	ProfileStatus            SeccompStatus            `json:"profileStatus"`
	SyscallsUsedCount        int                      `json:"syscallsUsedCount"`
	SyscallsUnusedCount      int                      `json:"syscallsUnusedCount"`
	SyscallsUsed             []string                 `json:"syscallsUsed"`
	SyscallUnused            []string                 `json:"syscallsUnused"`
	MissingRuntimeInfoReason MissingRuntimeInfoReason `json:"missingRuntimeInfoReason"`
}
