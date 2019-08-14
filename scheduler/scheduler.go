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

package scheduler

/*
Scheduler struct defines an interface of the default calls for the scheduler, as well as holding information
regarding the framework, client, and events handler.

End users may wish to make their own scheduler and satisfy this interface.

A default scheduler is provided for those who wish to keep the default implementations: All calls simply create
the protobuf required for the call and send it off to the client.

End users should only create their own scheduler if they wish to change the behavior of their calls.
*/
import (
	"errors"
	"net/http"
	"sync"

	"github.com/carlonelong/mesos-framework-sdk/client"
	mesos_v1 "github.com/carlonelong/mesos-framework-sdk/include/mesos/v1"
	sched "github.com/carlonelong/mesos-framework-sdk/include/mesos/v1/scheduler"
	"github.com/carlonelong/mesos-framework-sdk/logging"
	"github.com/carlonelong/mesos-framework-sdk/recordio"
)

type Scheduler interface {
	FrameworkInfo() *mesos_v1.FrameworkInfo

	// Default Calls for scheduler
	Subscribe(chan *sched.Event) (*http.Response, error)
	Teardown() (*http.Response, error)
	Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters) (*http.Response, error)
	Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters) (*http.Response, error)
	Revive() (*http.Response, error)
	Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID) (*http.Response, error)
	Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID) (*http.Response, error)
	Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte) (*http.Response, error)
	Reconcile(tasks []*mesos_v1.TaskInfo) (*http.Response, error)
	Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte) (*http.Response, error)
	SchedRequest(resources []*mesos_v1.Request) (*http.Response, error)
	Suppress() (*http.Response, error)
}

// Default Scheduler can be used as a higher-level construct.
type DefaultScheduler struct {
	frameworkInfo *mesos_v1.FrameworkInfo
	Client        client.Client
	logger        logging.Logger
	IsSuppressed  bool
	sync.RWMutex
}

func NewDefaultScheduler(c client.Client, info *mesos_v1.FrameworkInfo, logger logging.Logger) *DefaultScheduler {
	return &DefaultScheduler{
		Client:        c,
		frameworkInfo: info,
		logger:        logger,
		IsSuppressed:  false,
	}
}

func (c *DefaultScheduler) FrameworkInfo() *mesos_v1.FrameworkInfo {
	return c.frameworkInfo
}

// Make a subscription call to mesos.
// Channel passed is the channel for Event Controller.
func (c *DefaultScheduler) Subscribe(eventChan chan *sched.Event) (*http.Response, error) {
	call := &sched.Call{
		Type: sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: c.frameworkInfo,
		},
		FrameworkId: c.frameworkInfo.Id,
	}

	// If we disconnect we need to reset the stream ID. For this reason always start with a fresh stream ID.
	// Otherwise we'll never be able to reconnect.
	c.Client.SetStreamID("")

	resp, err := c.Client.Request(call)
	if err != nil {
		return nil, err
	}
	return resp, recordio.Decode(resp.Body, eventChan)
}

// Send a teardown request to mesos master.
func (c *DefaultScheduler) Teardown() (*http.Response, error) {
	teardown := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_TEARDOWN.Enum(),
	}
	resp, err := c.Client.Request(teardown)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}

	c.logger.Emit(logging.INFO, "Tearing down framework")
	return resp, err
}

// Accepts offers from mesos master
func (c *DefaultScheduler) Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters) (*http.Response, error) {
	accept := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      &sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.Client.Request(accept)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}

	c.logger.Emit(logging.INFO, "Accepting %d offers for %d tasks", len(offerIds), len(tasks))
	return resp, err
}

func (c *DefaultScheduler) Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters) (*http.Response, error) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     &sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.Client.Request(decline)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}

	l := len(offerIds)
	if l > 0 {
		c.logger.Emit(logging.INFO, "Declining %d offers", len(offerIds))
	}

	return resp, err
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *DefaultScheduler) Revive() (*http.Response, error) {
	c.RLock()
	if !c.IsSuppressed {
		c.RUnlock()
		return nil, nil
	}
	c.RUnlock()

	revive := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.Client.Request(revive)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	c.Lock()
	c.IsSuppressed = false
	c.Unlock()

	c.logger.Emit(logging.INFO, "Reviving offers")
	return resp, err
}

