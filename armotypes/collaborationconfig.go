package armotypes

// Config option type
// swagger:model CollaborationConfigOptionType
type CollaborationConfigOptionType struct {
	// Name of the type
	// Example: cloud
	Name string `json:"name"`

	// This is a mandatory option or not
	// Example: true
	Required bool `json:"required"`

	// Custom input available or not
	// Example: false
	CustomInput bool `json:"customInput"`
}

// Collaboration provider config option
// swagger:model CollaborationConfigOption
type CollaborationConfigOption struct {
	// Type of the option
	// Example: Project
	Type *CollaborationConfigOptionType `json:"type,omitempty"`

	// Name of the option
	// Example: jira-main-project
	Name string `json:"name"`

	// ID of the option
	// Example: 8313c5a0-bee1-4a3c-8f4f-71ce698259876
	ID string `json:"id"`

	// Icon url for the option. Optional
	// Example: https://site-admin-avatar-cdn.prod.public.atl-paas.net/avatars/240/triangle.png
	IconURL string `json:"iconURL,omitempty"`

	// Icon for the option encoded in base64. Optional
	IconBase64 string `json:"iconBase64,omitempty"`
}

// swagger:model CollaborationConfig
type CollaborationConfig struct {
	PortalBase `json:",inline"`

	// Provider name
	// Example: jira
	Provider string `json:"provider"`

	// Host name for private hosting
	// Example: http://example.com
	HostName string `json:"hostName,omitempty"`

	// The context of sharing (for example in jira it will be cloud, project, etc)
	Context map[string]CollaborationConfigOption `json:"context"`

	// Icon url for the option. Optional
	// Example: https://site-admin-avatar-cdn.prod.public.atl-paas.net/avatars/240/triangle.png
	IconURL string `json:"iconURL,omitempty"`

	// Icon for the option encoded in base64. Optional
	IconBase64 string `json:"iconBase64,omitempty"`
}
