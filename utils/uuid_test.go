package utils

import "testing"

// Make sure we return a valid v4 UUID with the proper length.
func TestUuid(t *testing.T) {
	u := Uuid()
	if len(u) != 16 {
		t.Fatal("The length of the UUID is incorrect")
	}
}
