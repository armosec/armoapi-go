package armotypes

// cloud severities
var CloudSeverityToInt = map[string]int{
	"critical": 500,
	"high":     400,
	"medium":   300,
	"low":      200,
	"info":     100,
}

var CloudIntToSeverity = map[int]string{
	UnknownScore:  "Unknown",
	InfoScore:     "info",
	LowScore:      "low",
	MediumScore:   "medium",
	HighScore:     "high",
	CriticalScore: "critical",
}

const (
	UnknownScore  = 1
	InfoScore     = 100
	LowScore      = 200
	MediumScore   = 300
	HighScore     = 400
	CriticalScore = 500
)

// cloud check statuses
var CloudCheckStatusToInt = map[string]int{
	"EMPTY":   -1,
	"MANUAL":  0,
	"FAIL":    1,
	"PASS":    2,
	"SKIPPED": 3,
}

var CloudIntToCheckStatus = map[int]string{
	-1: "EMPTY",
	0:  "MANUAL",
	1:  "FAIL",
	2:  "PASS",
	3:  "SKIPPED",
}

// cloud posture scans statuses
var CloudPostureScanStatusToInt = map[string]int{
	ScanFailed:     1,
	ScanInProgress: 2,
	ScanSuccess:    3,
}

var CloudPostureScanIntToStatus = map[int]string{
	ScanFailedScore:     ScanFailed,
	ScanInProgressScore: ScanInProgress,
	ScanSuccessScore:    ScanSuccess,
}

const (
	ScanFailed     = "FAILED"
	ScanInProgress = "INPROGRESS"
	ScanSuccess    = "SUCCESS"
)

const (
	ScanFailedScore     = 1
	ScanInProgressScore = 2
	ScanSuccessScore    = 3
)
