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
	"FAILED":     1,
	"INPROGRESS": 2,
	"SUCCESS":    3,
}
