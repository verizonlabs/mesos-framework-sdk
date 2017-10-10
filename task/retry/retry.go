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

package retry

import (
	"github.com/verizonlabs/mesos-framework-sdk/include/mesos_v1"
	"github.com/verizonlabs/mesos-framework-sdk/task"
	"time"
)

type (

	// Provides pluggable retry mechanisms.
	// Also used extensively for testing with mocks.
	Retry interface {
		AddPolicy(policy *task.TimeRetry, mesosTask *mesos_v1.TaskInfo) error
		CheckPolicy(mesosTask *mesos_v1.TaskInfo) *TaskRetry
		ClearPolicy(mesosTask *mesos_v1.TaskInfo) error
		RunPolicy(policy *TaskRetry, f func() error) error
	}

	// Primary retry mechanism used with policies in the task manager and persistence engine.
	TaskRetry struct {
		TotalRetries int
		MaxRetries   int
		RetryTime    time.Duration
		Backoff      bool
		Name         string
	}
)
