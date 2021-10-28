package armotypes

type RecordStatus int

const (
	RecordAlive        RecordStatus = 0
	RecordShouldDelete RecordStatus = 1
)
