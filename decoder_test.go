package scale

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type halfReader struct {
	rem []byte
}

func (t *halfReader) Read(buf []byte) (int, error) {
	n := copy(buf, t.rem[:1+len(t.rem)/2])
	t.rem = t.rem[n:]
	return n, nil
}

func testReadFull[T any](
	t *testing.T,
	name string,
	val T,
	encode func(e *Encoder, v T) (int, error),
	decode func(d *Decoder) (T, int, error),
) {
	t.Run(name, func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := NewEncoder(buf)
			_, err := encode(enc, val)
			require.NoError(t, err)

			dec := NewDecoder(buf)
			rst, _, err := decode(dec)
			require.NoError(t, err)
			require.Equal(t, val, rst)
		})
		t.Run("partial", func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := NewEncoder(buf)
			_, err := encode(enc, val)
			require.NoError(t, err)

			hr := &halfReader{rem: buf.Bytes()}
			dec := NewDecoder(hr)
			rst, _, err := decode(dec)
			require.NoError(t, err)
			require.Equal(t, val, rst)
		})
	})
}

func TestReadFull(t *testing.T) {
	testReadFull(
		t, "uint8", uint8(math.MaxUint8),
		EncodeCompact8, DecodeCompact8)
	testReadFull(
		t, "uint16", uint16(math.MaxUint16),
		EncodeCompact16, DecodeCompact16)
	testReadFull(t,
		"uint32", uint32(math.MaxUint32),
		EncodeCompact32, DecodeCompact32)
	testReadFull(
		t, "uint64", uint64(math.MaxUint64),
		EncodeCompact64, DecodeCompact64)
	testReadFull(
		t, "*uint8", intPtr[uint8](math.MaxUint8),
		EncodeCompact8Ptr, DecodeCompact8Ptr)
	testReadFull(
		t, "nil *uint8", (*uint8)(nil),
		EncodeCompact8Ptr, DecodeCompact8Ptr)
	testReadFull(
		t, "*uint16", intPtr[uint16](math.MaxUint16),
		EncodeCompact16Ptr, DecodeCompact16Ptr)
	testReadFull(
		t, "nil *uint16", (*uint16)(nil),
		EncodeCompact16Ptr, DecodeCompact16Ptr)
	testReadFull(
		t, "*uint32", intPtr[uint32](math.MaxUint32),
		EncodeCompact32Ptr, DecodeCompact32Ptr)
	testReadFull(
		t, "nil *uint32", (*uint32)(nil),
		EncodeCompact32Ptr, DecodeCompact32Ptr)
	testReadFull(
		t, "*uint64", intPtr[uint64](math.MaxUint64),
		EncodeCompact64Ptr, DecodeCompact64Ptr)
	testReadFull(
		t, "nil *uint64", (*uint64)(nil),
		EncodeCompact64Ptr, DecodeCompact64Ptr)
	testReadFull(
		t, "byte slice", []byte("dsa1232131312dsada123312"),
		EncodeByteSlice, DecodeByteSlice)
	testReadFull(
		t, "string", "dsa1232131312dsada123312",
		EncodeString, DecodeString)
	testReadFull(
		t, "string slice", []string{"qwe123", "dsa456"},
		EncodeStringSlice, DecodeStringSlice)
	testReadFull(
		t, "uint16 slice", []uint16{0, 1, 2, math.MaxUint8, math.MaxUint16},
		EncodeUint16Slice, DecodeUint16Slice)
	testReadFull(
		t, "uint32 slice", []uint32{0, 1, 2, math.MaxUint8, math.MaxUint16, math.MaxUint32},
		EncodeUint32Slice, DecodeUint32Slice)
	testReadFull(
		t, "uint64 slice", []uint64{0, 1, 2, math.MaxUint32, math.MaxUint64},
		EncodeUint64Slice, DecodeUint64Slice)
}

func testDecodeCompactIntegers[T any](
	t *testing.T,
	name string,
	tcs []encTestCase[T],
	decode func(d *Decoder) (T, int, error),
) {
	t.Run(name, func(t *testing.T) {
		for _, tc := range tcs {
			t.Run("", func(t *testing.T) {
				buf := bytes.NewBuffer(tc.expect)
				dec := NewDecoder(buf)
				rst, _, err := decode(dec)
				require.NoError(t, err)
				require.Equal(t, tc.value, rst)
			})
		}
	})
}

