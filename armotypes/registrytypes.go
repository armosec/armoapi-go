package armotypes

import "time"

type RegistryJobParams struct {
	Name            string `json:"name,omitempty"`
	ID              string `json:"id,omitempty"`
	ClusterName     string `json:"clusterName"`
	RegistryName    string `json:"registryName"`
	CronTabSchedule string `json:"cronTabSchedule,omitempty"`
	JobID           string `json:"jobID,omitempty"`
}

type RegistryInfo struct {
	RegistryName     string     `json:"registryName,omitempty" bson:"registryName"`
	RegistryProvider string     `json:"registryProvider,omitempty" bson:"registryProvider"`
	RegistryToken    string     `json:"registryToken,omitempty" bson:"registryToken"`
	Depth            *int       `json:"depth,omitempty" bson:"depth"`
	Include          []string   `json:"include,omitempty" bson:"include"`
	Exclude          []string   `json:"exclude,omitempty" bson:"exclude"`
	Kind             string     `json:"kind,omitempty" bson:"kind"`
	IsHTTPS          *bool      `json:"isHTTPS,omitempty" bson:"isHTTPS"`
	SkipTLSVerify    *bool      `json:"skipTLSVerify,omitempty" bson:"skipTLSVerify"`
	AuthMethod       AuthMethod `json:"authMethod,omitempty" bson:"authMethod"`
	SecretName       string     `json:"secretName,omitempty" bson:"secretName"`
}

type AuthMethod struct {
	Username string `json:"username,omitempty" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	Type     string `json:"type,omitempty" bson:"type"`
}

type Repository struct {
	RepositoryName string `json:"repositoryName"`
}

type RegistryProvider int

const (
	Quay RegistryProvider = iota
	Harbor
)

type BaseContainerImageRegistry struct {
	PortalBase    `json:",inline" bson:"inline"`
	Provider      RegistryProvider `json:"provider" bson:"provider"`
	ClusterName   string           `json:"clusterName" bson:"clusterName"`
	Repositories  []string         `json:"repositories" bson:"repositories"`
	LastScan      *time.Time       `json:"lastScan,omitempty" bson:"lastScan,omitempty"`
	ScanFrequency string           `json:"scanFrequency,omitempty" bson:"scanFrequency,omitempty"`
	ResourceHash  string           `json:"resourceHash,omitempty" bson:"resourceHash,omitempty"`
	AuthID        string           `json:"authID,omitempty" bson:"authID"`
}

type QuayImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	ContainerRegistryName      string `json:"containerRegistryName"`
	RobotAccountName           string `json:"RobotAccountName"`
	RobotAccountToken          string `json:"RobotAccountToken"`
}

type HarborImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	InstanceURL                string `json:"instanceURL"`
	Username                   string `json:"username"`
	Password                   string `json:"password"`
}
