package scale

import (
	"testing"
)

func Benchmark_UnsafeBytesToString(b *testing.B) {
	in := []byte("Hello World")
	var out string

	for i := 0; i < b.N; i++ {
		out = bytesToString(in)
	}
	_ = out
}

func Benchmark_SafeBytesToString(b *testing.B) {
	in := []byte("Hello World")
	var out string

	for i := 0; i < b.N; i++ {
		out = string(in)
	}
	_ = out
}

func Benchmark_UnsafeStringToBytes(b *testing.B) {
	in := "Hello World"
	var out []byte

	for i := 0; i < b.N; i++ {
		out = stringToBytes(in)
	}
	_ = out
}

func Benchmark_SafeStringToBytes(b *testing.B) {
	in := "Hello World"
	var out []byte
	for i := 0; i < b.N; i++ {
		out = []byte(in)
	}
	b.StopTimer()
	_ = out
}
