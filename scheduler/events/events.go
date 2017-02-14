package events

import (
	"fmt"
	"mesos-framework-sdk/include/mesos"
)

/*
The events package will hook in how an end user wants to deal with events received by the scheduler.
*/

// Define the behavior of how an end user will deal with events.
type SchedulerEvent interface {
	Subscribe()
	Offers([]*mesos_v1.Offer)
	Rescind()
	Update()
	Message()
	Failure()
	Error()
	Heartbeat()
}

// Mock type that satisfies interface.
type SchedEvent struct {
}

func NewSchedulerEvents() *SchedEvent {
	return &SchedEvent{}
}

func (s *SchedEvent) Subscribe() {

}
func (s *SchedEvent) Offers(offers []*mesos_v1.Offer) {
	fmt.Println("Offers event recieved.")
	for num, offer := range offers {
		fmt.Printf("Offer number: %v, Offer info: %v", num, offer)
	}

}
func (s *SchedEvent) Rescind() {
	fmt.Println("Rescind event recieved.")
}
func (s *SchedEvent) Update() {

}
func (s *SchedEvent) Message() {
	fmt.Println("Message event recieved.")
}
func (s *SchedEvent) Failure() {
	fmt.Println("Failured event recieved.")
}
func (s *SchedEvent) Error() {
	fmt.Println("Error event recieved.")
}
func (s *SchedEvent) Heartbeat() {
	fmt.Println("Heartbeat event recieved.")
}
