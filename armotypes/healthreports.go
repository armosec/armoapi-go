package armotypes

import "time"

// HealthReport represents a minimal health report emitted by the agent.
type HealthReport struct {
	// AgentVersion is the version of the agent emitting this report.
	AgentVersion string `json:"agent_version"`
	// SensorUpdated is the time the sensor was installed/updated.
	SensorUpdated time.Time `json:"sensor_updated"`
	// Timestamp is the time the report was generated.
	Timestamp time.Time `json:"timestamp"`
	// CloudMetadata contains enriched cloud provider metadata associated with this node.
	CloudMetadata CloudMetadata `json:"cloudMetadata"`
}
