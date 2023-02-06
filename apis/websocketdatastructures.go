package apis

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/docker/docker/api/types"
)

// Commands contains a collection of commands for the in-cluster components
type Commands struct {
	// A list of commands to execute
	//
	// Example: [ { "CommandName": "scanRegistry", "args": { "registryInfo-v1": { "registryName": "quay.io/armosec" } } } ]
	Commands []Command `json:"commands"`
}

// Command describes an individual command for the in-cluster components
type Command struct {
	// Name of the command
	//
	// Example: updateRules
	CommandName NotificationPolicyType `json:"commandName"`
	// ID of the response
	//
	// Example: 49cfe0a0-9fab-4e54-a6e4-7b27e566d3cd
	ResponseID string `json:"responseID,omitempty"`

	// Designators for the command
	//
	// Designators select the targets to which the command applies.
	Designators []armotypes.PortalDesignator `json:"designators,omitempty"`
	Wlid        string                       `json:"wlid,omitempty"`
	WildWlid    string                       `json:"wildWlid,omitempty"`
	Sid         string                       `json:"sid,omitempty"`
	WildSid     string                       `json:"wildSid,omitempty"`
	// Job tracking context for
	JobTracking JobTracking `json:"jobTracking,omitempty"`

	// Arguments for the command
	Args map[string]interface{} `json:"args,omitempty"`
}

// JobTracking describes a context in which the job is executing
// It is used to track job execution source and context: what spawned it, when and under what circumstances.
type JobTracking struct {
	// ID of the current job
	//
	// Example: 0f2c8611-ba99-40e5-af21-2bc3823e3283
	JobID string `json:"jobID,omitempty"`
	// ID of the parent job
	//
	// Example: 6ecfe560-104c-4e7b-8cd3-ee3cbc3b58fb
	ParentID string `json:"parentAction,omitempty"`
	// Number of the last action
	//
	// Example: 2
	LastActionNumber int `json:"numSeq,omitempty"`
	// Timestamp of the latest action
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// WebsocketScanCommand is a command that triggers a scan for vulnerabilities.
type WebsocketScanCommand struct {
	// Current session context
	//
	// Used for correlating requests in the logs.
	Session SessionChain `json:"session,omitempty"`
	// Tag of the image to scan
	//
	// Example: nginx:latest
	ImageTag string `json:"imageTag"`
	// ID of a workload that is running the image you want to scan
	//
	// Example: wlid://cluster-marina/namespace-default/deployment-nginx
	Wlid string `json:"wlid"`
	// Has the provided image been previously scanned or not?
	//
	// An image will only be scanned if it has not been scanned previously (value is `false`).
	// If an image has been previously scanned (value is `true`), it will not be scanned again.
	//
	// Example: false
	IsScanned bool `json:"isScanned"`
	// Name of the container that contains an image to be scanned
	//
	// Example: nginx
	ContainerName string `json:"containerName"`
	// ID of the scanning Job
	//
	// Example: 7b04592b-665a-4e47-a9c9-65b2b3cabb49
	JobID string `json:"jobID,omitempty"`
	// ID of the Parent Job — a job that initiated the current job
	//
	// Example: 825f0a9e-34a9-4727-b81a-6e1bf3a63725
	ParentJobID string `json:"parentJobID,omitempty"`
	// The last action received from the Websocket
	//
	// Example: 2
	LastAction int `json:"actionIDN"`
	// Hash of the image to scan
	//
	// Example: bcae378eacedab83da66079d9366c8f5df542d7ed9ab23bf487e3e1a8481375d
	ImageHash string `json:"imageHash"`
	// InstanceID for relevancy scan
	// namespace-<namespace>/<kind>-<name>/<resourceVersion>
	// Example: namespace-default/pod-nginx/75641
	InstanceID *string `json:"instanceID,omitempty"`
	// Deprecated: Credentials to the Container Registry that holds the image to be scanned
	//
	// Kept for backward compatibility
	Credentials *types.AuthConfig `json:"credentials,omitempty"`
	// A list of credentials for private Container Registries that store images to be scanned
	Credentialslist []types.AuthConfig `json:"credentialsList,omitempty"`
	// Arguments to pass to the scan command
	//
	// Example: {"useHTTP": true, "skipTLSVerify": true, "registryName": "", "repository": "", "tag": ""}
	Args map[string]interface{} `json:"args,omitempty"`
	// CustomerGUID string `json:"customerGUID"`
}

// CronJobParams parmas for cronJob
type CronJobParams struct {
	CronTabSchedule string `json:"cronTabSchedule"`
	JobName         string `json:"name,omitempty"`
}
