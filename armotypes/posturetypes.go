package armotypes

import (
	"encoding/json"
	"time"
)

const (
	PostureControlStatusUnknown    = 0
	PostureControlStatusPassed     = 1
	PostureControlStatusWarning    = 2 // deprecated
	PostureControlStatusFailed     = 3
	PostureControlStatusSkipped    = 4
	PostureControlStatusIrrelevant = 5 // deprecated
	PostureControlStatusError      = 6

	PostureResourceMaxCtrls = 6
)

// TODO: use CommonSummaryFields in all of the summaries

// swagger:model
type CommonSummaryFields struct {
	// The unique id of the report this summary belongs to
	ReportID GUID `json:"reportGUID"`

	// The designators of this summary
	Designators *PortalDesignator `json:"designators"`

	// Time of the scan that produced this summary
	Timestamp time.Time `json:"timestamp"`

	// swagger:ignore
	// Indication if this summary is marked for deletetion
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}

// -------- /api/v1/posture/clustersOvertime response datastructures
type PostureClusterOverTime struct {
	Designators  PortalDesignator           `json:"designators,omitempty"`
	ClusterName  string                     `json:"clusterName"`
	Frameworks   []PostureFrameworkOverTime `json:"frameworks"`
	DeleteStatus RecordStatus               `json:"deletionStatus,omitempty"`
}

// Used for elastic
type PostureFrameworksOverTime struct {
	ClusterName string `json:"clusterName"`

	ScoreValue float32   `json:"value"`
	ReportID   string    `json:"reportGUID"`
	Timestamp  time.Time `json:"timestamp"`
	Framework  string    `json:"frameworkName"`
}

// PostureFrameworkOverTime - the response structure
type PostureFrameworkOverTime struct {
	// "frameworkName": "MITRE",
	//                 "riskScore": 54,
	RiskScore float32                         `json:"riskScore"`
	Framework string                          `json:"frameworkName"`
	Coords    []PostureFrameworkOverTimeCoord `json:"cords"`
}

type PostureFrameworkOverTimeCoord struct {
	ScoreValue float32   `json:"value"`
	ReportID   string    `json:"reportGUID"`
	Timestamp  time.Time `json:"timestamp"`
}

//---- /api/v1/posture/frameworks

type PostureFrameworkSummary struct {
	Name             string           `json:"name"`
	Score            float32          `json:"value"`
	ImprovementScore float32          `json:"improvementScore"`
	TotalControls    int              `json:"totalControls"`
	FailedControls   int              `json:"failedControls"`
	SkippedControls  int              `json:"skippedControls,omitempty"`
	WarningControls  int              `json:"warningControls,omitempty"` // Deprecated
	ReportID         string           `json:"reportGUID"`
	Designators      PortalDesignator `json:"designators"`

	Timestamp    time.Time    `json:"timestamp"`
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}

