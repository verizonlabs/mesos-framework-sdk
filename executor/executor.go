package executor

/*
Executor interface and default executor implementation is defined here.
*/
import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	exec "mesos-framework-sdk/include/executor"
	"mesos-framework-sdk/include/mesos"
)

type Executor interface {
	FrameworkID() *mesos_v1.FrameworkID
	ExecutorID() *mesos_v1.ExecutorID
	Client() *client.Client
	Subscribe()
	Update(taskStatus *mesos_v1.TaskStatus)
	Message(data []byte)
}

type DefaultExecutor struct {
	frameworkID *mesos_v1.FrameworkID
	executorID  *mesos_v1.ExecutorID
	client      *client.Client
}

func NewDefaultExecutor(master string) *DefaultExecutor {
	return &DefaultExecutor{
		frameworkID: &mesos_v1.FrameworkID{Value: proto.String("")},
		executorID:  &mesos_v1.ExecutorID{Value: proto.String("")},
		client:      client.NewClient(master),
	}

}

func (d *DefaultExecutor) FrameworkID() {
	return d.frameworkID
}

func (d *DefaultExecutor) ExecutorID() {
	return d.executorID
}

func (d *DefaultExecutor) Subscribe() {
	// Both id's for framework and executor will be empty here.
	subscribe := &exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_SUBSCRIBE.Enum(),
	}
	d.client.Request(subscribe)
}

func (d *DefaultExecutor) Update(taskStatus *mesos_v1.TaskStatus) {
	update := exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_UPDATE.Enum(),
		Update: &exec.Call_Update{
			Status: taskStatus,
		},
	}
	d.client.Request(update)

}
func (d *DefaultExecutor) Message(data []byte) {
	message := exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_MESSAGE.Enum(),
		Message: &exec.Call_Message{
			Data: data,
		},
	}
	d.client.Request(message)

}
