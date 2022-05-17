package apis

// WebsocketScanCommand api
const (
	WebsocketScanCommandVersion string = "v1"
	WebsocketScanCommandPath    string = "scanImage"
	DBCommandPath               string = "DBCommand"
)

// commands send via websocket
// DEPRECATED
const (
	UPDATE            string = "update"
	ATTACH            string = "Attach"
	REMOVE            string = "remove"
	DETACH            string = "Detach"
	INCOMPATIBLE      string = "Incompatible"
	REPLACE_HEADERS   string = "ReplaceHeaders"
	IMAGE_UNREACHABLE string = "ImageUnreachable"
	SIGN              string = "sign"
	UNREGISTERED      string = "unregistered"
	INJECT            string = "inject"
	RESTART           string = "restart"
	ENCRYPT           string = "encryptSecret"
	DECRYPT           string = "decryptSecret"
	SCAN              string = "scan"
	SCAN_REGISTRY     string = "scanRegistry"
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
	TypeDecryptSecret              NotificationPolicyType = "encryptSecret"
	TypeScanImages                 NotificationPolicyType = "scan"
	TypeScanRegistry               NotificationPolicyType = "scanRegistry"
)

// KubescapeJobParams kubescape triggering cronJob params
type KubescapeJobParams struct {
	Name            string `json:"name,omitempty"`
	ID              string `json:"id,omitempty"`
	ClusterName     string `json:"clusterName"`
	FrameworkName   string `json:"frameworkName"`
	CronTabSchedule string `json:"cronTabSchedule,omitempty"`
	JobID           string `json:"jobID,omitempty"`
}
