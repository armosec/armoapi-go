package notifications

import (
	"time"

	"github.com/armosec/armoapi-go/containerscan"
	"github.com/armosec/armoapi-go/identifiers"
)

type WeeklyReport struct {
	ClustersScannedThisWeek             int                        `json:"clustersScannedThisWeek" bson:"clustersScannedThisWeek"`
	ClustersScannedPrevWeek             int                        `json:"clustersScannedPrevWeek" bson:"clustersScannedPrevWeek"`
	LinkToConfigurationScanningFiltered string                     `json:"linkToConfigurationScanningFiltered" bson:"linkToConfigurationScanningFiltered"`
	RepositoriesScannedThisWeek         int                        `json:"repositoriesScannedThisWeek" bson:"repositoriesScannedThisWeek"`
	RepositoriesScannedPrevWeek         int                        `json:"repositoriesScannedPrevWeek" bson:"repositoriesScannedPrevWeek"`
	LinkToRepositoriesScanningFiltered  string                     `json:"linkToRepositoriesScanningFiltered" bson:"linkToRepositoriesScanningFiltered"`
	RegistriesScannedThisWeek           int                        `json:"registriesScannedThisWeek" bson:"registriesScannedThisWeek"`
	RegistriesScannedPrevWeek           int                        `json:"registriesScannedPrevWeek" bson:"registriesScannedPrevWeek"`
	LinkToRegistriesScanningFiltered    string                     `json:"linkToRegistriesScanningFiltered" bson:"linkToRegistriesScanningFiltered"`
	Top5FailedControls                  []TopCtrlItem              `json:"top5FailedControls" bson:"top5FailedControls"`
	Top5FailedCVEs                      []containerscan.TopVulItem `json:"top5FailedCVEs" bson:"top5FailedCVEs"`
	ClustersScanned                     []ClusterResourceScanned   `json:"clustersScanned" bson:"clustersScanned"`
	RepositoriesScanned                 []RepositoryScanned        `json:"repositoriesScanned" bson:"repositoriesScanned"`
	RegistriesScanned                   []RegistryScanned          `json:"registriesScanned" bson:"registriesScanned"`
}
type PushNotification struct {
	Misconfigurations Misconfigurations
	NewClusterAdmins  NewClusterAdmins
}

type NewClusterAdmins []NewClusterAdmin
type NewClusterAdmin struct {
	Resource          string
	Link              string
	ClusterName       string
	ClusterFullName   string
	ResourceName      string
	ResourceKind      string
	ResourceNamespace string
}

type Misconfigurations []Misconfiguration
type Misconfiguration struct {
	Name                      string
	FullName                  string
	Type                      ScanType
	Link                      string
	PercentageIncrease        uint64
	FrameworksComplianceDrift map[string]int
	PercentageThreshold       uint8
}
type ScanType string

const (
	ScanTypePosture      ScanType = "posture"
	ScanTypeRepositories ScanType = "repository"
)

type NotificationConfigIdentifier struct {
	NotificationType NotificationType `json:"notificationType,omitempty" bson:"notificationType,omitempty"`
}
type AlertChannel struct {
	ChannelType             ChannelProvider `json:"channelType,omitempty" bson:"channelType,omitempty"`
	Scope                   []AlertScope    `json:"scope,omitempty" bson:"scope,omitempty"`
	CollaborationConfigGUID string          `json:"collaborationConfigId,omitempty" bson:"collaborationConfigId,omitempty"`
	Alerts                  []AlertConfig   `json:"notifications,omitempty" bson:"notifications,omitempty"`
}

type NotificationParams struct {
	DriftPercentage     *int     `json:"driftPercentage,omitempty" bson:"driftPercentage,omitempty"`
	MinSeverity         *int     `json:"minSeverity,omitempty" bson:"minSeverity,omitempty"` // To be DEPRECATED after workflows is live.
	IncidentPolicyGUIDs []string `json:"incidentPolicyGUIDs,omitempty" bson:"incidentPolicyGUIDs,omitempty"`

	// params for workflows
	Severities []string `json:"severities,omitempty" bson:"severities,omitempty"`

	// vulnerability params
	KnownExploited *bool    `json:"knownExploited,omitempty" bson:"knownExploited,omitempty"` // Known Exploited (CISA KEV)
	HighLikelihood *bool    `json:"highLikelihood,omitempty" bson:"highLikelihood,omitempty"` // High Likelihood (EPSS ≥ 10%)
	CVSS           *float32 `json:"cvss,omitempty" bson:"cvss,omitempty"`                     // CVSS (Common Vulnerability Scoring System)
	InUse          *bool    `json:"inUse,omitempty" bson:"inUse,omitempty"`                   // In Use (CISA IU)
	Fixable        *bool    `json:"fixable,omitempty" bson:"fixable,omitempty"`               // Fixable (CISA FX)

	// security risks params
	SecurityRiskIDs []string `json:"securityRiskIDs,omitempty" bson:"securityRiskIDs,omitempty"` // Security Risk ID

}

