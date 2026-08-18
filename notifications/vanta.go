package notifications

// Vanta CollaborationConfig attribute keys (stored on the embedded
// PortalBase.Attributes of a provider=="vanta" CollaborationConfig).
const (
	// VantaAttrResourceID holds the Vanta-assigned resource id the sync pushes
	// to. Its presence marks the connection "complete" (like Jira's siteID).
	VantaAttrResourceID = "vantaResourceId"
	// VantaAttrUserAccountResourceID / VulnResourceID / PackageVulnResourceID
	// hold the per-resource-type ids for the families ARMO syncs.
	VantaAttrUserAccountResourceID = "vantaUserAccountResourceId"
	VantaAttrVulnComponentResID    = "vantaVulnerableComponentResourceId"
	VantaAttrPackageVulnResID      = "vantaPackageVulnerabilityResourceId"

	VantaAttrEnabled        = "vantaEnabled"
	VantaAttrLastSyncAt     = "vantaLastSyncAt"
	VantaAttrLastSyncStatus = "vantaLastSyncStatus" // "ok" | "error"
	VantaAttrLastSyncError  = "vantaLastSyncError"
	VantaAttrSyncedCount    = "vantaSyncedCount"
)

// IsVantaEnabled reports whether a Vanta CollaborationConfig is enabled.
func IsVantaEnabled(cc CollaborationConfig) bool {
	v, ok := cc.Attributes[VantaAttrEnabled].(bool)
	return ok && v
}
