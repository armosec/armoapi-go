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

// WebsocketScanCommand trigger scan thru the websocket
type WebsocketScanCommand struct {
	// CustomerGUID string `json:"customerGUID"`
	Session         SessionChain           `json:"session,omitempty"`
	ImageTag        string                 `json:"imageTag"`
	Wlid            string                 `json:"wlid"`
	IsScanned       bool                   `json:"isScanned"`
	ContainerName   string                 `json:"containerName"`
	JobID           string                 `json:"jobID,omitempty"`
	ParentJobID     string                 `json:"parentJobID,omitempty"`
	LastAction      int                    `json:"actionIDN"`
	ImageHash       string                 `json:"imageHash"`
	Credentials     *types.AuthConfig      `json:"credentials,omitempty"`
	Credentialslist []types.AuthConfig     `json:"credentialsList,omitempty"`
	Args            map[string]interface{} `json:"args,omitempty"`
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
