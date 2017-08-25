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

import (
	"regexp"
	"testing"
)

// Make sure we return a UUID with the proper length.
func TestUuid(t *testing.T) {
	t.Parallel()

	u := Uuid()
	if len(u) != 16 {
		t.Fatal("The length of the UUID is incorrect")
	}
}

// Measures performance of generating UUID's.
func BenchmarkUuid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Uuid()
	}
}

// Make sure we get a valid v4 UUID back as a string.
func TestUuidAsString(t *testing.T) {
	t.Parallel()

	u := UuidAsString()
	m, err := regexp.MatchString("^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", u)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !m {
		t.Fatal("Not a valid v4 UUID")
	}

}

// Measures performance of converting getting a UUID as a string.
func BenchmarkUuidAsString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UuidAsString()
	}
}
