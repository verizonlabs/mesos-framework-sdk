package executor

import (
	exec "mesos-framework-sdk/include/executor"
	"mesos-framework-sdk/include/mesos"
)

type Executor interface {
	Subscribe(chan *exec.Event) error
	Update(*mesos_v1.TaskStatus)
	Message([]byte)
}
