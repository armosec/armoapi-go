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
	Scope []EnrichedScope `json:"scope"`

	// Jira ticket identifiers for automatic ticket creation
	// Example: [{collaborationGUID: "abc", siteId: "xyz", projectId: "PROJ", issueTypeId: "10001"}]
	JiraTicketIdentifiers []JiraTicketIdentifiers `json:"jiraTicketIdentifiers,omitempty"`

	// Linear ticket identifiers for automatic ticket creation
	// Example: [{workspaceId: "abc", teamId: "xyz"}]
	LinearTicketIdentifiers []LinearTicketIdentifiers `json:"linearTicketIdentifiers,omitempty"`
}
