package retry

import (
"mesos-framework-sdk/include/mesos_v1"
"mesos-framework-sdk/task"
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

