package armotypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/weeklyReport.json
var weeklyReport string

func TestWeeklyReport(t *testing.T) {
	report := WeeklyReport{
		ClustersScannedThisWeek:             1,
		ClustersScannedPrevWeek:             2,
		LinkToConfigurationScanningFiltered: "http://somelink1.com",
		RepositoriesScannedThisWeek:         3,
		RepositoriesScannedPrevWeek:         4,
		LinkToRepositoriesScanningFiltered:  "http://somelink2.com",
		RegistriesScannedThisWeek:           5,
		RegistriesScannedPrevWeek:           6,
		LinkToRegistriesScanningFiltered:    "http://somelink3.com",
		Top5FailedControls:                  []TopCtrlItem{{Name: "control1", TotalFailedResources: 1}},
	}
	got, err := json.Marshal(report)
	assert.NoError(t, err)
	assert.Equal(t, weeklyReport, string(got))

	report.Top5FailedControls = []TopCtrlItem{{Name: "control1", Clusters: []TopCtrlCluster{{Name: "cluster1", ResourcesCount: 10}}},
		{Name: "control2", Clusters: []TopCtrlCluster{{Name: "cluster1", ResourcesCount: 10}, {Name: "cluster1", ResourcesCount: 100}}}}
	assert.Equal(t, int64(10), report.Top5FailedControls[0].GetTotalFailedResources())
	assert.Equal(t, int64(110), report.Top5FailedControls[1].GetTotalFailedResources())
}

func TestAddLatestPushReport(t *testing.T) {
	ts := time.Now()
	type testCase struct {
		name                string
		notificationsConfig NotificationsConfig
		report              PushReport
		want                NotificationsConfig
	}
	testTable := []testCase{
		{
			name: "empty",
			want: NotificationsConfig{
				LatestPushReports: map[string]*PushReport{"_": {}},
			},
		},
		{
			name:                "first report",
			notificationsConfig: NotificationsConfig{},
			report: PushReport{
				Cluster:         "cluster1",
				ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
				ScanType:        ScanTypePosture,
				FailedResources: 3,
			},
			want: NotificationsConfig{
				LatestPushReports: map[string]*PushReport{"cluster1_posture": {
					Cluster:         "cluster1",
					ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
					ScanType:        ScanTypePosture,
					FailedResources: 3,
				}},
			},
		},
		{
			name: "add repository scan",
			notificationsConfig: NotificationsConfig{
				LatestPushReports: map[string]*PushReport{"cluster1_posture": {
					Cluster:         "cluster1",
					ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
					ScanType:        ScanTypePosture,
					FailedResources: 3,
				}},
			},
			report: PushReport{
				ScanType:        ScanTypeRepositories,
				Timestamp:       ts,
				FailedResources: 4,
			},
			want: NotificationsConfig{
				LatestPushReports: map[string]*PushReport{
					"cluster1_posture": {
						Cluster:         "cluster1",
						ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
						ScanType:        ScanTypePosture,
						FailedResources: 3,
					},
					"_repository": {
						ScanType:        ScanTypeRepositories,
						Timestamp:       ts,
						FailedResources: 4,
					}},
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.notificationsConfig.AddLatestPushReport(&test.report)
			assert.Equal(t, test.want, test.notificationsConfig, test.name)
		})
	}
}

func TestGetLatestPushReport(t *testing.T) {
	ts := time.Now()
	notificationsConfig := NotificationsConfig{
		LatestPushReports: map[string]*PushReport{
			"cluster1_posture": {
				Cluster:         "cluster1",
				ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
				ScanType:        ScanTypePosture,
				FailedResources: 3,
			},
			"_repository": {
				ScanType:        ScanTypeRepositories,
				Timestamp:       ts,
				FailedResources: 4,
			}},
	}

	type testCase struct {
		name                string
		notificationsConfig NotificationsConfig
		clusterName         string
		scanType            ScanType
		want                *PushReport
	}
	testTable := []testCase{
		{
			name: "empty",
		},
		{
			name:                "not found",
			notificationsConfig: notificationsConfig,
			clusterName:         "test",
		},
		{
			name:                "get repository scan",
			notificationsConfig: notificationsConfig,
			clusterName:         "cluster1",
			scanType:            ScanTypePosture,
			want: &PushReport{
				Cluster:         "cluster1",
				ReportGUID:      "0a801812-2777-4886-a64e-9a731a41c1c4",
				ScanType:        ScanTypePosture,
				FailedResources: 3,
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.notificationsConfig.GetLatestPushReport(test.clusterName, test.scanType), test.name)
		})
	}
}

