package notifications

// AlertChannelAPI An Alerting Channel configuration
// swagger:model AlertChannelAPI
type AlertChannelAPI struct {
	// Channel connection definition
	// Example: webhook connection
	Channel CollaborationConfig `json:"channel"`

	// Notifications configurations
	// Example: new cluster admin
	Notifications []AlertConfig `json:"notifications"`

	// Scope selected clusters/namespaces
	// Example cluster123, [nspace1, nspace2]
	Scope []AlertScope `json:"scope"`
}
