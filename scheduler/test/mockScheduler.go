package scheduler

import (
	"errors"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/include/mesos_v1_scheduler"
	"net/http"
)

type MockScheduler struct {
	events chan mesos_v1_scheduler.Event
}

func NewMockScheduler() *MockScheduler {
	return &MockScheduler{
		events: make(chan mesos_v1_scheduler.Event),
	}
}

func (m MockScheduler) FrameworkInfo() *mesos_v1.FrameworkInfo {
	return &mesos_v1.FrameworkInfo{}
}

func (m MockScheduler) Subscribe(chan *mesos_v1_scheduler.Event) (*http.Response, error) {

	return new(http.Response), nil
}

func (m MockScheduler) Teardown() (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Revive() (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Reconcile(tasks []*mesos_v1.TaskInfo) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) SchedRequest(resources []*mesos_v1.Request) (*http.Response, error) {
	return new(http.Response), nil
}

func (m MockScheduler) Suppress() (*http.Response, error) {
	return new(http.Response), nil
}

type MockBrokenScheduler struct{}

func (m MockBrokenScheduler) FrameworkInfo() *mesos_v1.FrameworkInfo {
	return nil
}

func (m MockBrokenScheduler) Subscribe(chan *mesos_v1_scheduler.Event) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Teardown() (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Accept(offerIds []*mesos_v1.OfferID, tasks []*mesos_v1.Offer_Operation, filters *mesos_v1.Filters) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Decline(offerIds []*mesos_v1.OfferID, filters *mesos_v1.Filters) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Revive() (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Kill(taskId *mesos_v1.TaskID, agentid *mesos_v1.AgentID) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Shutdown(execId *mesos_v1.ExecutorID, agentId *mesos_v1.AgentID) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Acknowledge(agentId *mesos_v1.AgentID, taskId *mesos_v1.TaskID, uuid []byte) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Reconcile(tasks []*mesos_v1.TaskInfo) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Message(agentId *mesos_v1.AgentID, executorId *mesos_v1.ExecutorID, data []byte) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) SchedRequest(resources []*mesos_v1.Request) (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}

func (m MockBrokenScheduler) Suppress() (*http.Response, error) {
	return new(http.Response), errors.New("Broken.")
}
