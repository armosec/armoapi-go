package operator

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/armosec/armoapi-go/armotypes/ecs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClusterResource_Success(t *testing.T) {
	cluster := types.Cluster{
		ClusterArn:                        aws.String("arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster"),
		ClusterName:                       aws.String("my-cluster"),
		Status:                            aws.String("ACTIVE"),
		RegisteredContainerInstancesCount: 5,
		RunningTasksCount:                 10,
		ActiveServicesCount:               3,
	}

	resource, err := NewClusterResource(cluster)

	require.NoError(t, err)
	assert.Equal(t, "arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster", resource.ResourceARN)
	assert.NotEmpty(t, resource.Payload)
	assert.NotEmpty(t, resource.Checksum)

	// Verify round-trip
	unmarshaled, err := resource.Unmarshal()
	require.NoError(t, err)
	assert.Equal(t, "my-cluster", *unmarshaled.ClusterName)
	assert.Equal(t, "ACTIVE", *unmarshaled.Status)
	assert.Equal(t, int32(10), unmarshaled.RunningTasksCount)
}

func TestNewClusterResource_NilArn(t *testing.T) {
	cluster := types.Cluster{
		ClusterName: aws.String("my-cluster"),
		Status:      aws.String("ACTIVE"),
	}

	resource, err := NewClusterResource(cluster)

	require.NoError(t, err)
	assert.Equal(t, "", resource.ResourceARN)
}

// ============================================================================
// NewServiceResource tests
// ============================================================================

func TestNewServiceResource_Success(t *testing.T) {
	service := types.Service{
		ServiceArn:   aws.String("arn:aws:ecs:us-west-2:123456789012:service/my-cluster/my-service"),
		ServiceName:  aws.String("my-service"),
		ClusterArn:   aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/my-cluster"),
		Status:       aws.String("ACTIVE"),
		DesiredCount: 3,
		RunningCount: 3,
		LaunchType:   types.LaunchTypeFargate,
	}

	resource, err := NewServiceResource(service)

	require.NoError(t, err)
	assert.Equal(t, "arn:aws:ecs:us-west-2:123456789012:service/my-cluster/my-service", resource.ResourceARN)
	assert.Equal(t, ecs.ResourceTypeService, resource.ResourceType)
	assert.NotEmpty(t, resource.Payload)
	assert.NotEmpty(t, resource.Checksum)

	// Verify round-trip
	unmarshaled, err := resource.UnmarshalService()
	require.NoError(t, err)
	assert.Equal(t, "my-service", *unmarshaled.ServiceName)
	assert.Equal(t, types.LaunchTypeFargate, unmarshaled.LaunchType)
}

func TestNewServiceResource_NilArn(t *testing.T) {
	service := types.Service{
		ServiceName: aws.String("my-service"),
	}

	resource, err := NewServiceResource(service)

	require.NoError(t, err)
	assert.Equal(t, "", resource.ResourceARN)
}

// ============================================================================
// NewTaskResource tests
// ============================================================================

func TestNewTaskResource_Success(t *testing.T) {
	task := types.Task{
		TaskArn:           aws.String("arn:aws:ecs:eu-west-1:123456789012:task/my-cluster/abc123"),
		TaskDefinitionArn: aws.String("arn:aws:ecs:eu-west-1:123456789012:task-definition/my-task:1"),
		ClusterArn:        aws.String("arn:aws:ecs:eu-west-1:123456789012:cluster/my-cluster"),
		LastStatus:        aws.String("RUNNING"),
		DesiredStatus:     aws.String("RUNNING"),
		LaunchType:        types.LaunchTypeEc2,
	}

	resource, err := NewTaskResource(task)

	require.NoError(t, err)
	assert.Equal(t, "arn:aws:ecs:eu-west-1:123456789012:task/my-cluster/abc123", resource.ResourceARN)
	assert.Equal(t, ecs.ResourceTypeTask, resource.ResourceType)
	assert.NotEmpty(t, resource.Payload)
	assert.NotEmpty(t, resource.Checksum)

	// Verify round-trip
	unmarshaled, err := resource.UnmarshalTask()
	require.NoError(t, err)
	assert.Equal(t, "RUNNING", *unmarshaled.LastStatus)
	assert.Equal(t, types.LaunchTypeEc2, unmarshaled.LaunchType)
}

func TestNewTaskResource_NilArn(t *testing.T) {
	task := types.Task{
		LastStatus: aws.String("RUNNING"),
	}

	resource, err := NewTaskResource(task)

	require.NoError(t, err)
	assert.Equal(t, "", resource.ResourceARN)
}

// ============================================================================
// NewTaskDefinitionResource tests
// ============================================================================

