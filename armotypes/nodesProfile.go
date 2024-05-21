package armotypes

// information of node-agent pod status can be taken from PodStatus table in postgres
type NodeProfile struct {
	PodStatuses []PodStatus `json:"podStatuses"`

	CurrentState string `json:"currentState"`

	NodeAgentRunning bool `json:"nodeAgentRunning"`

	RuntimeDetectionEnabled bool `json:"runtimeDetectionEnabled"`
}

type NodeStatus struct {
	CustomerGUID    string `json:"customerGUID"`
	Cluster         string `json:"cluster"`
	Name            string `json:"name"`
	K8sResourceHash string `json:"k8sResourceHash"`
	NodeProfile     `json:",inline"`
}

func (nc *NodeStatus) IsKDRMonitored() bool {
	return nc.NodeAgentRunning && nc.RuntimeDetectionEnabled
}

func (nc *NodeStatus) GetMonitoredNamespaces() []string {
	// Map to keep track of unique namespaces
	monitoredNamespaceMap := make(map[string]bool)
	var monitoredNamespaces []string

	if !nc.IsKDRMonitored() {
		return monitoredNamespaces
	}

	for _, pod := range nc.PodStatuses {
		if pod.IsKDRMonitored && pod.Phase == "Running" && !monitoredNamespaceMap[pod.Namespace] {
			// Add the namespace to the slice and mark it in the map
			monitoredNamespaces = append(monitoredNamespaces, pod.Namespace)
			monitoredNamespaceMap[pod.Namespace] = true
		}
	}
	return monitoredNamespaces
}

func (nc *NodeStatus) GetMonitoredPods() []PodStatus {
	var monitoredPods []PodStatus
	if !nc.IsKDRMonitored() {
		return monitoredPods

	}

	for _, pod := range nc.PodStatuses {
		if pod.IsKDRMonitored && pod.Phase == "Running" {
			monitoredPods = append(monitoredPods, pod)
		}
	}

	return monitoredPods
}

func (nc *NodeStatus) GetMonitoredContainers() map[string][]PodContainer {
	monitoredContainers := make(map[string][]PodContainer)
	if !nc.IsKDRMonitored() {
		return monitoredContainers
	}
	monitoredPods := nc.GetMonitoredPods()
	for _, pod := range monitoredPods {
		monitoredContainers[pod.Name] = pod.GetMonitoredContainers()
	}
	return monitoredContainers
}

func (nc *NodeStatus) CountMonitoredNamespaces() int {
	return len(nc.GetMonitoredNamespaces())
}

func (nc *NodeStatus) CountMonitoredPods() int {
	return len(nc.GetMonitoredPods())
}

func (nc *NodeStatus) GetRunningPods() []PodStatus {
	var runningPods []PodStatus
	for _, pod := range nc.PodStatuses {
		if pod.Phase == "Running" {
			runningPods = append(runningPods, pod)
		}
	}
	return runningPods
}

func (nc *NodeStatus) CountRunningPods() int {
	return len(nc.GetRunningPods())
}

func (nc *NodeStatus) CountRunningPodsContainers() int {
	var containersCount int
	runningPods := nc.GetRunningPods()
	for _, pod := range runningPods {
		containersCount += len(pod.Containers) + len(pod.InitContainers)
	}
	return containersCount
}

func (nc *NodeStatus) CountMonitoredContainers() int {
	count := 0
	monitoredContainers := nc.GetMonitoredContainers()

	for _, containers := range monitoredContainers {
		count += len(containers)
	}

	return count
}
