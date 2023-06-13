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
	nci2 := NotificationConfigIdentifier{NotificationType: NotificationTypePush}
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
}
