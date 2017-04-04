package scheduler

import (
	"mesos-framework-sdk/client"
	"mesos-framework-sdk/include/mesos"
	"mesos-framework-sdk/include/scheduler"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockClient struct{}

func (m *mockClient) Request(interface{}) (*http.Response, error) {
	return new(http.Response), nil
}

func (m *mockClient) StreamID() string {
	return "test"
}

func (m *mockClient) SetStreamID(string) client.Client {
	return m
}

type mockLogger struct{}

func (m *mockLogger) Emit(severity uint8, template string, args ...interface{}) {

}

var c = new(mockClient)
var i = &mesos_v1.FrameworkInfo{}
var l = new(mockLogger)

// Checks the internal state of a new scheduler.
func TestNewDefaultScheduler(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	if s.Client != c || s.logger != l || s.Info != i {
		t.Fatal("Scheduler does not have the right internal state")
	}
}

// Measures performance of creating a new scheduler.
func BenchmarkNewDefaultScheduler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewDefaultScheduler(c, i, l)
	}
}

// Tests our accept call to Mesos.
func TestDefaultScheduler_Accept(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	tasks := []*mesos_v1.Offer_Operation{}
	filters := &mesos_v1.Filters{}

	_, err := s.Accept(offerIds, tasks, filters)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our accept call to Mesos.
func BenchmarkDefaultScheduler_Accept(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	tasks := []*mesos_v1.Offer_Operation{}
	filters := &mesos_v1.Filters{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Accept(offerIds, tasks, filters)
	}
}

// Tests our acknowledge call to Mesos.
func TestDefaultScheduler_Acknowledge(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	agentId := &mesos_v1.AgentID{}
	taskId := &mesos_v1.TaskID{}
	uuid := []byte{}

	_, err := s.Acknowledge(agentId, taskId, uuid)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our acknowledge call to Mesos.
func BenchmarkDefaultScheduler_Acknowledge(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	agentId := &mesos_v1.AgentID{}
	taskId := &mesos_v1.TaskID{}
	uuid := []byte{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Acknowledge(agentId, taskId, uuid)
	}
}

// Tests our decline call to Mesos.
func TestDefaultScheduler_Decline(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	filters := &mesos_v1.Filters{}

	_, err := s.Decline(offerIds, filters)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our decline call to Mesos.
func BenchmarkDefaultScheduler_Decline(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	offerIds := []*mesos_v1.OfferID{}
	filters := &mesos_v1.Filters{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Decline(offerIds, filters)
	}
}

// Tests our kill call to Mesos.
func TestDefaultScheduler_Kill(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	taskId := &mesos_v1.TaskID{}
	agentId := &mesos_v1.AgentID{}

	_, err := s.Kill(taskId, agentId)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our kill call to Mesos.
func BenchmarkDefaultScheduler_Kill(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	taskId := &mesos_v1.TaskID{}
	agentId := &mesos_v1.AgentID{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Kill(taskId, agentId)
	}
}

// Tests our message call to Mesos.
func TestDefaultScheduler_Message(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	agentId := &mesos_v1.AgentID{}
	execId := &mesos_v1.ExecutorID{}
	data := []byte{}

	_, err := s.Message(agentId, execId, data)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our message call to Mesos.
func BenchmarkDefaultScheduler_Message(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	agentId := &mesos_v1.AgentID{}
	execId := &mesos_v1.ExecutorID{}
	data := []byte{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Message(agentId, execId, data)
	}
}

// Tests our reconcile call to Mesos.
func TestDefaultScheduler_Reconcile(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	tasks := []*mesos_v1.TaskInfo{}

	_, err := s.Reconcile(tasks)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our reconcile call to Mesos.
func BenchmarkDefaultScheduler_Reconcile(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	tasks := []*mesos_v1.TaskInfo{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Reconcile(tasks)
	}
}

// Tests our revive call to Mesos.
func TestDefaultScheduler_Revive(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)

	_, err := s.Revive()
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our revive call to Mesos.
func BenchmarkDefaultScheduler_Revive(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Revive()
	}
}

// Tests our request call to Mesos.
func TestDefaultScheduler_SchedRequest(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	res := []*mesos_v1.Request{}

	_, err := s.SchedRequest(res)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our request call to Mesos.
func BenchmarkDefaultScheduler_SchedRequest(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	res := []*mesos_v1.Request{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.SchedRequest(res)
	}
}

// Tests our shutdown call to Mesos.
func TestDefaultScheduler_Shutdown(t *testing.T) {
	t.Parallel()

	s := NewDefaultScheduler(c, i, l)
	execId := &mesos_v1.ExecutorID{}
	agentId := &mesos_v1.AgentID{}

	_, err := s.Shutdown(execId, agentId)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// Measures performance of our shutdown call to Mesos.
func BenchmarkDefaultScheduler_Shutdown(b *testing.B) {
	s := NewDefaultScheduler(c, i, l)
	execId := &mesos_v1.ExecutorID{}
	agentId := &mesos_v1.AgentID{}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.Shutdown(execId, agentId)
	}
}

// Tests our subscribe call to Mesos.
func TestDefaultScheduler_Subscribe(t *testing.T) {
	t.Parallel()

	ch := make(chan *mesos_v1_scheduler.Event)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c := client.NewClient(srv.URL, l)
	s := NewDefaultScheduler(c, i, l)
	c.Request(nil)

	_, err := s.Subscribe(ch)

	// We SHOULD get an error in this case; make sure that's true.
	if err == nil {
		t.Fatal(err.Error())
	}
}
