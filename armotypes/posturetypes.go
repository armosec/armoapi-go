package armotypes

import (
	"encoding/json"
	"time"
)

const (
	PostureControlStatusPassed = iota + 1
	PostureControlStatusWarning
	PostureControlStatusFailed

	PostureResourceMaxCtrls = 4
)

//-------- /api/v1/posture/clustersOvertime response datastructures
type PostureClusterOverTime struct {
	Designators  PortalDesignator           `json:"designators,omitempty"`
	ClusterName  string                     `json:"clusterName"`
	Frameworks   []PostureFrameworkOverTime `json:"frameworks"`
	DeleteStatus RecordStatus               `json:"deletionStatus,omitempty"`
}

//Used for elastic
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
	WarningControls  int              `json:"warningControls"`
	ReportID         string           `json:"reportGUID"`
	Designators      PortalDesignator `json:"designators"`

	Timestamp    time.Time    `json:"timestamp"`
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}

type PostureContainerSummary struct {
	ContainerName string `json:"containerName"`
	ImageTag      string `json:"image,omitempty"`
}

//----/api/v1/posture/controls
type PostureControlSummary struct {
	Designators                    PortalDesignator `json:"designators"`
	ControlID                      string           `json:"id"` // "C0001"
	ControlGUID                    string           `json:"guid"`
	Name                           string           `json:"name"`
	AffectedResourcesCount         int              `json:"affectedResourcesCount"`
	FailedResourcesCount           int              `json:"failedResourcesCount"`
	WarningResourcesCount          int              `json:"warningResourcesCount"`
	PreviousAffectedResourcesCount int              `json:"previousAffectedResourcesCount"`
	PreviousFailedResourcesCount   int              `json:"previousFailedResourcesCount"`
	PreviousWarningResourcesCount  int              `json:"previousWarningResourcesCount"`
	Framework                      string           `json:"frameworkName"`
	Remediation                    string           `json:"remediation"`
	Status                         int              `json:"status"`
	StatusText                     string           `json:"statusText"`
	Description                    string           `json:"description"`
	Section                        string           `json:"section"`
	Timestamp                      time.Time        `json:"timestamp"`
	ReportID                       string           `json:"reportGUID"`
	DeleteStatus                   RecordStatus     `json:"deletionStatus,omitempty"`
	Score                          float32          `json:"score"`
	ScoreFactor                    float32          `json:"baseScore"`
	ARMOImprovement                float32          `json:"ARMOimprovement"`
}

//---------/api/v1/posture/resources

//1 resource per 1 control
type PostureResource struct {
	UniqueResourceResult string                      `json:"uniqueResourceResult"` // vnf(name + framework + reportID) to allow fast search for aggregation
	Designators          PortalDesignator            `json:"designators"`
	Name                 string                      `json:"name"` // wlid/sid and etc.
	ControlName          string                      `json:"controlName"`
	ControlID            string                      `json:"controlID"`
	FrameworkName        string                      `json:"frameworkName"`
	ControlStatus        int                         `json:"controlStatus"` // it's rather resource status within the control, control might fail but on this specific resource it might be warning
	ControlStatusText    string                      `json:"controlStatusText"`
	RelatedExceptions    []PostureExceptionPolicy    `json:"relatedExceptions"` // configured in portal
	ExceptionApplied     []PostureExceptionPolicy    `json:"exceptionApplied"`  //actual ruleResponse
	ResourceKind         string                      `json:"kind"`
	ResourceNamespace    string                      `json:"namespace"`
	Remediation          string                      `json:"remediation"`
	Images               []PostureContainerSummary   `json:"containers,omitempty"`
	DeleteStatus         RecordStatus                `json:"deletionStatus,omitempty"`
	Recommendations      []RecommendationAssociation `json:"recommendations"`

	Timestamp time.Time `json:"timestamp"`
	ReportID  string    `json:"reportGUID"`
}

type PostureResourceSummary struct {
	Designators PortalDesignator `json:"designators"`
	Name        string           `json:"name"` // wlid/sid and etc.

	//gives upto PostureResourceMaxCtrls controls as an example
	FailedControl   []string `json:"failedControls"` // failed+warning controls
	WarningControls []string `json:"warningControls"`

	//totalcount (including the failed/warning controls slices)
	FailedControlCount  int `json:"failedControlsCount"`
	WarningControlCount int `json:"warningControlsCount"`

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

	Timestamp    time.Time    `json:"timestamp"`
	ReportID     string       `json:"reportGUID"`
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}

//--------/api/v1/posture/summary
type PostureSummary struct {
	RuntimeImprovementPercentage float32          `json:"runtimeImprovementPercentage"`
	LastRun                      time.Time        `json:"lastRun"`
	ReportID                     string           `json:"reportGUID"`
	Designators                  PortalDesignator `json:"designators"`

	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`
}

type RawResource struct {
	Designators         PortalDesignator          `json:"designators"`
	ResourceID          string                    `json:"resourceID"`
	PostureReportID     string                    `json:"postureReportID,omitempty"`
	SPIFFE              string                    `json:"spiffe"`
	Containers          []PostureContainerSummary `json:"containers,omitempty"`
	RelatedResourcesIDs []string                  `json:"relatedResourcesID,omitempty"`
	RAW                 json.RawMessage           `json:"object"`
}
