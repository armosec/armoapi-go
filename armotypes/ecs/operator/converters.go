package operator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/armosec/armoapi-go/armotypes/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// ============================================================================
// Marshal AWS SDK types to EcsResource (for ecs-operator)
// ============================================================================

// NewClusterResource creates a EcsResource from an AWS ECS Cluster
func NewClusterResource(cluster types.Cluster) (EcsClusterResource, error) {
	payload, err := json.Marshal(cluster)
	if err != nil {
		return EcsClusterResource{}, fmt.Errorf("failed to marshal cluster: %w", err)
	}

	arn := ""
	if cluster.ClusterArn != nil {
		arn = *cluster.ClusterArn
	}

	return EcsClusterResource{
		ResourceARN:  arn,
		ResourceType: "", // Clusters don't have a resource type
		Payload:      payload,
		Checksum:     ComputeChecksum(payload),
	}, nil
}

// NewServiceResource creates a ResourceData from an AWS ECS Service
func NewServiceResource(service types.Service) (EcsResource, error) {
	payload, err := json.Marshal(service)
	if err != nil {
		return EcsResource{}, fmt.Errorf("failed to marshal service: %w", err)
	}

	arn := ""
	if service.ServiceArn != nil {
		arn = *service.ServiceArn
	}

	return EcsResource{
		ResourceARN:  arn,
		ResourceType: ecs.ResourceTypeService,
		Payload:      payload,
		Checksum:     ComputeChecksum(payload),
	}, nil
}

// NewTaskResource creates a ResourceData from an AWS ECS Task
func NewTaskResource(task types.Task) (EcsResource, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return EcsResource{}, fmt.Errorf("failed to marshal task: %w", err)
	}

	arn := ""
	if task.TaskArn != nil {
		arn = *task.TaskArn
	}

	return EcsResource{
		ResourceARN:  arn,
		ResourceType: ecs.ResourceTypeTask,
		Payload:      payload,
		Checksum:     ComputeChecksum(payload),
	}, nil
}

// NewTaskDefinitionResource creates a ResourceData from an AWS ECS TaskDefinition
func NewTaskDefinitionResource(taskDef types.TaskDefinition) (EcsResource, error) {
	payload, err := json.Marshal(taskDef)
	if err != nil {
		return EcsResource{}, fmt.Errorf("failed to marshal task definition: %w", err)
	}

	arn := ""
	if taskDef.TaskDefinitionArn != nil {
		arn = *taskDef.TaskDefinitionArn
	}

	return EcsResource{
		ResourceARN:  arn,
		ResourceType: ecs.ResourceTypeTaskDefinition,
		Payload:      payload,
		Checksum:     ComputeChecksum(payload),
	}, nil
}

// ============================================================================
// Unmarshal ResourceData payload to AWS SDK types (for backend)
// ============================================================================

// UnmarshalCluster extracts the AWS Cluster from a ResourceData payload
func (r *EcsClusterResource) Unmarshal() (*types.Cluster, error) {
	var cluster types.Cluster
	if err := json.Unmarshal(r.Payload, &cluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cluster payload: %w", err)
	}
	return &cluster, nil
}

// UnmarshalService extracts the AWS Service from a ResourceData payload
func (r *EcsResource) UnmarshalService() (*types.Service, error) {
	if r.ResourceType != ecs.ResourceTypeService {
		return nil, fmt.Errorf("resource type mismatch: expected %s, got %s", ecs.ResourceTypeService, r.ResourceType)
	}

	var service types.Service
	if err := json.Unmarshal(r.Payload, &service); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service payload: %w", err)
	}
	return &service, nil
}

// UnmarshalTask extracts the AWS Task from a ResourceData payload
func (r *EcsResource) UnmarshalTask() (*types.Task, error) {
	if r.ResourceType != ecs.ResourceTypeTask {
		return nil, fmt.Errorf("resource type mismatch: expected %s, got %s", ecs.ResourceTypeTask, r.ResourceType)
	}

	var task types.Task
	if err := json.Unmarshal(r.Payload, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task payload: %w", err)
	}
	return &task, nil
}

// UnmarshalTaskDefinition extracts the AWS TaskDefinition from a ResourceData payload
func (r *EcsResource) UnmarshalTaskDefinition() (*types.TaskDefinition, error) {
	if r.ResourceType != ecs.ResourceTypeTaskDefinition {
		return nil, fmt.Errorf("resource type mismatch: expected %s, got %s", ecs.ResourceTypeTaskDefinition, r.ResourceType)
	}

	var taskDef types.TaskDefinition
	if err := json.Unmarshal(r.Payload, &taskDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task definition payload: %w", err)
	}
	return &taskDef, nil
}

// ============================================================================
// Utility functions
// ============================================================================

// ComputeChecksum computes a SHA-256 checksum of the payload for change detection
func ComputeChecksum(payload []byte) string {
	hash := sha256.Sum256(payload)
	return hex.EncodeToString(hash[:])
}

// ClassifyTaskType returns ecs.TaskTypeService if the task was launched by an
// ECS service, ecs.TaskTypeStandalone otherwise. AWS sets
// Task.Group="service:<name>" and Task.StartedBy="ecs-svc/<deployment-id>"
// for service-managed tasks; either match is sufficient.
func ClassifyTaskType(task *types.Task) ecs.TaskType {
	if task != nil {
		if task.Group != nil && strings.HasPrefix(*task.Group, "service:") {
			return ecs.TaskTypeService
		}
		if task.StartedBy != nil && strings.HasPrefix(*task.StartedBy, "ecs-svc/") {
			return ecs.TaskTypeService
		}
	}
	return ecs.TaskTypeStandalone
}