func TestNewTaskDefinitionResource_Success(t *testing.T) {
	taskDef := types.TaskDefinition{
		TaskDefinitionArn: aws.String("arn:aws:ecs:ap-southeast-1:123456789012:task-definition/my-task:5"),
		Family:            aws.String("my-task"),
		Revision:          5,
		Status:            types.TaskDefinitionStatusActive,
		Cpu:               aws.String("256"),
		Memory:            aws.String("512"),
		Compatibilities:   []types.Compatibility{types.CompatibilityFargate, types.CompatibilityEc2},
	}

	resource, err := NewTaskDefinitionResource(taskDef)

	require.NoError(t, err)
	assert.Equal(t, "arn:aws:ecs:ap-southeast-1:123456789012:task-definition/my-task:5", resource.ResourceARN)
	assert.Equal(t, ecs.ResourceTypeTaskDefinition, resource.ResourceType)
	assert.NotEmpty(t, resource.Payload)
	assert.NotEmpty(t, resource.Checksum)

	// Verify round-trip
	unmarshaled, err := resource.UnmarshalTaskDefinition()
	require.NoError(t, err)
	assert.Equal(t, "my-task", *unmarshaled.Family)
	assert.Equal(t, int32(5), unmarshaled.Revision)
	assert.Equal(t, "256", *unmarshaled.Cpu)
}

func TestNewTaskDefinitionResource_NilArn(t *testing.T) {
	taskDef := types.TaskDefinition{
		Family: aws.String("my-task"),
	}

	resource, err := NewTaskDefinitionResource(taskDef)

	require.NoError(t, err)
	assert.Equal(t, "", resource.ResourceARN)
}

// ============================================================================
// UnmarshalCluster tests
// ============================================================================

func TestUnmarshalCluster_InvalidPayload(t *testing.T) {
	resource := EcsClusterResource{
		Payload: json.RawMessage(`{invalid json`),
	}

	_, err := resource.Unmarshal()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// ============================================================================
// UnmarshalService tests
// ============================================================================

func TestUnmarshalService_TypeMismatch(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeTask,
		Payload:      json.RawMessage(`{}`),
	}

	_, err := resource.UnmarshalService()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource type mismatch")
}

func TestUnmarshalService_InvalidPayload(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeService,
		Payload:      json.RawMessage(`not valid json`),
	}

	_, err := resource.UnmarshalService()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// ============================================================================
// UnmarshalTask tests
// ============================================================================

func TestUnmarshalTask_TypeMismatch(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeService,
		Payload:      json.RawMessage(`{}`),
	}

	_, err := resource.UnmarshalTask()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource type mismatch")
}

func TestUnmarshalTask_InvalidPayload(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeTask,
		Payload:      json.RawMessage(`[broken`),
	}

	_, err := resource.UnmarshalTask()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// ============================================================================
// UnmarshalTaskDefinition tests
// ============================================================================

func TestUnmarshalTaskDefinition_TypeMismatch(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeService,
		Payload:      json.RawMessage(`{}`),
	}

	_, err := resource.UnmarshalTaskDefinition()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource type mismatch")
}

func TestUnmarshalTaskDefinition_InvalidPayload(t *testing.T) {
	resource := EcsResource{
		ResourceType: ecs.ResourceTypeTaskDefinition,
		Payload:      json.RawMessage(`}`),
	}

	_, err := resource.UnmarshalTaskDefinition()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// ============================================================================
// ComputeChecksum tests
// ============================================================================

func TestComputeChecksum_DeterministicOutput(t *testing.T) {
	payload := []byte(`{"key":"value"}`)

	checksum1 := ComputeChecksum(payload)
	checksum2 := ComputeChecksum(payload)

	assert.Equal(t, checksum1, checksum2)
	assert.Len(t, checksum1, 64) // SHA-256 produces 64 hex characters
}

func TestComputeChecksum_DifferentPayloads(t *testing.T) {
	payload1 := []byte(`{"key":"value1"}`)
	payload2 := []byte(`{"key":"value2"}`)

	checksum1 := ComputeChecksum(payload1)
	checksum2 := ComputeChecksum(payload2)

	assert.NotEqual(t, checksum1, checksum2)
}

func TestComputeChecksum_EmptyPayload(t *testing.T) {
	checksum := ComputeChecksum([]byte{})

	assert.NotEmpty(t, checksum)
	assert.Len(t, checksum, 64)
}

// ============================================================================
// Round-trip integration tests
// ============================================================================

func TestRoundTrip_Cluster(t *testing.T) {
	original := types.Cluster{
		ClusterArn:                        aws.String("arn:aws:ecs:us-east-1:123456789012:cluster/production"),
		ClusterName:                       aws.String("production"),
		Status:                            aws.String("ACTIVE"),
		RegisteredContainerInstancesCount: 10,
		RunningTasksCount:                 50,
		PendingTasksCount:                 2,
		ActiveServicesCount:               15,
	}

	resource, err := NewClusterResource(original)
	require.NoError(t, err)

	restored, err := resource.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, *original.ClusterArn, *restored.ClusterArn)
	assert.Equal(t, *original.ClusterName, *restored.ClusterName)
	assert.Equal(t, *original.Status, *restored.Status)
	assert.Equal(t, original.RegisteredContainerInstancesCount, restored.RegisteredContainerInstancesCount)
	assert.Equal(t, original.RunningTasksCount, restored.RunningTasksCount)
}

