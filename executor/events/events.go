package events

import (
	"fmt"
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
}

type DefaultExecutorEvents struct {
}

func (d *DefaultExecutorEvents) Subscribed(sub *mesos_v1_executor.Event_Subscribed) {
	sub.GetFrameworkInfo()
}

func (d *DefaultExecutorEvents) Launch(launch *mesos_v1_executor.Event_Launch) {
	fmt.Println(launch.GetTask())
}

func (d *DefaultExecutorEvents) LaunchGroup(launchGroup *mesos_v1_executor.Event_LaunchGroup) {
	fmt.Println(launchGroup.GetTaskGroup())
}

func (d *DefaultExecutorEvents) Kill(kill *mesos_v1_executor.Event_Kill) {
	fmt.Printf("%v, %v\n", kill.GetTaskId(), kill.GetKillPolicy())
}
func (d *DefaultExecutorEvents) Acknowledged(acknowledge *mesos_v1_executor.Event_Acknowledged) {
	fmt.Printf("%v\n", acknowledge.GetTaskId())
}
func (d *DefaultExecutorEvents) Message(message *mesos_v1_executor.Event_Message) {
	fmt.Printf("%v\n", message.GetData())
}
func (d *DefaultExecutorEvents) Shutdown() {

}
func (d *DefaultExecutorEvents) Error(error *mesos_v1_executor.Event_Error) {
	fmt.Printf("%v\n", error.GetMessage())
}
