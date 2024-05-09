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

	// Consider if this is necessary or this information should be taken by joining kubernetes resources table.
	HasApplicationProfile bool `json:"hasApplicationProfile,omitempty"`
}

type PodContainer struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
