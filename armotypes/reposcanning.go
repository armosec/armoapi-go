package armotypes

import (
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

// Kind of an entity. Can only be one of the following: `file` or `repo`
// Example: repo
type RepoEntityKind string

const (
	RepoEntityFile RepoEntityKind = "file"
	RepoEntityRepo RepoEntityKind = "repo"
)

// RepoEntitySummary summary of repo scanning entity.
type RepoEntitySummary struct {
	Designators identifiers.PortalDesignator `json:"designators"`

	// Name of this entity
	// Example: "my-repo"
	Name string `json:"name"`

	Kind RepoEntityKind `json:"kind"`

	// Number of children of the entity. For `file`s entity it would be
	// the amount of the resources inside this file, and for `repo`s -
	// the amount of scanned files
	// Example: 13
	ChildCount uint64 `json:"childCount"`

	// Status of the entity
	// Example: failed
	StatusText string `json:"statusText"`

	// Information about the controls that were run on this entity
	// The key is the status of the control (`failed`, `passed`, etc)
	ControlsInfo map[string][]ControlInfo `json:"controlsInfo"`

	// Statistics about the controls that were run
	// The key is the status of the control (`failed`, `passed`, etc).
	// The value is the number of controls
	// Example: {"failed": 3, "passed": 4}
	ControlsStats map[string]int `json:"controlsStats"`

	// Frameworks that were run.
	// In multi-frameworks-summary, this property is
	// taking the place of the `framework` property
	// Example: ["ArmoBest", "MITRE"]
	Frameworks []string `json:"frameworks,omitempty"`

	// Single framework this summary is for.
	// Example: ArmoBest
	Framework string `json:"framework,omitempty"`

	// Time of the scan that produced this result
	Timestamp time.Time `json:"timestamp"`
	ReportID  string    `json:"reportGUID"`

	// swagger:ignore
	// This record is marked for deletion or not
	DeleteStatus RecordStatus `json:"deletionStatus,omitempty"`

	//tickets opened for in this entity (repository or repository file)
	TicketManager TicketManager `json:"ticketManager,omitempty"`
	Tickets       []Ticket      `json:"tickets,omitempty"`
}
