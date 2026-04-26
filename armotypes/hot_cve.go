package armotypes

import (
	"fmt"
	"sort"
	"strings"
)

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

// HotCVEAffectedWorkload describes a workload (cluster/namespace/kind/name +
// the image tag) impacted by a hot CVE. Carried inline on
// HotCVEOnFinishedMessage so the UNS notification dispatcher can render
// per-workload context (cluster, namespace, image, ...) into the Slack/Teams
// alert body without an extra postgres or dashboard round-trip.
type HotCVEAffectedWorkload struct {
	Cluster      string `json:"cluster,omitempty" bson:"cluster,omitempty"`
	Namespace    string `json:"namespace,omitempty" bson:"namespace,omitempty"`
	WorkloadKind string `json:"workloadKind,omitempty" bson:"workloadKind,omitempty"`
	WorkloadName string `json:"workloadName,omitempty" bson:"workloadName,omitempty"`
	ImageTag     string `json:"imageTag,omitempty" bson:"imageTag,omitempty"`
}

// HotCVEOnFinishedMessage is the Pulsar message published for UNS after hot CVE processing.
//
// Title, References and AffectedWorkloads are all `omitempty` so existing
// publishers that haven't been upgraded continue to produce valid messages
// (UNS treats absent fields as the zero value).
//   - Title is sourced from HotCVE.Title (the admin-curated headline).
//   - References mirror HotCVE.References (admin-curated authoritative URLs:
//     NVD, vendor advisory, etc.). UNS uses References[0] as the CVELink that
//     the Slack template renders as a clickable hyperlink under the CVE id.
//     Hardcoding NVD on the consumer side would override the admin's curated
//     references — they belong to the publisher, not the consumer.
//   - AffectedWorkloads enumerates the workloads in *this customer's* tenant
//     that match the hot CVE — UNS uses them both for rendering context and
//     for filtering against per-workflow scope (cluster/namespace).
type HotCVEOnFinishedMessage struct {
	CustomerGUID      string                   `json:"customerGUID"`
	CVEID             string                   `json:"cveId"`
	Severity          string                   `json:"severity"`
	Title             string                   `json:"title,omitempty"`
	References        []string                 `json:"references,omitempty"`
	AffectedWorkloads []HotCVEAffectedWorkload `json:"affectedWorkloads,omitempty"`
}

// hotCVEValidSeverities is the allow-list of severity values the hot-CVE
// pipeline can propagate end-to-end. Values MUST be titlecase — the
// postgres-connector vulnerabilities_cves view upsert INNER-joins
// vulnerabilities_v1 with `severity = 'Critical' OR 'High' OR ...`
// (titlecase only), and that join silently drops anything else. Keep this
// list in lockstep with the SQL filter in postgres-connector/internal/
// dalhelpers/templates/vulnerabilitiesCVEsViewUpsert.sql. Empirical
// confirmation on 2026-04-21: 13 lowercase "critical" rows in dev
// vulnerabilities_v1 were filtered out of 260,278 vulnerabilities_cves
// rows, resulting in zero is_hot_cve=true rows cluster-wide.
//
// Unexported so external callers cannot mutate the allow-list at runtime;
// use IsValidHotCVESeverity for membership checks.
var hotCVEValidSeverities = map[string]struct{}{
	"Critical":   {},
	"High":       {},
	"Medium":     {},
	"Low":        {},
	"Unknown":    {},
	"Negligible": {},
}

// IsValidHotCVESeverity reports whether severity is accepted by the hot-CVE
// pipeline. Callers outside this package should use this helper rather than
// touching the underlying map, so the allow-list stays immutable at runtime.
func IsValidHotCVESeverity(severity string) bool {
	_, ok := hotCVEValidSeverities[severity]
	return ok
}

// hotCVEValidSeveritiesList returns the allow-list as a sorted slice — used
// to build deterministic error messages without duplicating the severity
// strings between the allow-list and the error text. Building the list at
// call time (rather than caching) avoids any mutation risk and keeps the
// validator output in lockstep with the map.
func hotCVEValidSeveritiesList() []string {
	out := make([]string, 0, len(hotCVEValidSeverities))
	for s := range hotCVEValidSeverities {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// Validate checks that the HotCVE has all required fields.
func (h *HotCVE) Validate() error {
	if h.CVEID == "" {
		return fmt.Errorf("cveId is required")
	}
	if h.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	if !IsValidHotCVESeverity(h.Severity) {
		return fmt.Errorf("invalid severity %q: must be one of %s (titlecase)", h.Severity, strings.Join(hotCVEValidSeveritiesList(), ", "))
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
