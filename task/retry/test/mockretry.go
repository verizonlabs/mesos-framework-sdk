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

package test

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/task"
	"mesos-framework-sdk/task/retry"
)

type MockRetry struct{}

func (r MockRetry) AddPolicy(policy *task.TimeRetry, mesosTask *mesos_v1.TaskInfo) error {
	return nil
}

func (r MockRetry) CheckPolicy(mesosTask *mesos_v1.TaskInfo) *retry.TaskRetry {
	return nil
}

func (r MockRetry) ClearPolicy(mesosTask *mesos_v1.TaskInfo) error {
	return nil
}

func (r MockRetry) RunPolicy(policy *retry.TaskRetry, f func() error) error {
	return f()
}

type MockBrokenRetry struct{}

func (r MockBrokenRetry) AddPolicy(policy *task.TimeRetry, mesosTask *mesos_v1.TaskInfo) error {
	return errors.New("Broken")
}

func (r MockBrokenRetry) CheckPolicy(mesosTask *mesos_v1.TaskInfo) *retry.TaskRetry {
	return nil
}

func (r MockBrokenRetry) ClearPolicy(mesosTask *mesos_v1.TaskInfo) error {
	return errors.New("Broken")
}

func (r MockBrokenRetry) RunPolicy(policy *retry.TaskRetry, f func() error) error {
	return errors.New("Broken")
}
