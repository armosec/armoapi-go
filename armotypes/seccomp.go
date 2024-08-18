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

func RemoveItems(list1, list2 []string) []string {
	// Create a map to track items in list2
	itemsToRemove := make(map[string]struct{})
	for _, item := range list2 {
		itemsToRemove[item] = struct{}{}
	}

	// Create a new list for the result
	result := make([]string, 0)
	for _, item := range list1 {
		if _, exists := itemsToRemove[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// Function to append unique strings from src to dst
func AppendUniqueStrings(dst, src []string) []string {
	// Create a map to track values already in dst
	valueMap := make(map[string]bool)
	for _, v := range dst {
		valueMap[v] = true
	}

	// Append only unique values from src to dst
	for _, v := range src {
		if !valueMap[v] {
			dst = append(dst, v)
			valueMap[v] = true
		}
	}
	return dst
}