func TestNotificationConfigIdentifier_Validate(t *testing.T) {
	// Test case 1: Valid NotificationType (NotificationTypeAll)
	nci1 := NotificationConfigIdentifier{NotificationType: NotificationTypeAll}
	err1 := nci1.Validate()
	if err1 != nil {
		t.Errorf("Test case 1 failed: Expected Validate to return nil error, but got %s", err1.Error())
	}

	// Test case 2: Valid NotificationType (NotificationTypePush)
	nci2 := NotificationConfigIdentifier{NotificationType: NotificationTypePushPosture}
	err2 := nci2.Validate()
	if err2 != nil {
		t.Errorf("Test case 2 failed: Expected Validate to return nil error, but got %s", err2.Error())
	}

	// Test case 3: Valid NotificationType (NotificationTypeWeekly)
	nci3 := NotificationConfigIdentifier{NotificationType: NotificationTypeWeekly}
	err3 := nci3.Validate()
	if err3 != nil {
		t.Errorf("Test case 3 failed: Expected Validate to return nil error, but got %s", err3.Error())
	}

	// Test case 4: Invalid NotificationType
	nci4 := NotificationConfigIdentifier{NotificationType: "invalidType"}
	err4 := nci4.Validate()
	expectedError := fmt.Errorf("invalid notification type: %s", nci4.NotificationType)
	if err4 == nil {
		t.Errorf("Test case 4 failed: Expected Validate to return non-nil error, but got nil")
	} else if err4.Error() != expectedError.Error() {
		t.Errorf("Test case 4 failed: Expected error %s, but got %s", expectedError.Error(), err4.Error())
	}

	// Test case 5: empty NotificationType
	nci5 := NotificationConfigIdentifier{NotificationType: ""}
	err5 := nci5.Validate()
	expectedError = fmt.Errorf("notification type is required")
	if err5 == nil {
		t.Errorf("Test case 5 failed: Expected Validate to return non-nil error, but got nil")
	} else if err5.Error() != expectedError.Error() {
		t.Errorf("Test case 5 failed: Expected error %s, but got %s", expectedError.Error(), err4.Error())
	}
}

func TestGetAlertConfig(t *testing.T) {
	alertChannel := AlertChannel{
		Alerts: []AlertConfig{
			{
				NotificationConfigIdentifier: NotificationConfigIdentifier{
					NotificationType: NotificationTypePushPosture,
				},
			},
			{
				NotificationConfigIdentifier: NotificationConfigIdentifier{
					NotificationType: NotificationTypeWeekly,
				},
			},
		},
	}

	// Test case where the alert config should be found
	config := alertChannel.GetAlertConfig(NotificationTypePushPosture)
	if config == nil {
		t.Errorf("Expected alert config, got nil")
	} else if config.NotificationType != NotificationTypePushPosture {
		t.Errorf("Expected NotificationType to be %s, got %s", NotificationTypePushPosture, config.NotificationType)
	}

	// Test case where the alert config should not be found
	config = alertChannel.GetAlertConfig(NotificationTypeAll)
	if config != nil {
		t.Errorf("Expected nil, got alert config")
	}
}

func TestGetAlertConfigurations(t *testing.T) {
	alertChannel1 := AlertChannel{
		Alerts: []AlertConfig{
			{
				NotificationConfigIdentifier: NotificationConfigIdentifier{
					NotificationType: NotificationTypePushPosture,
				},
			},
		},
	}

	alertChannel2 := AlertChannel{
		Alerts: []AlertConfig{
			{
				NotificationConfigIdentifier: NotificationConfigIdentifier{
					NotificationType: NotificationTypeWeekly,
				},
			},
		},
	}

	notificationsConfig := NotificationsConfig{
		AlertChannels: map[ChannelProvider][]AlertChannel{
			CollaborationTypeJira:  {alertChannel1},
			CollaborationTypeSlack: {alertChannel2},
		},
	}

	// Test case where the alert configs should be found
	alertConfigs := notificationsConfig.GetAlertConfigurations(NotificationTypePushPosture)
	if len(alertConfigs) != 1 {
		t.Errorf("Expected 1 alert config, got %d", len(alertConfigs))
	} else if alertConfigs[0].NotificationType != NotificationTypePushPosture {
		t.Errorf("Expected NotificationType to be %s, got %s", NotificationTypePushPosture, alertConfigs[0].NotificationType)
	}

	alertConfigs = notificationsConfig.GetAlertConfigurations(NotificationTypeWeekly)
	if len(alertConfigs) != 1 {
		t.Errorf("Expected 1 alert config, got %d", len(alertConfigs))
	} else if alertConfigs[0].NotificationType != NotificationTypeWeekly {
		t.Errorf("Expected NotificationType to be %s, got %s", NotificationTypeWeekly, alertConfigs[0].NotificationType)
	}

	// Test case where the alert configs should not be found
	alertConfigs = notificationsConfig.GetAlertConfigurations(NotificationTypeAll)
	if len(alertConfigs) != 0 {
		t.Errorf("Expected 0 alert configs, got %d", len(alertConfigs))
	}
}

