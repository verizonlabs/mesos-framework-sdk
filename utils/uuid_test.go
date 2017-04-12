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
