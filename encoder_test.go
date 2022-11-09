package scale

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkEncodeStrings_WithStringWriter(b *testing.B) {

	// strings.Builder implements the io.StringWriter interface.
	var buf strings.Builder
	enc := NewEncoder(&buf)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
	b.StopTimer()

	_ = buf.Len() // avoid optimizations by compiler
}

func BenchmarkEncodeStrings_WithWriterForStrings(b *testing.B) {

	// bytes.Buffer does not implement the io.StringWriter interface.
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
	b.StopTimer()

	_ = buf.Len() // avoid optimizations by compiler
}
