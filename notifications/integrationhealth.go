package notifications

import "time"

// Integration health is persisted as entries in CollaborationConfig.Attributes
// (map[string]interface{}) rather than as typed struct fields. Attributes
// entries survive read-modify-write round-trips through services built against
// older versions of this library, while unknown typed fields would be silently
// dropped on re-serialization. Absence of the health keys means healthy.
const (
	AttributeHealthStatus        = "healthStatus"
	AttributeHealthLastChecked   = "healthLastChecked"
	AttributeHealthLastError     = "healthLastError"
	AttributeHealthDegradedSince = "healthDegradedSince"
)

type IntegrationHealthStatus string

const (
	IntegrationHealthStatusHealthy  IntegrationHealthStatus = "healthy"
	IntegrationHealthStatusDegraded IntegrationHealthStatus = "degraded"
)

// IntegrationHealth is the typed view over the health attributes of a
// CollaborationConfig. Timestamps are RFC3339 strings, consistent with
// PortalBase.UpdatedTime and CollaborationConfig.CreationTime.
// swagger:model IntegrationHealth
type IntegrationHealth struct {
	Status IntegrationHealthStatus `json:"status"`
	// LastChecked is the time of the last health evaluation (successful or not)
	LastChecked string `json:"lastChecked,omitempty"`
	// LastError is the reason for the last failed health evaluation
	LastError string `json:"lastError,omitempty"`
	// DegradedSince is when the integration first transitioned into degraded
	DegradedSince string `json:"degradedSince,omitempty"`
}

// GetHealth returns the health recorded on the config's attributes.
// Missing/nil attributes mean the integration is considered healthy.
func (c *CollaborationConfig) GetHealth() IntegrationHealth {
	health := IntegrationHealth{Status: IntegrationHealthStatusHealthy}
	if c == nil || c.Attributes == nil {
		return health
	}
	if attributeString(c.Attributes, AttributeHealthStatus) == string(IntegrationHealthStatusDegraded) {
		health.Status = IntegrationHealthStatusDegraded
	}
	health.LastChecked = attributeString(c.Attributes, AttributeHealthLastChecked)
	health.LastError = attributeString(c.Attributes, AttributeHealthLastError)
	health.DegradedSince = attributeString(c.Attributes, AttributeHealthDegradedSince)
	return health
}

// SetHealthDegraded marks the integration as degraded with the given reason.
// Idempotent with respect to DegradedSince: re-marking an already-degraded
// config refreshes LastChecked/LastError but keeps the original transition time.
func (c *CollaborationConfig) SetHealthDegraded(reason string) {
	if c.Attributes == nil {
		c.Attributes = map[string]interface{}{}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if attributeString(c.Attributes, AttributeHealthStatus) != string(IntegrationHealthStatusDegraded) {
		c.Attributes[AttributeHealthDegradedSince] = now
	}
	c.Attributes[AttributeHealthStatus] = string(IntegrationHealthStatusDegraded)
	c.Attributes[AttributeHealthLastChecked] = now
	c.Attributes[AttributeHealthLastError] = reason
}

// SetHealthChecked records a successful health evaluation, clearing any
// degraded state.
func (c *CollaborationConfig) SetHealthChecked() {
	if c.Attributes == nil {
		c.Attributes = map[string]interface{}{}
	}
	c.ClearHealth()
	c.Attributes[AttributeHealthLastChecked] = time.Now().UTC().Format(time.RFC3339)
}

// ClearHealth removes all health attributes (e.g. after a successful
// reconnect), returning the config to the implicit healthy state.
func (c *CollaborationConfig) ClearHealth() {
	if c == nil || c.Attributes == nil {
		return
	}
	delete(c.Attributes, AttributeHealthStatus)
	delete(c.Attributes, AttributeHealthLastChecked)
	delete(c.Attributes, AttributeHealthLastError)
	delete(c.Attributes, AttributeHealthDegradedSince)
}

// attributeString reads a string attribute, tolerating the
// map[string]interface{} shapes produced by JSON/BSON round-trips.
func attributeString(attributes map[string]interface{}, key string) string {
	if v, ok := attributes[key].(string); ok {
		return v
	}
	return ""
}
