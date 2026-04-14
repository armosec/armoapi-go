package armotypes

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

const (
	HotCVESentinelLayerHash = "__hot_cve__"
	HotCVEStatusActive      = "active"
	HotCVEStatusInactive    = "inactive"
)