type PostureClusterSummary struct {
	Score           float32          `json:"score"`
	TotalControls   int              `json:"totalControls"`
	FailedControls  int              `json:"failedControls"`
	SkippedControls int              `json:"skippedControls,omitempty"`
	WarningControls int              `json:"warningControls,omitempty"` // Deprecated
	ReportID        string           `json:"reportGUID"`
	Designators     PortalDesignator `json:"designators"`

	Timestamp    time.Time    `json:"timestamp"`
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`

	Frameworks []string `json:"frameworks"`

	// Counters - Failed resources by severity
	CriticalSeverityResources int `json:"criticalSeverityResources"`
	HighSeverityResources     int `json:"highSeverityResources"`
	MediumSeverityResources   int `json:"mediumSeverityResources"`
	LowSeverityResources      int `json:"lowSeverityResources"`

	// Counters - Failed controls by severity
	CriticalSeverityControls int `json:"criticalSeverityControls"`
	HighSeverityControls     int `json:"highSeverityControls"`
	MediumSeverityControls   int `json:"mediumSeverityControls"`
	LowSeverityControls      int `json:"lowSeverityControls"`

	// Counters -  Resources by status
	PassedResources   int `json:"passedResources"`
	FailedResources   int `json:"failedResources"`
	SkippedResources  int `jsons:"skippedResources,omitempty"`
	ExcludedResources int `json:"excludedResources,omitempty"` // Deprecated

	// Metadata
	KubescapeVersion  string `json:"kubescapeVersion"`
	KubernetesVersion string `json:"kubernetesVersion"`
	WorkerNodeCount   int    `json:"workerNodeCount"`
	Location          string `json:"location"`
	CloudProvider     string `json:"cloudProvider"`

	// Information about the controls that were run on this entity
	// The key is the status of the control (`failed`, `passed`, etc)
	ControlsInfo map[string][]ControlInfo `json:"controlsInfo"`

	// Names of the cluster
	FullName   string `json:"clusterFullName"`
	ShortName  string `json:"clusterShortName"`
	PrefixName string `json:"clusterPrefixName"`
}

type PostureFrameworkSubsectionSummary struct {
	// The name (title) of the subsection
	// Example: General Policies
	Name string `json:"name"`

	// The name of the framework this subsection belongs to
	// Example: CIS
	Framework string `json:"framework"`

	// Unique id of the subsection inside its framework
	// Example: 5.7
	ID string `json:"id"`

	// Statistics about the controls that were run
	// The key is the status of the control (`failed`, `passed`, etc).
	// The value is the number of controls
	// Example: {"failed": 3, "passed": 4}
	ControlsStats map[string]uint `json:"controlsStats"`
}

type PostureContainerSummary struct {
	ContainerName string `json:"containerName"`
	ImageTag      string `json:"image,omitempty"`
}

type ControlInputs struct {
	Rulename string
	Inputs   []PostureAttributesList // Attribute = input list name, Values = list values
}

// ----/api/v1/posture/controls
type PostureControlSummary struct {
	Designators                    PortalDesignator `json:"designators"`
	ControlID                      string           `json:"id"` // "C0001"
	ControlGUID                    string           `json:"guid"`
	Name                           string           `json:"name"`
	AffectedResourcesCount         int              `json:"affectedResourcesCount"`
	FailedResourcesCount           int              `json:"failedResourcesCount"`
	SkippedResourcesCount          int              `json:"skippedResourcesCount"`
	WarningResourcesCount          int              `json:"warningResourcesCount"` // Deprecated
	PreviousAffectedResourcesCount int              `json:"previousAffectedResourcesCount"`
	PreviousFailedResourcesCount   int              `json:"previousFailedResourcesCount"`
	PreviousSkippedResourcesCount  int              `json:"previousSkippedResourcesCount"`
	PreviousWarningResourcesCount  int              `json:"previousWarningResourcesCount"` // Deprecated
	Framework                      string           `json:"frameworkName"`
	FrameworkSubSectionID          []string         `json:"frameworkSubsectionID,omitempty"`
	Remediation                    string           `json:"remediation"`
	Status                         int              `json:"status"`
	StatusText                     string           `json:"statusText"`
	SubStatusText                  string           `json:"subStatusText,omitempty"`
	Description                    string           `json:"description"`
	Section                        string           `json:"section"`
	Timestamp                      time.Time        `json:"timestamp"`
	ReportID                       string           `json:"reportGUID"`
	DeleteStatus                   RecordStatus     `json:"deletionStatus,omitempty"`
	Score                          float32          `json:"score"`
	ScoreFactor                    float32          `json:"baseScore"`
	ScoreWeight                    float32          `json:"scoreWeight"`
	ARMOImprovement                float32          `json:"ARMOimprovement"`
	RelevantCloudProvides          []string         `json:"relevantCloudProvides"`
	ControlInputs                  []ControlInputs  `json:"controlInputs"`
	IsLastScan                     int              `json:"isLastScan"`
	HighlightPathsCount            int64            `json:"highlightPathsCount"`
}

//---------/api/v1/posture/resources

// 1 resource per 1 control
type PostureResource struct {
	UniqueResourceResult string           `json:"uniqueResourceResult"` // FNV(customerGUID + cluster+resourceID+frameworkName + resource.ReportID) to allow fast search for aggregation
	Designators          PortalDesignator `json:"designators"`
	Name                 string           `json:"name"`       // wlid/sid and etc.
	ResourceID           string           `json:"resourceID"` //as given by kscape

	ControlName       string                      `json:"controlName"`
	HighlightPaths    []string                    `json:"highlightPaths"` // specifies "failedPath" - where exactly in the raw resources the control failed
	FixPaths          []FixPath                   `json:"fixPaths"`       // specifies "fixPaths" - what in the raw resources needs to be added by user
	ControlID         string                      `json:"controlID"`
	FrameworkName     string                      `json:"frameworkName"`
	ControlStatus     int                         `json:"controlStatus"` // it's rather resource status within the control, control might fail but on this specific resource it might be passed (exception)
	ControlStatusText string                      `json:"controlStatusText"`
	RelatedExceptions []PostureExceptionPolicy    `json:"relatedExceptions"` // configured in portal
	ExceptionApplied  []PostureExceptionPolicy    `json:"exceptionApplied"`  //actual ruleResponse
	ResourceKind      string                      `json:"kind"`
	ResourceNamespace string                      `json:"namespace"`
	Remediation       string                      `json:"remediation"`
	Images            []PostureContainerSummary   `json:"containers,omitempty"`
	DeleteStatus      RecordStatus                `json:"deletionStatus,omitempty"`
	Recommendations   []RecommendationAssociation `json:"recommendations"`

	Timestamp time.Time `json:"timestamp"`
	ReportID  string    `json:"reportGUID"`
}

type HighlightsByControl struct {
	ControlID  string    `json:"controlID"`
	Highlights []string  `json:"highlights"`
	FixPaths   []FixPath `json:"fixPaths"`
	FixCommand string    `json:"fixCommand"`
}

type PostureResourceSummary struct {
	Designators PortalDesignator `json:"designators"`
	Name        string           `json:"name"`       // wlid/sid and etc.
	ResourceID  string           `json:"resourceID"` //as given by kscape

	//gives upto PostureResourceMaxCtrls controls as an example
	FailedControl   []string `json:"failedControls"`
	WarningControls []string `json:"warningControls"` // Deprecated
	SkippedControls []string `json:"skippedControls"`
	//maps statusText 2 list of controlIDs
	StatusToControls map[string][]string `json:"statusToControls"`

	HighlightsPerCtrl []HighlightsByControl `json:"highlightsPerControl"`

	//totalcount (including the failed/warning controls slices)
	FailedControlCount     int                         `json:"failedControlsCount"`
	SkippedControlCount    int                         `json:"skippedControlsCount"`
	WarningControlCount    int                         `json:"warningControlsCount"` // Deprecated
	Status                 int                         `json:"status"`
	StatusText             string                      `json:"statusText"`
	Remediation            []string                    `json:"remediation"`
	ResourceKind           string                      `json:"resourceKind"`
	FrameworkName          string                      `json:"frameworkName"`
	ExceptionRecommendaion string                      `json:"exceptionRecommendaion"`
	RelatedExceptions      []PostureExceptionPolicy    `json:"relatedExceptions"` // configured in portal
	ExceptionApplied       []PostureExceptionPolicy    `json:"exceptionApplied"`  //actual ruleResponse
	Images                 []PostureContainerSummary   `json:"containers,omitempty"`
	Recommendations        []RecommendationAssociation `json:"recommendations"`

	Timestamp     time.Time    `json:"timestamp"`
	ReportID      string       `json:"reportGUID"`
	DeleteStatus  RecordStatus `json:"deletionStatus,omitempty"`
	ArmoBestScore int64        `json:"armoBestScore"`

	// Information about the controls that were run on this entity
	// The key is the status of the control (`failed`, `passed`, etc)
	ControlsInfo map[string][]ControlInfo `json:"controlsInfo"`

	// Counters - Failed controls by severity
	CriticalSeverityControls int `json:"criticalSeverityControls"`
	HighSeverityControls     int `json:"highSeverityControls"`
	MediumSeverityControls   int `json:"mediumSeverityControls"`
	LowSeverityControls      int `json:"lowSeverityControls"`
}

type PostureAttributesList struct {
	Attribute string   `json:"attributeName"`
	Values    []string `json:"values"`
}

// --------/api/v1/posture/summary
type PostureSummary struct {
	RuntimeImprovementPercentage float32               `json:"runtimeImprovementPercentage"`
	LastRun                      time.Time             `json:"lastRun"`
	ReportID                     string                `json:"reportGUID"`
	Designators                  PortalDesignator      `json:"designators"`
	PostureAttributes            PostureAttributesList `json:"postureAttributes"`
	ClusterCloudProvider         string                `json:"clusterCloudProvider"`

	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}
type PosturePaths struct {
	// must have FailedPath or FixPath, not both
	FailedPath string  `json:"failedPath,omitempty"`
	FixPath    FixPath `json:"fixPath,omitempty"`
	FixCommand string  `json:"fixCommand,omitempty"`
}
type FixPath struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}
type PostureReportResultRaw struct {
	Designators           PortalDesignator `json:"designators"`
	Timestamp             time.Time        `json:"timestamp"`
	ReportID              string           `json:"reportGUID"`
	ResourceID            string           `json:"resourceID"`
	ControlID             string           `json:"controlID"`
	ControlConfigurations []ControlInputs  `json:"controlConfigurations,omitempty"`
	HighlightsPaths       []PosturePaths   `json:"highlightsPaths"`
}
type RawResource struct {
	Designators  PortalDesignator `json:"designators"`
	Timestamp    time.Time        `json:"timestamp"`
	DeleteStatus RecordStatus     `json:"deletionStatus,omitempty"`

	ResourceID          string                    `json:"resourceID"`
	PostureReportID     string                    `json:"postureReportID,omitempty"`
	SPIFFE              string                    `json:"spiffe"`
	Containers          []PostureContainerSummary `json:"containers,omitempty"`
	RelatedResourcesIDs []string                  `json:"relatedResourcesID,omitempty"`
	RAW                 json.RawMessage           `json:"object"`
}

type PostureJobParams struct {
	Name            string `json:"name,omitempty"`
	ID              string `json:"id,omitempty"`
	ClusterName     string `json:"clusterName"`
	FrameworkName   string `json:"frameworkName"`
	CronTabSchedule string `json:"cronTabSchedule,omitempty"`
	JobID           string `json:"jobID,omitempty"`
}

// ControlInfo Basic information about a control
type ControlInfo struct {

	// ID of the control
	// Example: C-0034
	ID string `json:"id"`

	// How much this control is critical
	// Example: 6
	BaseScore float32 `json:"baseScore"`

	// How many failed resources for this control
	// Example: 3
	FailedResources int `json:"failedResources"`
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
