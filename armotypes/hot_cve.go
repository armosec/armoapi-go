package armotypes

import (
	"fmt"
	"strings"
)

// HotCVEAffectedPackage defines a package affected by a hot CVE
type HotCVEAffectedPackage struct {
	PackageName  string   `json:"packageName" bson:"packageName"`
	PackageTypes []string `json:"packageTypes,omitempty" bson:"packageTypes,omitempty"`
	VersionStart string   `json:"versionStart,omitempty" bson:"versionStart,omitempty"`
	VersionEnd   string   `json:"versionEnd,omitempty" bson:"versionEnd,omitempty"`
	VersionExact []string `json:"versionExact,omitempty" bson:"versionExact,omitempty"`
	FixedVersion string   `json:"fixedVersion,omitempty" bson:"fixedVersion,omitempty"`
}

// HotCVE defines a hot CVE from the external JSON feed
type HotCVE struct {
	CVEID            string                  `json:"cveId" bson:"cveId"`
	Title            string                  `json:"title" bson:"title"`
	Description      string                  `json:"description,omitempty" bson:"description,omitempty"`
	Severity         string                  `json:"severity" bson:"severity"`
	SeverityScore    int                     `json:"severityScore" bson:"severityScore"`
	IsRCE            bool                    `json:"isRce,omitempty" bson:"isRce,omitempty"`
	References       []string                `json:"references,omitempty" bson:"references,omitempty"`
	Status           string                  `json:"status" bson:"status"`
	AffectedPackages []HotCVEAffectedPackage `json:"affectedPackages" bson:"affectedPackages"`
}

// HotCVEEndpointResponse is the JSON response from the external hot CVE endpoint
type HotCVEEndpointResponse struct {
	Version string   `json:"version"`
	HotCVEs []HotCVE `json:"hotCves"`
}

// HotCVEOnFinishedMessage is the Pulsar message published for UNS after batch scan
type HotCVEOnFinishedMessage struct {
	CustomerGUID string   `json:"customerGUID"`
	ClusterName  string   `json:"clusterName"`
	Namespace    string   `json:"namespace"`
	Kind         string   `json:"kind"`
	WorkloadName string   `json:"workloadName"`
	CVEID        string   `json:"cveId"`
	Severity     string   `json:"severity"`
	Components   []string `json:"components"`
}

// validHotCVESeverities lists the accepted severity values for hot CVEs (lowercase for case-insensitive comparison).
// Subset of containerscan.KnownSeverities — cannot import directly due to circular dependency.
var validHotCVESeverities = map[string]bool{
	"critical": true, "high": true, "medium": true, "low": true,
}

// Validate checks that the HotCVE has all required fields and valid values.
func (h *HotCVE) Validate() error {
	if h.CVEID == "" {
		return fmt.Errorf("cveId is required")
	}
	if h.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	if !validHotCVESeverities[strings.ToLower(h.Severity)] {
		return fmt.Errorf("invalid severity %q: must be critical, high, medium, or low", h.Severity)
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
	// At least one version constraint must be specified
	if p.VersionStart == "" && p.VersionEnd == "" && len(p.VersionExact) == 0 {
		return fmt.Errorf("at least one version constraint (versionStart, versionEnd, or versionExact) is required for package %q", p.PackageName)
	}
	return nil
}

const (
	HotCVESentinelLayerHash = "__hot_cve__"
	HotCVEStatusActive      = "active"
	HotCVEStatusInactive    = "inactive"
)
