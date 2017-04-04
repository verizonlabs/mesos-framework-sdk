package logging

import "testing"

// Ensure we can use our default logger.
func TestNewDefaultLogger(t *testing.T) {
	t.Parallel()

	l := NewDefaultLogger()
	if _, ok := l.(Logger); !ok {
		t.Fatal("Default logger is of the wrong type")
	}
}

// Tests if there are any issues with emitting messages.
func TestDefaultLogger_Emit(t *testing.T) {
	t.Parallel()

	l := NewDefaultLogger()
	// Test code path.
	l.Emit(EVENT, "TEST %s", "VALUE")
}
