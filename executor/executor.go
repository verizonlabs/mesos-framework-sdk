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

package executor

import (
	"github.com/verizonlabs/mesos-framework-sdk/client"
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	exec "github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1_executor"
	"github.com/verizonlabs/mesos-framework-sdk/logging"
	"github.com/verizonlabs/mesos-framework-sdk/recordio"
)

type Executor interface {
	Subscribe(chan *exec.Event) error
	Update(*mesos_v1.TaskStatus) error
	Message([]byte) error
}

type DefaultExecutor struct {
	frameworkId    *mesos_v1.FrameworkID
	executorId     *mesos_v1.ExecutorID
	client         client.Client
	logger         logging.Logger
	unackedTasks   map[*mesos_v1.TaskID]mesos_v1.TaskInfo
	unackedUpdates map[string]exec.Call_Update
}

// Creates a new default executor
func NewDefaultExecutor(
	f *mesos_v1.FrameworkID,
	e *mesos_v1.ExecutorID,
	c client.Client,
	lgr logging.Logger) Executor {

	return &DefaultExecutor{
		frameworkId:    f,
		executorId:     e,
		client:         c,
		logger:         lgr,
		unackedTasks:   make(map[*mesos_v1.TaskID]mesos_v1.TaskInfo),
		unackedUpdates: make(map[string]exec.Call_Update),
	}
}

func (e *DefaultExecutor) unacknowledgedTasks() []*mesos_v1.TaskInfo {
	numTasks := len(e.unackedTasks)
	if numTasks == 0 {
		return nil
	}

	tasks := make([]*mesos_v1.TaskInfo, 0, numTasks)
	for task := range e.unackedTasks {
		t := e.unackedTasks[task]
		tasks = append(tasks, &t)
	}

	return tasks
}

func (e *DefaultExecutor) unacknowledgedUpdates() []*exec.Call_Update {
	numUpdates := len(e.unackedUpdates)
	if numUpdates == 0 {
		return nil
	}

	updates := make([]*exec.Call_Update, 0, numUpdates)
	for update := range e.unackedUpdates {
		u := e.unackedUpdates[update]
		updates = append(updates, &u)
	}

	return updates
}

func (e *DefaultExecutor) Subscribe(eventChan chan *exec.Event) error {
	subscribe := &exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_SUBSCRIBE.Enum(),
		Subscribe: &exec.Call_Subscribe{
			UnacknowledgedTasks:   e.unacknowledgedTasks(),
			UnacknowledgedUpdates: e.unacknowledgedUpdates(),
		},
	}

	resp, err := e.client.Request(subscribe)
	if err != nil {
		return err
	} else {
		return recordio.Decode(resp.Body, eventChan)
	}
}

func (e *DefaultExecutor) Update(taskStatus *mesos_v1.TaskStatus) error {
	update := &exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_UPDATE.Enum(),
		Update: &exec.Call_Update{
			Status: taskStatus,
		},
	}
	_, err := e.client.Request(update)

	return err
}

func (e *DefaultExecutor) Message(data []byte) error {
	message := &exec.Call{
		FrameworkId: e.frameworkId,
		ExecutorId:  e.executorId,
		Type:        exec.Call_MESSAGE.Enum(),
		Message: &exec.Call_Message{
			Data: data,
		},
	}
	_, err := e.client.Request(message)

	return err
}
