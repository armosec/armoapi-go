package apis

// WebsocketScanCommand api
const (
	WebsocketScanCommandVersion string = "v1"
	WebsocketScanCommandPath    string = "scanImage"
	DBCommandPath               string = "DBCommand"
	ServerReady                 string = "ready"
)

// Supported NotificationTypes
type NotificationPolicyType string

const (
	TypeRunKubescape              NotificationPolicyType = "kubescapeScan"
	TypeSetKubescapeCronJob       NotificationPolicyType = "setKubescapeCronJob"
	TypeUpdateKubescapeCronJob    NotificationPolicyType = "updateKubescapeCronJob"
	TypeDeleteKubescapeCronJob    NotificationPolicyType = "deleteKubescapeCronJob"
	TypeSetVulnScanCronJob        NotificationPolicyType = "setVulnScanCronJob"
	TypeUpdateVulnScanCronJob     NotificationPolicyType = "updateVulnScanCronJob"
	TypeDeleteVulnScanCronJob     NotificationPolicyType = "deleteVulnScanCronJob"
	TypeScanImages                NotificationPolicyType = "scan"
	TypeScanRegistry              NotificationPolicyType = "scanRegistry"
	TypeSetRegistryScanCronJob    NotificationPolicyType = "setRegistryScanCronJob"
	TypeUpdateRegistryScanCronJob NotificationPolicyType = "updateRegistryScanCronJob"
	TypeDeleteRegistryScanCronJob NotificationPolicyType = "deleteRegistryScanCronJob"
)
