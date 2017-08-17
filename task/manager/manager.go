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

package manager

import (
	"encoding/json"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/retry"
	"sync"
	"time"
)

// Consts for mesos states.
const (
	RUNNING          = mesos_v1.TaskState_TASK_RUNNING
	KILLED           = mesos_v1.TaskState_TASK_KILLED
	LOST             = mesos_v1.TaskState_TASK_LOST
	GONE             = mesos_v1.TaskState_TASK_GONE
	STAGING          = mesos_v1.TaskState_TASK_STAGING
	STARTING         = mesos_v1.TaskState_TASK_STARTING // Default executor never sends this, it sends RUNNING directly.
	UNKNOWN          = mesos_v1.TaskState_TASK_UNKNOWN
	UNREACHABLE      = mesos_v1.TaskState_TASK_UNREACHABLE
	FINISHED         = mesos_v1.TaskState_TASK_FINISHED
	DROPPED          = mesos_v1.TaskState_TASK_DROPPED
	FAILED           = mesos_v1.TaskState_TASK_FAILED
	ERROR            = mesos_v1.TaskState_TASK_ERROR
	GONE_BY_OPERATOR = mesos_v1.TaskState_TASK_GONE_BY_OPERATOR
	KILLING          = mesos_v1.TaskState_TASK_KILLING
)

// Task manager holds information about tasks coming into the framework from the API
// It can set the state of a task.  How the implementation holds/handles those tasks
// is up to the end user.
type TaskManager interface {
	Add(...*Task) error
	Delete(...*Task) error
	Get(*string) (*Task, error)
	GetById(id *mesos_v1.TaskID) (*Task, error)
	HasTask(*mesos_v1.TaskInfo) bool
	Update(...*Task) error
	AllByState(state mesos_v1.TaskState) ([]*Task, error)
	TotalTasks() int
	All() ([]*Task, error)
}

// Used to hold information about task states in the task manager.
// Task and its fields should be public so that we can encode/decode this.
type Task struct {
	lock      sync.Mutex
	Info      *mesos_v1.TaskInfo
	State     mesos_v1.TaskState
	Filters   []task.Filter
	Retry     *retry.TaskRetry
	Instances int
	IsKill    bool
	GroupInfo GroupInfo
}

type GroupInfo struct {
	GroupName string
	InGroup   bool
}

func NewTask(i *mesos_v1.TaskInfo, s mesos_v1.TaskState, f []task.Filter, r *retry.TaskRetry, n int, g GroupInfo) *Task {
	return &Task{
		Info:      i,
		State:     s,
		Filters:   f,
		Retry:     r,
		Instances: n,
		GroupInfo: g,
		IsKill:    false,
		lock:      sync.Mutex{},
	}
}

// TODO (tim): Create a serialize/deserialize mechanism from string <-> struct to avoid costly encoding?

func (t *Task) Reschedule(revive chan *Task) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.State = mesos_v1.TaskState_TASK_STAGING

	// Minimum is 1 seconds, max is 60.
	if t.Retry.RetryTime < 1*time.Second {
		t.Retry.RetryTime = 1 * time.Second
	} else if t.Retry.RetryTime > time.Minute {
		t.Retry.RetryTime = time.Minute
	}

	delay := t.Retry.RetryTime + t.Retry.RetryTime

	// Total backoff can't be greater than 5 minutes.
	if delay > 5*time.Minute {
		delay = 5 * time.Minute
	}

	t.Retry.RetryTime = delay // update with new time.

	reschedule := time.NewTimer(t.Retry.RetryTime)
	go func() {
		<-reschedule.C
		if t.Retry.TotalRetries >= t.Retry.MaxRetries {
			// kill itself.
			t.IsKill = true
		}
		t.State = mesos_v1.TaskState_TASK_UNKNOWN
		revive <- t               // Revive itself.
		t.Retry.TotalRetries += 1 // Increment retry counter.
	}()

}

// Encode encodes the task for transport.
func (t *Task) Encode() ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, err
}

// Decode Decodes task data back into the task type.
func (t *Task) Decode(data []byte) (*Task, error) {
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
