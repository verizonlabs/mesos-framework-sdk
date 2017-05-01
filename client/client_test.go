package client

import (
	"mesos-framework-sdk/include/mesos_v1_executor"
	"mesos-framework-sdk/include/mesos_v1"
	"mesos-framework-sdk/include/mesos_v1_scheduler"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLogger struct{}

func (m *mockLogger) Emit(severity uint8, template string, args ...interface{}) {

}

var l = new(mockLogger)

// Checks our newly created client to make sure it's in the right state.
func TestNewClient(t *testing.T) {
	t.Parallel()

	c := NewClient("test", l)
	if c.StreamID() != "" {
		t.Fatal("Stream ID should be empty")
	}
}

// Measures performance of creating a new client.
func BenchmarkNewClient(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewClient("test", l)
	}
}

// Ensure that our stream ID gets set correctly.
func TestDefaultClient_SetStreamID(t *testing.T) {
	t.Parallel()

	id := "id"
	c := NewClient("test", l)
	c.SetStreamID(id)
	if c.StreamID() != id {
		t.Fatal("Stream ID was not set correctly")
	}
}

// Measures performance of setting our stream ID.
func BenchmarkDefaultClient_StreamID(b *testing.B) {
	c := NewClient("test", l)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.SetStreamID("id")
	}
}

// Tests if we can make requests successfully or not.
func TestDefaultClient_Request(t *testing.T) {
	t.Parallel()

	val := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Mesos-Stream-Id", val)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, l)

	_, err := c.Request(nil)
	if err != nil {
		t.Fatal("Generic request could not be made successfully:" + err.Error())
	}

	if c.StreamID() == "" {
		t.Fatal("Mesos-Stream-Id header should have been set but it wasn't")
	}

	_, err = c.Request(&mesos_v1_scheduler.Call{})
	if err != nil {
		t.Fatal("Scheduler request could not be made successfully: " + err.Error())
	}

	_, err = c.Request(&mesos_v1_executor.Call{
		ExecutorId: &mesos_v1.ExecutorID{
			Value: &val,
		},
		FrameworkId: &mesos_v1.FrameworkID{
			Value: &val,
		},
	})

	if err != nil {
		t.Fatal("Executor request could not be made successfully: " + err.Error())
	}

	c = NewClient("", l)
	_, err = c.Request(nil)

	if err == nil {
		t.Fatal("Generic request should have failed but it didn't: " + err.Error())
	}

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts2.Close()

	c = NewClient(ts2.URL, l)
	_, err = c.Request(&mesos_v1_scheduler.Call{})

	if err == nil {
		t.Fatal("Response should have thrown a 400: " + err.Error())
	}

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPermanentRedirect)
	}))
	defer ts3.Close()

	c = NewClient(ts3.URL, l)
	_, err = c.Request(&mesos_v1_scheduler.Call{})

	if err == nil {
		t.Fatal("Redirect should have been encountered but it wasn't: " + err.Error())
	}
}

// Measures performance of creating and sending HTTP requests.
func BenchmarkDefaultClient_Request(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, l)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Request(nil)
	}
}
