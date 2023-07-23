package configservice

import (
	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/notifications"
)

type PortalCustomer struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	Description          string `json:"description,omitempty" bson:"description,omitempty,omitempty"`
	SubscriptionDate     string `json:"subscription_date,omitempty" bson:"subscription_date,omitempty"`
	LastLoginDate        string `json:"last_login_date,omitempty" bson:"last_login_date,omitempty"`
	Email                string `json:"email,omitempty" bson:"email,omitempty"`
	// customizable field that overrides the default max
	MaxFreeNodes int `json:"maxFreeNodes,omitempty" bson:"maxFreeNodes,omitempty"`

	// DEPRECATED - moved to subscription
	LicenseType string `json:"license_type,omitempty" bson:"license_type,omitempty"`

	// DEPRECATED - moved to subscription
	SubscriptionExpiration string `json:"subscription_expiration,omitempty" bson:"subscription_expiration,omitempty"`

	// DEPRECATED
	InitialLicenseType string `json:"initial_license_type,omitempty" bson:"initial_license_type,omitempty"`

	NotificationsConfig *notifications.NotificationsConfig `json:"notifications_config,omitempty" bson:"notifications_config,omitempty"`
	State               *armotypes.CustomerState           `json:"state,omitempty" bson:"state,omitempty"`

	OpenAiRequestCount int `json:"open_ai_request_count,omitempty" bson:"open_ai_request_count,omitempty"`

	// Paid/free subscriptions information
	ActiveSubscription      *armotypes.Subscription  `json:"activeSubscription,omitempty" bson:"activeSubscription,omitempty"`
	HistoricalSubscriptions []armotypes.Subscription `json:"historicalSubscriptions,omitempty" bson:"historicalSubscriptions,omitempty"`
}
