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

// Failure reason codes sent by scanners (kubevuln, node-agent) in FailureReason.
// These are short enum-like codes — the notification service (UNS) maps them to
// human-friendly text at render time, so notification wording can be changed
// without redeploying in-cluster scanners.
const (
	ReasonSBOMGenerationFailed = "sbom_generation_failed"
	ReasonImageTooLarge        = "image_too_large"
	ReasonSBOMTooLarge         = "sbom_too_large"
	ReasonSBOMIncomplete       = "sbom_incomplete"
	ReasonImageAuthFailed      = "image_auth_failed"
	ReasonImageNotFound        = "image_not_found"
	ReasonCVEMatchingFailed    = "cve_matching_failed"
	ReasonResultUploadFailed   = "result_upload_failed"
	ReasonSBOMStorageFailed    = "sbom_storage_failed"
	ReasonUnexpected           = "unexpected_error"
)

// reasonFriendlyText maps reason codes to human-friendly text for notifications.
var reasonFriendlyText = map[string]string{
	ReasonSBOMGenerationFailed: "Failed to generate software inventory (SBOM) for this image",
	ReasonImageTooLarge:        "Image exceeds the maximum size limit for vulnerability scanning",
	ReasonSBOMTooLarge:         "Generated software inventory (SBOM) exceeds the maximum size limit",
	ReasonSBOMIncomplete:       "SBOM generation was incomplete — the scan may have timed out or the image exceeded size limits",
	ReasonImageAuthFailed:      "Failed to authenticate when pulling the container image",
	ReasonImageNotFound:        "Container image manifest not found in registry",
	ReasonCVEMatchingFailed:    "Failed to match image components against vulnerability databases",
	ReasonResultUploadFailed:   "Scan completed but results could not be uploaded to the platform",
	ReasonSBOMStorageFailed:    "Failed to store the generated software inventory (SBOM)",
	ReasonUnexpected:           "An unexpected error occurred during vulnerability scanning",
}

// ReasonFriendlyText returns the human-friendly notification text for a reason code.
// If the code is unknown, returns the code itself as a fallback.
func ReasonFriendlyText(reasonCode string) string {
	if text, ok := reasonFriendlyText[reasonCode]; ok {
		return text
	}
	return reasonCode
}

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
	// Error holds the raw error string for R&D debugging. Not rendered in user-facing
	// notifications (Slack/Teams templates use FailureReason only). Producers must avoid
	// including secrets (tokens, credentials) — redact sensitive data before populating.
	Error     string    `json:"error,omitempty" bson:"error,omitempty"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	ImageHash string    `json:"imageHash,omitempty" bson:"imageHash,omitempty"`
	JobID     string    `json:"jobID,omitempty" bson:"jobID,omitempty"`

	// Registry scan context (no workloads).
	RegistryName   string `json:"registryName,omitempty" bson:"registryName,omitempty"`
	IsRegistryScan bool   `json:"isRegistryScan,omitempty" bson:"isRegistryScan,omitempty"`
}
