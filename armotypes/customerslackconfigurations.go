package armotypes

type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"
	AlertCritical AlertLevel = "critical"
	AlertError    AlertLevel = "error"
)

type SlackSettings struct {
	Token         string `json:"token" bson:"token"`
	Alert2Channel `json:",inline,omitempty" bson:"inline,omitempty"`
	Notifications `json:"notifications,omitempty" bson:"notifications,omitempty"`
}

type Alert2Channel struct {
	Critical []SlackChannel `json:"criticalChannels,omitempty" bson:"criticalChannels,omitempty"`
	Error    []SlackChannel `json:"errorChannels,omitempty" bson:"errorChannels,omitempty"`
	Info     []SlackChannel `json:"infoChannels,omitempty" bson:"infoChannels,omitempty"`
}

type SlackChannel struct {
	ChannelID   string     `json:"channelID" bson:"channelID"`
	ChannelName string     `json:"channelName" bson:"channelName"`
	AlertLevel  AlertLevel `json:"alertLevel" bson:"alertLevel"`
}

type SlackNotification struct {
	IsActive   bool                   `json:"isActive" bson:"isActive"`
	Channels   []SlackChannel         `json:"channels" bson:"channels"`
	Attributes map[string]interface{} `json:"attributes" bson:"attributes"`
}

type Notifications struct {
	PostureScan               []string `json:"postureScan,omitempty" bson:"postureScan,omitempty"` // bad approach kept till i see if can do something with mongo and old data
	PostureScoreAboveLastScan []string `json:"postureScoreAboveLastScan,omitempty" bson:"postureScoreAboveLastScan,omitempty"`

	PostureScanV1              []SlackNotification `json:"postureScanV1" bson:"postureScanV1"`
	PostureScanAboveLastScanV1 []SlackNotification `json:"postureScoreAboveLastScanV1" bson:"postureScoreAboveLastScanV1"`
	// PostureScanThresholdV1     []SlackNotification `json:"postureScanThreshold"`
}
