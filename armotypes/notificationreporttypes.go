package armotypes

type WeeklyReport struct {
	ClustersScannedThisWeek             int                 `json:"clustersScannedThisWeek"`
	ClustersScannedPrevWeek             int                 `json:"clustersScannedPrevWeek"`
	LinkToConfigurationScanningFiltered string              `json:"linkToConfigurationScanningFiltered"`
	RepositoriesScannedThisWeek         int                 `json:"repositoriesScannedThisWeek"`
	RepositoriesScannedPrevWeek         int                 `json:"repositoriesScannedPrevWeek"`
	LinkToRepositoriesScanningFiltered  string              `json:"linkToRepositoriesScanningFiltered"`
	RegistriesScannedThisWeek           int                 `json:"registriesScannedThisWeek"`
	RegistriesScannedPrevWeek           int                 `json:"registriesScannedPrevWeek"`
	LinkToRegistriesScanningFiltered    string              `json:"linkToRegistriesScanningFiltered"`
	Top5FailedControls                  []TopCtrlItem       `json:"top5FailedControls"`
	Top5FailedCVEs                      []TopVulItem        `json:"top5FailedCVEs"`
	ClustersScanned                     []ClusterScanned    `json:"clustersScanned"`
	RepositoriesScanned                 []RepositoryScanned `json:"repositoriesScanned"`
	RegistriesScanned                   []RegistryScanned   `json:"registriesScanned"`
}

type ClusterScanned struct {
	Cluster ClusterResourceScanned `json:"cluster"`
}

type RegistryScanned struct {
	Registry ResourceScanned `json:"registry"`
}

type RepositoryScanned struct {
	Repository ResourceScanned `json:"repository"`
}

type ClusterResourceScanned struct {
	FullName  string                     `json:"fullName"`
	ShortName string                     `json:"shortName"`
	Resource  map[string]ResourceDetails `json:"resourcesToDetails"`
}

type ResourceScanned struct {
	Name     string                     `json:"name"`
	Resource map[string]ResourceDetails `json:"resourcesToDetails"`
}

type ResourceDetails struct {
	FailedResourcesNumber int `json:"failedResourcesNumber"`
}

type Vulnerability struct {
	Name               string                         `json:"name"`
	ImgHash            string                         `json:"imageHash"`
	ImgTag             string                         `json:"imageTag"`
	RelatedPackageName string                         `json:"packageName"`
	PackageVersion     string                         `json:"packageVersion"`
	Link               string                         `json:"link"`
	Description        string                         `json:"description"`
	Severity           string                         `json:"severity"`
	SeverityScore      int                            `json:"severityScore"`
	Metadata           interface{}                    `json:"metadata"`
	Fixes              VulFixes                       `json:"fixedIn"`
	Relevancy          string                         `json:"relevant"`
	UrgentCount        int                            `json:"urgent"`
	NeglectedCount     int                            `json:"neglected"`
	HealthStatus       string                         `json:"healthStatus"`
	Categories         VulnerabilityCategory          `json:"categories"`
	ExceptionApplied   []VulnerabilityExceptionPolicy `json:"exceptionApplied,omitempty"`
}

type VulFixes []FixedIn

type VulnerabilityCategory struct {
	IsRCE bool `json:"isRce"`
}

type FixedIn struct {
	Name    string `json:"name"`
	ImgTag  string `json:"imageTag"`
	Version string `json:"version"`
}
type TopVulItem struct {
	Vulnerability   `json:",inline"`
	WorkloadsCount  int64 `json:"workloadsCount"`
	SeverityOverall int64 `json:"severityOverall"`
}

type TopCtrlItem struct {
	ControlID            string           `json:"id"`
	ControlGUID          string           `json:"guid"`
	Name                 string           `json:"name"`
	Remediation          string           `json:"remediation"`
	Description          string           `json:"description"`
	ClustersCount        int64            `json:"clustersCount"`
	SeverityOverall      int64            `json:"severityOverall"`
	BaseScore            int64            `json:"baseScore"`
	Clusters             []TopCtrlCluster `json:"clusters"`
	TotalFailedResources int64            `json:"-"`
}

type TopCtrlCluster struct {
	Name               string `json:"name"`
	ResourcesCount     int64  `json:"resourcesCount"`
	ReportGUID         string `json:"reportGUID"`
	TopFailedFramework string `json:"topFailedFramework"`
}

type CommonContainerScanSummaryResult struct {
	SeverityStats
	Designators     PortalDesignator `json:"designators"`
	Context         []ArmoContext    `json:"context"`
	JobIDs          []string         `json:"jobIDs"`
	CustomerGUID    string           `json:"customerGUID"`
	ContainerScanID string           `json:"containersScanID"`

	Timestamp     int64    `json:"timestamp"`
	WLID          string   `json:"wlid"`
	ImgTag        string   `json:"imageTag"`
	ImgHash       string   `json:"imageHash"`
	Cluster       string   `json:"cluster"`
	Namespace     string   `json:"namespace"`
	ContainerName string   `json:"containerName"`
	PackagesName  []string `json:"packages"`

	ListOfDangerousArtifcats []string `json:"listOfDangerousArtifcats"`

	Status string `json:"status"`

	Registry     string `json:"registry"`
	VersionImage string `json:"versionImage"`

	SeveritiesStats         []SeverityStats `json:"severitiesStats"`
	ExcludedSeveritiesStats []SeverityStats `json:"excludedSeveritiesStats,omitempty"`

	Version string `json:"version"`

	Vulnerabilities []ShortVulnerabilityResult `json:"vulnerabilities"`
}

type SeverityStats struct {
	Severity                     string `json:"severity,omitempty"`
	TotalCount                   int64  `json:"total"`
	RCEFixCount                  int64  `json:"rceFixCount"`
	FixAvailableOfTotalCount     int64  `json:"fixedTotal"`
	RelevantCount                int64  `json:"totalRelevant"`
	FixAvailableForRelevantCount int64  `json:"fixedRelevant"`
	RCECount                     int64  `json:"rceTotal"`
	UrgentCount                  int64  `json:"urgent"`
	NeglectedCount               int64  `json:"neglected"`
	HealthStatus                 string `json:"healthStatus"`
}

type ShortVulnerabilityResult struct {
	Name string `json:"name"`
}
