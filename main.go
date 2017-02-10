package main

import (
	//exec "mesos-framework-sdk/include/executor"
	sched "mesos-framework-sdk/include/scheduler"
	mesos "mesos-framework-sdk/include/mesos"
	"time"
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
)

func main() {
	frameworkInfo := mesos.FrameworkInfo{
		User: proto.String("root"),
		Name: proto.String("Sprint"),
		Id: &mesos.FrameworkID{Value: proto.String("")},
		FailoverTimeout: proto.Float64(1 * time.Second.Seconds()),
		Checkpoint: proto.Bool(true),
		Role: proto.String("*"),
		Hostname: proto.String(""),
		Principal: proto.String(""),
	}

	client.Subscribe_Call(&sched.Call{
		FrameworkId: frameworkInfo.Id,
		Type: sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: &frameworkInfo,
		},
	})


}