package task_manager

import (
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/structures"
)

type TaskManager interface {
	Add(*mesos_v1.TaskInfo) error
	Delete(*mesos_v1.TaskInfo)
	Get(*string) (*mesos_v1.TaskInfo, error)
	HasTask(*mesos_v1.TaskInfo) bool
	TotalTasks() int
	Tasks() *structures.ConcurrentMap
}
