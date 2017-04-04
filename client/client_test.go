package client

import (
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

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, l)
	_, err := c.Request(nil)
	if err != nil {
		t.Fatal("Request could not be made successfully")
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
