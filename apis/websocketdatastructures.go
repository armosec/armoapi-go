package apis

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/docker/docker/api/types"
)

// Commands list of commands received from websocket
type Commands struct {
	Commands []Command `json:"commands"`
}

// Command structure of command received from websocket
type Command struct {
	// basic command
	CommandName NotificationPolicyType `json:"commandName"`
	ResponseID  string                 `json:"responseID,omitempty"`

	// command designators
	Designators []armotypes.PortalDesignator `json:"designators,omitempty"`
	Wlid        string                       `json:"wlid,omitempty"`
	WildWlid    string                       `json:"wildWlid,omitempty"`
	Sid         string                       `json:"sid,omitempty"`
	WildSid     string                       `json:"wildSid,omitempty"`
	JobTracking JobTracking                  `json:"jobTracking,omitempty"`

	// command extra data
	Args map[string]interface{} `json:"args,omitempty"`
}

type JobTracking struct {
	JobID            string    `json:"jobID,omitempty"`
	ParentID         string    `json:"parentAction,omitempty"`
	LastActionNumber int       `json:"numSeq,omitempty"`
	Timestamp        time.Time `json:"timestamp,omitempty"`
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

type SafeMode struct {
	Reporter        string `json:"reporter"`                // "Agent"
	Action          string `json:"action,omitempty"`        // "action"
	Wlid            string `json:"wlid"`                    // CAA_WLID
	PodName         string `json:"podName"`                 // CAA_POD_NAME
	InstanceID      string `json:"instanceID"`              // CAA_POD_NAME
	ContainerName   string `json:"containerName,omitempty"` // CAA_CONTAINER_NAME
	ProcessName     string `json:"processName,omitempty"`
	ProcessID       int    `json:"processID,omitempty"`
	ProcessCMD      string `json:"processCMD,omitempty"`
	ComponentGUID   string `json:"componentGUID,omitempty"` // CAA_GUID
	StatusCode      int    `json:"statusCode"`              // 0/1/2
	ProcessExitCode int    `json:"processExitCode"`         // 0 +
	Timestamp       int64  `json:"timestamp"`
	Message         string `json:"message,omitempty"` // any string
	JobID           string `json:"jobID,omitempty"`   // any string
	Compatible      *bool  `json:"compatible,omitempty"`
}

// CronJobParams parmas for cronJob
type CronJobParams struct {
	CronTabSchedule string `json:"cronTabSchedule"`
	JobName         string `json:"name,omitempty"`
}
