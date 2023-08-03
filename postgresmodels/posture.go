package postgresmodels

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// TODO: add explicit column names, add validation

type Cluster struct {
	BaseModel
	GUID         string `gorm:"primaryKey"`
	CustomerGuid string
	//Location      string
	CloudProvider string
	FullName      string
	ShortName     string
	PrefixName    string
}

type ClusterPostureReport struct {
	ReportGUID        string `gorm:"primaryKey"`
	ClusterGUID       string
	Cluster           Cluster `gorm:"foreignKey:ClusterGUID"`
	Score             float32
	Timestamp         time.Time
	WorkerNodeCount   int
	KubescapeVersion  string
	KubernetesVersion string
	HelmChartVersion  string
	//RegoLibraryVersion  string
	ControlsSummary     ControlsSummary      `gorm:"foreignKey:ReportGUID"`
	ControlsScanResults []ControlScanResult  `gorm:"foreignKey:ReportGUID"`
	FrameworksSummary   []FrameworkSummary   `gorm:"foreignKey:ReportGUID"`
	ResourceScanResults []ResourceScanResult `gorm:"foreignKey:ReportGUID"`
}

/*
type Control struct {
	BaseModel
	ControlID          string `gorm:"primaryKey"`
	RegoLibraryVersion string `gorm:"primaryKey"`
	Name               string //TODO: primary key?
	Description        string
	Remediation        string
	BaseScore          float32
}*/

type ControlScanResult struct {
	BaseModel
	ControlID                  string `gorm:"primaryKey"`
	ReportGUID                 string `gorm:"primaryKey"`
	FrameworkName              string `gorm:"primaryKey"`
	Name                       string
	Status                     string
	SubStatus                  string
	StatusCode                 int
	ComplianceScore            float32
	AffectedResourcesCount     int
	FailedResourcesCount       int
	SkippedResourcesCount      int
	WarningResourcesCount      int
	TotalScannedResourcesCount int

	//HighlightPathsCount        int64
	//ControlInputs                  []ControlInputs
	//FrameworkSubSectionID          []string
}

type ControlsSummary struct {
	BaseModel
	BaseReport                      //we keep it here and not on the cluster summary for security framework scans
	ReportGUID               string `gorm:"primaryKey"`
	TotalControls            int
	FailedControls           int
	SkippedControls          int
	CriticalSeverityControls int
	HighSeverityControls     int
	MediumSeverityControls   int
	LowSeverityControls      int
}

type FrameworkSummary struct {
	BaseModel
	ReportGUID      string `gorm:"primaryKey"`
	FrameworkName   string `gorm:"primaryKey"`
	ComplianceScore float32
	TotalControls   int
	FailedControls  int
	SkippedControls int
	TypeTags        []string `gorm:"type:text[]"`
	//Designators      PortalDesignator `json:"designators"`
}

type Resource struct {
	BaseModel
	ResourceID string `gorm:"primaryKey"`
	ReportGUID string `gorm:"primaryKey"`
	//Designators                datatypes.JSON
	EventId           []byte
	Kind              string
	Name              string
	Namespace         string
	RelatedNames      pq.StringArray `gorm:"type:text[]"`
	RelatedKinds      pq.StringArray `gorm:"type:text[]"`
	RelatedNamespaces pq.StringArray `gorm:"type:text[]"`
}

type ResourceScanResult struct {
	BaseModel
	ResourceID      string         `gorm:"primaryKey"`
	ReportGUID      string         `gorm:"primaryKey"`
	FrameworkName   string         `gorm:"primaryKey"`
	Resource        Resource       `gorm:"foreignKey:ResourceID,ReportGUID"`
	FailedControl   pq.StringArray `gorm:"type:text[]"`
	WarningControls pq.StringArray `gorm:"type:text[]"`
	SkippedControls pq.StringArray `gorm:"type:text[]"`
	//maps statusText 2 list of controlIDs
	StatusToControls         datatypes.JSON `gorm:"type:json"`
	HighlightsPerCtrl        datatypes.JSON `gorm:"type:json"`
	FailedControlCount       int
	SkippedControlCount      int
	WarningControlCount      int
	Status                   int
	StatusText               string
	SubStatusText            string
	ExceptionApplied         datatypes.JSON `gorm:"type:json"`
	Images                   datatypes.JSON `gorm:"type:json"`
	ControlsInfo             datatypes.JSON `json:"controlsInfo"`
	CriticalSeverityControls int
	HighSeverityControls     int
	MediumSeverityControls   int
	LowSeverityControls      int
}
