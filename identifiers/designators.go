package identifiers

import (
	"fmt"
	"strings"
)

// Type of the designator
//
// swagger:enum DesignatorType
type DesignatorType string

func (dt DesignatorType) ToLower() DesignatorType {
	return DesignatorType(strings.ToLower(string(dt)))
}

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

// attributes
const (
	DesignatorsToken          = "designators"
	AttributeCustomerGUID     = "customerGUID"
	AttributeRegistryName     = "registryName"
	AttributeRepository       = "repository"
	AttributeTag              = "tag"
	AttributeCluster          = "cluster"
	AttributeClusterShortName = "clusterShortName"
	AttributeNamespace        = "namespace"
	AttributeKind             = "kind"
	AttributeName             = "name"
	AttributeContainerName    = "containerName"
	AttributeApiVersion       = "apiVersion"
	AttributeApiGroup         = "apiGroup"
	AttributeWorkloadHash     = "workloadHash"
	AttributeIsIncomplete     = "isIncomplete"
	AttributeSensor           = "sensor"
	AttributePath             = "path"
	AttributeResourceID       = "resourceID"
	AttributeContainerScanId  = "containerScanId"
	AttributeSyncKind         = "syncKind"
	AttributeSBOMToolVersion  = "sbomToolVersion"
	AttributeSecurityRiskID   = "securityRiskID"
	AttributeK8sResourceHash  = "k8sResourceHash"
	AttributeType             = "type"
	AttributeOwner            = "owner"
	AttributeRelated          = "relatedObjects"
	AttributeLayerHash        = "layerHash"
	AttributeImageRepository  = "imageRepository"
	AttributeResourceHash     = "resourceHash"
	AttributeComponentVersion = "componentVersion"
	AttributeComponent        = "component"
	AttributeSeverityScore    = "severityScore"
	AttributeSeverity         = "severity"
	AttributeCVEID            = "cveID"
	AttributeCVEName          = "cveName"
	AttributeControlID        = "controlID"
	AttributeBaseScore        = "baseScore"
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

// Worker nodes attribute related consts
const (
	AttributeWorkerNodes              = "workerNodes"
	WorkerNodesmax                    = "max"
	WorkerNodeslastReported           = "lastReported"
	WorkerNodeslastReportDate         = "lastReportDate"
	WorkerNodesmaxPerMonth            = "maxPerMonth"
	WorkerNodesmaxReportGUID          = "maxReportGUID"
	WorkerNodesmaxPerMonthReportGUID  = "maxPerMonthReportGUID"
	WorkerNodeslastReportGUID         = "lastReportGUID"
	LastPostureScanTriggered          = "lastPostureScanTriggered"
	LastTimeACEngineCompleted         = "lastTimeACEngineCompleted"
	LastTimeSecurityRiskScanCompleted = "lastTimeSecurityRiskScanCompleted"
)

type mapString2String map[string]string

var IgnoreLabels = []string{AttributeCluster, AttributeNamespace}

// AttributeDesignators describe a kubernetes object, with its labels.
type AttributesDesignators struct {
	cluster         string
	namespace       string
	kind            string
	name            string
	path            string
	labels          map[string]string
	resourceID      string
	k8sResourceHash string
}

func CalcResourceHash(customerGUID string, identifiers map[string]string) string {
	hash := (fmt.Sprintf("%s/%s/%s/%s/%s",
		customerGUID,
		strings.ToLower(identifiers[AttributeKind]),
		strings.ToLower(identifiers[AttributeName]),
		strings.ToLower(identifiers[AttributeNamespace]),
		strings.ToLower(identifiers[AttributeCluster])))

	return hash
}
