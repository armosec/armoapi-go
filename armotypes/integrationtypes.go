package armotypes

import "fmt"

type TicketManager string

const (
	TicketManagerJira TicketManager = "jira"
)

type Ticket struct {
	GUID          string              `json:"guid,omitempty"`      //ticket guid in armo
	TicketManager TicketManager       `json:"ticketManager"`       //ticket service provider
	Owner         EntityIdentifiers   `json:"ownerID,omitempty"`   //armo entity that owns the ticket
	Subjects      []EntityIdentifiers `json:"subjects,omitempty"`  //armo entities mentioned in the ticket
	Link          string              `json:"link,omitempty"`      //link to the ticket
	Status        string              `json:"status,omitempty"`    //status of the ticket
	LinkTitle     string              `json:"linkTitle,omitempty"` //title of the ticket
	Severity      string              `json:"severity,omitempty"`  //severity of the ticket
	Error         string              `json:"error,omitempty"`     //error message if any
	ErrorCode     int                 `json:"errorCode,omitempty"` //error code if any (e.g. http status code like 401)
}

// EntityIdentifiers is a struct that holds the identifiers of an entity (hard typed designators)
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
