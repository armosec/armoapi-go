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
	TypeValidateRules              NotificationPolicyType = "validateRules"
	TypeExecPostureScan            NotificationPolicyType = "execPostureScan"
	TypeUpdateRules                NotificationPolicyType = "updateRules"
	TypeRunKubescapeJob            NotificationPolicyType = "runKubescapeJob"
	TypeRunKubescape               NotificationPolicyType = "kubescapeScan"
	TypeSetKubescapeCronJob        NotificationPolicyType = "setKubescapeCronJob"
	TypeUpdateKubescapeCronJob     NotificationPolicyType = "updateKubescapeCronJob"
	TypeDeleteKubescapeCronJob     NotificationPolicyType = "deleteKubescapeCronJob"
	TypeSetVulnScanCronJob         NotificationPolicyType = "setVulnScanCronJob"
	TypeUpdateVulnScanCronJob      NotificationPolicyType = "updateVulnScanCronJob"
	TypeDeleteVulnScanCronJob      NotificationPolicyType = "deleteVulnScanCronJob"
	TypeUpdateWorkload             NotificationPolicyType = "update"
	TypeAttachWorkload             NotificationPolicyType = "Attach"
	TypeRemoveWorkload             NotificationPolicyType = "remove"
	TypeDetachWorkload             NotificationPolicyType = "Detach"
	TypeWorkloadIncompatible       NotificationPolicyType = "Incompatible"
	TypeSignWorkload               NotificationPolicyType = "sign"
	TypeClusterUnregistered        NotificationPolicyType = "unregistered"
	TypeReplaceHeadersInWorkload   NotificationPolicyType = "ReplaceHeaders"
	TypeImageUnreachableInWorkload NotificationPolicyType = "ImageUnreachable"
	TypeInjectToWorkload           NotificationPolicyType = "inject"
	TypeRestartWorkload            NotificationPolicyType = "restart"
	TypeEncryptSecret              NotificationPolicyType = "encryptSecret"
	TypeDecryptSecret              NotificationPolicyType = "decryptSecret"
	TypeScanImages                 NotificationPolicyType = "scan"
	TypeScanRegistry               NotificationPolicyType = "scanRegistry"
	TypeSetRegistryScanCronJob     NotificationPolicyType = "setRegistryScanCronJob"
	TypeUpdateRegistryScanCronJob  NotificationPolicyType = "updateRegistryScanCronJob"
	TypeDeleteRegistryScanCronJob  NotificationPolicyType = "deleteRegistryScanCronJob"
)
