package armotypes

// cloud severities
var CloudSeverityToInt = map[string]int{
	"critical": 500,
	"high":     400,
	"medium":   300,
	"low":      200,
	"info":     100,
	"unknown":  1,
}

var CloudIntToSeverity = map[int]string{
	UnknownScore:  "unknown",
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
	"FAIL":    10,
	"MANUAL":  20,
	"PASS":    30,
	"SKIPPED": 40,
}

var CloudIntToCheckStatus = map[int]string{
	-1: "EMPTY",
	10:  "FAIL",
	20:  "MANUAL",
	30:  "PASS",
	40:  "SKIPPED",
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
