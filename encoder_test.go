package scale

import (
	"bytes"
	"testing"
)

type customBuffer struct {
	buf bytes.Buffer
}

func (c *customBuffer) Write(b []byte) (int, error) {
	return c.buf.Write(b)
}

func (c *customBuffer) Len() int {
	return c.buf.Len()
}

func BenchmarkEncodeStrings_WithStringWriter(b *testing.B) {
	// bytes.Buffer implements the io.StringWriter interface.
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
	b.StopTimer()

	b.Log(buf.Len())
}

func BenchmarkEncodeStrings_WithWriterForStrings(b *testing.B) {
	// CustomBuffer does not implement the io.StringWriter interface.
	var buf customBuffer
	enc := NewEncoder(&buf)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
	b.StopTimer()

	b.Log(buf.Len())
}
