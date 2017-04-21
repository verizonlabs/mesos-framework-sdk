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
