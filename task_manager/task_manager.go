package task_manager

import (
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/structures"
)

type TaskManager interface {
	Add(*mesos_v1.Task)
	Delete(*mesos_v1.Task)
	Get(*mesos_v1.TaskID) *mesos_v1.Task
	SetTaskState(*mesos_v1.Task, *mesos_v1.TaskState) error
	IsTaskInState(*mesos_v1.Task, *mesos_v1.TaskState) (bool, error)
	HasTask(*mesos_v1.Task) bool
	HasQueuedTasks() bool
	TotalTasks() int
	Tasks() *structures.ConcurrentMap
}

type DefaultTaskManager struct {
	totalTasks int
	tasks      *structures.ConcurrentMap
}

func NewDefaultTaskManager() *DefaultTaskManager {
	return &DefaultTaskManager{tasks: structures.NewConcurrentMap(1)}
}

// Provision a task
func (m *DefaultTaskManager) Add(task *mesos_v1.Task) {
	m.tasks.Set(task.GetTaskId().GetValue(), *task)
}

// Delete a task
func (m *DefaultTaskManager) Delete(task *mesos_v1.Task) {
	m.tasks.Delete(task.GetTaskId())
}

// Set a task status
func (m *DefaultTaskManager) SetTaskState(task *mesos_v1.Task, state *mesos_v1.TaskState) error {
	task.State = state
	m.tasks.Set(task.GetTaskId().GetValue(), *task)
	return nil
}

func (m *DefaultTaskManager) Get(id *mesos_v1.TaskID) *mesos_v1.Task {
	ret := m.tasks.Get(id.GetValue()).(mesos_v1.Task)
	return &ret
}

// Check if a task has a particular status.
func (m *DefaultTaskManager) IsTaskInState(task *mesos_v1.Task, state *mesos_v1.TaskState) (bool, error) {
	m.tasks.Get(task.GetTaskId().GetValue())
	return task.GetState().Enum() == state, nil
}

// Check if the task is already in the task manager.
func (m *DefaultTaskManager) HasTask(task *mesos_v1.Task) bool {
	ret := m.tasks.Get(task.GetTaskId().GetValue())
	if ret == nil {
		return false
	}
	return true
}

// Check if we have tasks left to execute.
func (m *DefaultTaskManager) HasQueuedTasks() bool {
	// Check to see if we have any tasks in STAGING state still.
	for v := range m.tasks.Iterate() {
		task := v.Value.(mesos_v1.Task)
		if task.GetState() == mesos_v1.TaskState_TASK_STAGING {
			return true
		}
	}
	return false
}

func (m *DefaultTaskManager) Tasks() *structures.ConcurrentMap {
	return m.tasks
}

func (m *DefaultTaskManager) TotalTasks() int {
	return m.tasks.Length()
}
