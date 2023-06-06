package postgresmodels

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// TODO: add explicit column names, add validation

type Cluster struct {
	BaseModel
	GUID          string `gorm:"primaryKey"`
	CustomerGuid  string
	Location      string
	CloudProvider string
	FullName      string
	ShortName     string
	PrefixName    string
}

type ClusterPostureReport struct {
	ReportGUID          string `gorm:"primaryKey"`
	ClusterGUID         string
	Cluster             Cluster `gorm:"foreignKey:ClusterGUID"`
	Score               float32
	WorkerNodeCount     int
	KubescapeVersion    string
	KubernetesVersion   string
	HelmVersion         string
	RegoLibraryVersion  string
	ControlsSummary     ControlsSummary      `gorm:"foreignKey:ReportGUID"`
	ControlsScanResults []ControlScanResult  `gorm:"foreignKey:ReportGUID"`
	FrameworksSummary   []FrameworkSummary   `gorm:"foreignKey:ReportGUID"`
	ResourceScanResults []ResourceScanResult `gorm:"foreignKey:ReportGUID"`
}

type Control struct {
	BaseModel
	ControlID          string `gorm:"primaryKey"`
	RegoLibraryVersion string `gorm:"primaryKey"`
	Name               string //TODO: primary key?
	Description        string
	Remediation        string
	BaseScore          float32
}

type ControlScanResult struct {
	BaseModel
	ControlID          string    `gorm:"primaryKey"`
	ReportGUID         string    `gorm:"primaryKey"`
	FrameworkGUID      string    `gorm:"primaryKey"`
	Framework          Framework `gorm:"foreignKey:FrameworkGUID"`
	Control            Control   `gorm:"foreignKey:ControlID,RegoLibraryVersion"`
	RegoLibraryVersion string
	Status             string
	SubStatus          string
}

type ControlsSummary struct {
	ReportGUID               string `gorm:"primaryKey"`
	TotalControls            int
	FailedControls           int
	WarningControls          int
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

type Framework struct {
	BaseModel
	GUID         string `gorm:"primaryKey"`
	Name         string
	CustomerGuid string
}

type FrameworkControl struct {
	BaseModel
	ControlID          string    `gorm:"primaryKey"`
	FrameworkGUID      string    `gorm:"primaryKey"`
	Framework          Framework `gorm:"foreignKey:FrameworkGUID"`
	RegoLibraryVersion string
	Control            Control `gorm:"foreignKey:ControlID,RegoLibraryVersion"`
}

type FrameworkSummary struct {
	BaseModel
	ReportGUID       string    `gorm:"primaryKey"`
	FrameworkGUID    string    `gorm:"primaryKey"`
	Framework        Framework `gorm:"foreignKey:FrameworkGUID"`
	RiskScore        float32
	ImprovementScore float32
}

type ResourceScanResult struct {
	BaseModel
	ResourceID         string `gorm:"primaryKey"`
	ReportGUID         string `gorm:"primaryKey"`
	FrameworkGUID      string `gorm:"primaryKey"`
	ControlID          string `gorm:"primaryKey"`
	RegoLibraryVersion string
	Control            Control   `gorm:"foreignKey:ControlID,RegoLibraryVersion"`
	Resource           Resource  `gorm:"foreignKey:ResourceID,ReportGUID"`
	Framework          Framework `gorm:"foreignKey:FrameworkGUID"`
	Status             string
	ExceptionApplied   datatypes.JSON
	FixPaths           datatypes.JSON
}

type Designators struct {
}

/*
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
Ecosystem        Ecosystem `json:"ecosystem" sql:"type:ENUM('production', 'testsystem')"`*/
