package events

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/task_manager"
)

/*
The events package will hook in how an end user wants to deal with events received by the scheduler.
*/

// Define the behavior of how an end user will deal with events.
type SchedulerEvent interface {
	Subscribe(*mesos_v1_scheduler.Event_Subscribed)
	Offers(*mesos_v1_scheduler.Event_Offers)
	Rescind(*mesos_v1_scheduler.Event_Rescind)
	Update(*mesos_v1_scheduler.Event_Update)
	Message(*mesos_v1_scheduler.Event_Message)
	Failure(*mesos_v1_scheduler.Event_Failure)
	Error(*mesos_v1_scheduler.Event_Error)
	InverseOffer(*mesos_v1_scheduler.Event_InverseOffers)
	RescindInverseOffer(*mesos_v1_scheduler.Event_RescindInverseOffer)
}

// Mock type that satisfies interface.
type SchedEvent struct {
	frameworkId *mesos_v1.FrameworkID
	taskmanager task_manager.TaskManager
	channel     chan *mesos_v1_scheduler.Call // Channel used to talk to scheduler for making calls.
}

func NewSchedulerEvents(manager task_manager.TaskManager, callChan chan *mesos_v1_scheduler.Call) *SchedEvent {
	return &SchedEvent{
		taskmanager: manager,
		channel:     callChan,
	}
}

func (s *SchedEvent) Subscribe(subEvent *mesos_v1_scheduler.Event_Subscribed) {
	fmt.Printf("Subscribed event recieved: %v\n", *subEvent)
	s.frameworkId = subEvent.GetFrameworkId()
}

func (s *SchedEvent) Offers(offerEvent *mesos_v1_scheduler.Event_Offers) {
	fmt.Println("Offers event recieved.")
	var offerIDs []*mesos_v1.OfferID

	for num, offer := range offerEvent.GetOffers() {
		fmt.Printf("Offer number: %v, Offer info: %v\n", num, offer)
		offerIDs = append(offerIDs, offer.GetId())
	}

	// Check task manager for any active tasks.
	if s.taskmanager.HasQueuedTasks() && s.frameworkId.GetValue() != "" {
		var taskList []*mesos_v1.TaskInfo

		// Create the task infos
		// TODO have task manager convert Task -> TaskInfo
		for i, task := range s.taskmanager.Tasks() {
			t := &mesos_v1.TaskInfo{
				Name:      task.Name,
				TaskId:    task.TaskId,
				AgentId:   offerEvent.Offers[0].AgentId,
				Resources: offerEvent.Offers[0].Resources,
				Executor: &mesos_v1.ExecutorInfo{
					ExecutorId:  &mesos_v1.ExecutorID{Value: proto.String(i)},
					FrameworkId: s.frameworkId,
					Command:     &mesos_v1.CommandInfo{Value: proto.String("sleep 10")},
				},
			}
			taskList = append(taskList, t)
			s.taskmanager.Delete(&mesos_v1.TaskID{Value: proto.String(i)})
		}
		var operations []*mesos_v1.Offer_Operation
		offer := &mesos_v1.Offer_Operation{
			Type:   mesos_v1.Offer_Operation_LAUNCH.Enum(),
			Launch: &mesos_v1.Offer_Operation_Launch{TaskInfos: taskList}}

		operations = append(operations, offer)
		// Write the call to the scheduler.
		go func() {
			s.channel <- &mesos_v1_scheduler.Call{
				FrameworkId: s.frameworkId,
				Type:        mesos_v1_scheduler.Call_ACCEPT.Enum(),
				Accept: &mesos_v1_scheduler.Call_Accept{
					OfferIds:   offerIDs,
					Operations: operations,
				},
			}
		}()

	}
}
func (s *SchedEvent) Rescind(rescindEvent *mesos_v1_scheduler.Event_Rescind) {
	fmt.Printf("Rescind event recieved.: %v\n", *rescindEvent)
}
func (s *SchedEvent) Update(updateEvent *mesos_v1_scheduler.Event_Update) {
	fmt.Printf("Update recieved for: %v\n", *updateEvent.GetStatus())
	id := s.taskmanager.Get(updateEvent.GetStatus().GetTaskId())
	s.taskmanager.SetTaskState(id, updateEvent.GetStatus().State)
	go func() {
		s.channel <- &mesos_v1_scheduler.Call{
			FrameworkId: s.frameworkId,
			Type:        mesos_v1_scheduler.Call_ACKNOWLEDGE.Enum(),
			Acknowledge: &mesos_v1_scheduler.Call_Acknowledge{
				AgentId: updateEvent.GetStatus().GetAgentId(),
				TaskId:  updateEvent.GetStatus().GetTaskId(),
				Uuid:    updateEvent.GetStatus().GetUuid(),
			},
		}
	}()
}
func (s *SchedEvent) Message(msg *mesos_v1_scheduler.Event_Message) {
	fmt.Printf("Message event recieved: %v\n", *msg)
}
func (s *SchedEvent) Failure(fail *mesos_v1_scheduler.Event_Failure) {
	fmt.Printf("Failured event recieved: %v\n", *fail)
}
func (s *SchedEvent) Error(err *mesos_v1_scheduler.Event_Error) {
	fmt.Printf("Error event recieved: %v\n", err)
}
func (s *SchedEvent) InverseOffer(ioffers *mesos_v1_scheduler.Event_InverseOffers) {
	fmt.Printf("Inverse Offer event recieved: %v\n", ioffers)
}
func (s *SchedEvent) RescindInverseOffer(rioffers *mesos_v1_scheduler.Event_RescindInverseOffer) {
	fmt.Printf("Rescind Inverse Offer event recieved: %v\n", rioffers)
}
