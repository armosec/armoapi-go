package armotypes

type WeeklyReport struct {
	ClustersScannedThisWeek             int                      `json:"clustersScannedThisWeek" bson:"clustersScannedThisWeek"`
	ClustersScannedPrevWeek             int                      `json:"clustersScannedPrevWeek" bson:"clustersScannedPrevWeek"`
	LinkToConfigurationScanningFiltered string                   `json:"linkToConfigurationScanningFiltered" bson:"linkToConfigurationScanningFiltered"`
	RepositoriesScannedThisWeek         int                      `json:"repositoriesScannedThisWeek" bson:"repositoriesScannedThisWeek"`
	RepositoriesScannedPrevWeek         int                      `json:"repositoriesScannedPrevWeek" bson:"repositoriesScannedPrevWeek"`
	LinkToRepositoriesScanningFiltered  string                   `json:"linkToRepositoriesScanningFiltered" bson:"linkToRepositoriesScanningFiltered"`
	RegistriesScannedThisWeek           int                      `json:"registriesScannedThisWeek" bson:"registriesScannedThisWeek"`
	RegistriesScannedPrevWeek           int                      `json:"registriesScannedPrevWeek" bson:"registriesScannedPrevWeek"`
	LinkToRegistriesScanningFiltered    string                   `json:"linkToRegistriesScanningFiltered" bson:"linkToRegistriesScanningFiltered"`
	Top5FailedControls                  []TopCtrlItem            `json:"top5FailedControls" bson:"top5FailedControls"`
	Top5FailedCVEs                      []TopVulItem             `json:"top5FailedCVEs" bson:"top5FailedCVEs"`
	ClustersScanned                     []ClusterResourceScanned `json:"clustersScanned" bson:"clustersScanned"`
	RepositoriesScanned                 []RepositoryScanned      `json:"repositoriesScanned" bson:"repositoriesScanned"`
	RegistriesScanned                   []RegistryScanned        `json:"registriesScanned" bson:"registriesScanned"`
}

type NotificationsConfig struct {
	//Map of unsubscribed user id to notification config identifier
	UnsubscribedUsers  map[string]NotificationConfigIdentifier `json:"unsubscribedUsers,omitempty" bson:"unsubscribedUsers,omitempty"`
	LatestWeeklyReport *WeeklyReport                           `json:"latestWeeklyReport,omitempty" bson:"latestWeeklyReport,omitempty"`
}

type NotificationConfigIdentifier struct {
	NotificationType NotificationType `json:"notificationType,omitempty" bson:"notificationType,omitempty"`
}

type NotificationType string

const (
	NotificationTypeAll    NotificationType = "all"
	NotificationTypePush   NotificationType = "push"
	NotificationTypeWeekly NotificationType = "weekly"
)

type RegistryScanned struct {
	Registry ResourceScanned `json:"registry" bson:"registry"`
}

type RepositoryScanned struct {
	Repository ResourceScanned `json:"repository" bson:"repository"`
}

type ClusterResourceScanned struct {
	ShortName string          `json:"shortName" bson:"shortName"`
	Cluster   ResourceScanned `json:"cluster" bson:"cluster"`
}

type ResourceScanned struct {
	Name                         string                     `json:"name" bson:"name"`
	MapSeverityToResourceDetails map[string]ResourceDetails `json:"resourcesToDetails" bson:"resourcesToDetails"`
}

type ResourceDetails struct {
	FailedResourcesNumber int `json:"failedResourcesNumber" bson:"failedResourcesNumber"`
}

type Vulnerability struct {
	Name               string                         `json:"name" bson:"name"`
	ImgHash            string                         `json:"imageHash" bson:"imageHash"`
	ImgTag             string                         `json:"imageTag" bson:"imageTag"`
	RelatedPackageName string                         `json:"packageName" bson:"packageName"`
	PackageVersion     string                         `json:"packageVersion" bson:"packageVersion"`
	Link               string                         `json:"link" bson:"link"`
	Description        string                         `json:"description" bson:"description"`
	Severity           string                         `json:"severity" bson:"severity"`
	SeverityScore      int                            `json:"severityScore" bson:"severityScore"`
	Metadata           interface{}                    `json:"metadata" bson:"metadata"`
	Fixes              VulFixes                       `json:"fixedIn" bson:"fixedIn"`
	Relevancy          string                         `json:"relevant" bson:"relevant"`
	UrgentCount        int                            `json:"urgent" bson:"urgent"`
	NeglectedCount     int                            `json:"neglected" bson:"neglected"`
	HealthStatus       string                         `json:"healthStatus" bson:"healthStatus"`
	Categories         VulnerabilityCategory          `json:"categories" bson:"categories"`
	ExceptionApplied   []VulnerabilityExceptionPolicy `json:"exceptionApplied,omitempty" bson:"exceptionApplied,omitempty"`
}

