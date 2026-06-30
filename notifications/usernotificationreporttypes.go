package notifications

import (
	"encoding/json"
	"strconv"
	"strings"
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
type JiraTicketIdentifiers struct {
	CollaborationGUID string                 `json:"collaborationGUID,omitempty" bson:"collaborationGUID,omitempty"`
	SiteID            string                 `json:"siteId,omitempty" bson:"siteId,omitempty"`
	ProjectID         string                 `json:"projectId,omitempty" bson:"projectId,omitempty"`
	IssueTypeID       string                 `json:"issueTypeId,omitempty" bson:"issueTypeId,omitempty"`
	Fields            map[string]interface{} `json:"fields,omitempty" bson:"fields,omitempty"`
}

type LinearTicketIdentifiers struct {
	WorkspaceID string                 `json:"workspaceId,omitempty" bson:"workspaceId,omitempty"`
	TeamID      string                 `json:"teamId,omitempty" bson:"teamId,omitempty"`
	AssigneeID  string                 `json:"assigneeId,omitempty" bson:"assigneeId,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty" bson:"fields,omitempty"`
}

type GitHubTicketIdentifiers struct {
	CollaborationGUID string   `json:"collaborationGUID,omitempty" bson:"collaborationGUID,omitempty"`
	OrganizationName  string   `json:"organizationName,omitempty" bson:"organizationName,omitempty"`
	RepositoryName    string   `json:"repositoryName,omitempty" bson:"repositoryName,omitempty"`
	Labels            []string `json:"labels,omitempty" bson:"labels,omitempty"`
	Assignees         []string `json:"assignees,omitempty" bson:"assignees,omitempty"`
	MilestoneID       *int     `json:"milestoneId,omitempty" bson:"milestoneId,omitempty"`
}

// UnmarshalJSON decodes GitHubTicketIdentifiers tolerantly so the workflow /
// runtime-incident payload accepts the same select-option shape the UI submits
// for every provider, in addition to the historical bare forms:
//   - labels / assignees: bare strings (["bug"]) OR option objects
//     ([{"id":"bug"}] / [{"name":"bug"}]).
//   - milestone: the bare "milestoneId" integer (preferred) OR a "milestone"
//     select option ({"id":"12"} / {"id":12} / "12").
//
// Everything normalizes into the typed fields, so existing readers (e.g. the
// notification senders) and re-marshaling — which always emits the bare form —
// are unaffected. For GitHub the option id IS the value the API consumes (label
// name, assignee login, milestone number), mirroring the manual create path.
func (g *GitHubTicketIdentifiers) UnmarshalJSON(data []byte) error {
	var shadow struct {
		CollaborationGUID string          `json:"collaborationGUID"`
		OrganizationName  string          `json:"organizationName"`
		RepositoryName    string          `json:"repositoryName"`
		Labels            json.RawMessage `json:"labels"`
		Assignees         json.RawMessage `json:"assignees"`
		MilestoneID       json.RawMessage `json:"milestoneId"`
		Milestone         json.RawMessage `json:"milestone"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}

	g.CollaborationGUID = shadow.CollaborationGUID
	g.OrganizationName = shadow.OrganizationName
	g.RepositoryName = shadow.RepositoryName
	g.Labels = githubSelectStrings(shadow.Labels)
	g.Assignees = githubSelectStrings(shadow.Assignees)

	g.MilestoneID = nil
	if n, ok := githubMilestoneNumber(shadow.MilestoneID); ok {
		g.MilestoneID = &n
	} else if n, ok := githubMilestoneNumber(shadow.Milestone); ok {
		g.MilestoneID = &n
	}
	return nil
}

// githubSelectStrings coerces a labels/assignees JSON value into []string,
// accepting bare strings and option objects ({"id":...} / {"name":...}).
// Elements that resolve to "" are skipped; null/absent/empty input yields nil.
func githubSelectStrings(raw json.RawMessage) []string {
	if isEmptyJSON(raw) {
		return nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil || len(items) == 0 {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s := githubOptionString(item); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// githubOptionString extracts the string value of a single select-field option,
// reading "id" first then "name" for object options. The option id is the value
// GitHub consumes. Returns "" when nothing usable is found.
func githubOptionString(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}
	for _, key := range []string{"id", "name"} {
		v, ok := obj[key]
		if !ok {
			continue
		}
		var str string
		if err := json.Unmarshal(v, &str); err == nil && str != "" {
			return str
		}
		var num float64
		if err := json.Unmarshal(v, &num); err == nil { // numeric id, e.g. {"id": 1}
			return strconv.Itoa(int(num))
		}
	}
	return ""
}

// githubMilestoneNumber coerces a milestone value into the integer milestone
// number GitHub requires. Accepts a JSON number, a numeric string, or an option
// object ({"id": ...}). Returns ok=false for null/absent/non-numeric input.
func githubMilestoneNumber(raw json.RawMessage) (int, bool) {
	if isEmptyJSON(raw) {
		return 0, false
	}
	var num float64
	if err := json.Unmarshal(raw, &num); err == nil {
		return int(num), true
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if t := strings.TrimSpace(s); t != "" {
			if i, err := strconv.Atoi(t); err == nil {
				return i, true
			}
		}
		return 0, false
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err == nil {
		if v, ok := obj["id"]; ok {
			return githubMilestoneNumber(v)
		}
	}
	return 0, false
}

func isEmptyJSON(raw json.RawMessage) bool {
	return len(raw) == 0 || string(raw) == "null"
}

type AlertChannel struct {
	ChannelType             ChannelProvider           `json:"channelType,omitempty" bson:"channelType,omitempty"`
	Scope                   []AlertScope              `json:"scope,omitempty" bson:"scope,omitempty"`
	CollaborationConfigGUID string                    `json:"collaborationConfigId,omitempty" bson:"collaborationConfigId,omitempty"`
	Alerts                  []AlertConfig             `json:"notifications,omitempty" bson:"notifications,omitempty"`
	JiraTicketIdentifiers   []JiraTicketIdentifiers   `json:"jiraTicketIdentifiers,omitempty" bson:"jiraTicketIdentifiers,omitempty"`
	LinearTicketIdentifiers []LinearTicketIdentifiers `json:"linearTicketIdentifiers,omitempty" bson:"linearTicketIdentifiers,omitempty"`
	GitHubTicketIdentifiers []GitHubTicketIdentifiers `json:"githubTicketIdentifiers,omitempty" bson:"githubTicketIdentifiers,omitempty"`
}

type NotificationParams struct {
	DriftPercentage     *int     `json:"driftPercentage,omitempty" bson:"driftPercentage,omitempty"`
	MinSeverity         *int     `json:"minSeverity,omitempty" bson:"minSeverity,omitempty"` // To be DEPRECATED after workflows is live.
	IncidentPolicyGUIDs []string `json:"incidentPolicyGUIDs,omitempty" bson:"incidentPolicyGUIDs,omitempty"`

	// params for workflows
	Severities []string `json:"severities,omitempty" bson:"severities,omitempty"`

	// vulnerability params
	KnownExploited   *bool    `json:"knownExploited,omitempty" bson:"knownExploited,omitempty"`     // Known Exploited (CISA KEV)
	HighLikelihood   *bool    `json:"highLikelihood,omitempty" bson:"highLikelihood,omitempty"`     // High Likelihood (EPSS ≥ 10%)
	CVSS             *float32 `json:"cvss,omitempty" bson:"cvss,omitempty"`                         // CVSS (Common Vulnerability Scoring System) min threshold (≥)
	CVSSMax          *float32 `json:"cvssMax,omitempty" bson:"cvssMax,omitempty"`                   // CVSS max threshold (≤)
	InUse            *bool    `json:"inUse,omitempty" bson:"inUse,omitempty"`                       // In Use (CISA IU)
	Fixable          *bool    `json:"fixable,omitempty" bson:"fixable,omitempty"`                   // Fixable (CISA FX)
	NoExploitability *bool    `json:"noExploitability,omitempty" bson:"noExploitability,omitempty"` // No exploit intelligence
	RiskFactors      []string `json:"riskFactors,omitempty" bson:"riskFactors,omitempty"`           // Risk Factors

	// security risks params
	SecurityRiskIDs []string `json:"securityRiskIDs,omitempty" bson:"securityRiskIDs,omitempty"` // Security Risk ID

	// cluster status params
	ClusterStatus []string `json:"clusterStatus,omitempty" bson:"clusterStatus,omitempty"` // Cluster Status

	// system health — scan failure
	ScanFailure *bool `json:"scanFailure,omitempty" bson:"scanFailure,omitempty"`

	// hot CVE notification
	HotCVE *bool `json:"hotCVE,omitempty" bson:"hotCVE,omitempty"`
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
	NotificationTypeClusterStatusPush   NotificationType = "clusterStatusPush"   // cluster status
	NotificationTypeComplianceDrift     NotificationType = NotificationTypePush + ":complianceDrift"
	NotificationTypeNewClusterAdmin     NotificationType = NotificationTypePush + ":newClusterAdmin"
	NotificationTypeNewVulnerability    NotificationType = NotificationTypeContainerPush + ":newVulnerability"
	NotificationTypeVulnerabilityNewFix NotificationType = NotificationTypeContainerPush + ":vulnerabilityNewFix"

	NotificationTypeSecurityRiskNew    NotificationType = NotificationTypeSecurityRiskPush + ":newSecurityRisk"
	NotificationTypeRuntimeIncidentNew NotificationType = NotificationTypeRuntimeIncidentPush + ":newRuntimeIncident"

	NotificationTypeScanFailurePush NotificationType = "scanFailurePush"                                   // scan failure
	NotificationTypeScanFailureNew  NotificationType = NotificationTypeScanFailurePush + ":newScanFailure" // new scan failure event
)

var notificationTypes = []NotificationType{
	NotificationTypePush,
	NotificationTypeWeekly,
	NotificationTypeComplianceDrift,
	NotificationTypeNewClusterAdmin,
	NotificationTypeNewVulnerability,
	NotificationTypeVulnerabilityNewFix,
	NotificationTypeScanFailurePush,
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
