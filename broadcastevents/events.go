package broadcastevents

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

type EventBase struct {
	CustomerGUID string    `json:"customerGUID"`
	EventName    string    `json:"eventName"`
	EventTime    time.Time `json:"eventTime"`
	// The date (year,month,day of month) of the event
	EventDate          string `json:"eventDate,omitempty"`
	EventMonth         string `json:"eventMonth,omitempty"`
	EventWeekOfTheYear int    `json:"eventWeekOfTheYear,omitempty"`
}

type NetworkPolicyGenerated struct {
	EventBase     `json:",inline"`
	ClusterName   string `json:"clusterName"`
	WorkloadsAmount          string `json:"WorkloadsAmount"`
}

type AttackChainCreated struct {
	EventBase     `json:",inline"`
	ClusterName   string `json:"clusterName"`
	ACId          string `json:"ACId"`
	ACType        string `json:"ACType"`
	ACFirstSeeing string `json:"ACFirstSeeing"`
}

type AttackChainResolved struct {
	EventBase     `json:",inline"`
	ClusterName   string `json:"clusterName"`
	ACId          string `json:"ACId"`
	ACType        string `json:"ACType"`
	ACFirstSeeing string `json:"ACFirstSeeing"`
}

type AggregationEvent struct {
	EventBase        `json:",inline"`
	JobID            string `json:"jobID"`
	K8sVendor        string `json:"k8sVendor,omitempty"`
	K8sVersion       string `json:"k8sVersion,omitempty"`
	KSVersion        string `json:"kubescapeVersion,omitempty"`
	HelmChartVersion string `json:"helmChartVersion,omitempty"`
	ReportGUID       string `json:"reportGUID,omitempty"`
	ClusterName      string `json:"clusterName,omitempty"`
	WorkerNodesCount int    `json:"workerNodesCount,omitempty"`
}

type LoginEvent struct {
	EventBase         `json:",inline"`
	Email             string `json:"email"`
	UserName          string `json:"userName"`
	PreferredUserName string `json:"preferredUserName"`
}

type HelmInstalledEvent struct {
	EventBase
	ClusterName string `json:"clusterName"`
	*armotypes.InstallationData
}

type PodInTroubleEvent struct {
	EventBase
	ClusterName string `json:"clusterName"`
	ObjId       string `json:"podId"`
	Reason      string `json:"reason"`
	Message     string `json:"message"`
}

type PodInTroubleConditionEvent struct {
	PodInTroubleEvent
	Condition string `json:"condition"`
}

type PodInTroubleContainerEvent struct {
	PodInTroubleEvent
	ContainerName string `json:"containerName"`
	ExitCode      int32  `json:"exitCode"`
	RestartCount  int32  `json:"restartCount"`
}

type IgnoreRuleEvent struct {
	EventBase      `json:",inline"`
	IgnoreRuleType IgnoreRuleType           `json:"ignoreRuleType"`
	IgnoredIds     string                   `json:"ids"` // comma separated ids of controls or vulnerabilities
	Resources      int                      `json:"resources"`
	ExpirationType IgnoreRuleExpirationType `json:"expirationType"`
}

type AlertChannelEvent struct {
	EventBase        `json:",inline"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	AllClusters      *bool  `json:"allClusters,omitempty"`
	NewVulnerability string `json:"new vulnerability is detected in your cluster,omitempty"`
	NewFix           string `json:"new fix is available for vulnerability,omitempty"`
	Compliance       string `json:"compliance score has decreased,omitempty"`
	NewAdmin         string `json:"new cluster admin was added,omitempty"`
}

type ScanWithoutAccessKeyEvent struct {
	EventBase   `json:",inline"`
	ClusterName string `json:"clusterName"`
}
