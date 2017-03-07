package task_manager

import (
	"errors"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/structures"
)

type TaskManager interface {
	Add(*mesos_v1.TaskInfo)
	Delete(*mesos_v1.TaskInfo)
	Get(*string) (*mesos_v1.TaskInfo, error)
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
		tasks:         structures.NewConcurrentMap(100),
		launchedTasks: make(map[string]*mesos_v1.TaskInfo),
		queuedTasks:   make(map[string]*mesos_v1.TaskInfo),
	}
}

// Provision a task
func (m *DefaultTaskManager) Add(task *mesos_v1.TaskInfo) {
	m.tasks.Set(task.GetName(), *task)
	m.SetTaskQueued(task)
}

// Delete a task
func (m *DefaultTaskManager) Delete(task *mesos_v1.TaskInfo) {
	m.tasks.Delete(task.GetName())
	delete(m.queuedTasks, task.GetName())
	delete(m.launchedTasks, task.GetName())
}

func (m *DefaultTaskManager) Get(name *string) (*mesos_v1.TaskInfo, error) {
	ret := m.tasks.Get(*name)
	if ret != nil {
		r := ret.(mesos_v1.TaskInfo)
		return &r, nil
	}
	return &mesos_v1.TaskInfo{}, errors.New("Could not find task.")
}

// Check if the task is already in the task manager.
func (m *DefaultTaskManager) HasTask(task *mesos_v1.TaskInfo) bool {
	ret := m.tasks.Get(task.GetName())
	if ret == nil {
		return false
	}
	return true
}

func (m *DefaultTaskManager) SetTaskQueued(task *mesos_v1.TaskInfo) {
	m.queuedTasks[task.GetName()] = task
}

func (m *DefaultTaskManager) SetTaskLaunched(task *mesos_v1.TaskInfo) error {
	if _, ok := m.queuedTasks[task.GetName()]; ok {
		delete(m.queuedTasks, task.GetName())  // Delete it from queue.
		m.launchedTasks[task.GetName()] = task // Add to launched tasks.
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
		if _, ok := m.queuedTasks[task.GetName()]; ok {
			return true
		}
	}
	return false
}

func (m *DefaultTaskManager) GetById(id *mesos_v1.TaskID) *mesos_v1.TaskInfo {
	// Check to see if any tasks we have match the id passed in.
	for v := range m.tasks.Iterate() {
		task := v.Value.(mesos_v1.TaskInfo)
		if task.GetTaskId().GetValue() == id.GetValue() {
			return &task
		}
	}
	return nil
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
