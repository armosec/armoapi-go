package containerscan

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/identifiers"
)

type RelevantLabel string

const RelevantLabelYes RelevantLabel = "yes"
const RelevantLabelNo RelevantLabel = "no"
const RelevantLabelNotExists RelevantLabel = ""

type CommonContainerVulnerabilityResult struct {
	Designators       identifiers.PortalDesignator `json:"designators"`
	IntroducedInLayer string                       `json:"layerHash"`
	WLID              string                       `json:"wlid"`
	ContainerScanID   string                       `json:"containersScanID"`
	Vulnerability     `json:",inline"`
	Layers            []ESLayer                                `json:"layers"`
	LayersNested      []ESLayer                                `json:"layersNested"`
	Context           []identifiers.ArmoContext                `json:"context"`
	RelevantLinks     []string                                 `json:"links"`
	RelatedExceptions []armotypes.VulnerabilityExceptionPolicy `json:"relatedExceptions,omitempty"`
	Timestamp         int64                                    `json:"timestamp"`
	IsLastScan        int                                      `json:"isLastScan"`
	IsFixed           int                                      `json:"isFixed"`
	RelevantLabel     RelevantLabel                            `json:"relevantLabel"`
	ClusterShortName  string                                   `json:"clusterShortName"`
}

type TopVulItem struct {
	Vulnerability   `json:",inline"`
	WorkloadsCount  int64 `json:"workloadsCount"`
	SeverityOverall int64 `json:"severityOverall"`
}

type ESLayer struct {
	LayerHash       string `json:"layerHash"`
	ParentLayerHash string `json:"parentLayerHash"`
	*LayerInfo
}

type LayerInfo struct {
	CreatedBy   string     `json:"createdBy,omitempty"`
	CreatedTime *time.Time `json:"createdTime,omitempty"`
	LayerOrder  int        `json:"layerOrder,omitempty"` // order 0 is first layer in the list
}

type SeverityStats struct {
	Severity                 string `json:"severity,omitempty"`
	HealthStatus             string `json:"healthStatus"`
	TotalCount               int64  `json:"total"`
	RCEFixCount              int64  `json:"rceFixCount"`
	RelevantFixCount         int64  `json:"relevantFixCount"`
	FixAvailableOfTotalCount int64  `json:"fixedTotal"`
	RelevantCount            int64  `json:"relevantTotal"`
	RCECount                 int64  `json:"rceTotal"`
	UrgentCount              int64  `json:"urgent"`
	NeglectedCount           int64  `json:"neglected"`
	RelevancyScanCount       int64  `json:"relevancyScanCount"`
}

type ShortVulnerabilityResult struct {
	Name string `json:"name"`
}

type CommonContainerScanSeveritySummary struct {
	Designators identifiers.PortalDesignator `json:"designators"`
	SeverityStats
	ImgTag          string                    `json:"imageTag"`
	ContainerName   string                    `json:"containerName"`
	CustomerGUID    string                    `json:"customerGUID"`
	ContainerScanID string                    `json:"containersScanID"`
	DayDate         string                    `json:"dayDate"`
	WLID            string                    `json:"wlid"`
	Version         string                    `json:"version"`
	ImgHash         string                    `json:"imageHash"`
	Cluster         string                    `json:"cluster"`
	Namespace       string                    `json:"namespace"`
	VersionImage    string                    `json:"versionImage"`
	Status          string                    `json:"status"`
	Registry        string                    `json:"registry"`
	JobIDs          []string                  `json:"jobIDs"`
	Context         []identifiers.ArmoContext `json:"context"`
	Timestamp       int64                     `json:"timestamp"`
}

type CommonContainerScanSummaryResult struct {
	Designators identifiers.PortalDesignator `json:"designators"`
	SeverityStats
	Version                       string                     `json:"version"`
	Registry                      string                     `json:"registry"`
	CustomerGUID                  string                     `json:"customerGUID"`
	ContainerScanID               string                     `json:"containersScanID"`
	ImageSignatureValidationError string                     `json:"imageSignatureValidationError,omitempty"`
	WLID                          string                     `json:"wlid"`
	ImageID                       string                     `json:"imageHash"`
	ImageTag                      string                     `json:"imageTag"`
	ClusterName                   string                     `json:"cluster"`
	ClusterShortName              string                     `json:"clusterShortName"`
	Namespace                     string                     `json:"namespace"`
	ContainerName                 string                     `json:"containerName"`
	ImageTagSuffix                string                     `json:"versionImage"`
	Status                        string                     `json:"status"`
	ExcludedSeveritiesStats       []SeverityStats            `json:"excludedSeveritiesStats,omitempty"`
	PackagesName                  []string                   `json:"packages"`
	SeveritiesStats               []SeverityStats            `json:"severitiesStats"`
	JobIDs                        []string                   `json:"jobIDs"`
	Vulnerabilities               []ShortVulnerabilityResult `json:"vulnerabilities"`
	Context                       []identifiers.ArmoContext  `json:"context"`
	Timestamp                     int64                      `json:"timestamp"`
	ImageSignatureValid           bool                       `json:"imageSignatureValid,omitempty"`
	ImageHasSignature             bool                       `json:"imageHasSignature,omitempty"`
	RelevantLabel                 RelevantLabel              `json:"relevantLabel"`
	HasRelevancyData              bool                       `json:"hasRelevancyData"`
}

type CommonContainerScanSummaryResultStub struct {
	CommonContainerScanSummaryResult `json:",inline"`
	IsStub                           bool     `json:"isStub,omitempty"`
	ErrorsList                       []string `json:"errors,omitempty"`
}

type DesignatorsToVulnerabilityNames struct {
	Designators        identifiers.PortalDesignator `json:"designators"`
	VulnerabilityNames []string                     `json:"vulnerabilityNames"`
}
