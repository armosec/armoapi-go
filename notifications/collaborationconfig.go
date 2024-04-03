package notifications

import (
	"fmt"
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
	Provider             ChannelProvider `json:"provider,omitempty" bson:"provider,omitempty"`             //integration provider (e.g jira, slack, teams)
	ProviderData         interface{}     `json:"providerData,omitempty" bson:"providerData,omitempty"`     //integration provider data (e.g jira ticket data)
	Type                 ReferenceType   `json:"type,omitempty" bson:"type,omitempty"`                     //type of the reference (e.g cve-ticket, slack-message etc)
	Owner                *Entity         `json:"owner,omitempty" bson:"owner,omitempty"`                   //owner identifiers of this reference (e.g resourceHash, wlid)
	RelatedObjects       []Entity        `json:"relatedObjects,omitempty" bson:"relatedObjects,omitempty"` //related entities identifiers of this reference (e.g cves, controls)
	CreationTime         time.Time       `json:"creationTime" bson:"creationTime"`                         //creation time of the reference
}

type EntityType string

const (
	EntityTypePostureResource       EntityType = "postureResource"
	EntityTypeRepositoryResource    EntityType = "repositoryResource"
	EntityTypeContainerScanWorkload EntityType = "containerScanWorkload"
	EntityTypeImage                 EntityType = "image"
	EntityTypeImageLayer            EntityType = "imageLayer"
	EntityTypeVulanrability         EntityType = "vulnerability"
	EntityTypeControl               EntityType = "control"
)

type Entity struct {
	Type EntityType `json:"type,omitempty" bson:"type,omitempty"`

	Cluster  string `json:"cluster,omitempty" bson:"cluster,omitempty"`
	RepoHash string `json:"repoHash,omitempty" bson:"repoHash,omitempty"`

	Namespace    string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Kind         string `json:"kind,omitempty" bson:"kind,omitempty"`
	ResourceHash string `json:"resourceHash,omitempty" bson:"resourceHash,omitempty"`
	ResourceID   string `json:"resourceID,omitempty" bson:"resourceID,omitempty"`

	CVEID             string `json:"cveID,omitempty" bson:"cveID,omitempty"`
	VulnerabilityHash string `json:"vulnerabilityHash,omitempty" bson:"vulnerabilityHash,omitempty"`
	Severity          string `json:"severity,omitempty" bson:"severity,omitempty"`
	SeverityCode      int    `json:"severityCode,omitempty" bson:"severityCode,omitempty"`
	Component         string `json:"component,omitempty" bson:"component,omitempty"`
	ComponentVersion  string `json:"componentVersion,omitempty" bson:"componentVersion,omitempty"`

	ImageReposiotry string `json:"imageRepository,omitempty" bson:"imageRepository,omitempty"`
	LayerHash       string `json:"layerHash,omitempty" bson:"layerHash,omitempty"`

	ControlID string `json:"controlID,omitempty" bson:"controlID,omitempty"`
}

func (e *Entity) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("entity type is required")
	}
	switch e.Type {
	case EntityTypePostureResource, EntityTypeRepositoryResource, EntityTypeContainerScanWorkload:
		if e.Cluster == "" || e.Namespace == "" || e.Name == "" || e.Kind == "" || e.ResourceHash == "" || e.ResourceID == "" {
			return fmt.Errorf("namespace, name and kind are required for %s", e.Type)
		}
	case EntityTypeImage:
		if e.ImageReposiotry == "" {
			return fmt.Errorf("image repository is required for %s", e.Type)
		}
	case EntityTypeImageLayer:
		if e.ImageReposiotry == "" || e.LayerHash == "" {
			return fmt.Errorf("image repository and layer hash are required for %s", e.Type)
		}
	case EntityTypeVulanrability:
		if e.CVEID == "" || e.VulnerabilityHash == "" || e.Severity == "" || e.SeverityCode == 0 || e.Component == "" || e.ComponentVersion == ""  {
			return fmt.Errorf("cve, vulnerability hash, severity, component and component version are required for %s", e.Type)
		}
	case EntityTypeControl:
		if e.ControlID == "" || e.Severity == "" || e.SeverityCode == 0 {
			return fmt.Errorf("control id, severity and severity code are required for %s", e.Type)
		}
	default:
		return fmt.Errorf("entity type %s is not supported", e.Type)
	}

	// validate cluster and repository for relevant types
	switch e.Type {
	case EntityTypePostureResource, EntityTypeContainerScanWorkload:
		if e.ResourceHash == "" {
			return fmt.Errorf("resource hash is required for %s", e.Type)
		}
	case EntityTypeRepositoryResource:
		if e.RepoHash == "" {
			return fmt.Errorf("repo hash is required for %s", e.Type)
		}
	}
	return nil
}
