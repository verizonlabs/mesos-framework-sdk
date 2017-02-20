package events

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/scheduler"
	"mesos-framework-sdk/task_manager"
)

// Mock type that satisfies interface.
type EventController struct {
	scheduler   *scheduler.DefaultScheduler
	frameworkId *mesos_v1.FrameworkID
	taskmanager task_manager.TaskManager
	events      chan *sched.Event
}

func NewDefaultEventController(scheduler *scheduler.DefaultScheduler, manager task_manager.TaskManager, eventChan chan *sched.Event) *EventController {
	return &EventController{
		taskmanager: manager,
		scheduler:   scheduler,
		events:      eventChan,
	}
}

func (s *EventController) Run() {
	if s.scheduler.FrameworkInfo().GetId().GetValue() == "" {
		err := s.scheduler.Subscribe(s.events)
		if err != nil {
			log.Printf("Error: %v", err.Error())
		}
		s.frameworkId = s.scheduler.FrameworkInfo().GetId()
		s.Listen()
	} else {
		s.Listen()
	}

}

// Main event loop that listens on channels forever until framework terminates.
func (s *EventController) Listen() {
	for {
		select {
		case t := <-s.events:
			switch t.GetType() {
			case sched.Event_SUBSCRIBED:
				fmt.Println("Subscribe event.")
				go s.Subscribe(t.GetSubscribed())
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
			default:
			}
		}
	}
}

func (s *EventController) Subscribe(subEvent *sched.Event_Subscribed) {
	fmt.Printf("Subscribed event recieved: %v\n", *subEvent)
	s.frameworkId = subEvent.GetFrameworkId()
	info := s.scheduler.FrameworkInfo()
	s.scheduler.FramworkInfo = &mesos_v1.FrameworkInfo{
		Id:              s.frameworkId,
		Capabilities:    info.Capabilities,
		FailoverTimeout: info.FailoverTimeout,
		Checkpoint:      info.Checkpoint,
		Hostname:        info.Hostname,
		Labels:          info.Labels,
		Name:            info.Name,
		Principal:       info.Principal,
		Role:            info.Role,
		User:            info.User,
		WebuiUrl:        info.WebuiUrl,
	}
}

func (s *EventController) Offers(offerEvent *sched.Event_Offers) {
	fmt.Println("Offers event recieved.")
	var offerIDs []*mesos_v1.OfferID

	for num, offer := range offerEvent.GetOffers() {
		fmt.Printf("Offer number: %v, Offer info: %v\n", num, offer)
		offerIDs = append(offerIDs, offer.GetId())
	}

	// Check task manager for any active tasks.
	if s.taskmanager.HasQueuedTasks() && s.frameworkId.GetValue() != "" {
		var taskList []*mesos_v1.TaskInfo
		// TODO check if resources are available for this particular task before launch.
		for i, task := range s.taskmanager.Tasks() {
			t := &mesos_v1.TaskInfo{
				Name:      task.Name,
				TaskId:    task.TaskId,
				AgentId:   offerEvent.Offers[0].AgentId,
				Resources: offerEvent.Offers[0].Resources,
				Executor: &mesos_v1.ExecutorInfo{
					ExecutorId:  &mesos_v1.ExecutorID{Value: proto.String(i)},
					FrameworkId: s.frameworkId,
					Command:     &mesos_v1.CommandInfo{Value: proto.String("sleep 5")},
				},
			}
			taskList = append(taskList, t)
		}

		var operations []*mesos_v1.Offer_Operation

		offer := &mesos_v1.Offer_Operation{
			Type:   mesos_v1.Offer_Operation_LAUNCH.Enum(),
			Launch: &mesos_v1.Offer_Operation_Launch{TaskInfos: taskList}}

		operations = append(operations, offer)

		s.scheduler.Accept(offerIDs, operations, nil)

	} else {
		var ids []*mesos_v1.OfferID
		for _, v := range offerEvent.GetOffers() {
			ids = append(ids, v.GetId())
		}
		// decline offers.
		fmt.Println("Declining offers.")
		s.scheduler.Decline(ids, &mesos_v1.Filters{})
		s.scheduler.Suppress()
	}
}
func (s *EventController) Rescind(rescindEvent *sched.Event_Rescind) {
	fmt.Printf("Rescind event recieved.: %v\n", *rescindEvent)
	rescindEvent.GetOfferId().GetValue()
}
func (s *EventController) Update(updateEvent *sched.Event_Update) {
	fmt.Printf("Update recieved for: %v\n", *updateEvent.GetStatus())

	id := s.taskmanager.Get(updateEvent.GetStatus().GetTaskId())
	s.taskmanager.SetTaskState(id, updateEvent.GetStatus().State)

	status := updateEvent.GetStatus()
	s.scheduler.Acknowledge(status.GetAgentId(), status.GetTaskId(), status.GetUuid())
}
func (s *EventController) Message(msg *sched.Event_Message) {
	fmt.Printf("Message event recieved: %v\n", *msg)
}
func (s *EventController) Failure(fail *sched.Event_Failure) {
	fmt.Printf("Failured event recieved: %v\n", *fail)
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
