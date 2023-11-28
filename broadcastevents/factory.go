package broadcastevents

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/containerscan"
	"github.com/armosec/armoapi-go/notifications"
	"k8s.io/utils/ptr"
)

type IgnoreRuleType string
type IgnoreRuleExpirationType string

const (
	IgnoreRuleTypeMisconfiguration IgnoreRuleType = "misconfiguration"
	IgnoreRuleTypeVulnerability    IgnoreRuleType = "vulnerability"

	IgnoreRuleExpirationTypeNone IgnoreRuleExpirationType = "none"
	IgnoreRuleExpirationTypeDate IgnoreRuleExpirationType = "date"
	IgnoreRuleExpirationTypeFix  IgnoreRuleExpirationType = "fix"

	ignoreRuleEventPrefix = "IgnoreRule"

	AlertChannelPrefix = "AlertChannel"
)

func NewBaseEvent(customerGUID, eventName string, eventTime *time.Time) EventBase {
	var now time.Time
	if eventTime == nil {
		now = time.Now().UTC()
	} else {
		now = *eventTime
	}
	nowDate := now.Format("2006-01-02")
	nowMonth := now.Format("2006-01")
	_, nowWeekOfTheYear := now.ISOWeek()
	return EventBase{
		CustomerGUID:       customerGUID,
		EventName:          eventName,
		EventTime:          now,
		EventDate:          nowDate,
		EventMonth:         nowMonth,
		EventWeekOfTheYear: nowWeekOfTheYear,
	}
}

func NewAlertChannelDeletedEvent(customerGUID, name, provider string) AlertChannelEvent {
	return AlertChannelEvent{
		EventBase: NewBaseEvent(customerGUID, AlertChannelPrefix+"Deleted", nil),
		Name:      name,
		Type:      provider,
	}
}

func NewAlertChannelCreatedEvent(customerGUID, name string, channel notifications.AlertChannel) AlertChannelEvent {
	return newAlertChannelDetailedEvent(customerGUID, name, channel, "Created")
}

func NewAlertChannelUpdatedEvent(customerGUID, name string, channel notifications.AlertChannel) AlertChannelEvent {
	return newAlertChannelDetailedEvent(customerGUID, name, channel, "Updated")
}

func NewPostureExceptionEvent(customerGUID, changeMethod string, exception armotypes.PostureExceptionPolicy) IgnoreRuleEvent {
	expirationType := IgnoreRuleExpirationTypeNone
	if exception.ExpirationDate != nil && !exception.ExpirationDate.IsZero() {
		expirationType = IgnoreRuleExpirationTypeDate
	}
	resourcesCount := len(exception.Resources)
	var ids []string
	for _, policy := range exception.PosturePolicies {
		if policy.ControlID != "" {
			ids = append(ids, policy.ControlID)
		} else {
			ids = append(ids, policy.ControlName)
		}
	}
	return newIgnoreRuleEvent(customerGUID, changeMethod, IgnoreRuleTypeMisconfiguration, expirationType, ids, resourcesCount)
}

func NewVulnerabilityExceptionChangeEvent(customerGUID, changeMethod string, exception armotypes.VulnerabilityExceptionPolicy) IgnoreRuleEvent {
	expirationType := IgnoreRuleExpirationTypeNone
	if exception.ExpiredOnFix != nil && *exception.ExpiredOnFix {
		expirationType = IgnoreRuleExpirationTypeFix
	} else if exception.ExpirationDate != nil && !exception.ExpirationDate.IsZero() {
		expirationType = IgnoreRuleExpirationTypeDate
	}

	resourcesCount := len(exception.Designatores)
	var ids []string
	for _, policy := range exception.VulnerabilityPolicies {
		ids = append(ids, policy.Name)
	}
	return newIgnoreRuleEvent(customerGUID, changeMethod, IgnoreRuleTypeVulnerability, expirationType, ids, resourcesCount)
}

func NewFeatureFlagsEvent(customerGUID, userEmail, userName, userPreferredName string) LoginEvent {
	now := time.Now().UTC()
	nowDate := now.Format("2006-01-02")
	nowMonth := now.Format("2006-01")
	_, nowWeekOfTheYear := now.ISOWeek()
	return LoginEvent{
		EventBase: EventBase{
			CustomerGUID:       customerGUID,
			EventName:          "FeatureFlagsRequested",
			EventTime:          now,
			EventDate:          nowDate,
			EventMonth:         nowMonth,
			EventWeekOfTheYear: nowWeekOfTheYear,
		},
		Email:             userEmail,
		UserName:          userName,
		PreferredUserName: userPreferredName,
	}
}

