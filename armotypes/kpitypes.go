package armotypes

import "time"

type KPIPostureScanWrapper struct {
	ReportType    string         `json:"reportType"`
	ReportPayload KPIPostureScan `json:"reportPayload"`
}

type KPIPostureScan struct {
	CustomerGUID string    `json:"customerGUID"`
	ReportGUID   string    `json:"reportGUID"`
	Frameworks   []string  `json:"frameworks"`
	Timestamp    time.Time `json:"timestamp"`
	ClusterName  string    `json:"clusterName"`
	SourceType   string    `json:"sourceType"`   //yaml,helm,running - what we actually scanned
	OutputFormat string    `json:"outputFormat"` // json,junit,prettyprint - output prefered by customers
	CICD         int       `json:"CICD"`         // 0 -false , 1- true , 2- unknown
	IP           string    `json:"IP,omitempty"`
}

type KPIPLoginWrapper struct {
	ReportType    string         `json:"reportType"`
	ReportPayload KPIPostureScan `json:"reportPayload"`
}

type KPILogin struct {
	CustomerGUID string    `json:"tennantGUID"`
	Timestamp    time.Time `json:"timestamp"`
	Username     string    `json:"username"`
	Email        string    `json:"e-mail"`
	IP           string    `json:"IP,omitempty"`
}