func TestDecodeCompactIntegers(t *testing.T) {
	testDecodeCompactIntegers(t, "uint8", uint8CompactTestCases(), DecodeCompact8)
	testDecodeCompactIntegers(t, "uint16", uint16CompactTestCases(), DecodeCompact16)
	testDecodeCompactIntegers(t, "uint32", uint32CompactTestCases(), DecodeCompact32)
	testDecodeCompactIntegers(t, "uint64", uint64CompactTestCases(), DecodeCompact64)
	testDecodeCompactIntegers(t, "*uint8", uint8PtrTestCases(), DecodeCompact8Ptr)
	testDecodeCompactIntegers(t, "*uint16", uint16PtrTestCases(), DecodeCompact16Ptr)
	testDecodeCompactIntegers(t, "*uint32", uint32PtrTestCases(), DecodeCompact32Ptr)
	testDecodeCompactIntegers(t, "*uint64", uint64PtrTestCases(), DecodeCompact64Ptr)
	testDecodeCompactIntegers(t, "[]uint16", uint16SliceTestCases(), DecodeUint16Slice)
	testDecodeCompactIntegers(t, "[]uint32", uint32SliceTestCases(), DecodeUint32Slice)
	testDecodeCompactIntegers(t, "[]uint64", uint64SliceTestCases(), DecodeUint64Slice)
}

type boundTestCase struct {
	value []byte
}

func boundUint8Cases() []boundTestCase {
	return []boundTestCase{
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b0000_0010},
		},
		{
			value: []byte{0b0000_0001},
		},
		{
			value: []byte{0b1111_1101, 0b0000_0111},
		},
	}
}

func boundUint16Cases() []boundTestCase {
	return []boundTestCase{
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b0000_0010},
		},
		{
			value: []byte{0b0000_0001},
		},
		{
			value: []byte{0b0000_0001, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b1111_1110, 0b1111_1111, 0b0000_0000, 0b0000_0000},
		},
		{
			value: []byte{0b1111_1110, 0b1111_1111, 0b0000_0100, 0b0000_0000},
		},
	}
}

func boundUint32Cases() []boundTestCase {
	return []boundTestCase{
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b0000_0010},
		},
		{
			value: []byte{0b0000_0001},
		},
		{
			value: []byte{0b0000_0001, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b1111_1110, 0b1111_1111, 0b0000_0000, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0011, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111},
		},
	}
}

func boundUint64Cases() []boundTestCase {
	return []boundTestCase{
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b0000_0010},
		},
		{
			value: []byte{0b0000_0001},
		},
		{
			value: []byte{0b0000_0001, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0011},
		},
		{
			value: []byte{0b1111_1110, 0b1111_1111, 0b0000_0000, 0b0000_0000},
		},
		{
			value: []byte{0b0000_0011, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000, 0b0000_0000},
		},
		{
			value: []byte{0b0010_0011, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111, 0b1111_1111},
		},
	}
}

func testCompactIntegerBoundaries[T any](
	t *testing.T,
	name string,
	tcs []boundTestCase,
	decode func(d *Decoder) (T, int, error),
) {
	t.Run(name, func(t *testing.T) {
		for _, tc := range tcs {
			t.Run("", func(t *testing.T) {
				buf := bytes.NewBuffer(tc.value)
				dec := NewDecoder(buf)
				_, _, err := decode(dec)
				require.Error(t, err, fmt.Sprintf("%b", tc.value))
			})
		}
	})
}

func TestCompactIntegersBoundaries(t *testing.T) {
	testCompactIntegerBoundaries(t, "uint8", boundUint8Cases(), DecodeCompact8)
	testCompactIntegerBoundaries(t, "uint16", boundUint16Cases(), DecodeCompact16)
	testCompactIntegerBoundaries(t, "uint32", boundUint32Cases(), DecodeCompact32)
	testCompactIntegerBoundaries(t, "uint64", boundUint64Cases(), DecodeCompact64)
}

func testDecodeNonCompact[T any](
	t *testing.T,
	name string,
	val T,
	encoded []byte,
	decode func(d *Decoder) (T, int, error),
) {
	t.Run(name, func(t *testing.T) {
		buf := bytes.NewBuffer(encoded)
		dec := NewDecoder(buf)
		rst, _, err := decode(dec)
		require.NoError(t, err)
		require.Equal(t, val, rst)
	})
}

func TestDecodeUint(t *testing.T) {
	testDecodeNonCompact(
		t, "uint8",
		uint8(0x42),
		[]byte{0x42},
		DecodeByte,
	)
	testDecodeNonCompact(
		t, "uint16",
		uint16(0x1234),
		[]byte{0x34, 0x12},
		DecodeUint16,
	)
	testDecodeNonCompact(
		t, "uint32",
		uint32(0x12345678),
		[]byte{0x78, 0x56, 0x34, 0x12},
		DecodeUint32,
	)
	testDecodeNonCompact(
		t, "uint64",
		uint64(0x123456789abcdef0),
		[]byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12},
		DecodeUint64,
	)
}
