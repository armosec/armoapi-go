package armotypes

type ScanFrequency string

type CustomerConfig struct {
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes,omitempty"` // could be string
	Scope      PortalDesignator       `json:"scope"`
	Settings   Settings               `json:"settings"`
}

type Settings struct {
	PostureControlInputs    map[string][]string     `json:"postureControlInputs"`
	PostureScanConfig       PostureScanConfig       `json:"postureScanConfig"`
	VulnerabilityScanConfig VulnerabilityScanConfig `json:"vulnerabilityScanConfig"`
	SlackConfigurations     SlackSettings           `json:"slackConfigurations,omitempty"`
}

type PostureScanConfig struct {
	ScanFrequency ScanFrequency `json:"scanFrequency,omitempty"`
}

type VulnerabilityScanConfig struct {
	ScanFrequency             ScanFrequency `json:"scanFrequency,omitempty"`
	CriticalPriorityThreshold int           `json:"criticalPriorityThreshold,omitempty"`
	HighPriorityThreshold     int           `json:"highPriorityThreshold,omitempty"`
	MediumPriorityThreshold   int           `json:"mediumPriorityThreshold,omitempty"`
	ScanNewDeployment         bool          `json:"scanNewDeployment,omitempty"`
	AllowlistRegistries       []string      `json:"AllowlistRegistries,omitempty"`
	BlocklistRegistries       []string      `json:"BlocklistRegistries,omitempty"`
}
