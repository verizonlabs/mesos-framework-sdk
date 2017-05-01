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
	frameworkId *mesos_v1.FrameworkID
	executorId  *mesos_v1.ExecutorID
	client      client.Client
	logger      logging.Logger
}

// Creates a new default executor
func NewDefaultExecutor(
	f *mesos_v1.FrameworkID,
	e *mesos_v1.ExecutorID,
	c client.Client,
	lgr logging.Logger) Executor {

	return &DefaultExecutor{
		frameworkId: f,
		executorId:  e,
		client:      c,
		logger:      lgr,
	}
}

func (e *DefaultExecutor) Subscribe(eventChan chan *exec.Event) error {
	subscribe := &exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_SUBSCRIBE.Enum(),
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
