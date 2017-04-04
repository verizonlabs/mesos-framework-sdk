package logging

import "testing"

func TestNewDefaultLogger(t *testing.T) {
	l := NewDefaultLogger()
	if _, ok := l.(Logger); !ok {
		t.Fatal("Default logger is of the wrong type")
	}

	// Test code path.
	l.Emit(EVENT, "TEST %s", "VALUE")
}
