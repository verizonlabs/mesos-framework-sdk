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

// Measures performance of creating a new logger.
func BenchmarkNewDefaultLogger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewDefaultLogger()
	}
}

// Tests if there are any issues with emitting messages.
func TestDefaultLogger_Emit(t *testing.T) {
	t.Parallel()

	l := NewDefaultLogger()

	// Test code path.
	l.Emit(TEST, "TEST %s", "VALUE")
}

// Measures performance of emitting messages.
func BenchmarkDefaultLogger_Emit(b *testing.B) {
	l := NewDefaultLogger()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l.Emit(TEST, "TEST %s", "VALUE")
	}
}
