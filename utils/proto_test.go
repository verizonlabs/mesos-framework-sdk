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

package utils

import "testing"

func TestProtoFloat64(t *testing.T) {
	t.Parallel()

	f := 1.0
	v := ProtoFloat64(f)
	if *v != f {
		t.Fatal("Values don't match")
	}
}

func BenchmarkProtoFloat64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ProtoFloat64(1.0)
	}
}

func TestProtoInt64(t *testing.T) {
	t.Parallel()

	var i int64 = 64
	v := ProtoInt64(i)
	if *v != i {
		t.Fatal("Values don't match")
	}
}

func BenchmarkProtoInt64(b *testing.B) {
	var i int64 = 64
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ProtoInt64(i)
	}
}

func TestProtoString(t *testing.T) {
	t.Parallel()

	s := "test"
	v := ProtoString(s)
	if *v != s {
		t.Fatal("Values don't match")
	}
}

func BenchmarkProtoString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ProtoString("test")
	}
}
