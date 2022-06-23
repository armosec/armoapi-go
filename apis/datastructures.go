package apis

import (
	"net/http"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

/*
SessionChain
hold up all related jobIDs, earliest timestamp(command was issued/websocket started to handle it)
and action title

goal is to provide
for e.g in vulnerability scan context:

BE/cluster sends websocket request(With jobID ofc - jobid#1)
websocket takes all the cluster workloads and for each workload it creates jobID_i
for each container in workload_i it creates jobid_j

so when it sends the scan it sends the normal command object(pre sessionchain) to vulnscan
+
session: {
     jobIDs: [jobID#1,jobID_i,jobID_j]
	 timestamp: <jobID#1 timestamp>
	 rootJobID: jobID#1
}

WHy?
----
each scan will hold it's own unique sessionChain
rootJobID will allow customers to find their latest scans issues by cluster/other
jobIDs will allow them to take all specific workload related for that specific scan

*/
type SessionChain struct {
	JobIDs      []string  `json:"jobIDs"`              // all related JobIds in order: eg. grandparent,parent,this
	Timestamp   time.Time `json:"timestamp"`           //earliest/ timestamp
	RootJobID   string    `json:"rootJobID,omitempty"` //e,g grandparent
	ActionTitle string    `json:"action,omitempty"`    //e,g vulnerability-scan
}

type SessionChainWrapper struct {
	SessionChain `json:",inline"`
	Designators  armotypes.PortalDesignator `json:"designators"`
}

type DBCommand struct {
	Commands map[string]interface{} `json:"commands"`
}

//taken from BE
// ElasticRespTotal holds the total struct in Elastic array response
type ElasticRespTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// V2ListResponse holds the response of some list request with some metadata
type V2ListResponse struct {
	Total    ElasticRespTotal `json:"total"`
	Response interface{}      `json:"response"`
	// Cursor for quick access to the next page. Not supported yet
	Cursor string `json:"cursor"`
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

//PaginationMarks for split documents
type PaginationMarks struct {
	ReportNumber int  `json:"chunkNumber"` // serial number of report, used in pagination
	IsLastReport bool `json:"isLastChunk"` //specify this is the last report, used in pagination
}
