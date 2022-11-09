package scale

import (
	"testing"
)

func Benchmark_UnsafeBytesToString(b *testing.B) {
	in := []byte("Hello World")
	var out string

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out = bytesToString(in)
		// assert.Equal(b, "Hello World", out)
	}
	b.StopTimer()
	b.Log(out)
}

func Benchmark_SafeBytesToString(b *testing.B) {
	in := []byte("Hello World")
	var out string

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out = string(in)
		// assert.Equal(b, "Hello World", out)
	}
	b.StopTimer()
	b.Log(out)
}

func Benchmark_UnsafeStringToBytes(b *testing.B) {
	in := "Hello World"
	var out []byte

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out = stringToBytes(in)
		// assert.Equal(b, []byte("Hello World"), out)
	}
	b.StopTimer()
	b.Log(string(out))
}

func Benchmark_SafeStringToBytes(b *testing.B) {
	in := "Hello World"
	var out []byte

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out = []byte(in)
		// assert.Equal(b, []byte("Hello World"), out)
	}
	b.StopTimer()
	b.Log(string(out))
}
