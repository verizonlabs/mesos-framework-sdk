package main

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	//"mesos-framework-sdk/executor"
	mesos "mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/scheduler"
	"mesos-framework-sdk/scheduler/events"
	"mesos-framework-sdk/task_manager"
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
	// This channel handles events from the mesos master
	eventChan := make(chan *mesos_v1_scheduler.Event)

	// Channel listens for calls to make from the event handlers
	callChan := make(chan *mesos_v1_scheduler.Call)
	defer close(callChan)
	defer close(eventChan)

	manager := task_manager.NewDefaultTaskManager()

	// We can simply serve a file using the server here.
	//go server.NewServer("executor", ":8080", "/tmp/executor")

	// Create a http client for mesos, and create a new scheduler with the default handlers.
	c := client.NewClient("http://localhost:5050/api/v1/scheduler")

	// Create the new scheduler with both channels to listen for events and calls.
	// SchedulerEvents holds a reference to call chan so it can write to it.
	s := scheduler.NewDefaultScheduler(
		c,
		frameworkInfo,
		eventChan,
		callChan,
		manager,
		events.NewSchedulerEvents(manager, callChan),
	)
	// Run the scheduler.
	s.Run()
}
