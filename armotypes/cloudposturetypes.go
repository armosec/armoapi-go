package armotypes

var SeverityToInt = map[string]int{
	"critical": 500,
	"high":     400,
	"medium":   300,
	"low":      200,
	"info":     100,
}

var CheckStatusToInt = map[string]int{
	"EMPTY":   -1,
	"MANUAL":  0,
	"FAIL":    1,
	"PASS":    2,
	"SKIPPED": 3,
}

var ScanStatusToInt = map[string]int{
	ScanFailed:     1,
	ScanInProgress: 2,
	ScanSuccess:    3,
}

const (
	ScanFailed     = "FAILED"
	ScanInProgress = "INPROGRESS"
	ScanSuccess    = "SUCCESS"
)
