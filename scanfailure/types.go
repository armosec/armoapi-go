package scanfailure

import "time"

// ScanFailureCase enumerates the known reasons a scan can fail.
type ScanFailureCase int

const (
	// ScanFailureUnknown is the zero value; used when no specific case applies.
	ScanFailureUnknown ScanFailureCase = 0
	// ScanFailureCVE — have SBOM, can't match against vulnerability DBs.
	ScanFailureCVE ScanFailureCase = 1
	// ScanFailureSBOMGeneration — can't build SBOM from image.
	ScanFailureSBOMGeneration ScanFailureCase = 2
	// ScanFailureOOMKilled — scanner process was OOM-killed.
	ScanFailureOOMKilled ScanFailureCase = 3
	// ScanFailureBackendPost — scan succeeded but results couldn't be posted.
	ScanFailureBackendPost ScanFailureCase = 4
)

// String returns a human-readable description of the failure case.
func (f ScanFailureCase) String() string {
	switch f {
	case ScanFailureCVE:
		return "CVE matching failed"
	case ScanFailureSBOMGeneration:
		return "SBOM generation failed"
	case ScanFailureOOMKilled:
		return "Scanner process OOM-killed"
	case ScanFailureBackendPost:
		return "Backend post failed"
	default:
		return "Unknown failure"
	}
}

// WorkloadIdentifier identifies a single Kubernetes workload affected by a scan failure.
// A failed image may be used by multiple workloads, so the report carries a list of these.
type WorkloadIdentifier struct {
	ClusterName  string `json:"clusterName" bson:"clusterName"`
	Namespace    string `json:"namespace" bson:"namespace"`
	WorkloadKind  string `json:"workloadKind" bson:"workloadKind"`
	WorkloadName  string `json:"workloadName" bson:"workloadName"`
	ContainerName string `json:"containerName,omitempty" bson:"containerName,omitempty"`
}

// ScanFailureReport is emitted by the scanner when a scan fails.
// The scanner sends ONE report per failed image with all affected workloads listed.
// Downstream services (event-ingester, UNS) fan out per workload for notifications.
// For registry scans, Workloads is nil/empty and RegistryName is populated.
type ScanFailureReport struct {
	CustomerGUID  string               `json:"customerGUID" bson:"customerGUID"`
	Workloads     []WorkloadIdentifier `json:"workloads,omitempty" bson:"workloads,omitempty"`
	ImageTag      string               `json:"imageTag" bson:"imageTag"`
	FailureCase   ScanFailureCase      `json:"failureCase" bson:"failureCase"`
	FailureReason string               `json:"failureReason" bson:"failureReason"`
	Timestamp     time.Time            `json:"timestamp" bson:"timestamp"`
	ImageHash     string               `json:"imageHash,omitempty" bson:"imageHash,omitempty"`
	ContainerName string               `json:"containerName,omitempty" bson:"containerName,omitempty"`
	JobID         string               `json:"jobID,omitempty" bson:"jobID,omitempty"`

	// Registry scan context (no workloads).
	RegistryName   string `json:"registryName,omitempty" bson:"registryName,omitempty"`
	IsRegistryScan bool   `json:"isRegistryScan,omitempty" bson:"isRegistryScan,omitempty"`
}
