package task_manager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/structures"
)

type TaskManager interface {
	Add(*mesos_v1.TaskInfo)
	Delete(*mesos_v1.TaskInfo)
	Get(*mesos_v1.TaskID) *mesos_v1.TaskInfo
	HasTask(*mesos_v1.TaskInfo) bool
	TotalTasks() int
	Tasks() *structures.ConcurrentMap
}

type DefaultTaskManager struct {
	totalTasks    int
	tasks         *structures.ConcurrentMap
	launchedTasks map[string]*mesos_v1.TaskInfo
	queuedTasks   map[string]*mesos_v1.TaskInfo
}

func NewDefaultTaskManager() *DefaultTaskManager {
	return &DefaultTaskManager{
		tasks:         structures.NewConcurrentMap(1),
		totalTasks:    0,
		launchedTasks: make(map[string]*mesos_v1.TaskInfo),
		queuedTasks:   make(map[string]*mesos_v1.TaskInfo),
	}
}

// Provision a task
func (m *DefaultTaskManager) Add(task *mesos_v1.TaskInfo) {
	m.tasks.Set(task.GetTaskId().GetValue(), *task)
	m.SetTaskQueued(task)
}

// Delete a task
func (m *DefaultTaskManager) Delete(task *mesos_v1.TaskInfo) {
	m.tasks.Delete(task.GetTaskId())
}

func (m *DefaultTaskManager) Get(id *mesos_v1.TaskID) *mesos_v1.TaskInfo {
	ret := m.tasks.Get(id.GetValue()).(mesos_v1.TaskInfo)
	return &ret
}

// Check if the task is already in the task manager.
func (m *DefaultTaskManager) HasTask(task *mesos_v1.TaskInfo) bool {
	ret := m.tasks.Get(task.GetTaskId().GetValue())
	if ret == nil {
		return false
	}
	return true
}

func (m *DefaultTaskManager) SetTaskQueued(task *mesos_v1.TaskInfo) {
	m.queuedTasks[task.GetTaskId().GetValue()] = task
}

func (m *DefaultTaskManager) SetTaskLaunched(task *mesos_v1.TaskInfo) error {
	if _, ok := m.queuedTasks[task.GetTaskId().GetValue()]; ok {
		delete(m.queuedTasks, task.GetTaskId().GetValue())  // Delete it from queue.
		m.launchedTasks[task.GetTaskId().GetValue()] = task // Add to launched tasks.
		return nil
	}
	// This task isn't queued up, reject.
	return errors.New("Task is not in queue, cannot set to launching.")
}

// Check if we have tasks left to execute.
func (m *DefaultTaskManager) HasQueuedTasks() bool {
	// Check to see if we have any tasks in the queued tasks map.
	for v := range m.tasks.Iterate() {
		task := v.Value.(mesos_v1.TaskInfo)
		if _, ok := m.queuedTasks[task.GetTaskId().GetValue()]; ok {
			return true
		}
	}
	return false
}

func (m *DefaultTaskManager) QueuedTasks() map[string]*mesos_v1.TaskInfo {
	return m.queuedTasks
}

func (m *DefaultTaskManager) LaunchedTasks() map[string]*mesos_v1.TaskInfo {
	return m.launchedTasks
}

func (m *DefaultTaskManager) Tasks() *structures.ConcurrentMap {
	return m.tasks
}

func (m *DefaultTaskManager) TotalTasks() int {
	return m.tasks.Length()
}
