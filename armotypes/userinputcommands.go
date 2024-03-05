package armotypes

const (
	// msgKeys
	CommandKey            = "command"
	UserInputMessageIDKey = "userInputMessageID"
	ErrorKey              = "error"
	ReplyTopicKey         = "replyTopic"
	TimestampKey          = "timestamp"
)

type UserInputCommand string

const (
	UserInputCommandDeleteCluster                               = UserInputCommand("delete-cluster")
	UserInputCommandDeleteAllCustomerData                       = UserInputCommand("delete-all-customer-data")
	UserInputCommandMarkForDeletionRegistryScan                 = UserInputCommand("mark-for-deletion-registry-scan")
	UserInputCommandStoreContainerScanningSummaryStub           = UserInputCommand("store-container-scanning-summary-stub")
	UserInputCommandMarkForDeletionKubernetesResources          = UserInputCommand("mark-for-deletion-kubernetes-resources")
	UserInputCommandMarkForDeletionContainerScanSummary         = UserInputCommand("mark-for-deletion-container-scan-summary")
	UserInputCommandMarkForDeletionClusterPostureReport         = UserInputCommand("mark-for-deletion-cluster-posture-report")
	UserInputCommandMarkForDeletionRepositoryPostureReport      = UserInputCommand("mark-for-deletion-repository-posture-report")
	UserInputCommandMarkForDeletionRegistryScanVulnerabilities  = UserInputCommand("mark-for-deletion-registry-scan-vulnerabilities")
	UserInputCommandMarkForDeletionContainerScanVulnerabilities = UserInputCommand("mark-for-deletion-container-scan-vulnerabilities")
	UserInputCommandStoreExceptionSecurityRisk                  = UserInputCommand("store-exception-security-risk")
	UserInputCommandUpdateExceptionSecurityRisk                 = UserInputCommand("update-exception-security-risk")
	UserInputCommandMarkForDeletionExceptionSecurityRisk        = UserInputCommand("mark-for-deletion-exception-security-risk")
)