type AlertConfig struct {
	NotificationConfigIdentifier `json:",inline" bson:",inline"`
	Parameters                   NotificationParams `json:"parameters,omitempty" bson:"parameters,omitempty"`
	Disabled                     *bool              `json:"disabled,omitempty" bson:"disabled,omitempty"`
}

type AlertScope struct {
	Cluster    string   `json:"cluster,omitempty" bson:"cluster,omitempty"`
	Namespaces []string `json:"namespaces,omitempty" bson:"namespaces,omitempty"`
}

type EnrichedScope struct {
	AlertScope       `json:",inline"`
	ClusterShortName string `json:"clusterShortName,omitempty"`
}

type NotificationType string

const (
	NotificationTypeWeekly              NotificationType = "weekly"              //weekly report
	NotificationTypePush                NotificationType = "push"                //posture scan
	NotificationTypeContainerPush       NotificationType = "containerScanPush"   //container scan
	NotificationTypeSecurityRiskPush    NotificationType = "securityRiskPush"    //security risk
	NotificationTypeRuntimeIncidentPush NotificationType = "runtimeIncidentPush" // runtime incident (kdr)

	NotificationTypeComplianceDrift     NotificationType = NotificationTypePush + ":complianceDrift"
	NotificationTypeNewClusterAdmin     NotificationType = NotificationTypePush + ":newClusterAdmin"
	NotificationTypeNewVulnerability    NotificationType = NotificationTypeContainerPush + ":newVulnerability"
	NotificationTypeVulnerabilityNewFix NotificationType = NotificationTypeContainerPush + ":vulnerabilityNewFix"

	NotificationTypeSecurityRiskNew    NotificationType = NotificationTypeSecurityRiskPush + ":newSecurityRisk"
	NotificationTypeRuntimeIncidentNew NotificationType = NotificationTypeRuntimeIncidentPush + ":newRuntimeIncident"
)

var notificationTypes = []NotificationType{
	NotificationTypePush,
	NotificationTypeWeekly,
	NotificationTypeComplianceDrift,
	NotificationTypeNewClusterAdmin,
	NotificationTypeNewVulnerability,
	NotificationTypeVulnerabilityNewFix,
}

type PushReport struct {
	Cluster                   string             `json:"custer,omitempty" bson:"custer,omitempty"`
	ReportGUID                string             `json:"reportGUID,omitempty" bson:"reportGUID,omitempty"`
	ScanType                  ScanType           `json:"scanType" bson:"scanType"`
	Timestamp                 time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	FailedResources           uint64             `json:"failedResources,omitempty" bson:"failedResources,omitempty"`
	FrameworksComplianceScore map[string]float32 `json:"frameworksComplianceScore,omitempty" bson:"frameworksComplianceScore,omitempty"`
}

type RegistryScanned struct {
	Registry ResourceScanned `json:"registry" bson:"registry"`
}

type RepositoryScanned struct {
	ReportGUID string          `json:"reportGUID" bson:"reportGUID"`
	Repository ResourceScanned `json:"repository" bson:"repository"`
}

type ClusterResourceScanned struct {
	ShortName       string          `json:"shortName" bson:"shortName"`
	Cluster         ResourceScanned `json:"cluster" bson:"cluster"`
	ReportGUID      string          `json:"reportGUID" bson:"reportGUID"`
	FailedResources uint64          `json:"failedResources" bson:"failedResources"`
}

type ResourceScanned struct {
	Kind                         string                     `json:"kind" bson:"kind"`
	Name                         string                     `json:"name" bson:"name"`
	MapSeverityToSeverityDetails map[string]SeverityDetails `json:"mapSeverityToSeverityDetails" bson:"mapSeverityToSeverityDetails"`
}

type SeverityDetails struct {
	Severity              string `json:"severity" bson:"severity"`
	FailedResourcesNumber int    `json:"failedResourcesNumber" bson:"failedResourcesNumber"`
}

const (
	NotificationBeforeUpdateContainerScanEvent = "beforeUpdateContainerScan"
)

type NotificationPushEvent struct {
	EventName   string                       `json:"eventName"`
	EventTime   time.Time                    `json:"eventTime"`
	Designators identifiers.PortalDesignator `json:"designators,omitempty"`
}
