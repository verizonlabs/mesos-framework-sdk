package scheduler

import (
	"github.com/golang/protobuf/proto"
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"testing"
	"time"
)

const (
	clientUrl = "http://localhost:5050/api/v1/scheduler" // TODO remove this once client is mocked
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
	c := client.NewClient(clientUrl) // TODO mock this so it doesn't make a real HTTP call
	s := NewDefaultScheduler(c, frameworkInfo)
	ch := make(chan *sched.Event)

	if err := s.Subscribe(ch); err != nil {
		t.FailNow()
	}
}
