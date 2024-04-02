package notifications

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

// Config option type
// swagger:model CollaborationConfigOptionType
type CollaborationConfigOptionType struct {
	// Name of the type
	// Example: project
	Name string `json:"name" bson:"name,omitempty"`

	// Indicates if this option is a mandatory for collaboration configuration
	// Example: false
	ConfigRequired bool `json:"required" bson:"required"`

	// Indicates if this option is a mandatory for sharing
	// Example: true
	ShareRequired bool `json:"-"`
	//TODO set back to `json:"shareRequired"` after updating the schema in the portal

	// Custom input available or not
	// Example: false
	CustomInput bool `json:"customInput" bson:"customInput"`
}

// Collaboration provider config option
// swagger:model CollaborationConfigOption
type CollaborationConfigOption struct {
	// Type of the option
	// Example: Project
	Type *CollaborationConfigOptionType `json:"type,omitempty" bson:"type,omitempty"`

	// Name of the option
	// Example: jira-main-project
	Name string `json:"name" bson:"name,omitempty"`

	// ID of the option
	// Example: 8313c5a0-bee1-4a3c-8f4f-71ce698259876 or https://teams/mywebhook
	ID string `json:"id" bson:"id,omitempty"`

	// Icon url for the option. Optional
	// Example: https://site-admin-avatar-cdn.prod.public.atl-paas.net/avatars/240/triangle.png
	IconURL string `json:"iconURL,omitempty" bson:"iconURL,omitempty"`

	// Icon for the option encoded in base64. Optional
	IconBase64 string `json:"iconBase64,omitempty" bson:"iconBase64,omitempty"`
}

type ChannelProvider string

const (
	CollaborationTypeJira  ChannelProvider = "jira"
	CollaborationTypeSlack ChannelProvider = "slack"
	CollaborationTypeTeams ChannelProvider = "teams"
	CollaborationTypeEmail ChannelProvider = "email"
)

// swagger:model CollaborationConfig
type CollaborationConfig struct {
	armotypes.PortalBase `json:",inline" bson:",inline"`

	// Provider name
	// Example: jira
	Provider ChannelProvider `json:"provider,omitempty" bson:"provider,omitempty"`

	// Host name for private hosting
	// Example: http://example.com
	HostName string `json:"hostName,omitempty" bson:"hostName,omitempty"`

	// The context of sharing (for example in jira it will be cloud, project, etc)
	Context map[string]CollaborationConfigOption `json:"context" bson:"context,omitempty"`

	// Icon url for the option. Optional
	// Example: https://site-admin-avatar-cdn.prod.public.atl-paas.net/avatars/240/triangle.png
	IconURL string `json:"iconURL,omitempty" bson:"iconURL,omitempty"`

	// Icon for the option encoded in base64. Optional
	IconBase64 string `json:"iconBase64,omitempty" bson:"iconBase64,omitempty"`

	CreationTime string `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
}

type IntegrationConnectionStatus string

const (
	Connected    IntegrationConnectionStatus = "connected"
	Disconnected IntegrationConnectionStatus = "disconnected"
)

type IntegrationsConnectionStatus struct {
	Provider ChannelProvider             `json:"provider"`
	Status   IntegrationConnectionStatus `json:"status"`
}

type ReferenceType string

const (
	ReferenceTypeClusterControlTicket    ReferenceType = "ticket:cluster:control"
	ReferenceTypeRepositoryControlTicket ReferenceType = "ticket:repository:control"
	ReferenceTypeImageTicket             ReferenceType = "ticket:image"
	ReferenceTypeVulnerabilityTicket     ReferenceType = "ticket:vulnerability"
)

// Referance to external integration (e.g link to jira ticket)
type IntegrationReference struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	//integration provider (e.g jira, slack, teams)
	Provider ChannelProvider `json:"provider,omitempty" bson:"provider,omitempty"`
	//integration provider data (e.g jira ticket data)
	ProviderData map[string]interface{} `json:"providerData,omitempty" bson:"providerData,omitempty"`
	//type of the reference (e.g tickets kind)
	Type ReferenceType `json:"type,omitempty" bson:"type,omitempty"`
	//owner identifiers of this reference (e.g resourceHash, wlid)
	Owner map[string]string `json:"owner,omitempty" bson:"owner,omitempty"`
	//related entities identifiers of this reference (e.g cves, controls)
	RelatedObjects []map[string]string `json:"relatedObjects,omitempty" bson:"relatedObjects,omitempty"`
	//creation time of the reference
	CreationTime time.Time `json:"creationTime" bson:"creationTime"`
}
