package scheduler

/*
Scheduler struct defines an interface of the default calls for the scheduler, as well as holding information
regarding the framework, client, and events handler.

End users may wish to make their own scheduler and satisfy this interface.

A default scheduler is provided for those who wish to keep the default implementations: All calls simply create
the protobuf required for the call and send it off to the client.

End users should only create their own scheduler if they wish to change the behavior of their calls.
*/
import (
	"fmt"

	"io/ioutil"
	"log"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/recordio"
	"time"
)

const (
	subscribeRetry = 2
	executors      = 3
)

type Scheduler interface {
	// Scheduler must also hold framework information, a client and an event handler.
	FrameworkInfo() *mesos_v1.FrameworkInfo
	Client() *client.Client

	// Default Calls for scheduler
	Subscribe(chan *sched.Event) error
	Teardown()
	Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters)
	Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters)
	Revive()
	Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID)
	Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID)
	Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte)
	Reconcile(tasks []*mesos_v1.Task)
	Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte)
	SchedRequest(resources []*mesos_v1.Request)
	Suppress()
}

// Default Scheduler can be used as a higher-level construct.
type DefaultScheduler struct {
	Info     *mesos_v1.FrameworkInfo
	client   *client.Client
	executor int
}

func NewDefaultScheduler(c *client.Client, info *mesos_v1.FrameworkInfo) *DefaultScheduler {
	return &DefaultScheduler{
		client:   c,
		Info:     info,
		executor: executors,
	}
}

func (c *DefaultScheduler) FrameworkInfo() *mesos_v1.FrameworkInfo {
	return c.Info
}

// Make a subscription call to mesos.
// Channel passed is the "listener" channel for Event Controller.
func (c *DefaultScheduler) Subscribe(eventChan chan *sched.Event) error {
	call := &sched.Call{
		Type: sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: c.Info,
		},
	}
	go func() {
		for {
			resp, err := c.client.Request(call)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(recordio.Decode(resp.Body, eventChan))
			}

			// If we disconnect we need to reset the stream ID.
			// Otherwise we'll never be able to reconnect.
			c.client.StreamID = ""
			time.Sleep(time.Duration(subscribeRetry) * time.Second)
		}
	}()
	return nil
}

// Send a teardown request to mesos master.
func (c *DefaultScheduler) Teardown() {
	teardown := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_TEARDOWN.Enum(),
	}
	resp, err := c.client.Request(teardown)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(resp)
}

// Accepts offers from mesos master
func (c *DefaultScheduler) Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters) {
	accept := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      &sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.client.Request(accept)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// NOTE: If we get back an improper response, this will panic.
	// TODO: recover from panics.
	k, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(k))

}

func (c *DefaultScheduler) Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     &sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.client.Request(decline)
	if err != nil {
		log.Println(err.Error())
	}
	a, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(a))
	return
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *DefaultScheduler) Revive() {
	revive := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.client.Request(revive)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *DefaultScheduler) Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID) {
	// Probably want some validation that this is a valid task and valid agentid.
	kill := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_KILL.Enum(),
		Kill:        &sched.Call_Kill{TaskId: taskId, AgentId: agentid},
	}

	resp, err := c.client.Request(kill)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *DefaultScheduler) Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID) {
	shutdown := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_SHUTDOWN.Enum(),
		Shutdown: &sched.Call_Shutdown{
			ExecutorId: execId,
			AgentId:    agentId,
		},
	}
	resp, err := c.client.Request(shutdown)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// UUID should be a type
// TODO import extras uuid funcs.
func (c *DefaultScheduler) Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte) {
	fmt.Println("Acknowledge event.")
	acknowledge := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &sched.Call_Acknowledge{
			AgentId: agentId,
			TaskId:  taskId,
			Uuid:    uuid,
		},
	}
	resp, err := c.client.Request(acknowledge)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Reconcile(tasks []*mesos_v1.Task) {
	reconcile := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_RECONCILE.Enum(),
	}
	resp, err := c.client.Request(reconcile)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte) {
	message := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_MESSAGE.Enum(),
		Message: &sched.Call_Message{
			AgentId:    agentId,
			ExecutorId: executorId,
			Data:       data,
		},
	}
	resp, err := c.client.Request(message)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)

}

// Sent by the scheduler to request resources from the master/allocator.
// The built-in hierarchical allocator simply ignores this request but other allocators (modules) can interpret this in
// a customizable fashion.
func (c *DefaultScheduler) SchedRequest(resources []*mesos_v1.Request) {
	request := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_REQUEST.Enum(),
		Request: &sched.Call_Request{
			Requests: resources,
		},
	}
	resp, err := c.client.Request(request)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Suppress() {
	suppress := &sched.Call{
		FrameworkId: c.FrameworkInfo().GetId(),
		Type:        sched.Call_SUPPRESS.Enum(),
	}
	resp, err := c.client.Request(suppress)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}
