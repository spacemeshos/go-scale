package scale

import (
	"testing"
)

func Benchmark_UnsafeBytesToString(b *testing.B) {
	in := []byte("Hello World")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out := BytesToString(in)
		// assert.Equal(b, "Hello World", out)
		_ = out
	}
	b.StopTimer()
}

func Benchmark_SafeBytesToString(b *testing.B) {
	in := []byte("Hello World")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out := string(in)
		// assert.Equal(b, "Hello World", out)
		_ = out
	}
	b.StopTimer()
}

func Benchmark_UnsafeStringToBytes(b *testing.B) {
	in := "Hello World"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out := StringToBytes(in)
		// assert.Equal(b, []byte("Hello World"), out)
		_ = out
	}
	b.StopTimer()
}

func Benchmark_SafeStringToBytes(b *testing.B) {
	in := "Hello World"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		out := []byte(in)
		// assert.Equal(b, []byte("Hello World"), out)
		_ = out
	}
	b.StopTimer()
}
