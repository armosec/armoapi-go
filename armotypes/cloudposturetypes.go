package armotypes

// cloud severities
var CloudSeverityToInt = map[string]int{
	"critical": 500,
	"high":     400,
	"medium":   300,
	"low":      200,
	"info":     100,
	"none":     0,
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
	CloudCheckStatusEmpty    = "EMPTY"
	CloudCheckStatusFail     = "FAIL"
	CloudCheckStatusManual   = "MANUAL"
	CloudCheckStatusPass     = "PASS"
	CloudCheckStatusSkipped  = "SKIPPED"
	CloudCheckStatusAccepted = "ACCEPTED"
)

var CloudCheckStatusToInt = map[string]int{
	CloudCheckStatusEmpty:    -1,
	CloudCheckStatusFail:     10,
	CloudCheckStatusManual:   20,
	CloudCheckStatusPass:     30,
	CloudCheckStatusSkipped:  40,
	CloudCheckStatusAccepted: 50,
}

var CloudIntToCheckStatus = map[int]string{
	-1: CloudCheckStatusEmpty,
	10: CloudCheckStatusFail,
	20: CloudCheckStatusManual,
	30: CloudCheckStatusPass,
	40: CloudCheckStatusSkipped,
	50: CloudCheckStatusAccepted,
}

// cloud check types
const (
	CloudEmptyCheckType     = "EMPTY"
	CloudAutomatedCheckType = "AUTOMATED"
	CloudManualCheckType    = CloudCheckStatusManual
	CloudManualAndAutomated = CloudAutomatedCheckType + "/" + CloudManualCheckType
)

var CloudCheckTypeToInt = map[string]int{
	CloudEmptyCheckType:     -1,
	CloudAutomatedCheckType: 10,
	CloudManualCheckType:    20,
	CloudManualAndAutomated: 30,
}

var CloudIntToCheckType = map[int]string{
	-1: CloudEmptyCheckType,
	10: CloudAutomatedCheckType,
	20: CloudManualCheckType,
	30: CloudManualAndAutomated,
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

type CSPMExceptionPolicy struct {
	BaseExceptionPolicy `json:",inline"`
	Name                string   `json:"name"`     // rule name
	Controls            []string `json:"controls"` // affected controls
	Severity            string   `json:"severity"`
	SeverityScore       int      `json:"severityScore"`
	RuleHash            string   `json:"ruleHash"`
}
