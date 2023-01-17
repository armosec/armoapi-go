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
	UnsubscribedUsers  map[string][]NotificationConfigIdentifier `json:"unsubscribedUsers,omitempty" bson:"unsubscribedUsers,omitempty"`
	LatestWeeklyReport *WeeklyReport                             `json:"latestWeeklyReport,omitempty" bson:"latestWeeklyReport,omitempty"`
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
	MapSeverityToSeverityDetails map[string]SeverityDetails `json:"mapSeverityToSeverityDetails" bson:"mapSeverityToSeverityDetails"`
}

type SeverityDetails struct {
	Severity              string `json:"severity" bson:"severity"`
	FailedResourcesNumber int    `json:"failedResourcesNumber" bson:"failedResourcesNumber"`
}