func TestRoundTrip_Service(t *testing.T) {
	original := types.Service{
		ServiceArn:   aws.String("arn:aws:ecs:us-west-2:123456789012:service/prod/api-service"),
		ServiceName:  aws.String("api-service"),
		ClusterArn:   aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/prod"),
		Status:       aws.String("ACTIVE"),
		DesiredCount: 5,
		RunningCount: 5,
		LaunchType:   types.LaunchTypeFargate,
	}

	resource, err := NewServiceResource(original)
	require.NoError(t, err)

	restored, err := resource.UnmarshalService()
	require.NoError(t, err)

	assert.Equal(t, *original.ServiceArn, *restored.ServiceArn)
	assert.Equal(t, *original.ServiceName, *restored.ServiceName)
	assert.Equal(t, *original.ClusterArn, *restored.ClusterArn)
	assert.Equal(t, original.DesiredCount, restored.DesiredCount)
	assert.Equal(t, original.LaunchType, restored.LaunchType)
}

func TestRoundTrip_Task(t *testing.T) {
	original := types.Task{
		TaskArn:           aws.String("arn:aws:ecs:eu-central-1:123456789012:task/cluster/task123"),
		TaskDefinitionArn: aws.String("arn:aws:ecs:eu-central-1:123456789012:task-definition/web:10"),
		ClusterArn:        aws.String("arn:aws:ecs:eu-central-1:123456789012:cluster/cluster"),
		LastStatus:        aws.String("RUNNING"),
		DesiredStatus:     aws.String("RUNNING"),
		LaunchType:        types.LaunchTypeEc2,
		Cpu:               aws.String("1024"),
		Memory:            aws.String("2048"),
	}

	resource, err := NewTaskResource(original)
	require.NoError(t, err)

	restored, err := resource.UnmarshalTask()
	require.NoError(t, err)

	assert.Equal(t, *original.TaskArn, *restored.TaskArn)
	assert.Equal(t, *original.TaskDefinitionArn, *restored.TaskDefinitionArn)
	assert.Equal(t, *original.LastStatus, *restored.LastStatus)
	assert.Equal(t, original.LaunchType, restored.LaunchType)
	assert.Equal(t, *original.Cpu, *restored.Cpu)
}

func TestRoundTrip_TaskDefinition(t *testing.T) {
	original := types.TaskDefinition{
		TaskDefinitionArn: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/worker:25"),
		Family:            aws.String("worker"),
		Revision:          25,
		Status:            types.TaskDefinitionStatusActive,
		Cpu:               aws.String("512"),
		Memory:            aws.String("1024"),
		NetworkMode:       types.NetworkModeAwsvpc,
		Compatibilities:   []types.Compatibility{types.CompatibilityFargate},
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Name:  aws.String("app"),
				Image: aws.String("nginx:latest"),
			},
		},
	}

	resource, err := NewTaskDefinitionResource(original)
	require.NoError(t, err)

	restored, err := resource.UnmarshalTaskDefinition()
	require.NoError(t, err)

	assert.Equal(t, *original.TaskDefinitionArn, *restored.TaskDefinitionArn)
	assert.Equal(t, *original.Family, *restored.Family)
	assert.Equal(t, original.Revision, restored.Revision)
	assert.Equal(t, original.Status, restored.Status)
	assert.Equal(t, original.NetworkMode, restored.NetworkMode)
	assert.Len(t, restored.ContainerDefinitions, 1)
	assert.Equal(t, "app", *restored.ContainerDefinitions[0].Name)
}

func TestClassifyTaskType(t *testing.T) {
	tests := []struct {
		name     string
		task     *types.Task
		expected ecs.TaskType
	}{
		{
			name:     "nil task is standalone",
			task:     nil,
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "empty task is standalone",
			task:     &types.Task{},
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "Group=service:* is service",
			task:     &types.Task{Group: aws.String("service:my-svc")},
			expected: ecs.TaskTypeService,
		},
		{
			name:     "StartedBy=ecs-svc/* is service",
			task:     &types.Task{StartedBy: aws.String("ecs-svc/2637865999356093636")},
			expected: ecs.TaskTypeService,
		},
		{
			name: "both Group and StartedBy set to service markers",
			task: &types.Task{
				Group:     aws.String("service:my-svc"),
				StartedBy: aws.String("ecs-svc/abc"),
			},
			expected: ecs.TaskTypeService,
		},
		{
			name:     "Group=family:* is standalone",
			task:     &types.Task{Group: aws.String("family:my-family")},
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "StartedBy=user is standalone (RunTask)",
			task:     &types.Task{StartedBy: aws.String("user")},
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "StartedBy=events-rule/* is standalone (EventBridge scheduled)",
			task:     &types.Task{StartedBy: aws.String("events-rule/my-rule")},
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "non-service Group prefix wins over absent StartedBy",
			task:     &types.Task{Group: aws.String("family:foo")},
			expected: ecs.TaskTypeStandalone,
		},
		{
			name:     "service: substring not at start is standalone",
			task:     &types.Task{Group: aws.String("not-service:foo")},
			expected: ecs.TaskTypeStandalone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ClassifyTaskType(tt.task))
		})
	}
}
