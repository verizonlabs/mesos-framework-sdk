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
	"crypto/rand"
	"fmt"
)

// Generates a UUID using random bytes from a secure source.
func Uuid() []byte {
	b := make([]byte, 0, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to generate UUID")
	}

	// Generate a v4 UUID.
	b[6] = (b[6] | 0x40) & 0x4F
	b[8] = (b[8] | 0x80) & 0xBF

	return b
}

func UuidAsString() string {
	uuid := Uuid()
	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
