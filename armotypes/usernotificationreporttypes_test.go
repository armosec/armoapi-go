package armotypes

import (
	"encoding/json"
	"testing"
	"time"

	_ "embed"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/weeklyReport.json
var weeklyReport string

func TestWeeklyReport(t *testing.T) {
	from, _ := time.Parse(time.RFC3339,  "2023-01-07T00:00:00+00:00")
	to, _ := time.Parse(time.RFC3339, "2023-01-14T00:00:00+00:00")
	report :=WeeklyReport{
		From:                                from,
		To:                                  to,
		AccountName:                         "userAccount",
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
	b, _ := json.Marshal(report)
	assert.Equal(t, weeklyReport, string(b))
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
	notificationsConfig:=NotificationsConfig{
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
		clusterName          string
		scanType ScanType
		want                *PushReport
	}
	testTable := []testCase{
		{
			name: "empty",
		},
		{
			name: "not found",
			notificationsConfig: notificationsConfig,
			clusterName: "test",
		},
		{
			name: "get repository scan",
			notificationsConfig: notificationsConfig,
			clusterName: "cluster1",
			scanType: ScanTypePosture,
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
			assert.Equal(t, test.want,test.notificationsConfig.GetLatestPushReport(test.clusterName, test.scanType), test.name)
		})
	}
}
