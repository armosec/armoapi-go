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

// ScanFailureReport is emitted by the scanner when a scan fails.
// Each report covers a single workload-image pair (matching the container scan pattern).
// For registry scans, workload fields are empty and RegistryName is populated.
type ScanFailureReport struct {
	CustomerGUID  string          `json:"customerGUID" bson:"customerGUID"`
	ImageTag      string          `json:"imageTag" bson:"imageTag"`
	ImageHash     string          `json:"imageHash,omitempty" bson:"imageHash,omitempty"`
	ContainerName string          `json:"containerName,omitempty" bson:"containerName,omitempty"`
	FailureCase   ScanFailureCase `json:"failureCase" bson:"failureCase"`
	FailureReason string          `json:"failureReason" bson:"failureReason"`
	Timestamp     time.Time       `json:"timestamp" bson:"timestamp"`
	JobID         string          `json:"jobID,omitempty" bson:"jobID,omitempty"`

	// Workload context (single workload per report).
	ClusterName  string `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
	Namespace    string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	WorkloadKind string `json:"workloadKind,omitempty" bson:"workloadKind,omitempty"`
	WorkloadName string `json:"workloadName,omitempty" bson:"workloadName,omitempty"`

	// Registry scan context (mutually exclusive with workload fields).
	RegistryName   string `json:"registryName,omitempty" bson:"registryName,omitempty"`
	IsRegistryScan bool   `json:"isRegistryScan,omitempty" bson:"isRegistryScan,omitempty"`
}
