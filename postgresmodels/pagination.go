package postgresmodels

type BaseReport struct {
	// TotalChunksExpected and ReceivedChunkNumbers are used to track the progress of the report.

	// Total number of chunks expected. Will be populated with the (ReportNumber of the LastReport + 1) (IsLastReport == true)
	// If not known yet (i.e. IsLastReport not received yet), will be set to -1
	TotalChunksExpected int

	// A list of the numbers of the chunks received so far - will be updated after each chunk is received.
	ReceivedChunkNumbers []int

	// set to True when TotalChunksExpected == ReceivedChunkNumbers
	Completed bool
}

type Reports struct {
	BaseModel
	ReportID     string `gorm:"primaryKey; not null"`
	CustomerGUID string
	BaseReport
}
