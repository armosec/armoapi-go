package v1

import (
	"github.com/armosec/armoapi-go/apis"
	"github.com/armosec/armoapi-go/containerscan"
	"github.com/armosec/armoapi-go/identifiers"
)

type ScanResultReport struct {
	Designators     identifiers.PortalDesignator                       `json:"designators"`
	Summary         *containerscan.CommonContainerScanSummaryResult    `json:"summary,omitempty"`
	ContainerScanID string                                             `json:"containersScanID"`
	Vulnerabilities []containerscan.CommonContainerVulnerabilityResult `json:"vulnerabilities"`
	PaginationInfo  apis.PaginationMarks                               `json:"paginationInfo"`
	Timestamp       int64                                              `json:"timestamp"`
}