func NewLoginEvent(customerGUID, email, name, preferredName string) LoginEvent {
	return LoginEvent{
		EventBase:         NewBaseEvent(customerGUID, "UserLoggedIn", nil),
		Email:             email,
		UserName:          name,
		PreferredUserName: preferredName,
	}
}

func NewClusterImageScanSessionStartedEvent(jobId, clusterName, customerId string, timeStarted time.Time) AggregationEvent {
	return AggregationEvent{
		EventBase:   NewBaseEvent(customerId, "ClusterImageScanSessionStarted", &timeStarted),
		JobID:       jobId,
		ClusterName: clusterName,
	}
}

func NewRegistryImageScanSessionStartedEvent(jobId string, customerId string, timeStarted time.Time) AggregationEvent {
	aggEvent := AggregationEvent{
		EventBase: NewBaseEvent(customerId, "RegistryImageScanSessionStarted", &timeStarted),
		JobID:     jobId,
	}
	return aggEvent
}

func NewImageScanEventHookNotify(customerGUID, jobId string, scanTime time.Time) AggregationEvent {
	return AggregationEvent{
		EventBase: NewBaseEvent(customerGUID, "ContainerImageScanSubmitted", &scanTime),
		JobID:     jobId,
	}
}

func NewGitRepositoryRiskScanEvent(customerGUID, jodID, reportGUID, clusterName string, eventTime time.Time) AggregationEvent {
	return AggregationEvent{
		EventBase:   NewBaseEvent(customerGUID, "RepoRiskScanSubmitted", &eventTime),
		JobID:       jodID,
		ReportGUID:  reportGUID,
		ClusterName: clusterName,
	}
}

func NewClusterRiskScanV2Event(customerGUID, jobID, reportGUID, clusterName, kubescapeVersion, cloudProvider, K8sVersion, helmVersion string, numOfWorkerNodes int, scanTime time.Time) AggregationEvent {
	aggEvent := AggregationEvent{
		EventBase:   NewBaseEvent(customerGUID, "RiskScanSubmitted", &scanTime),
		JobID:       jobID,
		ReportGUID:  reportGUID,
		ClusterName: clusterName,
	}
	aggEvent.KSVersion = kubescapeVersion
	aggEvent.HelmChartVersion = helmVersion
	aggEvent.WorkerNodesCount = numOfWorkerNodes
	aggEvent.K8sVendor = cloudProvider
	aggEvent.K8sVersion = K8sVersion
	return aggEvent
}

func NewHelmInstalledEvent(clusterName, customerGUID string, installationData *armotypes.InstallationData) HelmInstalledEvent {
	aggEvent := HelmInstalledEvent{
		EventBase:        NewBaseEvent(customerGUID, "HelmInstalled", nil),
		ClusterName:      clusterName,
		InstallationData: installationData,
	}
	return aggEvent
}

func NewKubescapePodRunningConditionErrorEvent(clusterName, customerId, objId, condition, reason, message string) PodInTroubleConditionEvent {
	// TODO: add objId, reason, meassage to the event
	return PodInTroubleConditionEvent{
		PodInTroubleEvent: PodInTroubleEvent{
			EventBase:   NewBaseEvent(customerId, "KubescapePodRunningError", nil),
			ClusterName: clusterName,
			ObjId:       objId,
			Reason:      reason,
			Message:     message,
		},
		Condition: condition,
	}
}

func NewKubescapePodRunningContainerErrorEvent(clusterName, customerId, objId, containerName, reason, message string, exitCode, restartCount int32) PodInTroubleContainerEvent {
	// TODO: add objId, reason, meassage to the event
	return PodInTroubleContainerEvent{
		PodInTroubleEvent: PodInTroubleEvent{
			EventBase:   NewBaseEvent(customerId, "KubescapePodRunningError", nil),
			ClusterName: clusterName,
			ObjId:       objId,
			Reason:      reason,
			Message:     message,
		},
		ContainerName: containerName,
		ExitCode:      exitCode,
		RestartCount:  restartCount,
	}
}

