package task_manager

import (
	"mesos-framework-sdk/include/mesos"
)

type TaskManager interface {
	Add(task *mesos_v1.Task)
	Delete(id *mesos_v1.TaskID)
	SetTaskState(task *mesos_v1.Task, state *mesos_v1.TaskState) error
	IsTaskInState(task *mesos_v1.Task, state *mesos_v1.TaskState) (bool, error)
	HasTask(task *mesos_v1.Task) bool
	HasQueuedTasks() bool
	TotalTasks() int
	Tasks() map[string]mesos_v1.Task
}

type DefaultTaskManager struct {
	totalTasks int
	tasks      map[string]mesos_v1.Task
}

func NewDefaultTaskManager() *DefaultTaskManager {
	return &DefaultTaskManager{tasks: make(map[string]mesos_v1.Task)}
}

// Provision a task
func (m *DefaultTaskManager) Add(task *mesos_v1.Task) {
	m.tasks[task.GetTaskId().GetValue()] = *task
}

// Delete a task
func (m *DefaultTaskManager) Delete(id *mesos_v1.TaskID) {
	delete(m.tasks, id.GetValue())
}

// Set a task status
func (m *DefaultTaskManager) SetTaskState(task *mesos_v1.Task, state *mesos_v1.TaskState) error {
	task.State = state
	m.tasks[task.GetTaskId().GetValue()] = *task
	return nil
}

// Check if a task has a particular status.
func (m *DefaultTaskManager) IsTaskInState(task *mesos_v1.Task, state *mesos_v1.TaskState) (bool, error) {
	if _, ok := m.tasks[task.GetTaskId().GetValue()]; !ok {
		return false, nil
	}
	return task.GetState().Enum() == state, nil
}

// Check if the task is already in the task manager.
func (m *DefaultTaskManager) HasTask(task *mesos_v1.Task) bool {
	if _, ok := m.tasks[task.GetTaskId().GetValue()]; ok {
		return true
	}
	return false
}

// Check if we have tasks left to execute.
func (m *DefaultTaskManager) HasQueuedTasks() bool {
	return !(len(m.tasks) == 0)
}

func (m *DefaultTaskManager) Tasks() map[string]mesos_v1.Task {
	return m.tasks
}

func (m *DefaultTaskManager) TotalTasks() int {
	return len(m.tasks)
}
