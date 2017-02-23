package executor

/*
Executor interface and default executor implementation is defined here.
*/
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/executor/events"
	exec "mesos-framework-sdk/include/executor"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/recordio"
	"time"
)

const (
	subscribeRetry = 2
)

type Executor interface {
	FrameworkID() *mesos_v1.FrameworkID
	ExecutorID() *mesos_v1.ExecutorID
	Client() *client.Client
	Events() chan *exec.Event
	Subscribe()
	Update(taskStatus *mesos_v1.TaskStatus)
	Message(data []byte)
	Run()
}

type DefaultExecutor struct {
	frameworkID *mesos_v1.FrameworkID
	executorID  *mesos_v1.ExecutorID
	client      *client.Client
	events      chan *exec.Event
	handlers    events.ExecutorEvents
}

// Creates a new default executor
func NewDefaultExecutor(c *client.Client) *DefaultExecutor {
	return &DefaultExecutor{
		frameworkID: &mesos_v1.FrameworkID{Value: proto.String("")},
		executorID:  &mesos_v1.ExecutorID{Value: proto.String("")},
		client:      c,
		events:      make(chan *exec.Event),
	}

}

func (d *DefaultExecutor) Run() {
	if d.frameworkID.GetValue() == "" {
		d.Subscribe()
	}
	d.listen()
}

// Default listening method on the
func (d *DefaultExecutor) listen() {
	for {
		switch t := <-d.events; t.GetType() {
		case exec.Event_SUBSCRIBED:
			d.frameworkID = t.GetSubscribed().GetFrameworkInfo().GetId()
			go d.handlers.Subscribed(t.GetSubscribed())
			break
		case exec.Event_ACKNOWLEDGED:
			go d.handlers.Acknowledged(t.GetAcknowledged())
			break
		case exec.Event_MESSAGE:
			go d.handlers.Message(t.GetMessage())
			break
		case exec.Event_KILL:
			go d.handlers.Kill(t.GetKill())
			break
		case exec.Event_LAUNCH:
			go d.handlers.Launch(t.GetLaunch())
			break
		case exec.Event_LAUNCH_GROUP:
			go d.handlers.LaunchGroup(t.GetLaunchGroup())
			break
		case exec.Event_SHUTDOWN:
			go d.handlers.Shutdown()
			break
		case exec.Event_ERROR:
			go d.handlers.Error(t.GetError())
			break
		case exec.Event_UNKNOWN:
			fmt.Println("Unknown event caught.")
			break
		}
	}
}

func (d *DefaultExecutor) FrameworkID() *mesos_v1.FrameworkID {
	return d.frameworkID
}

func (d *DefaultExecutor) ExecutorID() *mesos_v1.ExecutorID {
	return d.executorID
}

func (d *DefaultExecutor) Subscribe() {
	// Both id's for framework and executor will be empty here.
	subscribe := &exec.Call{
		FrameworkId: d.frameworkID,
		ExecutorId:  d.executorID,
		Type:        exec.Call_SUBSCRIBE.Enum(),
	}

	go func() {
		for {
			resp, err := d.client.Request(subscribe)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(recordio.Decode(resp.Body, d.events))
			}

			// If we disconnect we need to reset the stream ID.
			// Otherwise we'll never be able to reconnect.
			d.client.StreamID = ""
			time.Sleep(time.Duration(subscribeRetry) * time.Second)
		}
	}()
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
