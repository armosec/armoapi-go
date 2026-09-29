package operator

import (
	"encoding/json"
	"time"

	"github.com/armosec/armoapi-go/armotypes/ecs"
)

// OperatorEventType represents the type of operator event
type OperatorEventType string

const (
	EventTypeKeepalive OperatorEventType = "keepalive"
	EventTypeResources OperatorEventType = "resources"
	EventTypeReconcile OperatorEventType = "reconcile"
)

// Status represents the status of an ECS operator
type Status string

const (
	StatusPending      Status = "PENDING"
	StatusHealthy      Status = "HEALTHY"
	StatusDegraded     Status = "DEGRADED"
	StatusUnhealthy    Status = "UNHEALTHY"
	StatusDisconnected Status = "DISCONNECTED"
)

// ChunkingInfo represents pagination/chunking information for splitting large requests
type ChunkingInfo struct {
	ChunkID     string `json:"chunkId"`     // UUID identifying a chunked request (same for all chunks)
	ChunkIndex  int    `json:"chunkIndex"`  // 0-based index of this chunk
	TotalChunks int    `json:"totalChunks"` // Total number of chunks
}

// OperatorEvent represents a generic ECS operator event
type OperatorEvent struct {
	MessageID    string            `json:"messageId"`          // UUID for debugging and tracing
	EventType    OperatorEventType `json:"eventType"`          // keepalive, resources, reconcile
	ClusterName  string            `json:"clusterName"`        // Required for all events
	Region       string            `json:"region"`             // AWS region (e.g., us-east-1)
	AWSAccountID string            `json:"awsAccountId"`       // AWS account ID
	Payload      json.RawMessage   `json:"payload"`            // Event-specific payload
	Chunking     *ChunkingInfo     `json:"chunking,omitempty"` // Pagination/chunking support (for future use)
}

// OperatorKeepAliveEvent represents the heartbeat payload along with cluster resource if it's the first event or if cluster has changed
type OperatorKeepAliveEvent struct {
	Timestamp                time.Time           `json:"timestamp"`
	Version                  string              `json:"version"`
	Status                   Status              `json:"status"` // HEALTHY, DEGRADED, UNHEALTHY
	Uptime                   string              `json:"uptime"` // Duration as string (e.g., "1h30m")
	LastSyncTimestamp        *time.Time          `json:"lastSyncTimestamp,omitempty"`
	LastReconcileTimestamp   *time.Time          `json:"lastReconcileTimestamp,omitempty"`
	InventoryCounts          InventoryCounts     `json:"inventoryCounts"`
	Errors                   ErrorCounts         `json:"errors"`
	Cluster                  *EcsClusterResource `json:"cluster,omitempty"`        // (optional) cluster resource will be set only if it's the first event or if cluster has changed
	OperatorVersionUpdatedAt time.Time           `json:"operatorVersionUpdatedAt"` // Time when the operator version was deployed on ECS
}

// InventoryCounts represents counts of discovered resources
type InventoryCounts struct {
	Services        int `json:"services"`
	Tasks           int `json:"tasks"`
	TaskDefinitions int `json:"taskDefinitions"`
}

// ErrorCounts represents error counts in the last hour
type ErrorCounts struct {
	SyncErrors int `json:"syncErrors"`
	APIErrors  int `json:"apiErrors"`
}

// OperatorResourcesEvent represents an event with resource information
type OperatorResourcesEvent struct {
	Resources []EcsResource `json:"resources"`
}

// EcsResource represents information about an ECS resource (for operator events)
type EcsResource struct {
	ResourceARN  string              `json:"resourceArn"`
	ResourceType ecs.EcsResourceType `json:"resourceType"` // Service, Task, TaskDefinition
	Payload      json.RawMessage     `json:"payload"`      // Raw AWS SDK JSON
	Checksum     string              `json:"checksum"`
}

// EcsClusterResource represents a cluster resource (for operator events)
// It's a type alias to EcsResource for JSON compatibility
type EcsClusterResource = EcsResource

// OperatorReconcileEvent represents a reconciliation event with all ARNs
type OperatorReconcileEvent struct {
	Services        []string `json:"services"`
	Tasks           []string `json:"tasks"`
	TaskDefinitions []string `json:"taskDefinitions"`
}
