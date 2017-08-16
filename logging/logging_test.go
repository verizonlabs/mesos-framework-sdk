// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
