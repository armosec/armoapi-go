package armotypes

type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"
	AlertCritical AlertLevel = "critical"
	AlertError    AlertLevel = "error"
)

type SlackSettings struct {
	Token         string `json:"token"`
	Alert2Channel `json:",inline,omitempty"`
	Notifications `json:"notifications,omitempty"`
}

type Alert2Channel struct {
	Critical []SlackChannel `json:"criticalChannels,omitempty"`
	Error    []SlackChannel `json:"errorChannels,omitempty"`
	Info     []SlackChannel `json:"infoChannels,omitempty"`
}

type SlackChannel struct {
	ChannelID   string     `json:"channelID"`
	ChannelName string     `json:"channelName"`
	AlertLevel  AlertLevel `json:"alertLevel"`
}

type SlackNotification struct {
	IsActive   bool                   `json:"isActive"`
	Channels   []SlackChannel         `json:"channels"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Notifications struct {
	PostureScan               []string `json:"postureScan,omitempty"` // bad approach kept till i see if can do something with mongo and old data
	PostureScoreAboveLastScan []string `json:"postureScoreAboveLastScan,omitempty"`

	PostureScanV1              []SlackNotification `json:"postureScanV1"`
	PostureScanAboveLastScanV1 []SlackNotification `json:"postureScoreAboveLastScanV1"`
	// PostureScanThresholdV1     []SlackNotification `json:"postureScanThreshold"`
}
