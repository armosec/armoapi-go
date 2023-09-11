package armotypes

import (
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

// PortalCluster holds cluster data from portal BE
type PortalCluster struct {
	PortalBase       `json:",inline" bson:"inline"`
	SubscriptionDate string            `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate    string            `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
	InstallationData *InstallationData `json:"installationData" bson:"installationData,omitempty"`
}

type ClusterAttackChainState struct {
	PortalBase               `json:",inline" bson:"inline"`
	CreationTime             string `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	ClusterName              string `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
	CustomerGUID             string `json:"customerGUID,omitempty" bson:"customerGUID,omitempty"`
	LastPostureScanTriggered string `json:"lastPostureScanTriggered,omitempty" bson:"lastPostureScanTriggered,omitempty"`
	LastTimeEngineCompleted  string `json:"lastTimeEngineCompleted,omitempty" bson:"lastTimeEngineCompleted,omitempty"`
}

type RelevantImageVulnerabilitiesConfiguration string

const (
	RelevantImageVulnerabilitiesConfigurationEnable  RelevantImageVulnerabilitiesConfiguration = "enable"
	RelevantImageVulnerabilitiesConfigurationDisable RelevantImageVulnerabilitiesConfiguration = "disable"
	RelevantImageVulnerabilitiesConfigurationDetect  RelevantImageVulnerabilitiesConfiguration = "detect"
)

type InstallationData struct {
	ClusterName                               string                                    `json:"clusterName,omitempty" bson:"clusterName,omitempty"`                                                             // cluster name defined manually or from the cluster context
	ClusterShortName                          string                                    `json:"clusterShortName,omitempty" bson:"clusterShortName,omitempty"`                                                   // cluster short name enriched from the cluster name by BE
	StorageEnabled                            *bool                                     `json:"storage,omitempty" bson:"storage,omitempty"`                                                                     // storage configuration (enabled/disabled)
	RelevantImageVulnerabilitiesEnabled       *bool                                     `json:"relevantImageVulnerabilitiesEnabled,omitempty" bson:"relevantImageVulnerabilitiesEnabled,omitempty"`             // relevancy actual state (enabled/disabled)
	RelevantImageVulnerabilitiesConfiguration RelevantImageVulnerabilitiesConfiguration `json:"relevantImageVulnerabilitiesConfiguration,omitempty" bson:"relevantImageVulnerabilitiesConfiguration,omitempty"` // relevancy configuration defined user
	Namespace                                 string                                    `json:"namespace,omitempty" bson:"namespace,omitempty"`                                                                 // namespace to deploy the components
	ImageVulnerabilitiesScanningEnabled       *bool                                     `json:"imageVulnerabilitiesScanningEnabled,omitempty" bson:"imageVulnerabilitiesScanningEnabled,omitempty"`             // image scanning configuration (enabled/disabled)
	PostureScanEnabled                        *bool                                     `json:"postureScanEnabled,omitempty" bson:"postureScanEnabled,omitempty"`                                               // posture configuration (enabled/disabled)
	OtelCollectorEnabled                      *bool                                     `json:"otelCollector,omitempty" bson:"otelCollector,omitempty"`                                                         // otel collector configuration (enabled/disabled)
	ClusterProvider                           string                                    `json:"clusterProvider,omitempty" bson:"clusterProvider,omitempty"`                                                     // cluster provider (aws/azure/gcp)

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
	Onboarding           *CustomerOnboarding      `json:"onboarding,omitempty" bson:"onboarding,omitempty"`
	GettingStarted       *GettingStartedChecklist `json:"gettingStarted,omitempty" bson:"gettingStarted,omitempty"`
	NodeUsage            *NodeUsage               `json:"nodeUsage,omitempty" bson:"nodeUsage,omitempty"`
	AttackChainsLastScan string                   `json:"attackChainsLastScan,omitempty" bson:"attackChainsLastScan,omitempty"`
}

type User struct {
	DismissedBanners map[string]Banner `json:"dismissedBanners,omitempty" bson:"dismissedBanners,omitempty"` // map of bannerID to Banner
}

type Banner struct {
	CustomerGUID string `json:"customerGUID,omitempty" bson:"customerGUID,omitempty"` // customerGUID of the account which clicked the banner
	ScanID       string `json:"scanID,omitempty" bson:"scanID,omitempty"`             // for detailed view, unique key for banner is combination of scanID and bannerID
}
