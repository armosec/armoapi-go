package armotypes

import (
	"strings"

	stripe "github.com/stripe/stripe-go/v74"
)

const (
	CustomerGuidQuery   = "customerGUID"
	ClusterNameQuery    = "cluster"
	DatacenterNameQuery = "datacenter"
	NamespaceQuery      = "namespace"
	ProjectQuery        = "project"
	WlidQuery           = "wlid"
	SidQuery            = "sid"
)

const (
	SubscriptionStatusIncomplete        = string(stripe.SubscriptionStatusIncomplete)
	SubscriptionStatusIncompleteExpired = string(stripe.SubscriptionStatusIncompleteExpired)
	SubscriptionStatusTrialing          = string(stripe.SubscriptionStatusTrialing)
	SubscriptionStatusActive            = string(stripe.SubscriptionStatusActive)
	SubscriptionStatusPastDue           = string(stripe.SubscriptionStatusPastDue)
	SubscriptionStatusCanceled          = string(stripe.SubscriptionStatusCanceled)
	SubscriptionStatusUnpaid            = string(stripe.SubscriptionStatusUnpaid)
)

type LicenseType string

const (
	LicenseTypeFree       LicenseType = "Free"
	LicenseTypeTeam       LicenseType = "Team"
	LicenseTypeEnterprise LicenseType = "Enterprise"
)

type CustomerAccessStatus string

const (
	PayingCustomer  CustomerAccessStatus = "paying"
	FreeCustomer    CustomerAccessStatus = "free"
	TrialCustomer   CustomerAccessStatus = "trial"
	BlockedCustomer CustomerAccessStatus = "blocked"
)

var ActiveSubscriptionStatuses = []string{SubscriptionStatusIncomplete, SubscriptionStatusTrialing, SubscriptionStatusActive}

// PortalBase holds basic items data from portal BE
type PortalBase struct {
	GUID        string                 `json:"guid" bson:"guid"`
	Name        string                 `json:"name" bson:"name"`
	Attributes  map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"` // could be string
	UpdatedTime string                 `json:"updatedTime,omitempty" bson:"updatedTime,omitempty"`
}

// Type of the designator
//
// swagger:enum DesignatorType
type DesignatorType string

// Supported designators
const (
	DesignatorAttributes DesignatorType = "Attributes"
	DesignatorAttribute  DesignatorType = "Attribute" // Deprecated
	// WorkloadID format.
	//
	// Has two formats:
	//  1. Kubernetes format: wlid://cluster-<cluster>/namespace-<namespace>/<kind>-<name>
	//  2. Native format: wlid://datacenter-<datacenter>/project-<project>/native-<name>
	DesignatorWlid DesignatorType = "Wlid"
	// A WorkloadID wildcard expression.
	//
	// A wildcard expression that includes a cluster:
	//
	//  wlid://cluster-<cluster>/
	//
	// An expression that includes a cluster and namespace (filters out all other namespaces):
	//
	//  wlid://cluster-<cluster>/namespace-<namespace>/
	DesignatorWildWlid      DesignatorType = "WildWlid"
	DesignatorWlidContainer DesignatorType = "WlidContainer"
	DesignatorWlidProcess   DesignatorType = "WlidProcess"
	DesignatorSid           DesignatorType = "Sid" // secret id
)

func (dt DesignatorType) ToLower() DesignatorType {
	return DesignatorType(strings.ToLower(string(dt)))
}

// attributes
const (
	DesignatorsToken       = "designators"
	AttributeCustomerGUID  = "customerGUID"
	AttributeRegistryName  = "registryName"
	AttributeRepository    = "repository"
	AttributeTag           = "tag"
	AttributeCluster       = "cluster"
	AttributeNamespace     = "namespace"
	AttributeKind          = "kind"
	AttributeName          = "name"
	AttributeContainerName = "containerName"
	AttributeApiVersion    = "apiVersion"
	AttributeWorkloadHash  = "workloadHash"
	AttributeIsIncomplete  = "isIncomplete"
	AttributeSensor        = "sensor"
	AttributePath          = "path"
	AttributeResourceID    = "resourceID"
)

