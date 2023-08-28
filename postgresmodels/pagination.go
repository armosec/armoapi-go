package postgresmodels

import "github.com/lib/pq"

type BaseReport struct {
	// TotalChunksExpected and TotalChunksRecieved are used to track the progress of the report.

	// Total number of chunks expected. Will be populated with the (ReportNumber of the LastReport + 1) (IsLastReport == true)
	// If not known yet (i.e. IsLastReport not recieved yet), will be set to -1
	TotalChunksExpected int

	//specify the total number of chunks recieved so far - will be increment by one on each chunk recieved.
	TotalChunksRecieved int

	// set to True when TotalChunksExpected == TotalChunksRecieved
	Completed bool
}

type ReportStatus struct {
	BaseModel
	ReportIdentifier      string `gorm:"primaryKey"` // reportGUID or imageScanID+timestamp
	TotalChunksExpected   int
	ReceivedChunksNumbers pq.Int64Array `gorm:"type:int[]"` // list of chunk numbers received so far
	Completed             bool
}