package armotypes

type ScanFrequency string

type CustomerConfig struct {
	Name       string                 `json:"name" bson:"name"`
	Attributes map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"` // could be string
	Scope      PortalDesignator       `json:"scope" bson:"scope"`
	Settings   Settings               `json:"settings" bson:"settings"`
}

type Settings struct {
	PostureControlInputs    map[string][]string     `json:"postureControlInputs" bson:"postureControlInputs"`
	PostureScanConfig       PostureScanConfig       `json:"postureScanConfig" bson:"postureScanConfig"`
	VulnerabilityScanConfig VulnerabilityScanConfig `json:"vulnerabilityScanConfig" bson:"vulnerabilityScanConfig"`
	SlackConfigurations     SlackSettings           `json:"slackConfigurations,omitempty" bson:"slackConfigurations,omitempty"`
}

type PostureScanConfig struct {
	ScanFrequency ScanFrequency `json:"scanFrequency,omitempty" bson:"scanFrequency,omitempty"`
}

type VulnerabilityScanConfig struct {
	ScanFrequency             ScanFrequency `json:"scanFrequency,omitempty" bson:"scanFrequency,omitempty"`
	CriticalPriorityThreshold int           `json:"criticalPriorityThreshold,omitempty" bson:"criticalPriorityThreshold,omitempty"`
	HighPriorityThreshold     int           `json:"highPriorityThreshold,omitempty" bson:"highPriorityThreshold,omitempty"`
	MediumPriorityThreshold   int           `json:"mediumPriorityThreshold,omitempty" bson:"mediumPriorityThreshold,omitempty"`
	ScanNewDeployment         bool          `json:"scanNewDeployment,omitempty" bson:"scanNewDeployment,omitempty"`
	AllowlistRegistries       []string      `json:"AllowlistRegistries,omitempty" bson:"AllowlistRegistries,omitempty"`
	BlocklistRegistries       []string      `json:"BlocklistRegistries,omitempty" bson:"BlocklistRegistries,omitempty"`
}
