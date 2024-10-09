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

func uint8CompactTestCases() []compactTestCase[uint8] {
	return []compactTestCase[uint8]{
		{0, []byte{0b0000_0000}},
		{1, []byte{0b0000_0100}},
		{maxUint6, []byte{0b1111_1100}},
		{maxUint8, []byte{0b1111_1101, 0b0000_0011}},
	}
}

func uint16CompactTestCases() []compactTestCase[uint16] {
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

func uint32CompactTestCases() []compactTestCase[uint32] {
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

func uint64CompactTestCases() []compactTestCase[uint64] {
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

func intPtr[T ~uint8 | ~uint16 | ~uint32 | ~uint64](v T) *T {
	return &v
}

func uint8PtrTestCases() []compactTestCase[*uint8] {
	return []compactTestCase[*uint8]{
		{nil, []byte{0}},
		{intPtr[uint8](1), []byte{1, 0b0000_0100}},
		{intPtr[uint8](maxUint6), []byte{1, 0b1111_1100}},
		{intPtr[uint8](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
	}
}

func uint16PtrTestCases() []compactTestCase[*uint16] {
	return []compactTestCase[*uint16]{
		{nil, []byte{0}},
		{intPtr[uint16](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{intPtr[uint16](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
	}
}

func uint32PtrTestCases() []compactTestCase[*uint32] {
	return []compactTestCase[*uint32]{
		{nil, []byte{0}},
		{intPtr[uint32](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{intPtr[uint32](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{intPtr[uint32](maxUint30), []byte{1, 0b1111_1110, 0b1111_1111, 0b1111_1111, 0b1111_1111}},
		{intPtr[uint32](math.MaxUint32), []byte{
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
		{intPtr[uint64](maxUint8), []byte{1, 0b1111_1101, 0b0000_0011}},
		{intPtr[uint64](maxUint16), []byte{1, 0b1111_1110, 0b1111_1111, 0b0000_0011, 0b0000_0000}},
		{intPtr[uint64](math.MaxUint32), []byte{
			1,
			0b0000_0011,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
			0b1111_1111,
		}},
		{intPtr[uint64](math.MaxUint64), []byte{
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

func testEncodeCompactIntegers[T any](
	t *testing.T,
	name string,
	tcs []compactTestCase[T],
	encode func(enc *Encoder, value T) (int, error),
) {
	t.Run(name, func(t *testing.T) {
		for _, tc := range tcs {
			t.Run("", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				enc := NewEncoder(buf)
				_, err := encode(enc, tc.value)
				require.NoError(t, err)
				require.Equal(t, tc.expect, buf.Bytes())
			})
		}
	})
}

func TestEncodeCompactIntegers(t *testing.T) {
	testEncodeCompactIntegers(t, "uint8", uint8CompactTestCases(), EncodeCompact8)
	testEncodeCompactIntegers(t, "uint16", uint16CompactTestCases(), EncodeCompact16)
	testEncodeCompactIntegers(t, "uint32", uint32CompactTestCases(), EncodeCompact32)
	testEncodeCompactIntegers(t, "uint64", uint64CompactTestCases(), EncodeCompact64)
	testEncodeCompactIntegers(t, "*uint8", uint8PtrTestCases(), EncodeCompact8Ptr)
	testEncodeCompactIntegers(t, "*uint16", uint16PtrTestCases(), EncodeCompact16Ptr)
	testEncodeCompactIntegers(t, "*uint32", uint32PtrTestCases(), EncodeCompact32Ptr)
	testEncodeCompactIntegers(t, "*uint64", uint64PtrTestCases(), EncodeCompact64Ptr)
	testEncodeCompactIntegers(t, "[]uint16", uint16SliceTestCases(), EncodeUint16Slice)
	testEncodeCompactIntegers(t, "[]uint32", uint32SliceTestCases(), EncodeUint32Slice)
	testEncodeCompactIntegers(t, "[]uint64", uint64SliceTestCases(), EncodeUint64Slice)
}

func testEncodeNonCompact[T any](
	t *testing.T,
	name string,
	val T,
	encoded []byte,
	encode func(enc *Encoder, value T) (int, error),
) {
	t.Run(name, func(t *testing.T) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		_, err := encode(enc, val)
		require.NoError(t, err)
		require.Equal(t, encoded, buf.Bytes())
	})
}

func TestEncodeUint(t *testing.T) {
	testEncodeNonCompact(
		t, "uint8",
		uint8(0x42),
		[]byte{0x42},
		EncodeByte,
	)
	testEncodeNonCompact(
		t, "uint16",
		uint16(0x1234),
		[]byte{0x34, 0x12},
		EncodeUint16)
	testEncodeNonCompact(
		t, "uint32",
		uint32(0x12345678),
		[]byte{0x78, 0x56, 0x34, 0x12},
		EncodeUint32)
	testEncodeNonCompact(
		t, "uint64",
		uint64(0x123456789abcdef0),
		[]byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12},
		EncodeUint64,
	)
}
