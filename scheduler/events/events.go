package events

import (
	sched "mesos-framework-sdk/include/scheduler"
	"net/http"
)

type SchedEvents struct {
	Handlers map[sched.Event_Type]http.HandlerFunc
}

func (e *SchedEvents) Subscribed() {

}

func (e *SchedEvents) Offers() {

}

func (e *SchedEvents) Rescind() {

}

func (e *SchedEvents) Update() {

}

func (e *SchedEvents) Message() {

}

func (e *SchedEvents) Failure() {

}

func (e *SchedEvents) Error() {

}

func (e *SchedEvents) Heartbeat() {

}