func TestAlertParams(t *testing.T) {
	ap := make(NotificationParams)

	// Test SetDriftPercentage and GetDriftPercentage
	ap.SetDriftPercentage(10)
	percentage, err := ap.GetDriftPercentage()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if percentage != 10 {
		t.Errorf("Expected 10, got %d", percentage)
	}

	// Test SetMinSeverity and GetMinSeverity
	ap.SetMinSeverity(5)
	severity, err := ap.GetMinSeverity()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if severity != 5 {
		t.Errorf("Expected 5, got %d", severity)
	}
}

func TestNotificationsConfigChannels(t *testing.T) {
	nc := NotificationsConfig{
		AlertChannels: make(map[ChannelProvider][]AlertChannel),
	}
	ac := AlertChannel{
		Alerts: []AlertConfig{
			{
				Scope: []AlertScope{
					{
						Cluster:    "testCluster",
						Namespaces: []string{"testNamespace"},
					},
				},
			},
		},
	}
	nc.AlertChannels["testProvider"] = []AlertChannel{ac}

	// Test GetProviderChannels
	channels := nc.GetProviderChannels("testProvider")
	if len(channels) != 1 {
		t.Errorf("Expected 1, got %d", len(channels))
	}

	// Test IsInScope for NotificationsConfig, AlertChannel, AlertConfig, and AlertScope
	if !nc.IsInScope("testCluster", "testNamespace") {
		t.Errorf("Expected true, got false")
	}
	if !channels[0].IsInScope("testCluster", "testNamespace") {
		t.Errorf("Expected true, got false")
	}
	if !channels[0].Alerts[0].IsInScope("testCluster", "testNamespace") {
		t.Errorf("Expected true, got false")
	}
	if !channels[0].Alerts[0].Scope[0].IsInScope("testCluster", "testNamespace") {
		t.Errorf("Expected true, got false")
	}

	// Test AddAlertConfig and GetAlertConfig
	err := channels[0].AddAlertConfig(AlertConfig{
		NotificationConfigIdentifier: NotificationConfigIdentifier{
			NotificationType: "testType",
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	config := channels[0].GetAlertConfig("testType")
	if config == nil {
		t.Errorf("Expected non-nil, got nil")
	}
	if config.NotificationType != "testType" {
		t.Errorf("Expected 'testType', got '%s'", config.NotificationType)
	}
}

func TestNotificationsConfigChannelsNegative(t *testing.T) {
	nc := NotificationsConfig{
		AlertChannels: make(map[ChannelProvider][]AlertChannel),
	}
	ac := AlertChannel{
		Alerts: []AlertConfig{
			{
				NotificationConfigIdentifier: NotificationConfigIdentifier{
					NotificationType: "testType",
				},
				Scope: []AlertScope{
					{
						Cluster:    "testCluster",
						Namespaces: []string{"testNamespace"},
					},
				},
			},
		},
	}
	nc.AlertChannels["testProvider"] = []AlertChannel{ac}

	// Test GetProviderChannels with non-existing provider
	channels := nc.GetProviderChannels("nonExistingProvider")
	if len(channels) != 0 {
		t.Errorf("Expected 0, got %d", len(channels))
	}

	// Test IsInScope for NotificationsConfig, AlertChannel, AlertConfig, and AlertScope with non-existing cluster and namespace
	if nc.IsInScope("nonExistingCluster", "nonExistingNamespace") {
		t.Errorf("Expected false, got true")
	}
	if nc.AlertChannels["testProvider"][0].IsInScope("nonExistingCluster", "nonExistingNamespace") {
		t.Errorf("Expected false, got true")
	}
	if nc.AlertChannels["testProvider"][0].Alerts[0].IsInScope("nonExistingCluster", "nonExistingNamespace") {
		t.Errorf("Expected false, got true")
	}
	if nc.AlertChannels["testProvider"][0].Alerts[0].Scope[0].IsInScope("nonExistingCluster", "nonExistingNamespace") {
		t.Errorf("Expected false, got true")
	}

	// Test AddAlertConfig with existing notification type
	err := nc.AlertChannels["testProvider"][0].AddAlertConfig(AlertConfig{
		NotificationConfigIdentifier: NotificationConfigIdentifier{
			NotificationType: "testType",
		},
	})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test GetAlertConfig with non-existing notification type
	config := nc.AlertChannels["testProvider"][0].GetAlertConfig("nonExistingType")
	if config != nil {
		t.Errorf("Expected nil, got non-nil")
	}
}
