package armotypes

import "time"

// SynchronizerClient represents a client which is connected to the synchronizer server
type SynchronizerClient struct {
	CustomerGUID        string    `json:"customerGUID"`
	Cluster             string    `json:"cluster"`
	Replica             string    `json:"replica"`
	LastKeepAlive       time.Time `json:"lastKeepAlive"`
	ConnectionTime      time.Time `json:"connectionTime"`
	HelmVersion         string    `json:"helmVersion"`
	SynchronizerVersion string    `json:"synchronizerVersion"`
	ConnectionId        string    `json:"connectionId"`
}
