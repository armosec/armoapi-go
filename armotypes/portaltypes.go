package armotypes

import (
	"strings"
)

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
	GUID        string                 `json:"guid" bson:"guid"`
	Name        string                 `json:"name" bson:"name"`
	Attributes  map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"` // could be string
	UpdatedTime string                 `json:"updatedTime,omitempty" bson:"updatedTime,omitempty"`
}

// Type of the designator
//
// swagger:enum DesignatorType
type DesignatorType string

// Supported designators
const (
	DesignatorAttributes DesignatorType = "Attributes"
	DesignatorAttribute  DesignatorType = "Attribute" // Deprecated
	// WorkloadID format.
	//
	// Has two formats:
	//  1. Kubernetes format: wlid://cluster-<cluster>/namespace-<namespace>/<kind>-<name>
	//  2. Native format: wlid://datacenter-<datacenter>/project-<project>/native-<name>
	DesignatorWlid DesignatorType = "Wlid"
	// A WorkloadID wildcard expression.
	//
	// A wildcard expression that includes a cluster:
	//
	//  wlid://cluster-<cluster>/
	//
	// An expression that includes a cluster and namespace (filters out all other namespaces):
	//
	//  wlid://cluster-<cluster>/namespace-<namespace>/
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
	AttributeSensor        = "sensor"
	AttributePath          = "path"
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

	AttributeFilePath          = "filePath"
	AttributeFileType          = "fileType"
	AttributeFileDir           = "fileDirectory"
	AttributeFileUrl           = "fileUrl"
	AttributeFileHelmChartName = "fileHelmChartName"

	AttributeLastFileCommitHash     = "lastFileCommitHash"
	AttributeLastFileCommitterName  = "lastFileCommitterName"
	AttributeLastFileCommitterEmail = "LastFileCommitterEmail"
	AttributeLastFileCommitTime     = "lastFileCommitTime"

	AttributeUseHTTP       = "useHTTP"
	AttributeSkipTLSVerify = "skipTLSVerify"
)

// rego-library attributes
const (
	AttributeImageScanRelated     = "imageScanRelated"
	AttributeImageRelatedControls = "imageRelatedControls"
	AttributeHostSensorRule       = "hostSensorRule"
	AttributeHostSensor           = "hostSensor"
)

// PortalDesignator represents a single designation option
type PortalDesignator struct {
	DesignatorType DesignatorType `json:"designatorType" bson:"designatorType"`
	// A specific Workload ID
	WLID string `json:"wlid,omitempty" bson:"wlid,omitempty"`
	// An expression that describes applicable workload IDs
	WildWLID string `json:"wildwlid,omitempty" bson:"wildwlid,omitempty"`
	// A specific Secret ID
	SID string `json:"sid,omitempty" bson:"sid,omitempty"`
	// Attributes that describe the targets
	Attributes map[string]string `json:"attributes" bson:"attributes"`
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

// PortalCluster holds cluster data from portal BE
type PortalCluster struct {
	PortalBase       `json:",inline" bson:"inline"`
	SubscriptionDate string `json:"subscription_date" bson:"subscription_date"`
	LastLoginDate    string `json:"last_login_date" bson:"last_login_date"`
}

type PortalCustomer struct {
	PortalBase       `json:",inline" bson:"inline"`
	Description      string `json:"description" bson:"description"`
	SubscriptionDate string `json:"subscription_date" bson:"subscription_date"`
	LastLoginDate    string `json:"last_login_date" bson:"last_login_date"`
	Email            string `json:"email" bson:"email"`
	//License
	LicenseType            string               `json:"license_type" bson:"license_type"`
	SubscriptionExpiration string               `json:"subscription_expiration" bson:"subscription_expiration"`
	InitialLicenseType     string               `json:"initial_license_type" bson:"initial_license_type"`
	NotificationsConfig    *NotificationsConfig `json:"notifications_config,omitempty" bson:"notifications_config,omitempty"`
}

type PortalRepository struct {
	PortalBase   `json:",inline" bson:"inline"`
	CreationDate string `json:"creationDate" bson:"creationDate"`
	Provider     string `json:"provider" bson:"provider"`
	Owner        string `json:"owner" bson:"owner"`
	RepoName     string `json:"repoName" bson:"repoName"`
	BranchName   string `json:"branchName" bson:"branchName"`
}

type PortalRegistryCronJob struct {
	PortalBase      `json:",inline" bson:"inline"`
	RegistryInfo    `json:",inline" bson:"inline"`
	CreationDate    string       `json:"creationDate" bson:"creationDate"`
	ID              string       `json:"id" bson:"id"`
	ClusterName     string       `json:"clusterName" bson:"clusterName"`
	CronTabSchedule string       `json:"cronTabSchedule" bson:"cronTabSchedule"`
	Repositories    []Repository `json:"repositories" bson:"repositories"`
}
