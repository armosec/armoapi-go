package armotypes

import "time"

type PodStatus struct {
	CustomerGUID        string         `json:"customerGUID"`
	Cluster             string         `json:"cluster"`
	ResourceHash        string         `json:"resourceHash"`
	ResourceVersion     string         `json:"resourceVersion"`
	Name                string         `json:"name"`
	Namespace           string         `json:"namespace"`
	NodeName            string         `json:"nodeName"`
	App                 string         `json:"app"`
	Phase               string         `json:"phase"`
	CurrentState        string         `json:"currentState"`
	LastStateExitCode   int            `json:"lastStateExitCode"`
	LastStateFinishedAt time.Time      `json:"lastStateFinishedAt"`
	LastStateStartedAt  time.Time      `json:"lastStateStartedAt"`
	LastStateReason     string         `json:"lastStateReason"`
	LastStateMessage    string         `json:"lastStateMessage"`
	RestartCount        int            `json:"restartCount"`
	CreationTimestamp   time.Time      `json:"creationTimestamp"`
	Containers          []PodContainer `json:"podContainers,omitempty"`
	InitContainers      []PodContainer `json:"initContainers,omitempty"`

	HasFinalApplicationProfile bool `json:"hasFinalApplicationProfile"`

	HasApplicableRuleBindings bool `json:"hasApplicableRuleBindings"`

	HasRelevancyCalculating bool `json:"hasRelevancyCalculating"`
}

type PodContainer struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	IsMonitored bool   `json:"isMonitored"`
}

func (ps *PodStatus) IsMonitored() bool {
	return ps.HasFinalApplicationProfile && ps.Phase == "Running" && ps.HasApplicableRuleBindings && ps.HasRelevancyCalculating
}

func (ps *PodStatus) GetMonitoredContainers() []PodContainer {
	var monitoredContainers []PodContainer
	if ps.IsMonitored() {
		for _, container := range ps.Containers {
			if container.IsMonitored {
				monitoredContainers = append(monitoredContainers, container)
			}
		}

		for _, container := range ps.InitContainers {
			if container.IsMonitored {
				monitoredContainers = append(monitoredContainers, container)
			}
		}

	}
	return monitoredContainers
}
