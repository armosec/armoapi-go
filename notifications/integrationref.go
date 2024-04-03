package notifications

import (
	"fmt"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

type ReferenceType string //type of the reference (e.g cve-ticket, slack-message etc)

const (
	//tickets types
	ReferenceTypeClusterControlTicket    ReferenceType = "ticket:cluster:control"
	ReferenceTypeRepositoryControlTicket ReferenceType = "ticket:repository:control"
	ReferenceTypeImageTicket             ReferenceType = "ticket:image"
	ReferenceTypeVulnerabilityTicket     ReferenceType = "ticket:vulnerability"
)

// Referance to external integration (e.g link to jira ticket)
type IntegrationReference struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	Provider             ChannelProvider     `json:"provider,omitempty" bson:"provider,omitempty"`             //integration provider (e.g jira, slack, teams)
	ProviderData         interface{}         `json:"providerData,omitempty" bson:"providerData,omitempty"`     //integration provider data (e.g jira ticket data)
	Type                 ReferenceType       `json:"type,omitempty" bson:"type,omitempty"`                     //type of the reference (e.g cve-ticket, slack-message etc)
	Owner                *EntityIdentifiers  `json:"owner,omitempty" bson:"owner,omitempty"`                   //owner identifiers of this reference (e.g resourceHash, wlid)
	RelatedObjects       []EntityIdentifiers `json:"relatedObjects,omitempty" bson:"relatedObjects,omitempty"` //related entities identifiers of this reference (e.g cves, controls)
	CreationTime         time.Time           `json:"creationTime" bson:"creationTime"`                         //creation time of the reference
}


// EntityIdentifiers is a struct that holds the identifiers of an entity
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

type EntityIdentifiers struct {
	Type EntityType `json:"type,omitempty" bson:"type,omitempty"`

	Cluster  string `json:"cluster,omitempty" bson:"cluster,omitempty"`
	RepoHash string `json:"repoHash,omitempty" bson:"repoHash,omitempty"`

	Namespace    string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Kind         string `json:"kind,omitempty" bson:"kind,omitempty"`
	ResourceHash string `json:"resourceHash,omitempty" bson:"resourceHash,omitempty"`
	ResourceID   string `json:"resourceID,omitempty" bson:"resourceID,omitempty"`

	CVEName          string `json:"cveName,omitempty" bson:"cveName,omitempty"`
	CVEID            string `json:"cveID,omitempty" bson:"cveID,omitempty"`
	Severity         string `json:"severity,omitempty" bson:"severity,omitempty"`
	SeverityScore    int    `json:"severityScore,omitempty" bson:"severityScore,omitempty"`
	Component        string `json:"component,omitempty" bson:"component,omitempty"`
	ComponentVersion string `json:"componentVersion,omitempty" bson:"componentVersion,omitempty"`

	ImageReposiotry string `json:"imageRepository,omitempty" bson:"imageRepository,omitempty"`
	LayerHash       string `json:"layerHash,omitempty" bson:"layerHash,omitempty"`

	ControlID string  `json:"controlID,omitempty" bson:"controlID,omitempty"`
	BaseScore float32 `json:"baseScore,omitempty" bson:"baseScore,omitempty"`
}


const (
	//IntegrationReference fileds names	
	AttributeProvider     = "provider"
	AttributeType         = "type"
	AttributeOwner        = "owner"
	AttributeRelated      = "relatedObjects"

	//Entity Identifiers fields names (that are not defined in identifiers.go)
	AttributeLayerHash        = "layerHash"
	AttributeImageRepository  = "imageRepository"
	AttributeResourceHash     = "resourceHash"
	AttributeResourceID       = "resourceID"
	AttributeComponentVersion = "componentVersion"
	AttributeComponent        = "component"
	AttributeSeverityScore    = "severityScore"
	AttributeSeverity         = "severity"
	AttributeCVEID            = "cveID"
	AttributeCVEName          = "cveName"
	AttributeControlID        = "controlID"
	AttributeBaseScore        = "baseScore"
)

func (e *EntityIdentifiers) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("entity type is required")
	}
	switch e.Type {
	case EntityTypePostureResource:
		if e.Cluster == "" || e.Namespace == "" || e.Name == "" || e.Kind == "" || e.ResourceHash == "" || e.ResourceID == "" {
			return fmt.Errorf("namespace, name, kind, resource hash, cluster and resource id are required for %s", e.Type)
		}
	case EntityTypeContainerScanWorkload:
		if e.Cluster == "" || e.Namespace == "" || e.Name == "" || e.Kind == "" || e.ResourceHash == "" {
			return fmt.Errorf("namespace, name, kind, resource hash and cluster are required for %s", e.Type)
		}
	case EntityTypeRepositoryResource:
		if e.RepoHash == "" || e.Namespace == "" || e.Name == "" || e.Kind == "" || e.ResourceID == "" {
			return fmt.Errorf("namespace, name, kind, resource hash, repo hash and resource id are required for %s", e.Type)
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
		if e.CVEName == "" || e.CVEID == "" || e.Severity == "" || e.SeverityScore == 0 || e.Component == "" || e.ComponentVersion == "" {
			return fmt.Errorf("cveName, cveID, severity, severity score, component and component version are required for %s", e.Type)
		}
	case EntityTypeControl:
		if e.ControlID == "" || e.Severity == "" || e.BaseScore == 0 {
			return fmt.Errorf("control id, severity and base score are required for %s", e.Type)
		}
	default:
		return fmt.Errorf("entity type %s is not supported", e.Type)
	}
	return nil
}