// Repository scan related attributes
const (
	AttributeRepoName      = "repoName"
	AttributeRepoOwner     = "repoOwner"
	AttributeRepoHash      = "repoHash"
	AttributeBranchName    = "branch"
	AttributeDefaultBranch = "defaultBranch"
	AttributeProvider      = "provider"
	AttributeRemoteURL     = "remoteURL"

	AttributeLastCommitHash     = "lastCommitHash"
	AttributeLastCommitterName  = "lastCommitterName"
	AttributeLastCommitterEmail = "lastCommitterEmail"
	AttributeLastCommitTime     = "lastCommitTime"

	AttributeFilePath          = "filePath"
	AttributeFileType          = "fileType"
	AttributeFileDir           = "fileDirectory"
	AttributeFileUrl           = "fileUrl"
	AttributeFileHelmChartName = "fileHelmChartName"

	AttributeLastFileCommitHash     = "lastFileCommitHash"
	AttributeLastFileCommitterName  = "lastFileCommitterName"
	AttributeLastFileCommitterEmail = "LastFileCommitterEmail"
	AttributeLastFileCommitTime     = "lastFileCommitTime"

	AttributeUseHTTP       = "useHTTP"
	AttributeSkipTLSVerify = "skipTLSVerify"
)

// rego-library attributes
const (
	AttributeImageScanRelated     = "imageScanRelated"
	AttributeImageRelatedControls = "imageRelatedControls"
	AttributeHostSensorRule       = "hostSensorRule"
	AttributeHostSensor           = "hostSensor"
)

// PortalDesignator represents a single designation option
type PortalDesignator struct {
	DesignatorType DesignatorType `json:"designatorType" bson:"designatorType"`
	// A specific Workload ID
	WLID string `json:"wlid,omitempty" bson:"wlid,omitempty"`
	// An expression that describes applicable workload IDs
	WildWLID string `json:"wildwlid,omitempty" bson:"wildwlid,omitempty"`
	// A specific Secret ID
	SID string `json:"sid,omitempty" bson:"sid,omitempty"`
	// Attributes that describe the targets
	Attributes map[string]string `json:"attributes" bson:"attributes"`
}

// Worker nodes attribute related consts
const (
	AttributeWorkerNodes             = "workerNodes"
	WorkerNodesmax                   = "max"
	WorkerNodeslastReported          = "lastReported"
	WorkerNodeslastReportDate        = "lastReportDate"
	WorkerNodesmaxPerMonth           = "maxPerMonth"
	WorkerNodesmaxReportGUID         = "maxReportGUID"
	WorkerNodesmaxPerMonthReportGUID = "maxPerMonthReportGUID"
	WorkerNodeslastReportGUID        = "lastReportGUID"
)

// PortalCluster holds cluster data from portal BE
type PortalCluster struct {
	PortalBase       `json:",inline" bson:"inline"`
	ClusterShortName string            `json:"clusterShortName,omitempty" bson:"clusterShortName,omitempty"`
	SubscriptionDate string            `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate    string            `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
	InstallationData *InstallationData `json:"installationData" bson:"installationData,omitempty"`
}

type InstallationData struct {
	ClusterName                         string `json:"clusterName,omitempty" bson:"clusterName,omitempty"`                                                 // cluster name defined manually or from the cluster context
	StorageEnabled                      *bool  `json:"storage,omitempty" bson:"storage,omitempty"`                                                         // storage configuration (enabled/disabled)
	RelevantImageVulnerabilitiesEnabled *bool  `json:"relevantImageVulnerabilitiesEnabled,omitempty" bson:"relevantImageVulnerabilitiesEnabled,omitempty"` // relevancy configuration (enabled/disabled)
	Namespace                           string `json:"namespace,omitempty" bson:"namespace,omitempty"`                                                     // namespace to deploy the components
	ImageVulnerabilitiesScanningEnabled *bool  `json:"imageVulnerabilitiesScanningEnabled,omitempty" bson:"imageVulnerabilitiesScanningEnabled,omitempty"` // image scanning configuration (enabled/disabled)
	PostureScanEnabled                  *bool  `json:"postureScanEnabled,omitempty" bson:"postureScanEnabled,omitempty"`                                   // posture configuration (enabled/disabled)
	OtelCollectorEnabled                *bool  `json:"otelCollector,omitempty" bson:"otelCollector,omitempty"`                                             // otel collector configuration (enabled/disabled)
}