type VulFixes []FixedIn

type VulnerabilityCategory struct {
	IsRCE bool `json:"isRce" bson:"isRce"`
}

type FixedIn struct {
	Name    string `json:"name" bson:"name"`
	ImgTag  string `json:"imageTag" bson:"imageTag"`
	Version string `json:"version" bson:"version"`
}
type TopVulItem struct {
	Vulnerability   `json:",inline"`
	WorkloadsCount  int64 `json:"workloadsCount" bson:"workloadsCount"`
	SeverityOverall int64 `json:"severityOverall" bson:"severityOverall"`
}

type TopCtrlItem struct {
	ControlID            string           `json:"id" bson:"id"`
	ControlGUID          string           `json:"guid" bson:"guid"`
	Name                 string           `json:"name" bson:"name"`
	Remediation          string           `json:"remediation" bson:"remediation"`
	Description          string           `json:"description" bson:"description"`
	ClustersCount        int64            `json:"clustersCount" bson:"clustersCount"`
	SeverityOverall      int64            `json:"severityOverall" bson:"severityOverall"`
	BaseScore            int64            `json:"baseScore" bson:"baseScore"`
	Clusters             []TopCtrlCluster `json:"clusters" bson:"clusters"`
	TotalFailedResources int64            `json:"-"`
}

type TopCtrlCluster struct {
	Name               string `json:"name" bson:"name"`
	ResourcesCount     int64  `json:"resourcesCount" bson:"resourcesCount"`
	ReportGUID         string `json:"reportGUID" bson:"reportGUID"`
	TopFailedFramework string `json:"topFailedFramework" bson:"topFailedFramework"`
}

type CommonContainerScanSummaryResult struct {
	SeverityStats
	Designators     PortalDesignator `json:"designators" bson:"designators"`
	Context         []ArmoContext    `json:"context" bson:"context"`
	JobIDs          []string         `json:"jobIDs" bson:"jobIDs"`
	CustomerGUID    string           `json:"customerGUID" bson:"customerGUID"`
	ContainerScanID string           `json:"containersScanID" bson:"containersScanID"`

	Timestamp     int64    `json:"timestamp" bson:"timestamp"`
	WLID          string   `json:"wlid" bson:"wlid"`
	ImgTag        string   `json:"imageTag" bson:"imageTag"`
	ImgHash       string   `json:"imageHash" bson:"imageHash"`
	Cluster       string   `json:"cluster" bson:"cluster"`
	Namespace     string   `json:"namespace" bson:"namespace"`
	ContainerName string   `json:"containerName" bson:"containerName"`
	PackagesName  []string `json:"packages" bson:"packages"`

	ListOfDangerousArtifcats []string `json:"listOfDangerousArtifcats" bson:"listOfDangerousArtifcats"`

	Status string `json:"status" bson:"status"`

	Registry     string `json:"registry" bson:"registry"`
	VersionImage string `json:"versionImage" bson:"versionImage"`

	SeveritiesStats         []SeverityStats `json:"severitiesStats" bson:"severitiesStats"`
	ExcludedSeveritiesStats []SeverityStats `json:"excludedSeveritiesStats,omitempty" bson:"excludedSeveritiesStats,omitempty"`

	Version string `json:"version" bson:"version"`

	Vulnerabilities []ShortVulnerabilityResult `json:"vulnerabilities" bson:"vulnerabilities"`
}

type SeverityStats struct {
	Severity                     string `json:"severity,omitempty" bson:"severity,omitempty"`
	TotalCount                   int64  `json:"total" bson:"total"`
	RCEFixCount                  int64  `json:"rceFixCount" bson:"rceFixCount"`
	FixAvailableOfTotalCount     int64  `json:"fixedTotal" bson:"fixedTotal"`
	RelevantCount                int64  `json:"totalRelevant" bson:"totalRelevant"`
	FixAvailableForRelevantCount int64  `json:"fixedRelevant" bson:"fixedRelevant"`
	RCECount                     int64  `json:"rceTotal" bson:"rceTotal"`
	UrgentCount                  int64  `json:"urgent" bson:"urgent"`
	NeglectedCount               int64  `json:"neglected" bson:"neglected"`
	HealthStatus                 string `json:"healthStatus" bson:"healthStatus"`
}

type ShortVulnerabilityResult struct {
	Name string `json:"name" bson:"name"`
}
