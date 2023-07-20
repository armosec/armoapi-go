package armotypes

import "github.com/armosec/armoapi-go/identifiers"

type AttackChainStatus string
type ProcessingStatus string

const (
	StatusActive AttackChainStatus = "active"
	StatusFixed  AttackChainStatus = "fixed"
	// StatusFixedSeen AttackChainStatus = "fixedSeen"

	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusDone       ProcessingStatus = "done"
)

type AttackChainType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AttackChain struct {
	Type             *AttackChainType             `json:"type"`
	ClusterName      string                       `json:"clusterName"`
	Resource         identifiers.PortalDesignator `json:"resource"`
	AttackChainID    string                       `json:"attackChainID"` // name/cluster/resourceID
	CustomerGUID     string                       `json:"customerGUID"`
	AttackChainNodes AttackChainNode              `json:"attackChainNodes"`
	UIStatus         *AttackChainUIStatus         `json:"uiStatus"`
	LatestReportGUID string                       `json:"latestReportGUID"` // latest reportGUID in which this attack chain was identified
}

type AttackChainNode struct {
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	ControlIDs       []string                       `json:"controlIDs,omitempty"` // failed/ignored controls that are associated to this attack chain node
	Vulnerabilities  []Vulnerabilities              `json:"vulnerabilities,omitempty"`
	RelatedResources []identifiers.PortalDesignator `json:"relatedResources"`
	NextNodes        []AttackChainNode              `json:"nextNodes,omitempty"`
}

type Vulnerabilities struct {
	ContainersScanID string   `json:"containersScanID"`
	ContainerName    string   `json:"containerName"`
	Names            []string `json:"names"` // CVE names
}

// struct for UI support. All strings are timestamps
type AttackChainUIStatus struct {
	// fields updated by the BE
	FirstSeen string `json:"firstSeen"` // timestamp of first scan in which the attack chain was identified
	// fields updated by the UI
	ViewedMainScreen string `json:"wasViewedMainScreen"` // if the attack chain was viewed by the user// New badge
	ProcessingStatus string `json:"processingStatus"`    // "processing"/ "done"
}
