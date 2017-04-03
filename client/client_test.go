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