// hold information of a single subscription.
type Subscription struct {

	// -------- Stripe specific properties ------- //

	// Stripe internal customer ID, usually generated on subscription creation.
	StripeCustomerID string `json:"stripeCustomerID,omitempty" bson:"stripeCustomerID,omitempty"`

	// Stripe subscription id.
	StripeSubscriptionID string `json:"stripeSubscriptionID,omitempty" bson:"stripeSubscriptionID,omitempty"`

	// -------- Stripe 'borrowed' properties, to be used also for none-stripe plans ------- //

	// Stripe subscription status, optional values: incomplete, incomplete_expired, trialing, active, past_due, canceled, or unpaid.
	SubscriptionStatus string `json:"subscriptionStatus,omitempty" bson:"subscriptionStatus,omitempty"`

	// Date when the subscription was first created. The date might differ from the created date due to backdating
	StartDate int64 `json:"startDate,omitempty" bson:"startDate,omitempty"`

	// Stripe The most recent invoice this subscription has generated.
	LatestInvoice string `json:"latestInvoice,omitempty" bson:"latestInvoice,omitempty"`

	// determine whether a subscription that has a status of active is scheduled to be canceled at the end of the current period.
	CancelAtPeriodEnd *bool `json:"cancelAtPeriodEnd,omitempty" bson:"cancelAtPeriodEnd,omitempty"`

	// End of the current period that the subscription has been invoiced for. At the end of this period, a new invoice will be created.
	CurrentPeriodStart int64 `json:"currentPeriodStart,omitempty" bson:"currentPeriodStart,omitempty"`

	// End of the current period that the subscription has been invoiced for. At the end of this period, a new invoice will be created.
	CurrentPeriodEnd int64 `json:"currentPeriodEnd,omitempty" bson:"currentPeriodEnd,omitempty"`

	// If the subscription has a trial, the end of that trial.
	TrialEnd int64 `json:"trialEnd,omitempty" bson:"trialEnd,omitempty"`

	// ------- Internal properties associated with plan --------- //

	// monthly average of daily sum of max scanned Worker Nodes per cluster per day
	NumNodes int `json:"numNodes,omitempty" bson:"numNodes,omitempty"`

	// can be "free", "team" or "enterprise"
	LicenseType LicenseType `json:"licenseType,omitempty" bson:"licenseType,omitempty"`
}

