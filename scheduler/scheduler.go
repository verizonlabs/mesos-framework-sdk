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
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"mesos-framework-sdk/client"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/recordio"
	"mesos-framework-sdk/scheduler/events"
	"mesos-framework-sdk/task_manager"
	"strconv"
	"time"
)

const (
	subscribeRetry = 2
)

type Scheduler interface {
	// Scheduler must also hold framework information, a client and an event handler.
	FrameworkInfo()
	Client()
	Events()
	// Default Calls for scheduler
	Subscribe() error
	Teardown()
	Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters)
	Decline(offerIds []*mesos.OfferID, filters *mesos.Filters)
	Revive()
	Kill(taskId *mesos.TaskID, agentid *mesos.AgentID)
	Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID)
	Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte)
	Reconcile(tasks []*mesos.Task)
	Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte)
	SchedRequest(resources []*mesos.Request)
	Suppress()
}

// Default Scheduler can be used as a higher-level construct.
type DefaultScheduler struct {
	FramworkInfo *mesos.FrameworkInfo
	client       *client.Client
	events       chan *sched.Event
	calls        chan *sched.Call
	handlers     events.SchedulerEvent
	manager      task_manager.TaskManager
}

func NewDefaultScheduler(c *client.Client, info *mesos.FrameworkInfo,
	event chan *sched.Event, callChan chan *sched.Call,
	manager task_manager.TaskManager, handlers events.SchedulerEvent) *DefaultScheduler {

	return &DefaultScheduler{
		client:       c,
		events:       event,
		calls:        callChan,
		FramworkInfo: info,
		handlers:     handlers,
		manager:      manager,
	}
}

// Events
func (c *DefaultScheduler) Run() {
	// If we don't have a framework id, subscribe.
	if c.FramworkInfo.GetId().GetValue() == "" {
		fmt.Println("subscribing...")
		err := c.Subscribe()
		if err != nil {
			log.Println(err.Error())
		}
		c.launchExecutors(1)
		c.listen()
		return
	}
	// Otherwise we're already connected. Just listen for events.
	c.listen()
}

// Create n default executors and launch them.
func (c *DefaultScheduler) launchExecutors(num int) {
	for i := 0; i < num; i++ {
		// Add tasks to task manager
		task := mesos.Task{
			Name:    proto.String("Sprint_" + strconv.Itoa(i)),
			TaskId:  &mesos.TaskID{Value: proto.String(strconv.Itoa(i))},
			AgentId: &mesos.AgentID{Value: proto.String("")},
			State:   mesos.TaskState_TASK_STAGING.Enum(),
		}
		c.manager.Add(&task)
	}
}

// Main event loop that listens on channels forever until framework terminates.
func (c *DefaultScheduler) listen() {
	for {
		select {
		case t := <-c.events:
			switch t.GetType() {
			case sched.Event_SUBSCRIBED:
				fmt.Println("Subscribe event.")
				go c.handlers.Subscribe(t.GetSubscribed())
			case sched.Event_ERROR:
				go c.handlers.Error(t.GetError())
			case sched.Event_FAILURE:
				go c.handlers.Failure(t.GetFailure())
			case sched.Event_INVERSE_OFFERS:
				go c.handlers.InverseOffer(t.GetInverseOffers())
			case sched.Event_MESSAGE:
				go c.handlers.Message(t.GetMessage())
			case sched.Event_OFFERS:
				log.Println("Offers...")
				go c.handlers.Offers(t.GetOffers())
			case sched.Event_RESCIND:
				go c.handlers.Rescind(t.GetRescind())
			case sched.Event_RESCIND_INVERSE_OFFER:
				go c.handlers.RescindInverseOffer(t.GetRescindInverseOffer())
			case sched.Event_UPDATE:
				go c.handlers.Update(t.GetUpdate())
			case sched.Event_HEARTBEAT:
			case sched.Event_UNKNOWN:
				fmt.Println("Unknown event recieved.")
			default:
			}
		case k := <-c.calls:
			switch k.GetType() {
			case sched.Call_ACCEPT:
				accept := k.GetAccept()
				c.FramworkInfo.Id = k.FrameworkId
				c.Accept(accept.OfferIds, accept.Operations, accept.Filters)
			case sched.Call_ACCEPT_INVERSE_OFFERS:
			case sched.Call_ACKNOWLEDGE:
				ack := k.GetAcknowledge()
				c.Acknowledge(ack.GetAgentId(), ack.GetTaskId(), ack.GetUuid())
			case sched.Call_DECLINE:
				decline := k.GetDecline()
				c.Decline(decline.GetOfferIds(), decline.GetFilters())
			case sched.Call_DECLINE_INVERSE_OFFERS:
			case sched.Call_MESSAGE:
			case sched.Call_KILL:
			case sched.Call_RECONCILE:
			case sched.Call_REVIVE:
			case sched.Call_SUBSCRIBE:
				// TODO decide on how we want to set framework info during subscribe call.
				c.FramworkInfo = k.GetSubscribe().GetFrameworkInfo()
			case sched.Call_SUPPRESS:
				c.Suppress()
			case sched.Call_SHUTDOWN:
			case sched.Call_TEARDOWN:
			default:

			}
		}
	}
}

// Make a subscription call to mesos.
func (c *DefaultScheduler) Subscribe() error {
	call := &sched.Call{
		Type: sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: c.FramworkInfo,
		},
	}
	go func() {
		for {
			resp, err := c.client.Request(call)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(recordio.Decode(resp.Body, c.events))
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
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_TEARDOWN.Enum(),
	}
	resp, err := c.client.Request(teardown)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(resp)
}

// Accepts offers from mesos master
func (c *DefaultScheduler) Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters) {
	accept := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
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

func (c *DefaultScheduler) Decline(offerIds []*mesos.OfferID, filters *mesos.Filters) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     &sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.client.Request(decline)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *DefaultScheduler) Revive() {
	revive := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.client.Request(revive)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *DefaultScheduler) Kill(taskId *mesos.TaskID, agentid *mesos.AgentID) {
	// Probably want some validation that this is a valid task and valid agentid.
	kill := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
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

func (c *DefaultScheduler) Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID) {
	shutdown := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
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
func (c *DefaultScheduler) Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte) {
	fmt.Println("acknowledge event.")
	acknowledge := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &sched.Call_Acknowledge{AgentId: agentId, TaskId: taskId, Uuid: uuid},
	}
	resp, err := c.client.Request(acknowledge)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Reconcile(tasks []*mesos.Task) {
	reconcile := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_RECONCILE.Enum(),
	}
	resp, err := c.client.Request(reconcile)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte) {
	message := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
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
func (c *DefaultScheduler) SchedRequest(resources []*mesos.Request) {
	request := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
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
	supress := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_SUPPRESS.Enum(),
	}
	resp, err := c.client.Request(supress)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}
