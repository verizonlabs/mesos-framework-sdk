package main

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/scheduler"
	eventCtrl "mesos-framework-sdk/scheduler/events/default_event_controller"
	"mesos-framework-sdk/task_manager"
	"time"
)

func main() {
	frameworkInfo := &mesos_v1.FrameworkInfo{
		User:            proto.String("root"),
		Name:            proto.String("Sprint"),
		FailoverTimeout: proto.Float64(1 * time.Second.Seconds()),
		Checkpoint:      proto.Bool(true),
		Role:            proto.String("*"),
		Hostname:        proto.String(""),
		Principal:       proto.String(""),
	}
	// This channel handles events from the mesos master
	eventChan := make(chan *mesos_v1_scheduler.Event)
	defer close(eventChan)

	manager := task_manager.NewDefaultTaskManager()

	// Create a http client for mesos, and create a new scheduler with the default handlers.
	c := client.NewClient("http://localhost:5050/api/v1/scheduler")
	s := scheduler.NewDefaultScheduler(c, frameworkInfo)
	e := eventCtrl.NewDefaultEventController(s, manager, eventChan)

	// We can simply serve a file using the server here.
	//go server.NewServer("executor", ":8080", "/tmp/executor")

	// Run the scheduler.
	e.Run()
}
