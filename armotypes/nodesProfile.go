package armotypes

// information of node-agent pod status can be taken from PodStatus table in postgres
type NodeProfile struct {
	CustomerGUID string      `json:"customerGUID"`
	Cluster      string      `json:"cluster"`
	Name         string      `json:"name"`
	PodStatuses  []PodStatus `json:"podStatuses"`

	// consider join on PodStatus table to get missingRuntimeInfo
	NodeAgentRunning bool `json:"nodeAgentRunning"`

	// from configMap?
	RuntimeDetectionEnabled bool `json:"runtimeDetectionEnabled"`
}

func (nc *NodeProfile) GetMonitoredPods() []PodStatus {
	var monitoredPods []PodStatus
	if nc.NodeAgentRunning && nc.RuntimeDetectionEnabled {
		for _, pod := range nc.PodStatuses {
			if pod.IsKDRMonitored {
				monitoredPods = append(monitoredPods, pod)
			}
		}
	}

	return monitoredPods
}

func (nc *NodeProfile) GetMonitoredContainers() map[string][]PodContainer {
	monitoredContainers := make(map[string][]PodContainer)
	monitoredPods := nc.GetMonitoredPods()
	for _, pod := range monitoredPods {
		monitoredContainers[pod.Name] = pod.GetMonitoredContainers()
	}
	return monitoredContainers
}

func (nc *NodeProfile) CountMonitoredPods() int {
	return len(nc.GetMonitoredPods())
}

func (nc *NodeProfile) GetRunningPods() []PodStatus {
	var runningPods []PodStatus
	for _, pod := range nc.PodStatuses {
		if pod.Phase == "Running" {
			runningPods = append(runningPods, pod)
		}
	}
	return runningPods
}

func (nc *NodeProfile) CountRunningPods() int {
	return len(nc.GetRunningPods())
}

func (nc *NodeProfile) CountRunningPodsContainers() int {
	var containersCount int
	runningPods := nc.GetRunningPods()
	for _, pod := range runningPods {
		containersCount += len(pod.Containers) + len(pod.InitContainers)
	}
	return containersCount
}

func (nc *NodeProfile) CountMonitoredContainers() int {
	return len(nc.GetMonitoredContainers())
}
