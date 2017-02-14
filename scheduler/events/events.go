package events

import (
	"fmt"
)

/*
The events package will hook in how an end user wants to deal with events received by the scheduler.
*/

// Define the behavior of how an end user will deal with events.
type SchedulerEvent interface {
	Subscribe()
	Offers()
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
func (s *SchedEvent) Offers() {
	fmt.Println("Offers event recieved.")
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
