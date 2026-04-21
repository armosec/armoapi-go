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

// HotCVEValidSeverities is the allow-list of severity values the hot-CVE
// pipeline can propagate end-to-end. Values MUST be titlecase — the
// postgres-connector vulnerabilities_cves view upsert INNER-joins
// vulnerabilities_v1 with `severity = 'Critical' OR 'High' OR ...`
// (titlecase only), and that join silently drops anything else. Keep this
// list in lockstep with the SQL filter in postgres-connector/internal/
// dalhelpers/templates/vulnerabilitiesCVEsViewUpsert.sql. Empirical
// confirmation on 2026-04-21: 13 lowercase "critical" rows in dev
// vulnerabilities_v1 were filtered out of 260,278 vulnerabilities_cves
// rows, resulting in zero is_hot_cve=true rows cluster-wide.
var HotCVEValidSeverities = map[string]struct{}{
	"Critical":   {},
	"High":       {},
	"Medium":     {},
	"Low":        {},
	"Unknown":    {},
	"Negligible": {},
}

// Validate checks that the HotCVE has all required fields.
func (h *HotCVE) Validate() error {
	if h.CVEID == "" {
		return fmt.Errorf("cveId is required")
	}
	if h.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	if _, ok := HotCVEValidSeverities[h.Severity]; !ok {
		return fmt.Errorf("invalid severity %q: must be one of Critical, High, Medium, Low, Unknown, Negligible (titlecase)", h.Severity)
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

const (
	HotCVESentinelLayerHash = "__hot_cve__"
	HotCVEStatusActive      = "active"
	HotCVEStatusInactive    = "inactive"
)
