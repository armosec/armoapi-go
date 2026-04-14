package armotypes

// HotCVEAffectedPackage defines a package affected by a hot CVE
type HotCVEAffectedPackage struct {
	PackageName  string   `json:"package_name" bson:"package_name"`
	PackageTypes []string `json:"package_types,omitempty" bson:"package_types,omitempty"`
	VersionStart string   `json:"version_start,omitempty" bson:"version_start,omitempty"`
	VersionEnd   string   `json:"version_end,omitempty" bson:"version_end,omitempty"`
	VersionExact []string `json:"version_exact,omitempty" bson:"version_exact,omitempty"`
	FixedVersion string   `json:"fixed_version,omitempty" bson:"fixed_version,omitempty"`
}

// HotCVE defines a hot CVE from the external JSON feed
type HotCVE struct {
	CVEID            string                  `json:"cve_id" bson:"cve_id"`
	Title            string                  `json:"title" bson:"title"`
	Description      string                  `json:"description,omitempty" bson:"description,omitempty"`
	Severity         string                  `json:"severity" bson:"severity"`
	SeverityScore    int                     `json:"severity_score" bson:"severity_score"`
	IsRCE            bool                    `json:"is_rce,omitempty" bson:"is_rce,omitempty"`
	References       []string                `json:"references,omitempty" bson:"references,omitempty"`
	Status           string                  `json:"status" bson:"status"`
	AffectedPackages []HotCVEAffectedPackage `json:"affected_packages" bson:"affected_packages"`
}

// HotCVEEndpointResponse is the JSON response from the external hot CVE endpoint
type HotCVEEndpointResponse struct {
	Version string   `json:"version"`
	HotCVEs []HotCVE `json:"hot_cves"`
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

const (
	HotCVESentinelLayerHash = "__hot_cve__"
	HotCVEStatusActive      = "active"
	HotCVEStatusInactive    = "inactive"
)