func NewKubescapePodPendingConditionErrorEvent(clusterName, customerId, objId, condition, reason, message string) PodInTroubleConditionEvent {
	return PodInTroubleConditionEvent{
		PodInTroubleEvent: PodInTroubleEvent{
			EventBase:   NewBaseEvent(customerId, "KubescapePodPendingError", nil),
			ClusterName: clusterName,
			ObjId:       objId,
			Reason:      reason,
			Message:     message,
		},
		Condition: condition,
	}
}

func NewKubescapePodPendingContainerErrorEvent(clusterName, customerId, objId, conainerName, reason, meassage string, exitCode, restartCount int32) PodInTroubleContainerEvent {
	return PodInTroubleContainerEvent{
		PodInTroubleEvent: PodInTroubleEvent{
			EventBase:   NewBaseEvent(customerId, "KubescapePodPendingError", nil),
			ClusterName: clusterName,
			ObjId:       objId,
			Reason:      reason,
			Message:     meassage,
		},
		ContainerName: conainerName,
		ExitCode:      exitCode,
		RestartCount:  restartCount,
	}
}

// helper
func newIgnoreRuleEvent(customerGUID string, changeMethod string, ruleType IgnoreRuleType, expirationType IgnoreRuleExpirationType, ids []string, resourcesCount int) IgnoreRuleEvent {
	var eventName string
	switch changeMethod {
	case http.MethodPost:
		eventName = ignoreRuleEventPrefix + "Created"
	case http.MethodPut:
		eventName = ignoreRuleEventPrefix + "Updated"
	case http.MethodDelete:
		eventName = ignoreRuleEventPrefix + "Deleted"
	default:
		eventName = ignoreRuleEventPrefix + "unknown"
	}
	now := time.Now().UTC()
	nowDate := now.Format("2006-01-02")
	nowMonth := now.Format("2006-01")
	_, nowWeekOfTheYear := now.ISOWeek()
	return IgnoreRuleEvent{
		EventBase: EventBase{
			CustomerGUID:       customerGUID,
			EventName:          eventName,
			EventTime:          now,
			EventDate:          nowDate,
			EventMonth:         nowMonth,
			EventWeekOfTheYear: nowWeekOfTheYear,
		},
		IgnoreRuleType: ruleType,
		IgnoredIds:     strings.Join(ids, ","),
		Resources:      resourcesCount,
		ExpirationType: expirationType,
	}
}

func newAlertChannelDetailedEvent(customerGUID, name string, channel notifications.AlertChannel, eventOp string) AlertChannelEvent {
	event := AlertChannelEvent{
		EventBase:   NewBaseEvent(customerGUID, AlertChannelPrefix+eventOp, nil),
		Name:        name,
		Type:        string(channel.ChannelType),
		AllClusters: ptr.To(len(channel.Scope) == 0),
	}
	alert := channel.GetAlertConfig(notifications.NotificationTypeComplianceDrift)
	if alert != nil && alert.IsEnabled() {
		drift := 0
		if alert.Parameters.DriftPercentage != nil {
			drift = *alert.Parameters.DriftPercentage
		}
		event.Compliance = fmt.Sprintf("%d%%", drift)
	}
	alert = channel.GetAlertConfig(notifications.NotificationTypeNewVulnerability)
	if alert != nil && alert.IsEnabled() {
		severityScore := 0
		if alert.Parameters.MinSeverity != nil {
			severityScore = *alert.Parameters.MinSeverity
		}
		event.NewVulnerability = containerscan.SeverityScoreToString(severityScore) + " and above"
	}
	alert = channel.GetAlertConfig(notifications.NotificationTypeVulnerabilityNewFix)
	if alert != nil && alert.IsEnabled() {
		severityScore := 0
		if alert.Parameters.MinSeverity != nil {
			severityScore = *alert.Parameters.MinSeverity
		}
		event.NewFix = containerscan.SeverityScoreToString(severityScore) + " and above"
	}
	alert = channel.GetAlertConfig(notifications.NotificationTypeNewClusterAdmin)
	if alert != nil && alert.IsEnabled() {
		event.NewAdmin = "true"
	}
	return event
}

func NewScanWithoutAccessKeyEvent(customerGUID, clusterName string) ScanWithoutAccessKeyEvent {
	return ScanWithoutAccessKeyEvent{
		EventBase:   NewBaseEvent(customerGUID, "ClusterScanWithoutAccessKey", nil),
		ClusterName: clusterName,
	}
}
