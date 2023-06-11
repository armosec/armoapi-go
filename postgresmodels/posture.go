package postgresmodels

import (
	"database/sql/driver"
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
	ControlGUID                string
	Name                       string
	Status                     string
	SubStatus                  string
	ComplianceScore            float32
	AffectedResourcesCount     int
	FailedResourcesCount       int
	SkippedResourcesCount      int
	TotalScannedResourcesCount int
}

type ControlsSummary struct {
	ReportGUID               string `gorm:"primaryKey"`
	TotalControls            int
	FailedControls           int
	SkippedControls          int
	CriticalSeverityControls int
	HighSeverityControls     int
	MediumSeverityControls   int
	LowSeverityControls      int
}

type Resource struct {
	BaseModel
	ResourceID        string `gorm:"primaryKey"`
	ReportGUID        string `gorm:"primaryKey"`
	EventId           []byte
	Kind              string
	Name              string
	Namespace         string
	RelatedNames      pq.StringArray `gorm:"type:text[]"`
	RelatedKinds      pq.StringArray `gorm:"type:text[]"`
	RelatedNamespaces pq.StringArray `gorm:"type:text[]"`
}

/*
type Framework struct {
	BaseModel
	GUID         string `gorm:"primaryKey"`
	Name         string
	CustomerGuid string
}*/

type FrameworkSummary struct {
	BaseModel
	ReportGUID    string `gorm:"primaryKey"`
	FrameworkName string `gorm:"primaryKey"`
	//Framework       Framework `gorm:"foreignKey:FrameworkGUID"`
	ComplianceScore float32
	TotalControls   int `json:"totalControls"`
	FailedControls  int `json:"failedControls"`
	SkippedControls int `json:"skippedControls,omitempty"`

	//Designators      PortalDesignator `json:"designators"`
}

type ResourceScanResult struct {
	BaseModel
	ResourceID    string `gorm:"primaryKey"`
	ReportGUID    string `gorm:"primaryKey"`
	FrameworkName string `gorm:"primaryKey"`
	ControlID     string `gorm:"primaryKey"`
	ControlName   string
	//RegoLibraryVersion string
	//Control            Control   `gorm:"foreignKey:ControlID,RegoLibraryVersion"`
	Resource Resource `gorm:"foreignKey:ResourceID,ReportGUID"`
	//Framework        Framework `gorm:"foreignKey:FrameworkGUID"`
	Status           string
	ExceptionApplied datatypes.JSON
	FixPaths         datatypes.JSON
}

type Designators struct {
}

type Ecosystem string

const (
	Production Ecosystem = "production"
	TestSystem Ecosystem = "testsystem"
)

func (e *Ecosystem) Scan(value interface{}) error {
	*e = Ecosystem(value.([]byte))
	return nil
}

func (e Ecosystem) Value() (driver.Value, error) {
	return string(e), nil
}

// model define
//Ecosystem        Ecosystem `json:"ecosystem" sql:"type:ENUM('production', 'testsystem')"`*/
