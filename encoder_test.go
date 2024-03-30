package scale

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
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
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
}

func BenchmarkEncodeStrings_WithWriterForStrings(b *testing.B) {
	// CustomBuffer does not implement the io.StringWriter interface.
	var buf customBuffer
	enc := NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		EncodeString(enc, "Hello World")
	}
}

type compactTestCase[T any] struct {
	value  T
	expect []byte
}

func uint8TestCases() []compactTestCase[uint8] {
	return []compactTestCase[uint8]{
		{0, []byte{0b0000_0000}},
		{1, []byte{0b0000_0100}},
		{maxUint6, []byte{0b1111_1100}},
		{maxUint8, []byte{0b1111_1101, 0b0000_0011}},
	}
}

func uint16TestCases() []compactTestCase[uint16] {
	return []compactTestCase[uint16]{
		{0, []byte{0b0000_0000}},
		{1, []byte{0b0000_0100}},
		{maxUint6, []byte{0b1111_1100}},
		{maxUint8, []byte{0b1111_1101, 0b0000_0011}},
		{maxUint14, []byte{0b1111_1101, 0b1111_1111}},
		{maxUint14 + 1, []byte{0b0000_0010, 0b0000_0000, 0b0000_0001, 0b0000_0000}},
		{maxUint16, []byte{0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
	}
}

func uint32TestCases() []compactTestCase[uint32] {
	return []compactTestCase[uint32]{
		{0, []byte{0b0000_0000}},
		{1, []byte{0b0000_0100}},
		{maxUint6, []byte{0b1111_1100}},
		{maxUint8, []byte{0b1111_1101, 0b0000_0011}},
		{maxUint14, []byte{0b1111_1101, 0b1111_1111}},
		{maxUint14 + 1, []byte{0b0000_0010, 0b0000_0000, 0b0000_0001, 0b0000_0000}},
		{maxUint16, []byte{0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{maxUint30, []byte{0b1111_1110, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
		{maxUint30 + 1, []byte{0b0000_0011, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0100_0000}},
		{math.MaxUint32, []byte{0b0000_0011, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
	}
}

func uint64TestCases() []compactTestCase[uint64] {
	return []compactTestCase[uint64]{
		{0, []byte{0b0000_0000}},
		{1, []byte{0b0000_0100}},
		{maxUint6, []byte{0b1111_1100}},
		{maxUint8, []byte{0b1111_1101, 0b0000_0011}},
		{maxUint14, []byte{0b1111_1101, 0b1111_1111}},
		{maxUint14 + 1, []byte{0b0000_0010, 0b0000_0000, 0b0000_0001, 0b0000_0000}},
		{maxUint16, []byte{0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{maxUint30, []byte{0b1111_1110, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
		{maxUint30 + 1, []byte{0b0000_0011, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0100_0000}},
		{math.MaxUint32, []byte{0b0000_0011, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
		{math.MaxUint32 + 1, []byte{0b0000_0111, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0001}},
		{1 << 40, []byte{0b0000_1011, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0001}},
		{1 << 48, []byte{
			0b0000_1111,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0001,
		}},
		{1 << 56, []byte{
			0b0001_0011,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0000,
			0b0000_0001,
		}},
		{math.MaxUint64, []byte{
			0b0001_0011,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
		}},
	}
}

func newInteger[T ~uint8 | ~uint16 | ~uint32 | ~uint64](v T) *T {
	return &v
}

func uint8PtrTestCases() []compactTestCase[*uint8] {
	return []compactTestCase[*uint8]{
		{nil, []byte{0}},
		{newInteger[uint8](1), []byte{1, 1}},
		{newInteger[uint8](math.MaxUint8), []byte{1, math.MaxUint8}},
	}
}

func uint16PtrTestCases() []compactTestCase[*uint16] {
	return []compactTestCase[*uint16]{
		{nil, []byte{0}},
		{newInteger[uint16](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{newInteger[uint16](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
	}
}

func uint32PtrTestCases() []compactTestCase[*uint32] {
	return []compactTestCase[*uint32]{
		{nil, []byte{0}},
		{newInteger[uint32](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{newInteger[uint32](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{newInteger[uint32](maxUint30), []byte{1, 0b1111_1110, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
		{newInteger[uint32](math.MaxUint32), []byte{
			1,
			0b0000_0011,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
		}},
	}
}

func uint64PtrTestCases() []compactTestCase[*uint64] {
	return []compactTestCase[*uint64]{
		{nil, []byte{0}},
		{newInteger[uint64](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{newInteger[uint64](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{newInteger[uint64](math.MaxUint32), []byte{
			1,
			0b0000_0011,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
		}},
		{newInteger[uint64](math.MaxUint64), []byte{
			1,
			0b0001_0011,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
		}},
	}
}

func mustDecodeHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return b
}

func uint16SliceTestCases() []compactTestCase[[]uint16] {
	return []compactTestCase[[]uint16]{
		{[]uint16{4, 15, 23, math.MaxUint16}, mustDecodeHex("10103c5cfeff0300")},
	}
}

func uint32SliceTestCases() []compactTestCase[[]uint32] {
	return []compactTestCase[[]uint32]{
		{[]uint32{4, 15, 23, math.MaxUint32}, mustDecodeHex("10103c5c03ffffffff")},
	}
}

func uint64SliceTestCases() []compactTestCase[[]uint64] {
	return []compactTestCase[[]uint64]{
		{[]uint64{4, 15, 23, math.MaxUint64}, mustDecodeHex("10103c5c13ffffffffffffffff")},
	}
}

func encodeTest[T any](t *testing.T, value T, expect []byte) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	var err error
	switch typed := any(value).(type) {
	case uint8:
		_, err = EncodeCompact8(enc, typed)
	case uint16:
		_, err = EncodeCompact16(enc, typed)
	case uint32:
		_, err = EncodeCompact32(enc, typed)
	case uint64:
		_, err = EncodeCompact64(enc, typed)
	case *uint8:
		_, err = EncodeBytePtr(enc, typed)
	case *uint16:
		_, err = EncodeCompact16Ptr(enc, typed)
	case *uint32:
		_, err = EncodeCompact32Ptr(enc, typed)
	case *uint64:
		_, err = EncodeCompact64Ptr(enc, typed)
	case []uint16:
		_, err = EncodeUint16Slice(enc, typed)
	case []uint32:
		_, err = EncodeUint32Slice(enc, typed)
	case []uint64:
		_, err = EncodeUint64Slice(enc, typed)
	default:
		t.Fatal("unsupported type")
	}
	require.NoError(t, err)
	require.Equal(t, expect, buf.Bytes())
}

func TestEncodeCompactIntegers(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		for _, tc := range uint8TestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("uint16", func(t *testing.T) {
		for _, tc := range uint16TestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("uint32", func(t *testing.T) {
		for _, tc := range uint32TestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("uint64", func(t *testing.T) {
		for _, tc := range uint64TestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("*uint8", func(t *testing.T) {
		for _, tc := range uint8PtrTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("*uint16", func(t *testing.T) {
		for _, tc := range uint16PtrTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("*uint32", func(t *testing.T) {
		for _, tc := range uint32PtrTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("*uint64", func(t *testing.T) {
		for _, tc := range uint64PtrTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("[]uint16", func(t *testing.T) {
		for _, tc := range uint16SliceTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("[]uint32", func(t *testing.T) {
		for _, tc := range uint32SliceTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
	t.Run("[]uint64", func(t *testing.T) {
		for _, tc := range uint64SliceTestCases() {
			t.Run("", func(t *testing.T) {
				encodeTest(t, tc.value, tc.expect)
			})
		}
	})
}
