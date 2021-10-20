package armotypes

type SlackSettings struct {
	Token         string `json:"token"`
	Alert2Channel `json:",inline"`
	Notifications `json:"notifications"`
}

type Alert2Channel struct {
	Critical []SlackChannel `json:"criticalChannels"`
	Error    []SlackChannel `json:"errorChannels"`
	Info     []SlackChannel `json:"infoChannels"`
}

type SlackChannel struct {
	ChannelID   string `json:"channelID"`
	ChannelName string `json:"channelName"`
}

type Notifications struct {
	PostureScan               []string `json:"postureScan"`
	PostureScoreAboveLastScan []string `json:"postureScoreAboveLastScan"`
}
