package postgresmodels

type BaseReport struct {
	// TotalChunksExpected and TotalChunksRecieved are used to track the progress of the report.
	// Once TotalChunksExpected == TotalChunksRecieved, compeleted is set to True.

	//specify the total number of chunks expected. Will be populated with the ReportNumber of the LastReport (IsLastReport == true)
	TotalChunksExpected int

	//specify the total number of chunks recieved so far - will be increment by one on each chunk recieved.
	TotalChunksRecieved int

	Completed bool
}