type PortalCustomer struct {
	PortalBase       `json:",inline" bson:"inline"`
	Description      string `json:"description,omitempty" bson:"description,omitempty,omitempty"`
	SubscriptionDate string `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate    string `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`
	// customizable field that overrides the default max
	MaxFreeNodes int `json:"maxFreeNodes,omitempty" bson:"maxFreeNodes,omitempty"`

	// DEPRECATED - moved to subscription
	LicenseType string `json:"license_type,omitempty" bson:"license_type,omitempty"`

	// DEPRECATED - moved to subscription
	SubscriptionExpiration string `json:"subscription_expiration,omitempty" bson:"subscription_expiration,omitempty"`

	// DEPRECATED
	InitialLicenseType string `json:"initial_license_type,omitempty" bson:"initial_license_type,omitempty"`

	NotificationsConfig *NotificationsConfig `json:"notifications_config,omitempty" bson:"notifications_config,omitempty"`
	State               *CustomerState       `json:"state,omitempty" bson:"state,omitempty"`

	OpenAiRequestCount int `json:"open_ai_request_count,omitempty" bson:"open_ai_request_count,omitempty"`

	// Paid/free subscriptions information
	ActiveSubscription      *Subscription  `json:"activeSubscription,omitempty" bson:"activeSubscription,omitempty"`
	HistoricalSubscriptions []Subscription `json:"historicalSubscriptions,omitempty" bson:"historicalSubscriptions,omitempty"`
}

type PortalRepository struct {
	PortalBase   `json:",inline" bson:"inline"`
	CreationDate string `json:"creationDate,omitempty" bson:"creationDate,omitempty"`
	Provider     string `json:"provider,omitempty" bson:"provider,omitempty"`
	Owner        string `json:"owner,omitempty" bson:"owner,omitempty"`
	RepoName     string `json:"repoName,omitempty" bson:"repoName,omitempty"`
	BranchName   string `json:"branchName,omitempty" bson:"branchName,omitempty"`
}

type PortalRegistryCronJob struct {
	PortalBase      `json:",inline" bson:"inline"`
	RegistryInfo    `json:",inline" bson:"inline"`
	CreationDate    string       `json:"creationDate,omitempty" bson:"creationDate,omitempty"`
	ID              string       `json:"id,omitempty" bson:"id,omitempty"`
	ClusterName     string       `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
	CronTabSchedule string       `json:"cronTabSchedule,omitempty" bson:"cronTabSchedule,omitempty"`
	Repositories    []Repository `json:"repositories,omitempty" bson:"repositories,omitempty"`
}

type CustomerOnboarding struct {
	Completed   *bool    `json:"completed,omitempty" bson:"completed,omitempty"`     // user completed the onboarding
	CompanySize *string  `json:"companySize,omitempty" bson:"companySize,omitempty"` // user company size
	Role        *string  `json:"role,omitempty" bson:"role,omitempty"`               // user role
	OrgName     *string  `json:"orgName,omitempty" bson:"orgName,omitempty"`         // user organization name
	Interests   []string `json:"interests,omitempty" bson:"interests,omitempty"`     // user interests
}

type GettingStartedChecklist struct {
	// indicates if the user has dismissed the checklist
	GettingStartedDismissed *bool `json:"gettingStartedDismissed,omitempty" bson:"gettingStartedDismissed,omitempty"`
	// checklist items
	EverConnectedCluster   *bool `json:"everConnectedCluster,omitempty" bson:"everConnectedCluster,omitempty"`
	EverScannedRepository  *bool `json:"everScannedRepository,omitempty" bson:"everScannedRepository,omitempty"`
	EverScannedRegistry    *bool `json:"everScannedRegistry,omitempty" bson:"everScannedRegistry,omitempty"`
	EverCollaborated       *bool `json:"everCollaborated,omitempty" bson:"everCollaborated,omitempty"`
	EverInvitedTeammate    *bool `json:"everInvitedTeammate,omitempty" bson:"everInvitedTeammate,omitempty"`
	EverUsedRbacVisualizer *bool `json:"everUsedRbacVisualizer,omitempty" bson:"everUsedRbacVisualizer,omitempty"`
}

type NodeUsage struct {
	// max sum of nodes across all clusters ever scanned on one day
	MaxNodesSumEver int `json:"maxNodesSumEver,omitempty" bson:"maxNodesSumEver,omitempty"`
	// date of MaxNodesSumEver
	MaxNodesSumDate string `json:"maxNodesSumDate,omitempty" bson:"maxNodesSumDate,omitempty"`
}

// CustomerState holds the state of the customer, used for UI purposes
type CustomerState struct {
	Onboarding     *CustomerOnboarding      `json:"onboarding,omitempty" bson:"onboarding,omitempty"`
	GettingStarted *GettingStartedChecklist `json:"gettingStarted,omitempty" bson:"gettingStarted,omitempty"`
	NodeUsage      *NodeUsage               `json:"nodeUsage,omitempty" bson:"nodeUsage,omitempty"`
}
