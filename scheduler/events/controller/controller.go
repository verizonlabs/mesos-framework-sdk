package controller

import (
	"fmt"
	"log"
	"mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/logging"
	"mesos-framework-sdk/resources"
	"mesos-framework-sdk/resources/manager"
	"mesos-framework-sdk/scheduler"
	"mesos-framework-sdk/task/manager"
)

// Mock type that satisfies interface.
type EventController struct {
	scheduler       *scheduler.DefaultScheduler
	taskmanager     *task_manager.DefaultTaskManager
	resourcemanager *manager.DefaultResourceManager
	events          chan *sched.Event
	logger          logging.Logger
}

func NewDefaultEventController(scheduler *scheduler.DefaultScheduler, manager *task_manager.DefaultTaskManager, resourceManager *manager.DefaultResourceManager, eventChan chan *sched.Event, logger logging.Logger) *EventController {
	return &EventController{
		taskmanager:     manager,
		scheduler:       scheduler,
		events:          eventChan,
		resourcemanager: resourceManager,
		logger:          logger,
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
			case sched.Event_ERROR:
				go s.Error(t.GetError())
			case sched.Event_FAILURE:
				go s.Failure(t.GetFailure())
			case sched.Event_INVERSE_OFFERS:
				go s.InverseOffer(t.GetInverseOffers())
			case sched.Event_MESSAGE:
				go s.Message(t.GetMessage())
			case sched.Event_OFFERS:
				go s.Offers(t.GetOffers())
			case sched.Event_RESCIND:
				go s.Rescind(t.GetRescind())
			case sched.Event_RESCIND_INVERSE_OFFER:
				go s.RescindInverseOffer(t.GetRescindInverseOffer())
			case sched.Event_UPDATE:
				go s.Update(t.GetUpdate())
			case sched.Event_HEARTBEAT:
			case sched.Event_UNKNOWN:
				fmt.Println("Unknown event recieved.")
			}
		}
	}
}

func (s *EventController) Offers(offerEvent *sched.Event_Offers) {

	// Check task manager for any active tasks.
	if s.taskmanager.HasQueuedTasks() {
		// Update our resources in the manager
		s.resourcemanager.AddOffers(offerEvent.GetOffers())

		offerIDs := []*mesos_v1.OfferID{}
		operations := []*mesos_v1.Offer_Operation{}

		for _, mesosTask := range s.taskmanager.QueuedTasks() {
			// See if we have resources.
			if s.resourcemanager.HasResources() {

				offer, err := s.resourcemanager.Assign(mesosTask)
				if err != nil {

					// It didn't match any offers.
					s.logger.Emit(logging.ERROR, err.Error())
					continue // We should decline.
				}

				t := &mesos_v1.TaskInfo{
					Name:      mesosTask.Name,
					TaskId:    mesosTask.GetTaskId(),
					AgentId:   offer.GetAgentId(),
					Command:   mesosTask.GetCommand(),
					Container: mesosTask.GetContainer(),
					Resources: mesosTask.GetResources(),
				}

				s.taskmanager.SetTaskLaunched(t)

				offerIDs = append(offerIDs, offer.Id)
				operations = append(operations, resources.LaunchOfferOperation([]*mesos_v1.TaskInfo{t}))
			}
		}
		s.scheduler.Accept(offerIDs, operations, nil)
	} else {
		var ids []*mesos_v1.OfferID
		for _, v := range offerEvent.GetOffers() {
			ids = append(ids, v.GetId())
		}

		// Decline and suppress offers until we're ready again.
		s.logger.Emit(logging.INFO, "Declining %d offers", len(ids))
		s.scheduler.Decline(ids, nil) // We want to make sure all offers are declined.
		s.scheduler.Suppress()
	}
}

func (s *EventController) Rescind(rescindEvent *sched.Event_Rescind) {
	s.logger.Emit(logging.INFO, "Rescind event recieved.: %v", *rescindEvent)
	rescindEvent.GetOfferId().GetValue()
}

func (s *EventController) Update(updateEvent *sched.Event_Update) {
	s.logger.Emit(logging.INFO, "Update recieved for: %v\n", *updateEvent.GetStatus())
	task := s.taskmanager.GetById(updateEvent.GetStatus().GetTaskId())
	// TODO: Handle more states in regard to tasks.
	if updateEvent.GetStatus().GetState() != mesos_v1.TaskState_TASK_FAILED {
		// Only set the task to "launched" if it didn't fail.
		s.taskmanager.SetTaskLaunched(task)
	}

	status := updateEvent.GetStatus()
	s.scheduler.Acknowledge(status.GetAgentId(), status.GetTaskId(), status.GetUuid())
}

func (s *EventController) Message(msg *sched.Event_Message) {
	s.logger.Emit(logging.INFO, "Message event recieved: %v", *msg)
}

func (s *EventController) Failure(fail *sched.Event_Failure) {
	s.logger.Emit(logging.ERROR, "Executor %s failed with status %d", fail.GetExecutorId().GetValue(), fail.GetStatus())
}

func (s *EventController) Error(err *sched.Event_Error) {
	s.logger.Emit(logging.ERROR, "Error event recieved: %v", err)
}

func (s *EventController) InverseOffer(ioffers *sched.Event_InverseOffers) {
	s.logger.Emit(logging.INFO, "Inverse Offer event recieved: %v", ioffers)
}

func (s *EventController) RescindInverseOffer(rioffers *sched.Event_RescindInverseOffer) {
	s.logger.Emit(logging.INFO, "Rescind Inverse Offer event recieved: %v", rioffers)
}
