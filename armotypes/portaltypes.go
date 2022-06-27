package armotypes

import "strings"

const (
	CustomerGuidQuery   = "customerGUID"
	ClusterNameQuery    = "cluster"
	DatacenterNameQuery = "datacenter"
	NamespaceQuery      = "namespace"
	ProjectQuery        = "project"
	WlidQuery           = "wlid"
	SidQuery            = "sid"
)

// PortalBase holds basic items data from portal BE
type PortalBase struct {
	GUID       string                 `json:"guid"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes,omitempty"` // could be string
}

type DesignatorType string

// Supported designators
const (
	DesignatorAttributes DesignatorType = "Attributes"
	DesignatorAttribute  DesignatorType = "Attribute" // Deprecated
	/*
		WorkloadID format.
		k8s format: wlid://cluster-<cluster>/namespace-<namespace>/<kind>-<name>
		native format: wlid://datacenter-<datacenter>/project-<project>/native-<name>
	*/
	DesignatorWlid DesignatorType = "Wlid"
	/*
		Wild card - subset of wlid. e.g.
		1. Include cluster:
			wlid://cluster-<cluster>/
		2. Include cluster and namespace (filter out all other namespaces):
			wlid://cluster-<cluster>/namespace-<namespace>/
	*/
	DesignatorWildWlid      DesignatorType = "WildWlid"
	DesignatorWlidContainer DesignatorType = "WlidContainer"
	DesignatorWlidProcess   DesignatorType = "WlidProcess"
	DesignatorSid           DesignatorType = "Sid" // secret id
)

func (dt DesignatorType) ToLower() DesignatorType {
	return DesignatorType(strings.ToLower(string(dt)))
}

// attributes
const (
	DesignatorsToken       = "designators"
	AttributeCustomerGUID  = "customerGUID"
	AttributeRegistryName  = "registryName"
	AttributeRepository    = "repository"
	AttributeTag           = "tag"
	AttributeCluster       = "cluster"
	AttributeNamespace     = "namespace"
	AttributeKind          = "kind"
	AttributeName          = "name"
	AttributeContainerName = "containerName"
	AttributeApiVersion    = "apiVersion"
	AttributeWorkloadHash  = "workloadHash"
	AttributeIsIncomplete  = "isIncomplete"
)

// Repository scan related attributes
const (
	AttributeRepoName      = "repoName"
	AttributeRepoOwner     = "repoOwner"
	AttributeRepoHash      = "repoHash"
	AttributeBranchName    = "branch"
	AttributeDefaultBranch = "defaultBranch"
	AttributeProvider      = "provider"
	AttributeRemoteURL     = "remoteURL"

	AttributeLastCommitHash     = "lastCommitHash"
	AttributeLastCommitterName  = "lastCommitterName"
	AttributeLastCommitterEmail = "lastCommitterEmail"
	AttributeLastCommitTime     = "lastCommitTime"

	AttributeFilePath = "filePath"
	AttributeFileType = "fileType"
	AttributeFileDir  = "fileDirectory"

	AttributeLastFileCommitHash     = "lastFileCommitHash"
	AttributeLastFileCommitterName  = "lastFileCommitterName"
	AttributeLastFileCommitterEmail = "LastFileCommitterEmail"
	AttributeLastFileCommitTime     = "lastFileCommitTime"
)

// PortalDesignator represented single designation options
type PortalDesignator struct {
	DesignatorType DesignatorType    `json:"designatorType"`
	WLID           string            `json:"wlid,omitempty"`
	WildWLID       string            `json:"wildwlid,omitempty"`
	SID            string            `json:"sid,omitempty"`
	Attributes     map[string]string `json:"attributes"`
}

// Worker nodes attribute related consts
const (
	AttributeWorkerNodes             = "workerNodes"
	WorkerNodesmax                   = "max"
	WorkerNodeslastReported          = "lastReported"
	WorkerNodeslastReportDate        = "lastReportDate"
	WorkerNodesmaxPerMonth           = "maxPerMonth"
	WorkerNodesmaxReportGUID         = "maxReportGUID"
	WorkerNodesmaxPerMonthReportGUID = "maxPerMonthReportGUID"
	WorkerNodeslastReportGUID        = "lastReportGUID"
)
