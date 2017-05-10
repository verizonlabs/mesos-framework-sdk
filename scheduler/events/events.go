package events

import "mesos-framework-sdk/include/mesos_v1_scheduler"

/*
The events package will hook in how an end user wants to deal with events received by the scheduler.
*/

// Define the behavior of how an end user will deal with events.
type SchedulerEvent interface {
	Subscribed(*mesos_v1_scheduler.Event_Subscribed)
	Offers(*mesos_v1_scheduler.Event_Offers)
	Rescind(*mesos_v1_scheduler.Event_Rescind)
	Update(*mesos_v1_scheduler.Event_Update)
	Message(*mesos_v1_scheduler.Event_Message)
	Failure(*mesos_v1_scheduler.Event_Failure)
	Error(*mesos_v1_scheduler.Event_Error)
	InverseOffer(*mesos_v1_scheduler.Event_InverseOffers)
	RescindInverseOffer(*mesos_v1_scheduler.Event_RescindInverseOffer)
	Run()
	Listen()
}
