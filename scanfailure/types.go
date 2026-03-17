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

// Human-friendly failure reasons displayed in Slack/Teams notifications.
// Scanners (kubevuln, node-agent) classify errors and pick the matching constant.
// These are rendered directly as {{ .FailureReason }} in notification templates.
const (
	ReasonSBOMGenerationFailed = "Failed to generate software inventory (SBOM) for this image"
	ReasonImageTooLarge        = "Image exceeds the maximum size limit for vulnerability scanning"
	ReasonSBOMTooLarge         = "Generated software inventory (SBOM) exceeds the maximum size limit"
	ReasonSBOMIncomplete       = "SBOM generation was incomplete — the scan may have timed out or the image exceeded size limits"
	ReasonImageAuthFailed      = "Failed to authenticate when pulling the container image"
	ReasonImageNotFound        = "Container image manifest not found in registry"
	ReasonCVEMatchingFailed    = "Failed to match image components against vulnerability databases"
	ReasonResultUploadFailed   = "Scan completed but results could not be uploaded to the platform"
	ReasonSBOMStorageFailed    = "Failed to store the generated software inventory (SBOM)"
	ReasonUnexpected           = "An unexpected error occurred during vulnerability scanning"
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
	ClusterName   string `json:"clusterName" bson:"clusterName"`
	Namespace     string `json:"namespace" bson:"namespace"`
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
	Error         string               `json:"error,omitempty" bson:"error,omitempty"`
	Timestamp     time.Time            `json:"timestamp" bson:"timestamp"`
	ImageHash string `json:"imageHash,omitempty" bson:"imageHash,omitempty"`
	JobID     string `json:"jobID,omitempty" bson:"jobID,omitempty"`

	// Registry scan context (no workloads).
	RegistryName   string `json:"registryName,omitempty" bson:"registryName,omitempty"`
	IsRegistryScan bool   `json:"isRegistryScan,omitempty" bson:"isRegistryScan,omitempty"`
}
