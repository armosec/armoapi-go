package containerscan

import (
	"github.com/armosec/armoapi-go/apis"
	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/identifiers"
)

type ScanReport interface {
	IsLastReport() bool
	GetDesignators() identifiers.PortalDesignator
	GetContainerScanID() string
	GetTimestamp() int64
	GetWorkloadHash() string
	GetCustomerGUID() string
	GetSummary() ContainerScanSummaryResult
	GetVulnerabilities() []ContainerScanVulnerabilityResult
	GetVersion() string
	GetPaginationInfo() apis.PaginationMarks
	Validate() bool

	SetDesignators(identifiers.PortalDesignator)
	SetContainerScanID(string)
	SetTimestamp(int64)
	SetWorkloadHash(string)
	SetCustomerGUID(string)
}

type ContainerScanSummaryResult interface {
	GetDesignators() identifiers.PortalDesignator
	GetContext() []identifiers.ArmoContext
	GetWLID() string
	GetImageTag() string
	GetImageID() string
	GetSeverityStats() SeverityStats
	GetSeveritiesStats() []SeverityStats
	GetClusterName() string
	GetClusterShortName() string
	GetNamespace() string
	GetContainerName() string
	GetStatus() string
	GetRegistry() string
	GetRepository() string
	GetImageTageSuffix() string
	GetVersion() string
	GetCustomerGUID() string
	GetContainerScanID() string
	GetTimestamp() int64
	GetJobIDs() []string
	GetRelevantLabel() RelevantLabel
	Validate() bool
	GetHasRelevancyData() bool

	SetDesignators(identifiers.PortalDesignator)
	SetContext([]identifiers.ArmoContext)
	SetWLID(string)
	SetImageTag(string)
	SetImageID(string)
	SetSeverityStats(SeverityStats)
	SetSeveritiesStats([]SeverityStats)
	SetClusterName(string)
	SetClusterShortName(string)
	SetNamespace(string)
	SetContainerName(string)
	SetStatus(string)
	SetRegistry(string)
	SetImageTageSuffix(string)
	SetVersion(string)
	SetCustomerGUID(string)
	SetContainerScanID(string)
	SetTimestamp(int64)
	SetRelevantLabel(RelevantLabel)
	SetHasRelevancyData(bool)
}

type ContainerScanVulnerabilityResult interface {
	GetDesignators() identifiers.PortalDesignator
	GetContext() []identifiers.ArmoContext
	GetWLID() string
	GetContainerScanID() string
	GetLayers() []ESLayer
	GetLayersNested() []ESLayer
	GetTimestamp() int64
	GetIsLastScan() int
	GetIsFixed() int
	GetIntroducedInLayer() string
	GetRelevantLinks() []string
	GetRelatedExceptions() []armotypes.VulnerabilityExceptionPolicy
	GetVulnerability() VulnerabilityResult
	GetRelevantLabel() RelevantLabel
	GetClusterShortName() string

	SetDesignators(designators identifiers.PortalDesignator)
	SetContext(context []identifiers.ArmoContext)
	SetWLID(wlid string)
	SetContainerScanID(containerScanID string)
	SetLayers(layers []ESLayer)
	SetLayersNested(layersNested []ESLayer)
	SetTimestamp(timestamp int64)
	SetIsLastScan(isLastScan int)
	SetIsFixed(isFixed int)
	SetIntroducedInLayer(introducedInLayer string)
	SetLink(link string)
	SetRelevantLinks(relevantLinks []string)
	SetRelatedExceptions(relatedExceptions []armotypes.VulnerabilityExceptionPolicy)
	SetRelevantLabel(relevantLabel RelevantLabel)
	SetClusterShortName(clusterShortName string)
}

type VulnerabilityResult interface {
	GetName() string
	GetImageID() string
	GetImageTag() string
	GetRelatedPackageName() string
	GetPackageType() string
	GetPackageVersion() string
	GetLink() string
	GetDescription() string
	GetSeverity() string
	GetSeverityScore() int
	GetFixes() VulFixes
	GetIsRelevant() *bool
	GetUrgentCount() int
	GetNeglectedCount() int
	GetHealthStatus() string
	GetCategories() VulnerabilityCategory
	GetExceptionApplied() []armotypes.VulnerabilityExceptionPolicy

	SetName(string)
	SetImageID(string)
	SetImageTag(string)
	SetRelatedPackageName(string)
	SetPackageType(string)
	SetPackageVersion(string)
	SetLink(string)
	SetDescription(string)
	SetSeverity(string)
	SetSeverityScore(int)
	SetFixes(VulFixes)
	SetIsRelevant(*bool)
	SetUrgentCount(int)
	SetNeglectedCount(int)
	SetHealthStatus(string)
	SetCategories(VulnerabilityCategory)
	SetExceptionApplied([]armotypes.VulnerabilityExceptionPolicy)
}
