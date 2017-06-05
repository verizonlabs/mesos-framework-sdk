package executor

import (
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos_v1"
	exec "mesos-framework-sdk/include/mesos_v1_executor"
	"mesos-framework-sdk/logging"
	"mesos-framework-sdk/recordio"
)

type Executor interface {
	Subscribe(chan *exec.Event) error
	Update(*mesos_v1.TaskStatus)
	Message([]byte)
}

type DefaultExecutor struct {
	frameworkId    *mesos_v1.FrameworkID
	executorId     *mesos_v1.ExecutorID
	client         client.Client
	logger         logging.Logger
	unackedTasks   map[*mesos_v1.TaskID]mesos_v1.TaskInfo
	unackedUpdates map[string]exec.Call_Update
}

// Creates a new default executor
func NewDefaultExecutor(
	f *mesos_v1.FrameworkID,
	e *mesos_v1.ExecutorID,
	c client.Client,
	lgr logging.Logger) Executor {

	return &DefaultExecutor{
		frameworkId:    f,
		executorId:     e,
		client:         c,
		logger:         lgr,
		unackedTasks:   make(map[*mesos_v1.TaskID]mesos_v1.TaskInfo),
		unackedUpdates: make(map[string]exec.Call_Update),
	}
}

func (e *DefaultExecutor) unacknowledgedTasks() []*mesos_v1.TaskInfo {
	numTasks := len(e.unackedTasks)
	if numTasks == 0 {
		return nil
	}

	tasks := make([]*mesos_v1.TaskInfo, 0, numTasks)
	for task := range e.unackedTasks {
		t := e.unackedTasks[task]
		tasks = append(tasks, &t)
	}

	return tasks
}

func (e *DefaultExecutor) unacknowledgedUpdates() []*exec.Call_Update {
	numUpdates := len(e.unackedUpdates)
	if numUpdates == 0 {
		return nil
	}

	updates := make([]*exec.Call_Update, 0, numUpdates)
	for update := range e.unackedUpdates {
		u := e.unackedUpdates[update]
		updates = append(updates, &u)
	}

	return updates
}

func (e *DefaultExecutor) Subscribe(eventChan chan *exec.Event) error {
	subscribe := &exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_SUBSCRIBE.Enum(),
		Subscribe: &exec.Call_Subscribe{
			UnacknowledgedTasks:   e.unacknowledgedTasks(),
			UnacknowledgedUpdates: e.unacknowledgedUpdates(),
		},
	}

	// If we disconnect we need to reset the stream ID. For this reason always start with a fresh stream ID.
	// Otherwise we'll never be able to reconnect.
	e.client.SetStreamID("")

	resp, err := e.client.Request(subscribe)
	if err != nil {
		return err
	} else {
		return recordio.Decode(resp.Body, eventChan)
	}
}

func (e *DefaultExecutor) Update(taskStatus *mesos_v1.TaskStatus) {
	update := exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_UPDATE.Enum(),
		Update: &exec.Call_Update{
			Status: taskStatus,
		},
	}
	e.client.Request(update)
}

func (e *DefaultExecutor) Message(data []byte) {
	message := exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_MESSAGE.Enum(),
		Message: &exec.Call_Message{
			Data: data,
		},
	}
	e.client.Request(message)
}
