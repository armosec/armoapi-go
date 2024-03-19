package notifications

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (ap *NotificationParams) SetDriftPercentage(percentage int) {
	ap.DriftPercentage = &percentage
}

func (ap *NotificationParams) SetMinSeverity(severity int) {
	ap.MinSeverity = &severity
}

func (ac *AlertChannel) GetAlertConfig(notificationType NotificationType) *AlertConfig {
	for _, alert := range ac.Alerts {
		if alert.NotificationType == notificationType {
			return &alert
		}
	}
	return nil
}

func (ac *AlertChannel) IsEqualOrGreaterThanMinSeverity(severity int, notificationType NotificationType) bool {
	if ac.Alerts == nil {
		return true
	}

	for _, alert := range ac.Alerts {
		if alert.IsEnabled() && alert.NotificationType == notificationType {
			if alert.Parameters.MinSeverity != nil {
				if *alert.Parameters.MinSeverity > severity {
					return false
				}
			}
		}

	}

	return true
}

func (ac *AlertChannel) AddAlertConfig(config AlertConfig) error {
	if config.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	for _, alert := range ac.Alerts {
		if alert.NotificationType == config.NotificationType {
			return fmt.Errorf("alert config for notification type %s already exists", config.NotificationType)
		}
	}
	ac.Alerts = append(ac.Alerts, config)
	return nil
}

func (ac *AlertChannel) IsInScope(cluster, namespace string) bool {
	if ac.Scope == nil {
		return true
	}
	for _, scope := range ac.Scope {
		if scope.IsInScope(cluster, namespace) {
			return true
		}
	}
	return false
}

func (ac *AlertChannel) IsNotificationTypeEnabled(notificationType NotificationType) bool {
	if ac.Alerts == nil {
		return false
	}

	config := ac.GetAlertConfig(notificationType)
	return config != nil && config.IsEnabled()
}

func (ac *AlertConfig) IsEnabled() bool {
	if ac.Disabled == nil {
		return true
	}
	return !*ac.Disabled
}

func (ac *AlertScope) IsInScope(cluster, namespace string) bool {
	if ac.Cluster == "" {
		//no scope defined
		return true
	}
	if ac.Cluster != cluster {
		return false
	}
	if namespace == "" {
		return true
	}
	if len(ac.Namespaces) == 0 {
		return true
	}
	return slices.Contains(ac.Namespaces, namespace)
}

func (nc *NotificationsConfig) IsInScope(cluster, namespace string) bool {
	for _, typesChannels := range nc.AlertChannels {
		for _, alertChannel := range typesChannels {
			if alertChannel.IsInScope(cluster, namespace) {
				return true
			}
		}
	}
	return false
}

func (nc *NotificationsConfig) GetProviderChannels(provider ChannelProvider) []AlertChannel {
	if nc.AlertChannels == nil {
		return nil
	}
	return nc.AlertChannels[provider]
}

func (nc *NotificationsConfig) GetAllChannels() []AlertChannel {
	if len(nc.AlertChannels) == 0 {
		return nil
	}
	var channels []AlertChannel
	for i := range nc.AlertChannels {
		channels = append(channels, nc.AlertChannels[i]...)
	}
	return channels
}

func (nc *NotificationsConfig) GetAlertConfigurations(notificationType NotificationType) []AlertConfig {
	alerts := make([]AlertConfig, 0)
	for _, typesChannels := range nc.AlertChannels {
		for _, alertChannel := range typesChannels {
			if config := alertChannel.GetAlertConfig(notificationType); config != nil {
				alerts = append(alerts, *config)
			}
		}
	}
	return alerts
}

func (nc *NotificationsConfig) AddLatestPushReport(report *PushReport) {
	if report == nil {
		return
	}
	if nc.LatestPushReports == nil {
		nc.LatestPushReports = make(map[string]*PushReport, 0)
	}
	nc.LatestPushReports[fmt.Sprintf("%s_%s", report.Cluster, report.ScanType)] = report
}

func (nc *NotificationsConfig) GetLatestPushReport(cluster string, scanType ScanType) *PushReport {
	if val, ok := nc.LatestPushReports[fmt.Sprintf("%s_%s", cluster, scanType)]; ok {
		return val
	}
	return nil
}

func (nc *NotificationsConfig) GetAlertChannelByCollaborationID(collaborationId string) (*AlertChannel, error) {
	providerToChannels := nc.AlertChannels
	for _, alertChannels := range providerToChannels {
		for _, alertChannel := range alertChannels {
			if alertChannel.CollaborationConfigGUID == collaborationId {
				return &alertChannel, nil
			}
		}
	}
	return nil, fmt.Errorf("alert channel with collaboration id %s not found", collaborationId)
}

func (nc *NotificationsConfig) RemoveAlertChannel(collaborationId string) error {
	for key, alertChannels := range nc.AlertChannels {
		for i, alertChannel := range alertChannels {
			if alertChannel.CollaborationConfigGUID == collaborationId {
				nc.AlertChannels[key] = append(alertChannels[:i], alertChannels[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("alert channel with collaboration id %s not found", collaborationId)
}

func (nc *NotificationsConfig) RemoveProviderConfig(provider ChannelProvider) error {
	if _, exists := nc.AlertChannels[provider]; exists {
		nc.AlertChannels[provider] = make([]AlertChannel, 0)
		return nil
	}
	return fmt.Errorf("provider with identifier %v not found", provider)
}

func (nci *NotificationConfigIdentifier) Validate() error {
	if slices.Contains(notificationTypes, nci.NotificationType) {
		return nil
	}
	if nci.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	return fmt.Errorf("invalid notification type: %s", nci.NotificationType)
}

type TopCtrlItem struct {
	ControlID            string           `json:"id" bson:"id"`
	ControlGUID          string           `json:"guid" bson:"guid"`
	Name                 string           `json:"name" bson:"name"`
	Remediation          string           `json:"remediation" bson:"remediation"`
	Description          string           `json:"description" bson:"description"`
	ClustersCount        int64            `json:"clustersCount" bson:"clustersCount"`
	SeverityOverall      int64            `json:"severityOverall" bson:"severityOverall"`
	BaseScore            int64            `json:"baseScore" bson:"baseScore"`
	Clusters             []TopCtrlCluster `json:"clusters" bson:"clusters"`
	TotalFailedResources int64            `json:"-"`
}

type TopCtrlCluster struct {
	Name               string `json:"name" bson:"name"`
	ResourcesCount     int64  `json:"resourcesCount" bson:"resourcesCount"`
	ReportGUID         string `json:"reportGUID" bson:"reportGUID"`
	TopFailedFramework string `json:"topFailedFramework" bson:"topFailedFramework"`
}

func (t *TopCtrlItem) GetTotalFailedResources() int64 {
	var totalFailedResources int64
	for _, c := range t.Clusters {
		totalFailedResources += c.ResourcesCount
	}
	return totalFailedResources
}
