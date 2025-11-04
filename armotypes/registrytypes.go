package armotypes

import "time"

const (
	RegistryResourcePrefix      = "kubescape-registry-scan"
	RegistryAuthFieldInSecret   = "registriesAuth"
	RegistryCommandBody         = "request-body.json"
	RegistryCronjobTemplateName = "cronjobTemplate"
	RegistryRequestVolumeName   = "request-body-volume"
)

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

type RegistryProvider string

const (
	AWS    RegistryProvider = "aws"
	Azure  RegistryProvider = "azure"
	Google RegistryProvider = "google"
	Harbor RegistryProvider = "harbor"
	Quay   RegistryProvider = "quay"
	Nexus  RegistryProvider = "nexus"
	Gitlab RegistryProvider = "gitlab"
)

type RegistryManageStatus string
type RegistryScanStatus string

const (
	Empty   RegistryManageStatus = ""
	Created RegistryManageStatus = "Created"
	Updated RegistryManageStatus = "Updated"
	Error   RegistryManageStatus = "Error"

	// Scan statuses
	Failed     RegistryScanStatus = "Failed"
	InProgress RegistryScanStatus = "In Progress"
	Completed  RegistryScanStatus = "Completed"
)

type ContainerImageRegistry interface {
	MaskSecret()
	ExtractSecret() interface{}
	FillSecret(interface{}) error
	GetBase() *BaseContainerImageRegistry
	SetBase(*BaseContainerImageRegistry)
	Validate() error
	GetDisplayName() string
}

type BaseContainerImageRegistry struct {
	PortalBase          `json:",inline" bson:"inline"`
	Provider            RegistryProvider     `json:"provider" bson:"provider"`
	ClusterName         string               `json:"clusterName" bson:"clusterName"`
	Repositories        []string             `json:"repositories" bson:"repositories"`
	LastScan            *time.Time           `json:"lastScan,omitempty" bson:"lastScan,omitempty"`
	ScanFrequency       string               `json:"scanFrequency,omitempty" bson:"scanFrequency"`
	NextScan            *time.Time           `json:"nextScan,omitempty" bson:"nextScan,omitempty"`
	ResourceName        string               `json:"resourceName,omitempty" bson:"resourceName,omitempty"`
	AuthID              string               `json:"authID,omitempty" bson:"authID"`
	ManageStatus        RegistryManageStatus `json:"manageStatus,omitempty" bson:"manageStatus"`
	ManageStatusMessage string               `json:"manageStatusMessage,omitempty" bson:"manageStatusMessage"`
	ScanStatus          RegistryScanStatus   `json:"scanStatus,omitempty" bson:"scanStatus"`
	ScanStatusMessage   string               `json:"scanStatusMessage,omitempty" bson:"scanStatusMessage"`
}

const RegistryScanStatusesKindPath = "registrystatuses"
const RegistryScanStatusesKind = "RegistryStatuses"

type ContainerImageRegistryScanStatusUpdate struct {
	GUID              string             `json:"guid"`
	ScanStatus        RegistryScanStatus `json:"scanStatus"`
	ScanStatusMessage string             `json:"scanStatusMessage,omitempty"`
	ScanTime          time.Time          `json:"scanTime"`
}

type QuayImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	ContainerRegistryName      string `json:"containerRegistryName"`
	RobotAccountName           string `json:"robotAccountName"`
	RobotAccountToken          string `json:"robotAccountToken,omitempty"`
}

type HarborImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	InstanceURL                string `json:"instanceURL"`
	Username                   string `json:"username"`
	Password                   string `json:"password,omitempty"`
}

type AzureImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	LoginServer                string `json:"loginServer"`
	Username                   string `json:"username"`
	AccessToken                string `json:"accessToken,omitempty"`
}

type AWSImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	RegistryURI                string `json:"registryURI"`
	RegistryRegion             string `json:"registryRegion"`
	AccessKeyID                string `json:"accessKeyID,omitempty"`
	SecretAccessKey            string `json:"secretAccessKey,omitempty"`
	RoleARN                    string `json:"roleARN,omitempty"`
}

type GoogleImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	RegistryURI                string                 `json:"registryURI"`
	ProjectID                  string                 `json:"projectID"`
	Key                        map[string]interface{} `json:"key,omitempty"`
}

type NexusImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	RegistryURL                string `json:"registryURL"`
	Username                   string `json:"username"`
	Password                   string `json:"password,omitempty"`
}

type GitlabImageRegistry struct {
	BaseContainerImageRegistry `json:",inline"`
	Username                   string `json:"username"`
	AccessToken                string `json:"accessToken,omitempty"`
}

type CheckRegistryResp struct {
	Repositories []string `json:"repositories,omitempty"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
}
