package scheduler

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/scheduler/events"
	"testing"
	"time"
)

const (
	clientUrl = "http://localhost:5050/api/v1/scheduler"
)

var (
	frameworkInfo = &mesos_v1.FrameworkInfo{
		User:            proto.String("root"),
		Name:            proto.String("Test"),
		FailoverTimeout: proto.Float64(1 * time.Second.Seconds()),
		Checkpoint:      proto.Bool(true),
		Role:            proto.String("*"),
		Hostname:        proto.String(""),
		Principal:       proto.String(""),
	}
)

// Tests if the scheduler can be created.
func TestNewScheduler(t *testing.T) {
	mockClient := client.NewClient(clientUrl)
	mockEventHandler := events.NewDefaultEventController()
	mockScheduler := NewScheduler(mockClient, frameworkInfo, mockEventHandler)
	if err := mockScheduler.Subscribe(); err != nil {
		t.FailNow()
	}
}

//
func TestScheduler_Kill(t *testing.T) {
	mockClient := client.NewClient(clientUrl)
	mockEventHandler := events.NewDefaultEventController()
	mockScheduler := NewScheduler(mockClient, frameworkInfo, mockEventHandler)
	if err := mockScheduler.Subscribe(); err != nil {
		t.FailNow()
	}
}
