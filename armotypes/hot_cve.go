package armotypes

import "fmt"

// HotCVEAffectedPackage defines a package affected by a hot CVE.
type HotCVEAffectedPackage struct {
	PackageName  string   `json:"packageName" bson:"packageName"`
	PackageTypes []string `json:"packageTypes,omitempty" bson:"packageTypes,omitempty"`
	VersionStart string   `json:"versionStart,omitempty" bson:"versionStart,omitempty"`
	VersionEnd   string   `json:"versionEnd,omitempty" bson:"versionEnd,omitempty"`
	VersionExact []string `json:"versionExact,omitempty" bson:"versionExact,omitempty"`
	FixedVersion string   `json:"fixedVersion,omitempty" bson:"fixedVersion,omitempty"`
}

// HotCVE defines a hot CVE definition published via admin API.
type HotCVE struct {
	CVEID            string                  `json:"cveId" bson:"cveId"`
	Title            string                  `json:"title" bson:"title"`
	Description      string                  `json:"description,omitempty" bson:"description,omitempty"`
	Severity         string                  `json:"severity" bson:"severity"`
	SeverityScore    int                     `json:"severityScore" bson:"severityScore"`
	References       []string                `json:"references,omitempty" bson:"references,omitempty"`
	Status           string                  `json:"status" bson:"status"`
	AffectedPackages []HotCVEAffectedPackage `json:"affectedPackages" bson:"affectedPackages"`
}

// HotCVEEndpointResponse is the JSON response from the external hot CVE endpoint.
type HotCVEEndpointResponse struct {
	Version string   `json:"version"`
	HotCVEs []HotCVE `json:"hotCves"`
}

// HotCVEOnFinishedMessage is the Pulsar message published for UNS after hot CVE processing.
type HotCVEOnFinishedMessage struct {
	CustomerGUID string `json:"customerGUID"`
	CVEID        string `json:"cveId"`
	Severity     string `json:"severity"`
}

// Validate checks that the HotCVE has all required fields.
func (h *HotCVE) Validate() error {
	if h.CVEID == "" {
		return fmt.Errorf("cveId is required")
	}
	if h.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	if len(h.AffectedPackages) == 0 {
		return fmt.Errorf("affectedPackages is required and must not be empty")
	}
	for i, pkg := range h.AffectedPackages {
		if err := pkg.Validate(); err != nil {
			return fmt.Errorf("affectedPackages[%d]: %w", i, err)
		}
	}
	if h.Status != "" && h.Status != HotCVEStatusActive && h.Status != HotCVEStatusInactive {
		return fmt.Errorf("invalid status %q: must be %q or %q", h.Status, HotCVEStatusActive, HotCVEStatusInactive)
	}
	return nil
}

// Validate checks that the HotCVEAffectedPackage has required fields.
func (p *HotCVEAffectedPackage) Validate() error {
	if p.PackageName == "" {
		return fmt.Errorf("packageName is required")
	}
	if p.VersionStart == "" && p.VersionEnd == "" && len(p.VersionExact) == 0 {
		return fmt.Errorf("at least one version constraint (versionStart, versionEnd, or versionExact) is required for package %q", p.PackageName)
	}
	return nil
}

// HotCVEMatchMessage is published to the KDR ingester topic when a hot CVE
// matches a customer's SBOM component. The HotCVEMessageProcessor creates
// runtime incidents from these messages.
type HotCVEMatchMessage struct {
	CustomerGUID     string `json:"customerGUID"`
	CVEID            string `json:"cveId"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Severity         string `json:"severity"`
	SeverityScore    int    `json:"severityScore"`
	Component        string `json:"component"`
	ComponentVersion string `json:"componentVersion"`
	PackageType      string `json:"packageType"`
	FixedVersion     string `json:"fixedVersion"`
	ResourceHash     string `json:"resourceHash"`
	ClusterName      string `json:"clusterName"`
	Namespace        string `json:"namespace"`
	WorkloadName     string `json:"workloadName"`
	WorkloadKind     string `json:"workloadKind"`
	IsInUse          bool   `json:"isInUse"`
}

const (
	HotCVESentinelLayerHash = "__hot_cve__"
	HotCVEStatusActive      = "active"
	HotCVEStatusInactive    = "inactive"
)