func (c *DefaultScheduler) Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID) (*http.Response, error) {
	kill := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_KILL.Enum(),
		Kill:        &sched.Call_Kill{TaskId: taskId, AgentId: agentid},
	}

	resp, err := c.Client.Request(kill)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	// Kill returns a 202 accepted.
	if resp.StatusCode == 202 {
		c.logger.Emit(logging.INFO, "Killing task %s", taskId.GetValue())
	}
	return resp, err
}

func (c *DefaultScheduler) Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID) (*http.Response, error) {
	shutdown := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_SHUTDOWN.Enum(),
		Shutdown: &sched.Call_Shutdown{
			ExecutorId: execId,
			AgentId:    agentId,
		},
	}
	resp, err := c.Client.Request(shutdown)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	c.logger.Emit(logging.INFO, "Shutting down")
	return resp, err
}

func (c *DefaultScheduler) Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte) (*http.Response, error) {

	// Note that with the new API, schedulers are responsible for explicitly acknowledging the receipt of status
	// updates that have “status.uuid()” set.
	// These status updates are reliably retried until they are acknowledged by the scheduler.
	// The scheduler must not acknowledge status updates that do not have "status.uuid()" set as they are not retried.
	if uuid == nil {
		return nil, errors.New("No uuid passed in to ACK.")
	}

	acknowledge := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &sched.Call_Acknowledge{
			AgentId: agentId,
			TaskId:  taskId,
			Uuid:    uuid,
		},
	}

	resp, err := c.Client.Request(acknowledge)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	return resp, err
}

func (c *DefaultScheduler) Reconcile(tasks []*mesos_v1.TaskInfo) (*http.Response, error) {
	reconcileTasks := make([]*sched.Call_Reconcile_Task, 0, len(tasks))
	for _, task := range tasks {
		reconcileTasks = append(reconcileTasks, &sched.Call_Reconcile_Task{
			AgentId: task.GetAgentId(),
			TaskId:  task.GetTaskId(),
		})
	}

	reconcile := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_RECONCILE.Enum(),
		Reconcile: &sched.Call_Reconcile{
			Tasks: reconcileTasks,
		},
	}
	resp, err := c.Client.Request(reconcile)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}

	c.logger.Emit(logging.INFO, "Reconciling %d tasks", len(tasks))
	return resp, err
}

func (c *DefaultScheduler) Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte) (*http.Response, error) {
	message := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_MESSAGE.Enum(),
		Message: &sched.Call_Message{
			AgentId:    agentId,
			ExecutorId: executorId,
			Data:       data,
		},
	}
	resp, err := c.Client.Request(message)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	c.logger.Emit(logging.INFO, "Message received from agent %s and executor %s", agentId.GetValue(), executorId.GetValue())
	return resp, err
}

// NOTE: This method is only kept to conform to official Mesos codebase.  This does nothing.
func (c *DefaultScheduler) SchedRequest(resources []*mesos_v1.Request) (*http.Response, error) {
	request := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_REQUEST.Enum(),
		Request: &sched.Call_Request{
			Requests: resources,
		},
	}
	resp, err := c.Client.Request(request)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
		return nil, err
	}
	c.logger.Emit(logging.INFO, "Requesting resources")
	return resp, err
}

// Makes a call to Mesos to suppress any further offers.
func (c *DefaultScheduler) Suppress() (*http.Response, error) {
	c.RLock()
	if c.IsSuppressed {
		c.RUnlock()
		return nil, nil
	}
	c.RUnlock()

	suppress := &sched.Call{
		FrameworkId: c.frameworkInfo.GetId(),
		Type:        sched.Call_SUPPRESS.Enum(),
	}
	resp, err := c.Client.Request(suppress)
	if err != nil {
		c.logger.Emit(logging.ERROR, err.Error())
	} else {
		c.Lock()
		c.IsSuppressed = true
		c.Unlock()

		c.logger.Emit(logging.INFO, "Suppressing offers")
	}

	return resp, err
}
