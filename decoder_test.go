package scale

import (
	"bytes"
	"fmt"
	"io"
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

// encodeCompact is a generic encode for primitive types using compact representation.
// other types are derived from this types.
func encodeCompact(tb testing.TB, value any) []byte {
	var (
		buf = bytes.NewBuffer(nil)
		enc = NewEncoder(buf)
		err error
	)
	switch val := value.(type) {
	case uint8:
		_, err = EncodeCompact8(enc, val)
	case uint16:
		_, err = EncodeCompact16(enc, val)
	case uint32:
		_, err = EncodeCompact32(enc, val)
	case uint64:
		_, err = EncodeCompact64(enc, val)
	case *uint8:
		_, err = EncodeCompact8Ptr(enc, val)
	case *uint16:
		_, err = EncodeCompact16Ptr(enc, val)
	case *uint32:
		_, err = EncodeCompact32Ptr(enc, val)
	case *uint64:
		_, err = EncodeCompact64Ptr(enc, val)
	case []byte:
		_, err = EncodeByteSlice(enc, val)
	case []uint16:
		_, err = EncodeUint16Slice(enc, val)
	case []uint32:
		_, err = EncodeUint32Slice(enc, val)
	case []uint64:
		_, err = EncodeUint64Slice(enc, val)
	case string:
		_, err = EncodeString(enc, val)
	case []string:
		_, err = EncodeStringSlice(enc, val)
	}
	require.NoError(tb, err)
	return buf.Bytes()
}

func expectEqual_Compact(tb testing.TB, value any, r io.Reader) {
	var (
		dec = NewDecoder(r)
		err error
		rst any
	)
	switch value.(type) {
	case uint8:
		rst, _, err = DecodeCompact8(dec)
	case uint16:
		rst, _, err = DecodeCompact16(dec)
	case uint32:
		rst, _, err = DecodeCompact32(dec)
	case uint64:
		rst, _, err = DecodeCompact64(dec)
	case *uint8:
		rst, _, err = DecodeCompact8Ptr(dec)
	case *uint16:
		rst, _, err = DecodeCompact16Ptr(dec)
	case *uint32:
		rst, _, err = DecodeCompact32Ptr(dec)
	case *uint64:
		rst, _, err = DecodeCompact64Ptr(dec)
	case []byte:
		rst, _, err = DecodeByteSlice(dec)
	case []uint16:
		rst, _, err = DecodeUint16Slice(dec)
	case []uint32:
		rst, _, err = DecodeUint32Slice(dec)
	case []uint64:
		rst, _, err = DecodeUint64Slice(dec)
	case string:
		rst, _, err = DecodeString(dec)
	case []string:
		rst, _, err = DecodeStringSlice(dec)
	}
	require.NoError(tb, err)
	require.Equal(tb, value, rst)
}

func TestReadFull(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		expect any
	}{
		{
			desc:   "uint8",
			expect: uint8(math.MaxUint8),
		},
		{
			desc:   "uint16",
			expect: uint16(math.MaxUint16),
		},
		{
			desc:   "uint32",
			expect: uint32(math.MaxUint32),
		},
		{
			desc:   "uint64",
			expect: uint64(math.MaxUint64),
		},
		{
			desc:   "*uint8",
			expect: intPtr[uint8](math.MaxUint8),
		},
		{
			desc:   "nil *uint8",
			expect: (*uint8)(nil),
		},
		{
			desc:   "*uint16",
			expect: intPtr[uint16](math.MaxUint8),
		},
		{
			desc:   "nil *uint16",
			expect: (*uint16)(nil),
		},
		{
			desc:   "*uint32",
			expect: intPtr[uint32](math.MaxUint8),
		},
		{
			desc:   "nil *uint32",
			expect: (*uint32)(nil),
		},
		{
			desc:   "*uint64",
			expect: intPtr[uint64](math.MaxUint8),
		},
		{
			desc:   "nil *uint64",
			expect: (*uint64)(nil),
		},
		{
			desc:   "byte slice",
			expect: []byte("dsa1232131312dsada123312"),
		},
		{
			desc:   "string",
			expect: "dsa1232131312dsada123312",
		},
		{
			desc:   "string slice",
			expect: []string{"qwe123", "dsa456"},
		},
		{
			desc:   "uint16 slice",
			expect: []uint16{0, 1, 2, math.MaxUint8, math.MaxUint16},
		},
		{
			desc:   "uint32 slice",
			expect: []uint32{0, 1, 2, math.MaxUint8, math.MaxUint16, math.MaxUint32},
		},
		{
			desc:   "uint64 slice",
			expect: []uint64{0, 1, 2, math.MaxUint32, math.MaxUint64},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Run("full", func(t *testing.T) {
				expectEqual_Compact(t, tc.expect, bytes.NewReader(encodeCompact(t, tc.expect)))
			})
			t.Run("partial", func(t *testing.T) {
				expectEqual_Compact(t, tc.expect, &halfReader{
					rem: encodeCompact(t, tc.expect),
				})
			})
		})
	}
}

