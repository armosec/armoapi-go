package postgresmodels

import "github.com/lib/pq"

type ReportStatus struct {
	BaseModel
	ReportIdentifier      string `gorm:"primaryKey"` // reportGUID or imageScanID+timestamp
	TotalChunksExpected   int
	ReceivedChunksNumbers pq.Int64Array `gorm:"type:int[]"` // list of chunk numbers received so far
	Completed             bool
}
