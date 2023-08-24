package broadcastevents

import (
	"net/http"
	"testing"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/notifications"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseEvent(t *testing.T) {
	event := NewBaseEvent("testGUID", "testEvent")
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "testEvent", event.EventName)
}

func TestNewAlertChannelDeletedEvent(t *testing.T) {
	event := NewAlertChannelDeletedEvent("testGUID", "testName", "testProvider")
	assert.Equal(t, "testName", event.Name)
	assert.Equal(t, "testProvider", event.Type)
}

func TestNewAlertChannelCreatedEvent(t *testing.T) {
	channel := notifications.AlertChannel{
		ChannelType: notifications.CollaborationTypeTeams,
	}
	event := NewAlertChannelCreatedEvent("testGUID", "testName", channel)
	assert.Equal(t, "testName", event.Name)
	assert.Equal(t, "teams", event.Type)
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
	eventTime := time.Now()
	event := NewClusterImageScanSessionStartedEvent("testJobId", "testClusterName", "testCustomerId", eventTime)
	assert.Equal(t, "testJobId", event.JobID)
	assert.Equal(t, "testClusterName", event.ClusterName)
}

// ... Add more tests for other methods ...

func TestNewIgnoreRuleEvent(t *testing.T) {
	event := newIgnoreRuleEvent("testGUID", http.MethodPost, IgnoreRuleTypeMisconfiguration, IgnoreRuleExpirationTypeNone, []string{"id1", "id2"}, 2)
	assert.Equal(t, "testGUID", event.CustomerGUID)
	assert.Equal(t, "IgnoreRuleCreated", event.EventName)
	assert.Equal(t, "id1,id2", event.IgnoredIds)
	assert.Equal(t, 2, event.Resources)
}
