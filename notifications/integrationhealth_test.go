package notifications

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHealthDefaultsToHealthy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		config *CollaborationConfig
	}{
		{name: "nil config", config: nil},
		{name: "nil attributes", config: &CollaborationConfig{}},
		{name: "unrelated attributes", config: configWithAttributes(map[string]interface{}{"workspaceId": "w-1"})},
		{name: "non-string health attribute", config: configWithAttributes(map[string]interface{}{AttributeHealthStatus: 42})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			health := tt.config.GetHealth()
			assert.Equal(t, IntegrationHealthStatusHealthy, health.Status)
			assert.Empty(t, health.LastChecked)
			assert.Empty(t, health.LastError)
			assert.Empty(t, health.DegradedSince)
		})
	}
}

func TestSetHealthDegraded(t *testing.T) {
	t.Parallel()
	config := &CollaborationConfig{}
	config.SetHealthDegraded("token_refresh_failed: invalid_grant")

	health := config.GetHealth()
	assert.Equal(t, IntegrationHealthStatusDegraded, health.Status)
	assert.Equal(t, "token_refresh_failed: invalid_grant", health.LastError)
	assertRFC3339(t, health.LastChecked)
	assertRFC3339(t, health.DegradedSince)
}

func TestSetHealthDegradedKeepsDegradedSince(t *testing.T) {
	t.Parallel()
	config := &CollaborationConfig{}
	config.SetHealthDegraded("first failure")
	firstSince := config.GetHealth().DegradedSince

	config.SetHealthDegraded("second failure")
	health := config.GetHealth()
	assert.Equal(t, firstSince, health.DegradedSince, "re-marking must not reset the transition time")
	assert.Equal(t, "second failure", health.LastError)
}

func TestSetHealthDegradedPreservesOtherAttributes(t *testing.T) {
	t.Parallel()
	config := configWithAttributes(map[string]interface{}{"workspaceId": "w-1"})
	config.SetHealthDegraded("boom")
	assert.Equal(t, "w-1", config.Attributes["workspaceId"])
}

func TestSetHealthChecked(t *testing.T) {
	t.Parallel()
	config := &CollaborationConfig{}
	config.SetHealthDegraded("boom")

	config.SetHealthChecked()
	health := config.GetHealth()
	assert.Equal(t, IntegrationHealthStatusHealthy, health.Status)
	assert.Empty(t, health.LastError)
	assert.Empty(t, health.DegradedSince)
	assertRFC3339(t, health.LastChecked)
}

func TestClearHealth(t *testing.T) {
	t.Parallel()
	config := configWithAttributes(map[string]interface{}{"workspaceId": "w-1"})
	config.SetHealthDegraded("boom")

	config.ClearHealth()
	for _, key := range []string{AttributeHealthStatus, AttributeHealthLastChecked, AttributeHealthLastError, AttributeHealthDegradedSince} {
		assert.NotContains(t, config.Attributes, key)
	}
	assert.Equal(t, "w-1", config.Attributes["workspaceId"], "unrelated attributes must survive")
	assert.NotPanics(t, func() { (&CollaborationConfig{}).ClearHealth() })
}

// Simulates the config-service round-trip: the config is marshalled to JSON
// and back (Attributes become map[string]interface{} with plain strings) and
// the helpers must still read the health correctly.
func TestHealthSurvivesJSONRoundTrip(t *testing.T) {
	t.Parallel()
	config := configWithAttributes(map[string]interface{}{"workspaceId": "w-1"})
	config.SetHealthDegraded("token_refresh_failed")

	raw, err := json.Marshal(config)
	require.NoError(t, err)
	var roundTripped CollaborationConfig
	require.NoError(t, json.Unmarshal(raw, &roundTripped))

	health := roundTripped.GetHealth()
	assert.Equal(t, IntegrationHealthStatusDegraded, health.Status)
	assert.Equal(t, "token_refresh_failed", health.LastError)
	assertRFC3339(t, health.DegradedSince)
	assert.Equal(t, "w-1", roundTripped.Attributes["workspaceId"])
}

func configWithAttributes(attributes map[string]interface{}) *CollaborationConfig {
	config := &CollaborationConfig{}
	config.Attributes = attributes
	return config
}

func assertRFC3339(t *testing.T, value string) {
	t.Helper()
	_, err := time.Parse(time.RFC3339, value)
	assert.NoError(t, err, "expected RFC3339 timestamp, got %q", value)
}
