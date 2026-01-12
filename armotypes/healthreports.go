package armotypes

import "time"

// HealthReport represents a minimal health report emitted by the agent.
type HealthReport struct {
	// Timestamp is the time the report was generated.
	Timestamp time.Time `json:"timestamp"`
	// AgentVersion is the version of the agent emitting this report.
	AgentVersion string `json:"agent_version"`
	// CloudMetadata contains enriched cloud provider metadata associated with this node.
	CloudMetadata CloudMetadata `json:"cloudMetadata"`
}
