package events

import (
	"mesos-framework-sdk/include/executor"
)

/*
Sent by the agent whenever it needs to assign a new task to the executor. The executor is required to send an
UPDATE message back to the agent indicating the success or failure of the task initialization.

The executor must maintain a list of unacknowledged tasks (see SUBSCRIBE in Calls section).
If for some reason, the executor is disconnected from the agent,
these tasks must be sent as part of SUBSCRIBE request in the tasks field.

*/
// Interface for all events sent to a custom executor.
type ExecutorEvents interface {
	Subscribed(*mesos_v1_executor.Event_Subscribed)
	Launch(*mesos_v1_executor.Event_Launch)
	LaunchGroup(*mesos_v1_executor.Event_LaunchGroup)
	Kill(*mesos_v1_executor.Event_Kill)
	Acknowledged(*mesos_v1_executor.Event_Acknowledged)
	Message(*mesos_v1_executor.Event_Message)
	Shutdown()
	Error(*mesos_v1_executor.Event_Error)
	Run()
}
