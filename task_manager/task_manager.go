package task_manager

import (
	"fmt"
	"mesos-framework-sdk/include/mesos"
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
	fmt.Println(task.GetTaskId().GetValue())
	m.tasks[task.GetTaskId().GetValue()] = *task
}

// Delete a task
func (m *DefaultTaskManager) Delete(task *mesos_v1.Task) {
	delete(m.tasks, task.GetTaskId().GetValue())
}

// Set a task status
func (m *DefaultTaskManager) SetTaskState(task *mesos_v1.Task, state *mesos_v1.TaskState) error {
	task.State = state
	m.tasks[task.GetTaskId().GetValue()] = *task
	return nil
}

func (m *DefaultTaskManager) Get(id *mesos_v1.TaskID) *mesos_v1.Task {
	ret := m.tasks[id.GetValue()]
	return &ret
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
	// Check to see if we have any tasks in STAGING state still.
	for _, v := range m.tasks {
		fmt.Println(v.State.Enum())
		fmt.Println(v.GetState() == mesos_v1.TaskState_TASK_STAGING)
		if v.GetState() == mesos_v1.TaskState_TASK_STAGING {
			return true
		}
	}
	return false
}

func (m *DefaultTaskManager) Tasks() map[string]mesos_v1.Task {
	return m.tasks
}

func (m *DefaultTaskManager) TotalTasks() int {
	return len(m.tasks)
}
