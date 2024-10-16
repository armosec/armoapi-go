package armotypes

// cloud severities
var CloudSeverityToInt = map[string]int{
	"critical": 500,
	"high":     400,
	"medium":   300,
	"low":      200,
	"info":     100,
	"none":  0,
}

var CloudIntToSeverity = map[int]string{
	UnknownScore:  "none",
	InfoScore:     "info",
	LowScore:      "low",
	MediumScore:   "medium",
	HighScore:     "high",
	CriticalScore: "critical",
}

const (
	UnknownScore  = 0
	InfoScore     = 100
	LowScore      = 200
	MediumScore   = 300
	HighScore     = 400
	CriticalScore = 500
)

// cloud check statuses
const (
	CloudCheckStatusEmpty   = "EMPTY"
	CloudCheckStatusFail    = "FAIL"
	CloudCheckStatusManual  = "MANUAL"
	CloudCheckStatusPass    = "PASS"
	CloudCheckStatusSkipped = "SKIPPED"

	CloudAutomatedCheckType = "AUTOMATED"
	CloudManualCheckType    = CloudCheckStatusManual
	CloudManualAndAutomated = CloudAutomatedCheckType + "/" + CloudManualCheckType
)

var CloudCheckStatusToInt = map[string]int{
	CloudCheckStatusEmpty:   -1,
	CloudCheckStatusFail:    10,
	CloudCheckStatusManual:  20,
	CloudCheckStatusPass:    30,
	CloudCheckStatusSkipped: 40,
}

var CloudIntToCheckStatus = map[int]string{
	-1: CloudCheckStatusEmpty,
	10: CloudCheckStatusFail,
	20: CloudCheckStatusManual,
	30: CloudCheckStatusPass,
	40: CloudCheckStatusSkipped,
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
