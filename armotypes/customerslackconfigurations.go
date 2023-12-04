package armotypes

type AlertLevel string

type SlackChannels struct {
	Channels []SlackChannel `json:"channels"`
}

type SlackChannel struct {
	ChannelID   string `json:"channelID"   bson:"channelID"`
	ChannelName string `json:"channelName" bson:"channelName"`
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