func decodeCompactTest[T any](t *testing.T, value []byte, expect T) {
	buf := bytes.NewBuffer(value)
	dec := NewDecoder(buf)
	switch typed := any(expect).(type) {
	case uint8:
		rst, _, err := DecodeCompact8(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint16:
		rst, _, err := DecodeCompact16(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint32:
		rst, _, err := DecodeCompact32(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint64:
		rst, _, err := DecodeCompact64(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	}
}

func TestDecodeCompactIntegers(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		for _, tc := range uint8CompactTestCases() {
			t.Run("", func(t *testing.T) {
				decodeCompactTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint16", func(t *testing.T) {
		for _, tc := range uint16CompactTestCases() {
			t.Run("", func(t *testing.T) {
				decodeCompactTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint32", func(t *testing.T) {
		for _, tc := range uint32CompactTestCases() {
			t.Run("", func(t *testing.T) {
				decodeCompactTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint64", func(t *testing.T) {
		for _, tc := range uint64CompactTestCases() {
			t.Run("", func(t *testing.T) {
				decodeCompactTest(t, tc.expect, tc.value)
			})
		}
	})
}

type boundsTestCase struct {
	value []byte
}

func boundsDecodeTest(t *testing.T, tc boundsTestCase, decodeFunc func(*Decoder) error) {
	buf := bytes.NewBuffer(tc.value)
	dec := NewDecoder(buf)
	require.Error(t, decodeFunc(dec), fmt.Sprintf("%b", tc.value))
}

func boundsUint8Cases() []boundsTestCase {
	return []boundsTestCase{
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

func boundsUint16Cases() []boundsTestCase {
	return []boundsTestCase{
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

func boundsUint32Cases() []boundsTestCase {
	return []boundsTestCase{
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

func boundsUint64Cases() []boundsTestCase {
	return []boundsTestCase{
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

func TestCompactIntegersBoundaries(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		for _, tc := range boundsUint8Cases() {
			t.Run("", func(t *testing.T) {
				boundsDecodeTest(t, tc, func(dec *Decoder) error {
					_, _, err := DecodeCompact8(dec)
					return err
				})
			})
		}
	})
	t.Run("uint16", func(t *testing.T) {
		for _, tc := range boundsUint16Cases() {
			t.Run("", func(t *testing.T) {
				boundsDecodeTest(t, tc, func(dec *Decoder) error {
					_, _, err := DecodeCompact16(dec)
					return err
				})
			})
		}
	})
	t.Run("uint32", func(t *testing.T) {
		for _, tc := range boundsUint32Cases() {
			t.Run("", func(t *testing.T) {
				boundsDecodeTest(t, tc, func(dec *Decoder) error {
					_, _, err := DecodeCompact32(dec)
					return err
				})
			})
		}
	})
	t.Run("uint64", func(t *testing.T) {
		for _, tc := range boundsUint64Cases() {
			t.Run("", func(t *testing.T) {
				boundsDecodeTest(t, tc, func(dec *Decoder) error {
					_, _, err := DecodeCompact64(dec)
					return err
				})
			})
		}
	})
}

func decodeTest[T any](t *testing.T, value []byte, expect T) {
	buf := bytes.NewBuffer(value)
	dec := NewDecoder(buf)
	switch typed := any(expect).(type) {
	case uint8:
		rst, _, err := DecodeByte(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint16:
		rst, _, err := DecodeUint16(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint32:
		rst, _, err := DecodeUint32(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	case uint64:
		rst, _, err := DecodeUint64(dec)
		require.NoError(t, err)
		require.Equal(t, typed, rst)
	}
}

func TestDecodeUint(t *testing.T) {
	for _, tc := range nonCompactTestCases() {
		t.Run("", func(t *testing.T) {
			decodeTest(t, tc.expect, tc.value)
		})
	}
}
