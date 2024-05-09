package armotypes

// information of node-agent pod status can be taken from PodStatus table in postgres
type NodeProfile struct {
	CustomerGUID string      `json:"customerGUID"`
	Cluster      string      `json:"cluster"`
	Name         string      `json:"name"`
	PodsStatus   []PodStatus `json:"podsStatus"`
}

func (nc *NodeProfile) GetRunningPods() []PodStatus {
	var runningPods []PodStatus
	for _, pod := range nc.PodsStatus {
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
		containersCount += pod.ContainersCount
	}
	return containersCount
}
