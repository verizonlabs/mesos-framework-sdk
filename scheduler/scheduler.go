package scheduler

/*
Scheduler:

*/
import (
	"fmt"
	"log"
	"mesos-framework-sdk/client"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/recordio"
	"mesos-framework-sdk/scheduler/events"
	"time"
)

const (
	subscribeRetry = 2
)

type Scheduler struct {
	FramworkInfo *mesos.FrameworkInfo
	client       *client.Client
	events       chan *sched.Event
	handlers     events.SchedulerEvent
}

func NewScheduler(c *client.Client, info *mesos.FrameworkInfo, handlers events.SchedulerEvent) *Scheduler {
	return &Scheduler{
		client:       c,
		events:       make(chan *sched.Event),
		FramworkInfo: info,
		handlers:     handlers,
	}
}

// Events
func (c *Scheduler) Run() {
	// If we don't have a framework id, subscribe.
	if c.FramworkInfo.GetId().GetValue() == "" {
		_, err := c.Subscribe()
		if err != nil {
			log.Println(err.Error())
		}
		c.listen()
	}
	// Otherwise we're already connected. Just listen for events.
	c.listen()
}

// Main event loop that listens on channels forever until framework terminates.
func (c *Scheduler) listen() {
	for {
		switch t := <-c.events; t.GetType() {
		case sched.Event_SUBSCRIBED:
			go c.handlers.Subscribe(t.GetSubscribed())
			break
		case sched.Event_ERROR:
			go c.handlers.Error(t.GetError())
			break
		case sched.Event_FAILURE:
			go c.handlers.Failure(t.GetFailure())
			break
		case sched.Event_INVERSE_OFFERS:
			go c.handlers.InverseOffer(t.GetInverseOffers())
			break
		case sched.Event_MESSAGE:
			go c.handlers.Message(t.GetMessage())
			break
		case sched.Event_OFFERS:
			go c.handlers.Offers(t.GetOffers())
			break
		case sched.Event_RESCIND:
			go c.handlers.Rescind(t.GetRescind())
			break
		case sched.Event_RESCIND_INVERSE_OFFER:
			go c.handlers.RescindInverseOffer(t.GetRescindInverseOffer())
			break
		case sched.Event_UPDATE:
			go c.handlers.Update(t.GetUpdate())
			break
		case sched.Event_HEARTBEAT:
			break
		case sched.Event_UNKNOWN:
			fmt.Println("Unknown event recieved.")
			break
		}
	}

}

// Make a subscription call to mesos.
func (c *Scheduler) Subscribe() (<-chan *sched.Event, error) {
	// We really want the ID after the call...
	call := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_SUBSCRIBE.Enum(),
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

	return c.events, nil
}

// Send a teardown request to mesos master.
func (c *Scheduler) Teardown() {
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

// Skeleton funcs for the rest of the calls.

// Accepts offers from mesos master
func (c *Scheduler) Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters) {
	accept := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      &sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.client.Request(accept)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *Scheduler) Decline(offerIds []*mesos.OfferID, filters *mesos.Filters) {
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
func (c *Scheduler) Revive() {

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

func (c *Scheduler) Kill(taskId *mesos.TaskID, agentid *mesos.AgentID) {
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

func (c *Scheduler) Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID) {
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
func (c *Scheduler) Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte) {
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

func (c *Scheduler) Reconcile(tasks []*mesos.Task) {
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

func (c *Scheduler) Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte) {
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
func (c *Scheduler) SchedRequest(resources []*mesos.Request) {
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
