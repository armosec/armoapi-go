package ecs

// EcsResourceType represents the type of ECS resource (shared between operator and backend)
type EcsResourceType string

func (r EcsResourceType) String() string {
	return string(r)
}

func (r EcsResourceType) Values() []EcsResourceType {
	return []EcsResourceType{
		ResourceTypeService,
		ResourceTypeTask,
		ResourceTypeTaskDefinition,
	}
}

const (
	ResourceTypeService        EcsResourceType = "Service"
	ResourceTypeTask           EcsResourceType = "Task"
	ResourceTypeTaskDefinition EcsResourceType = "TaskDefinition"
)

// TaskType distinguishes ECS tasks launched by a service from standalone tasks
// (launched via RunTask, EventBridge schedules, etc.). The K8s analogue is the
// distinction between a controller-owned Pod and an orphan Pod.
type TaskType string

const (
	TaskTypeService    TaskType = "service"
	TaskTypeStandalone TaskType = "standalone"
)
