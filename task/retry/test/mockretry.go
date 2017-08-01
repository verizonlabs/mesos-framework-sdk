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
