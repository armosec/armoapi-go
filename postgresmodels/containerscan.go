package postgresmodels

import (
	"time"

	"github.com/lib/pq"
)

// TODO: add explicit column names, add validation

type Vulnerability struct {
	BaseModel
	Name          string `gorm:"primaryKey"`
	Severity      string
	SeverityScore int
	IsRCE         bool
	Links         pq.StringArray `gorm:"type:text[]"`
	Description   string
}

type VulnerabilityFinding struct {
	BaseModel
	VulnerabilityName string        `gorm:"primaryKey"`
	Vulnerability     Vulnerability `gorm:"foreignKey:VulnerabilityName"`
	ImageScanId       string        `gorm:"primaryKey"`
	Component         string        `gorm:"primaryKey"`
	ComponentVersion  string        `gorm:"primaryKey"`
	LayerHash         string        `gorm:"primaryKey"`
	FixAvailable      *bool
	FixedInVersion    string
	LayerIndex        *int
	LayerCommand      string
	// TODO: add applied exceptions
}

type VulnerabilityScanSummary struct {
	BaseModel
	ScanKind        string
	ImageScanId     string `gorm:"primaryKey"`
	Timestamp       time.Time
	CustomerGuid    string
	Wlid            string
	Designators     JSONB
	ImageRegistry   string
	ImageRepository string
	ImageTag        string
	ImageHash       string
	JobIds          pq.StringArray `gorm:"type:text[]"`
	Status          string
	Errors          pq.StringArray         `gorm:"type:text[]"`
	Findings        []VulnerabilityFinding `gorm:"foreignKey:ImageScanId"`
	IsStub          *bool                  // if true, this is a stub scan summary, and the actual scan summary is not yet available. Should be deleted once we have the real one.
}