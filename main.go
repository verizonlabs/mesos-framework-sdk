package main

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	mesos "mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/scheduler"
	"mesos-framework-sdk/scheduler/events"
	"time"
)

func main() {
	frameworkInfo := &mesos.FrameworkInfo{
		User:            proto.String("root"),
		Name:            proto.String("Sprint"),
		FailoverTimeout: proto.Float64(1 * time.Second.Seconds()),
		Checkpoint:      proto.Bool(true),
		Role:            proto.String("*"),
		Hostname:        proto.String(""),
		Principal:       proto.String(""),
	}

	c := client.NewClient("http://localhost:5050/api/v1/scheduler")
	s := scheduler.NewScheduler(c, frameworkInfo, events.NewSchedulerEvents())
	s.Run()
	//s.Subscribe(frameworkInfo)
}
