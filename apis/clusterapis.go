package apis

// WebsocketScanCommand api
const (
	WebsocketScanCommandVersion string = "v1"
	WebsocketScanCommandPath    string = "scanImage"
	DBCommandPath               string = "DBCommand"
	ServerReady                 string = "ready"
)

// Supported NotificationTypes
//
// swagger:enum NotificationPolicyType
type NotificationPolicyType string

const (
	TypeValidateRules NotificationPolicyType = "validateRules"
	// Execute a posture scan
	TypeExecPostureScan NotificationPolicyType = "execPostureScan"
	TypeUpdateRules     NotificationPolicyType = "updateRules"
	TypeRunKubescapeJob NotificationPolicyType = "runKubescapeJob"
	// Trigger a Kubescape scan
	TypeRunKubescape NotificationPolicyType = "kubescapeScan"
	// Create a CronJob that runs a Kubescape scan
	TypeSetKubescapeCronJob NotificationPolicyType = "setKubescapeCronJob"
	// Update a CronJob that runs a Kubescape scan
	TypeUpdateKubescapeCronJob NotificationPolicyType = "updateKubescapeCronJob"
	// Delete a CronJob that runs a Kubescape scan
	TypeDeleteKubescapeCronJob NotificationPolicyType = "deleteKubescapeCronJob"
	// Create a CronJob that runs a Vulnerability Scan
	TypeSetVulnScanCronJob NotificationPolicyType = "setVulnScanCronJob"
	// Update a CronJob that runs a Vulnerability Scan
	TypeUpdateVulnScanCronJob NotificationPolicyType = "updateVulnScanCronJob"
	// Delete a CronJob that runs a Vulnerability Scan
	TypeDeleteVulnScanCronJob NotificationPolicyType = "deleteVulnScanCronJob"
	// Trigger an image scan
	TypeScanImages NotificationPolicyType = "scan"
	// Trigger a relevancy image scan
	TypeScanRelevancyImages NotificationPolicyType = "scanRelevancy"
	// Trigger a registry scan
	TypeScanRegistry NotificationPolicyType = "scanRegistry"
	// Create a CronJob that runs registry scans
	TypeSetRegistryScanCronJob NotificationPolicyType = "setRegistryScanCronJob"
	// Update a CronJob that runs registry scans
	TypeUpdateRegistryScanCronJob NotificationPolicyType = "updateRegistryScanCronJob"
	// Delete a CronJob that runs registry scans
	TypeDeleteRegistryScanCronJob NotificationPolicyType = "deleteRegistryScanCronJob"
	TypeTestRegistryConnectivity  NotificationPolicyType = "testRegistryConnectivity"
)
