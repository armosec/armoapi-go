package apis

import (
	"net/http"
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

/*
SessionChain provides the context of a given job.

The goal is to provide context for a given job: its parent jobs, a chain of how the jobs were spawned and some metadata.

Consider a vulnerability scan, for example:
  - The Backend or cluster sends a websocket request with a Job ID, e.g. jobID_1.
  - The Websocket takes all the cluster workloads and for each workload it creates a job with ID `jobID_i`.
  - Then, for each container in `workload_i` it creates a job with ID `jobID_j`.

So when the Websocket sends the scan command, it sends the normal command object (pre Session Chain) to the Vulnerability Scanner

	session: {
	  "jobIDs": ["jobID_1", "jobID_i", "jobID_j"],
	  "timestamp": "<jobID#1 timestamp>",
	  "rootJobID": "jobID_1"
	}

This Session Chain is needed so that:
  - each scan will hold it's own unique sessionChain.
  - `rootJobID` will allow customers to find their latest scans issues by cluster/other.
  - `jobID`s will allow customers to take all specific workload related for that specific scan.
*/
type SessionChain struct {
	// All related job IDs in order from the most distant to the closes relative.
	//
	// For instance: grandparent → parent → current.
	//
	// Example: ["825f0a9e-34a9-4727-b81a-6e1bf3a63725", "c188de09-c6ec-4814-b36a-722dcccea64b"]
	JobIDs []string `json:"jobIDs"`
	// The timestamp of the earliest job
	Timestamp time.Time `json:"timestamp"`
	// ID of the job that started this chain.
	//
	// Example: 825f0a9e-34a9-4727-b81a-6e1bf3a63725
	// swagger:strfmt uuid4
	RootJobID string `json:"rootJobID,omitempty"`
	// Title of the current action being performed
	//
	// Example: vulnerability-scan
	ActionTitle string `json:"action,omitempty"`
}

type SessionChainWrapper struct {
	SessionChain `json:",inline"`
	Designators  identifiers.PortalDesignator `json:"designators"`
}

type DBCommand struct {
	Commands map[string]interface{} `json:"commands"`
}

// Oauth2Customer returns inside the "ca_groups" field in claims section of
// Oauth2 verification process
type Oauth2Customer struct {
	CustomerName string `json:"customerName"`
	CustomerGUID string `json:"customerGUID"`
}

type LoginObject struct {
	Authorization string `json:"authorization"`
	GUID          string
	Cookies       []*http.Cookie
	Expires       string
}

// PaginationMarks for split documents
type PaginationMarks struct {
	ReportNumber int  `json:"chunkNumber"` // serial number of report, used in pagination
	IsLastReport bool `json:"isLastChunk"` //specify this is the last report, used in pagination

	// TotalChunksExpected and TotalChunksRecieved are used to track the progress of the report.
	// Once TotalChunksExpected == TotalChunksRecieved, fields values will be set to -1 and reported is considered as completed.

	//specify the total number of chunks expected. Will be populated with the ReportNumber of the LastReport (IsLastReport == true)
	TotalChunksExpected int `json:"totalChunksExpected"`

	//specify the total number of chunks recieved so far - will be increment by one on each chunk recieved.
	TotalChunksRecieved int `json:"totalChunksRecieved"`
}
