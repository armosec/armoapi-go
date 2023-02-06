package armotypes

import (
	"strings"
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

type LicenseType string

const (
	LicenseTypeFree       LicenseType = "Free"
	LicenseTypeTeam                   = "Team"
	LicenseTypeEnterprise             = "Enterprise"
)

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
	SubscriptionDate string `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate    string `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
}

// hold information of a single subscription.
type Subscription struct {

	// -------- Stripe specific properties ------- //

	// Stripe internal customer ID, usually generated on subscription creation.
	StripeCustomerID string `json:"stripe_customer_id,omitempty" bson:"stripe_customer_id,omitempty"`

	// Stripe subscription id.
	StripeSubscriptionID string `json:"stripe_subscription_id,omitempty" bson:"stripe_subscription_id,omitempty"`

	// -------- Stripe 'borrowed' properties, to be used also for none-stripe plans ------- //

	// Stripe subscription status, optional values: incomplete, incomplete_expired, trialing, active, past_due, canceled, or unpaid.
	SubscriptionStatus string `json:"subscription_status,omitempty" bson:"subscription_status,omitempty"`

	// Stripe The most recent invoice this subscription has generated.
	LatestInvoice string `json:"latest_invoice,omitempty" bson:"latest_invoice,omitempty"`

	// determine whether a subscription that has a status of active is scheduled to be canceled at the end of the current period.
	CancelAtPeriodEnd bool `json:"cancel_at_period_end,omitempty" bson:"cancel_at_period_end,omitempty"`

	// End of the current period that the subscription has been invoiced for. At the end of this period, a new invoice will be created.
	CurrentPeriodStart int64 `json:"current_period_start,omitempty" bson:"current_period_start,omitempty"`

	// End of the current period that the subscription has been invoiced for. At the end of this period, a new invoice will be created.
	CurrentPeriodEnd int64 `json:"current_period_end,omitempty" bson:"current_period_end,omitempty"`

	// If the subscription has a trial, the end of that trial.
	TrialEnd int64 `json:"trial_end,omitempty" bson:"trial_end,omitempty"`

	// ------- Internal properties associated with plan --------- //

	// monthly average of daily sum of max scanned Worker Nodes per cluster per day
	NumNodes int `json:"num_nodes,omitempty" bson:"num_nodes,omitempty"`

	// can be "free", "team" or "enterprise"
	LicenseType LicenseType `json:"LicenseType,omitempty" bson:"LicenseType,omitempty"`
}

type PortalCustomer struct {
	PortalBase       `json:",inline" bson:"inline"`
	Description      string `json:"description,omitempty" bson:"description,omitempty,omitempty"`
	SubscriptionDate string `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate    string `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
	Email            string `json:"email,omitempty" bson:"email,omitempty"`

	// DEPRECATED - moved to subscription
	LicenseType string `json:"license_type,omitempty" bson:"license_type,omitempty"`

	// DEPRECATED - moved to subscription
	SubscriptionExpiration string `json:"subscription_expiration,omitempty" bson:"subscription_expiration,omitempty"`

	// DEPRECATED
	InitialLicenseType string `json:"initial_license_type,omitempty" bson:"initial_license_type,omitempty"`

	NotificationsConfig *NotificationsConfig `json:"notifications_config,omitempty" bson:"notifications_config,omitempty"`
	State               *CustomerState       `json:"state,omitempty" bson:"state,omitempty"`

	// Paid/free subscriptions information
	ActiveSubscription      *Subscription  `json:"active_subscription,omitempty" bson:"active_subscription,omitempty"`
	HistoricalSubscriptions []Subscription `json:"historical_subscriptions,omitempty" bson:"historical_subscriptions,omitempty"`
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

// CustomerState holds the state of the customer, used for UI purposes
type CustomerState struct {
	Onboarding     *CustomerOnboarding      `json:"onboarding,omitempty" bson:"onboarding,omitempty"`
	GettingStarted *GettingStartedChecklist `json:"gettingStarted,omitempty" bson:"gettingStarted,omitempty"`
}
