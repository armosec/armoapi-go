package broadcastevents

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/notifications"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseEvent(t *testing.T) {
	event := NewBaseEvent("testGUID", "testEvent", nil)
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "testEvent", event.EventName)
}

//go:embed fixtures/alertChannel.json
var alertChannelJSON []byte

func TestNewAlertChannelDeletedEvent(t *testing.T) {
	event := NewAlertChannelDeletedEvent("testGUID", "testName", "testProvider")
	assert.Equal(t, "testName", event.Name)
	assert.Equal(t, "testProvider", event.Type)
}

func TestNewAlertChannelCreatedEvent(t *testing.T) {
	var channel notifications.AlertChannel
	err := json.Unmarshal(alertChannelJSON, &channel)
	if err != nil {
		t.Error(err)
	}
	event := NewAlertChannelCreatedEvent("testGUID", "testName", channel)
	assert.Equal(t, "AlertChannelCreated", event.EventName)
	assert.Equal(t, "testName", event.Name)
	assert.Equal(t, "teams", event.Type)
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "Low and above", event.NewVulnerability)
	assert.Equal(t, "High and above", event.NewFix)
	assert.Equal(t, "true", event.NewAdmin)
}

func TestNewAlertChannelUpdatedEvent(t *testing.T) {
	var channel notifications.AlertChannel
	err := json.Unmarshal(alertChannelJSON, &channel)
	if err != nil {
		t.Error(err)
	}
	event := NewAlertChannelUpdatedEvent("testGUID", "testName", channel)
	assert.Equal(t, "testName", event.Name)
	assert.Equal(t, "teams", event.Type)
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "AlertChannelUpdated", event.EventName)
	assert.Equal(t, "Low and above", event.NewVulnerability)
	assert.Equal(t, "High and above", event.NewFix)
	assert.Equal(t, "true", event.NewAdmin)
	assert.Equal(t, "25%", event.Compliance)

}

func TestNewPostureExceptionEvent(t *testing.T) {
	exception := armotypes.PostureExceptionPolicy{}
	event := NewPostureExceptionEvent("testGUID", "POST", exception)
	assert.Equal(t, "testGUID", event.CustomerGUID)
}

func TestNewVulnerabilityExceptionChangeEvent(t *testing.T) {
	exception := armotypes.VulnerabilityExceptionPolicy{}
	event := NewVulnerabilityExceptionChangeEvent("testGUID", "POST", exception)
	assert.Equal(t, "testGUID", event.CustomerGUID)
}

func TestNewFeatureFlagsEvent(t *testing.T) {
	event := NewFeatureFlagsEvent("testGUID", "testEmail", "testName", "testPreferredName")
	assert.Equal(t, "testEmail", event.Email)
	assert.Equal(t, "testName", event.UserName)
	assert.Equal(t, "testPreferredName", event.PreferredUserName)
}

func TestNewLoginEvent(t *testing.T) {
	event := NewLoginEvent("testGUID", "testEmail", "testName", "testPreferredName")
	assert.Equal(t, "testEmail", event.Email)
	assert.Equal(t, "testName", event.UserName)
	assert.Equal(t, "testPreferredName", event.PreferredUserName)
}

func TestNewClusterImageScanSessionStartedEvent(t *testing.T) {
	eventTime := time.Now().Add(-time.Hour)
	event := NewClusterImageScanSessionStartedEvent("testJobId", "testClusterName", "testCustomerId", eventTime)
	assert.Equal(t, "testJobId", event.JobID)
	assert.Equal(t, "testClusterName", event.ClusterName)
	assert.Equal(t, eventTime, event.EventTime)
}

func TestNewClusterRiskScanV2Event(t *testing.T) {
	eventTime := time.Now().Add(-time.Hour)
	event := NewClusterRiskScanV2Event("customerGUID", "testJobId", "testReportGUID", "testClusterName", "testKubescapeVersion", "testCloudProvider", "testK8sVer", "testHelBVersion", 10, eventTime)
	assert.Equal(t, "testJobId", event.JobID)
	assert.Equal(t, "testClusterName", event.ClusterName)
	assert.Equal(t, eventTime, event.EventTime)
	assert.Equal(t, "testReportGUID", event.ReportGUID)
	assert.Equal(t, "testKubescapeVersion", event.KSVersion)
	assert.Equal(t, "testCloudProvider", event.K8sVendor)
	assert.Equal(t, "testK8sVer", event.K8sVersion)
	assert.Equal(t, "testHelBVersion", event.HelmChartVersion)
	assert.Equal(t, 10, event.WorkerNodesCount)

}

func TestNewIgnoreRuleEvent(t *testing.T) {
	event := newIgnoreRuleEvent("testGUID", http.MethodPost, IgnoreRuleTypeMisconfiguration, IgnoreRuleExpirationTypeNone, []string{"id1", "id2"}, 2)
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "IgnoreRuleCreated", event.EventName)
	assert.Equal(t, "id1,id2", event.IgnoredIds)
	assert.Equal(t, 2, event.Resources)
}
