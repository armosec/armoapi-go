package postgresmodels

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type ClusterPostureReport struct {
	BaseModel
	ReportGUID               string `gorm:"primaryKey"`
	ClusterGUID              string
	ClusterName              string
	CustomerGUID             string
	Score                    float32
	Timestamp                time.Time
	WorkerNodeCount          int
	KubescapeVersion         string
	KubernetesVersion        string
	HelmChartVersion         string
	RegoLibraryVersion       string
	TotalControls            int
	FailedControls           int
	SkippedControls          int
	CriticalSeverityControls int
	HighSeverityControls     int
	MediumSeverityControls   int
	LowSeverityControls      int
}

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
}

// We need this table for quicker queries although it could be calculated from ControlScanResult
type FrameworkSummary struct {
	BaseModel
	ReportGUID      string `gorm:"primaryKey"`
	FrameworkName   string `gorm:"primaryKey"`
	ComplianceScore float32
	TotalControls   int
	FailedControls  int
	SkippedControls int
	TypeTags        pq.StringArray `gorm:"type:text[]"`
}

type Resource struct {
	BaseModel
	ResourceID        string         `gorm:"primaryKey"`
	ReportGUID        string         `gorm:"primaryKey"`
	Designators       datatypes.JSON `gorm:"type:json"` //Portal designators
	ResourceObjectRef string         //external storage ref(e.g. S3 bucket:key) to the resource file
}

type ResourceFixPath struct {
	ResourceID string `gorm:"primaryKey"`
	ReportGUID string `gorm:"primaryKey"`
	ControlID  string `gorm:"primaryKey"`
	FailedPath string `gorm:"primaryKey"`
	FixCommand string `gorm:"primaryKey"`
	FixPath    string `gorm:"primaryKey"`
	FixValue   string
}

type ResourceControlResult struct {
	BaseModel
	ResourceID          string   `gorm:"primaryKey"`
	ReportGUID          string   `gorm:"primaryKey"`
	FrameworkName       string   `gorm:"primaryKey"`
	ControlID           string   `gorm:"primaryKey"`
	Resource            Resource `gorm:"foreignKey:ResourceID,ReportGUID"`
	StatusCode          int
	StatusText          string
	SubStatusText       string
	IgnoreRulesIDs      pq.StringArray `gorm:"type:text[]"`
	SystemRulesNames    pq.StringArray `gorm:"type:text[]"`
	RelatedResourcesIDs pq.StringArray `gorm:"type:text[]"`
}

type ResourceContainer struct {
	BaseModel
	ResourceID    string `gorm:"primaryKey"`
	ReportGUID    string `gorm:"primaryKey"`
	ContainerName string `gorm:"primaryKey"`
	Image        string
	ImageHash     string
}
