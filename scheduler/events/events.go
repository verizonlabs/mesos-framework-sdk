package events

import (
	"fmt"
	"mesos-framework-sdk/include/scheduler"
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
}

func NewSchedulerEvents() *SchedEvent {
	return &SchedEvent{}
}

func (s *SchedEvent) Subscribe(subEvent *mesos_v1_scheduler.Event_Subscribed) {
	fmt.Printf("Subscribed event recieved: %v\n", *subEvent)
}
func (s *SchedEvent) Offers(offerEvent *mesos_v1_scheduler.Event_Offers) {
	fmt.Println("Offers event recieved.")
	for num, offer := range offerEvent.GetOffers() {
		fmt.Printf("Offer number: %v, Offer info: %v\n", num, offer)
	}

}
func (s *SchedEvent) Rescind(rescindEvent *mesos_v1_scheduler.Event_Rescind) {
	fmt.Printf("Rescind event recieved.: %v\n", *rescindEvent)
}
func (s *SchedEvent) Update(updateEvent *mesos_v1_scheduler.Event_Update) {
	fmt.Printf("Update recieved for: %v\n", *updateEvent.GetStatus())
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
