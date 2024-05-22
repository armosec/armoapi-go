package armotypes

import (
	"slices"
	"time"
)

type PodStatus struct {
	CustomerGUID               string         `json:"customerGUID"`
	Cluster                    string         `json:"cluster"`
	ResourceHash               string         `json:"resourceHash"`
	ResourceVersion            string         `json:"resourceVersion"`
	Name                       string         `json:"name"`
	Namespace                  string         `json:"namespace"`
	NodeName                   string         `json:"nodeName"`
	App                        string         `json:"app"`
	Phase                      string         `json:"phase"`
	CurrentState               string         `json:"currentState"`
	LastStateReason            string         `json:"lastStateReason"`
	LastStateMessage           string         `json:"lastStateMessage"`
	LastStateTransitionTime    time.Time      `json:"lastStateTransitionTime"`
	CreationTimestamp          time.Time      `json:"creationTimestamp"`
	Containers                 []PodContainer `json:"containers,omitempty"`
	InitContainers             []PodContainer `json:"initContainers,omitempty"`
	EphemeralContainers        []PodContainer `json:"ephemeralContainers,omitempty"`
	HasFinalApplicationProfile bool           `json:"hasFinalApplicationProfile"`
	HasApplicableRuleBindings  bool           `json:"hasApplicableRuleBindings"`
	HasRelevancyCalculating    bool           `json:"hasRelevancyCalculating"`
	IsKDRMonitored             bool           `json:"isKDRMonitored"`
}

type PodContainer struct {
	Name                string    `json:"name"`
	Image               string    `json:"image"`
	IsKDRMonitored      bool      `json:"isKDRMonitored"`
	CurrentState        string    `json:"currentState"`
	LastStateExitCode   int       `json:"lastStateExitCode"`
	LastStateFinishedAt time.Time `json:"lastStateFinishedAt"`
	LastStateStartedAt  time.Time `json:"lastStateStartedAt"`
	RestartCount        int       `json:"restartCount"`
}

func (ps *PodStatus) GetMonitoredContainers() []PodContainer {
	var monitoredContainers []PodContainer
	if ps.IsKDRMonitored {
		for _, container := range slices.Concat(ps.Containers, ps.InitContainers, ps.EphemeralContainers) {
			if container.IsKDRMonitored {
				monitoredContainers = append(monitoredContainers, container)
			}
		}
	}
	return monitoredContainers
}
