package armotypes

import "time"

type ProfileKind string

const (
	ContainerProfileKind    ProfileKind = "ContainerProfile"
	TSContainerProfileKind  ProfileKind = "TSContainerProfile"
	ApplicationProfileKind  ProfileKind = "ApplicationProfile"
	NetworkNeighborhoodKind ProfileKind = "NetworkNeighborhood"
)

// ProfileScope identifies the platform (hostType) and location (cluster, namespace, awsAccountID, region, hostID) of a storage resource.
type ProfileScope struct {
	HostType     HostType `json:"hostType"`
	Cluster      string   `json:"cluster"`
	Namespace    string   `json:"namespace"`
	AWSAccountID string   `json:"awsAccountID"`
	Region       string   `json:"region"`
	HostID       string   `json:"hostID"`
}

// ProfileIdentifier uniquely identifies a profile resource by combining
// its scope with its name. Used for storage key building/parsing.
type ProfileIdentifier struct {
	ProfileScope
	Name string `json:"name"`
}

type TimeSeriesContainerProfileObject struct {
	CustomerGUID string `json:"customerGUID"`
	ProfileScope
	Name                    string `json:"name"`
	SeriesID                string `json:"seriesID"`
	TSSuffix                string `json:"tsSuffix"`
	ReportTimestamp         string `json:"reportTimestamp"`
	Status                  string `json:"status"`
	Completion              string `json:"completion"`
	PreviousReportTimestamp string `json:"previousReportTimestamp"`
	ResourceObjectRef       string `json:"resourceObjectRef"`
	HasData                 bool   `json:"hasData"`
}

// AgentsProfileObject represents a platform-agnostic storage resource.
type AgentsProfileObject struct {
	// Identity
	CustomerGUID string `json:"customerGUID"`
	ResourceHash string `json:"resourceHash"`
	Kind         string `json:"kind"`
	Name         string `json:"name"`

	// Scope (platform + location)
	ProfileScope

	// Resource metadata
	ResourceObjectRef string    `json:"resourceObjectRef"`
	ResourceVersion   string    `json:"resourceVersion,omitempty"`
	Checksum          string    `json:"checksum"`
	CreationTimestamp  time.Time `json:"creationTimestamp"`
	SyncKind          string    `json:"syncKind,omitempty"`
	APIVersion        string    `json:"apiVersion,omitempty"`

	// Related resource info
	RelatedName            string `json:"relatedName"`
	RelatedKind            string `json:"relatedKind"`
	RelatedResourceType    string `json:"relatedResourceType"`
	RelatedAPIGroup        string `json:"relatedAPIGroup"`
	RelatedNamespace       string `json:"relatedNamespace"`
	RelatedAPIVersion      string `json:"relatedAPIVersion"`
	RelatedResourceVersion string `json:"relatedResourceVersion"`

	// Status
	Status           string `json:"status"`
	CompletionStatus string `json:"completionStatus"`

	// Storage
	RelatedContainerProfiles map[string]string `json:"relatedContainerProfiles,omitempty"`
	AdditionalProps          map[string]string `json:"additionalProps,omitempty"`
	Containers               []string          `json:"containers,omitempty"`
	InitContainers           []string          `json:"initContainers,omitempty"`
	EphemeralContainers      []string          `json:"ephemeralContainers,omitempty"`
	ResourceSize             int               `json:"resourceSize"`
}
