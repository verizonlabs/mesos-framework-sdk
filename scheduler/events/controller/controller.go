package controller

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/resources/manager"
	"mesos-framework-sdk/scheduler"
	"mesos-framework-sdk/task_manager"
	"mesos-framework-sdk/utils"
	"strconv"
)

// Mock type that satisfies interface.
type EventController struct {
	scheduler       *scheduler.DefaultScheduler
	taskmanager     *task_manager.DefaultTaskManager
	resourcemanager *manager.DefaultResourceManager
	events          chan *sched.Event
}

func NewDefaultEventController(scheduler *scheduler.DefaultScheduler, manager *task_manager.DefaultTaskManager, resourceManager *manager.DefaultResourceManager, eventChan chan *sched.Event) *EventController {
	return &EventController{
		taskmanager:     manager,
		scheduler:       scheduler,
		events:          eventChan,
		resourcemanager: resourceManager,
	}
}

func (s *EventController) Subscribe(subEvent *sched.Event_Subscribed) {
	id := subEvent.GetFrameworkId()
	idVal := id.GetValue()
	s.scheduler.Info.Id = id
	log.Printf("Subscribed with an ID of %s", idVal)
}

func (s *EventController) Run() {
	if s.scheduler.FrameworkInfo().GetId() == nil {
		err := s.scheduler.Subscribe(s.events)
		if err != nil {
			log.Printf("Failed to subscribe: %s", err.Error())
		}

		// Wait here until we have our framework ID.
		select {
		case e := <-s.events:
			s.Subscribe(e.GetSubscribed())
		}
	}
	s.Listen()
}

// Main event loop that listens on channels forever until framework terminates.
func (s *EventController) Listen() {
	for {
		select {
		case t := <-s.events:
			switch t.GetType() {
			case sched.Event_SUBSCRIBED:
				log.Println("Subscribe event.")
			case sched.Event_ERROR:
				go s.Error(t.GetError())
			case sched.Event_FAILURE:
				go s.Failure(t.GetFailure())
			case sched.Event_INVERSE_OFFERS:
				go s.InverseOffer(t.GetInverseOffers())
			case sched.Event_MESSAGE:
				go s.Message(t.GetMessage())
			case sched.Event_OFFERS:
				log.Println("Offers...")
				go s.Offers(t.GetOffers())
			case sched.Event_RESCIND:
				go s.Rescind(t.GetRescind())
			case sched.Event_RESCIND_INVERSE_OFFER:
				go s.RescindInverseOffer(t.GetRescindInverseOffer())
			case sched.Event_UPDATE:
				go s.Update(t.GetUpdate())
			case sched.Event_HEARTBEAT:
				fmt.Println("Heart beat.")
			case sched.Event_UNKNOWN:
				fmt.Println("Unknown event recieved.")
			}
		}
	}
}

func (s *EventController) Offers(offerEvent *sched.Event_Offers) {
	fmt.Println("Offers event recieved.")
	//Reconcile any tasks.
	var reconcileTasks []*mesos_v1.Task
	s.scheduler.Reconcile(reconcileTasks)

	var offerIDs []*mesos_v1.OfferID

	for num, offer := range offerEvent.GetOffers() {
		fmt.Printf("Offer number: %v, Offer info: %v\n", num, offer)
		offerIDs = append(offerIDs, offer.GetId())
	}

	// Check task manager for any active tasks.
	if s.taskmanager.HasQueuedTasks() {
		// Update our resources in the manager
		s.resourcemanager.AddOffers(offerEvent.GetOffers())

		for _, item := range s.taskmanager.QueuedTasks() {
			// See if we have resources.
			if s.resourcemanager.HasResources() {
				taskList := []*mesos_v1.TaskInfo{} // Clear it out every time.
				operations := []*mesos_v1.Offer_Operation{}
				mesosTask := item

				offer, err := s.resourcemanager.Assign(mesosTask)
				if err != nil {
					// It didn't match any offers.
					log.Println(err.Error())
				}

				t := &mesos_v1.TaskInfo{
					Name:    mesosTask.Name,
					TaskId:  mesosTask.TaskId,
					AgentId: offer.AgentId,
					Command: &mesos_v1.CommandInfo{
						User:  proto.String("root"),
						Value: proto.String("/bin/sleep 10"),
					},
					Resources: []*mesos_v1.Resource{
						resources.CreateCpu(0.1, ""),
						resources.CreateMem(64.0, ""),
					},
				}

				taskList = append(taskList, t)

				operations = append(operations, resources.LaunchOfferOperation(taskList))

				log.Printf("Launching task %v\n", taskList)
				s.scheduler.Accept(offerIDs, operations, nil)
			}
		}
	} else {
		var ids []*mesos_v1.OfferID
		for _, v := range offerEvent.GetOffers() {
			ids = append(ids, v.GetId())
		}
		// decline offers.
		fmt.Println("Declining offers.")
		s.scheduler.Decline(ids, nil)
		s.scheduler.Suppress()
	}
}

func (s *EventController) Rescind(rescindEvent *sched.Event_Rescind) {
	fmt.Printf("Rescind event recieved.: %v\n", *rescindEvent)
	rescindEvent.GetOfferId().GetValue()
}

func (s *EventController) Update(updateEvent *sched.Event_Update) {
	fmt.Printf("Update recieved for: %v\n", *updateEvent.GetStatus())
	task := s.taskmanager.Get(updateEvent.GetStatus().GetTaskId())
	// TODO: Handle more states in regard to tasks.
	if updateEvent.GetStatus().GetState() != mesos_v1.TaskState_TASK_FAILED {
		// Only set the task to "launched" if it didn't fail.
		s.taskmanager.SetTaskLaunched(task)
	}

	status := updateEvent.GetStatus()
	s.scheduler.Acknowledge(status.GetAgentId(), status.GetTaskId(), status.GetUuid())
}

func (s *EventController) Message(msg *sched.Event_Message) {
	fmt.Printf("Message event recieved: %v\n", *msg)
}

func (s *EventController) Failure(fail *sched.Event_Failure) {
	log.Println("Executor " + fail.GetExecutorId().GetValue() + " failed with status " + strconv.Itoa(int(fail.GetStatus())))
}

func (s *EventController) Error(err *sched.Event_Error) {
	fmt.Printf("Error event recieved: %v\n", err)
}

func (s *EventController) InverseOffer(ioffers *sched.Event_InverseOffers) {
	fmt.Printf("Inverse Offer event recieved: %v\n", ioffers)
}

func (s *EventController) RescindInverseOffer(rioffers *sched.Event_RescindInverseOffer) {
	fmt.Printf("Rescind Inverse Offer event recieved: %v\n", rioffers)
}
