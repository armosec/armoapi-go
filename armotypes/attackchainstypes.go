package armotypes

import (
	"reflect"
	"sort"

	"github.com/armosec/armoapi-go/identifiers"
)

type AttackChainStatus string
type ProcessingStatus string

const (
	ViewedMainScreenField = "viewedMainScreen"
	ProcessingStatusField = "processingStatus"
	//AttackChainStatuss
	StatusActive AttackChainStatus = "active"
	StatusFixed  AttackChainStatus = "fixed"
	// StatusFixedSeen AttackChainStatus = "fixedSeen"

	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusDone       ProcessingStatus = "done"
	ProcessingStatusFailed     ProcessingStatus = "failed"
	ProcessingStatusTimeout    ProcessingStatus = "timeout"
)

type AttackChain struct {
	AttackChainNodes AttackChainNode `json:"attackChainNodes,omitempty" bson:"attackChainNodes,omitempty"`
	AttackChainConfig
}

type AttackChainConfig struct {
	PortalBase       `json:",inline" bson:",inline"`
	Resource         identifiers.PortalDesignator `json:"resource,omitempty" bson:"resource,omitempty"`
	Description      string                       `json:"description,omitempty" bson:"description,omitempty"`
	CreationTime     string                       `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	AttackChainID    string                       `json:"attackChainID,omitempty" bson:"attackChainID,omitempty"` // name/cluster/resourceID
	ClusterName      string                       `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
	CustomerGUID     string                       `json:"customerGUID,omitempty" bson:"customerGUID,omitempty"`
	LatestReportGUID string                       `json:"latestReportGUID,omitempty" bson:"latestReportGUID,omitempty"` // latest reportGUID in which this attack chain was identified
	UIStatus         *AttackChainUIStatus         `json:"uiStatus,omitempty" bson:"uiStatus,omitempty"`
	Status           AttackChainStatus            `json:"status,omitempty" bson:"status,omitempty"` // "active"/ "fixed"
}

type AttackChainNode struct {
	Name                           string             `json:"name" bson:"name,omitempty"`
	Description                    string             `json:"description" bson:"description,omitempty"`
	ControlIDs                     []string           `json:"controlIDs,omitempty" bson:"controlIDs,omitempty"` // failed/ignored controls that are associated to this attack chain node
	Vulnerabilities                []Vulnerabilities  `json:"vulnerabilities,omitempty" bson:"vulnerabilities,omitempty"`
	RelatedResources               []RelatedResources `json:"relatedResources" bson:"relatedResources,omitempty"`
	NextNodes                      []AttackChainNode  `json:"nextNodes,omitempty" bson:"nextNodes,omitempty"`
	FlattenRelatedResourcesDisplay bool               `json:"flattenRelatedResourcesDisplay,omitempty" bson:"flattenRelatedResourcesDisplay,omitempty"`
}

type RelatedResources struct {
	identifiers.PortalDesignator `json:",inline" bson:",inline"`
	Clickable                    bool               `json:"clickable,omitempty" bson:"clickable,omitempty"`
	EdgeText                     []string           `json:"edgeText,omitempty" bson:"edgeText,omitempty"`
	RelatedResources             []RelatedResources `json:"relatedResources" bson:"relatedResources,omitempty"`
}

type Vulnerabilities struct {
	ContainerName string   `json:"containerName" bson:"containerName,omitempty"`
	ImageScanID   string   `json:"imageScanID" bson:"imageScanID,omitempty"`
	Names         []string `json:"names" bson:"names,omitempty"` // CVE names
}

// struct for UI support. All strings are timestamps
type AttackChainUIStatus struct {
	// fields updated by the BE
	FirstSeen string `json:"firstSeen,omitempty" bson:"firstSeen,omitempty"` // timestamp of first scan in which the attack chain was identified
	// fields updated by the UI
	ViewedMainScreen string `json:"viewedMainScreen,omitempty" bson:"viewedMainScreen,omitempty"` // if the attack chain was viewed by the user// New badge
	ProcessingStatus string `json:"processingStatus,omitempty" bson:"processingStatus,omitempty"` // "processing"/ "done"
}

func (a *AttackChainNode) Equals(b *AttackChainNode) bool {
	// Sort string slices
	sort.Strings(a.ControlIDs)
	sort.Strings(b.ControlIDs)

	// Sort Vulnerabilities slice (assuming Vulnerabilities has a defined order)
	sort.Slice(a.Vulnerabilities, func(i, j int) bool {
		// Provide logic for sorting Vulnerabilities here, e.g.,
		return a.Vulnerabilities[i].ContainerName < a.Vulnerabilities[j].ContainerName
	})
	sort.Slice(b.Vulnerabilities, func(i, j int) bool {
		// Provide logic for sorting Vulnerabilities here
		return b.Vulnerabilities[i].ContainerName < b.Vulnerabilities[j].ContainerName
	})

	if a.Description != b.Description || a.Name != b.Name {
		return false
	}

	// Recursively sort and compare NextNodes
	for i := range a.NextNodes {
		if !a.NextNodes[i].Equals(&b.NextNodes[i]) {
			return false
		}
	}

	return reflect.DeepEqual(a, b)
}
