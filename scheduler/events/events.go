// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Listen(chan *mesos_v1_scheduler.Event)
}
