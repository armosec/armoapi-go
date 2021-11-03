package armotypes

import "time"

type KPIPostureScan struct {
	Client           string    `json:"client"`
	ClientVersion    string    `json:"clientVersion"`
	Framework        string    `json:"framework"`
	FrameworkVersion string    `json:"frameworkVersion"`
	Timestamp        time.Time `json:"timestamp"`
	Target           string    `json:"target"` //yaml,helm,running - what we actually scanned
	ClientIP         string    `json:"clientIP"`
}

type KPILogin struct {
	CustomerGUID string    `json:"tennantGUID"`
	Timestamp    time.Time `json:"timestamp"`
	Username     string    `json:"username"`
	Email        string    `json:"e-mail"`
	IP           string    `json:"IP,omitempty"`
}
