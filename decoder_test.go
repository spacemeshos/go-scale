package scale

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/spacemeshos/go-scale/compat"
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

// testEncode is a generic encode for primitive types.
// other types are derived from this types.
func testEncode(tb testing.TB, value any) []byte {
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
	case []byte:
		_, err = EncodeByteSlice(enc, val)
	case string:
		_, err = EncodeString(enc, val)
	case []string:
		_, err = EncodeStringSlice(enc, val)
	}
	require.NoError(tb, err)
	return buf.Bytes()
}

func expectEqual(tb testing.TB, value any, r io.Reader) {
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
	case []byte:
		rst, _, err = DecodeByteSlice(dec)
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
	} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Run("full", func(t *testing.T) {
				expectEqual(t, tc.expect, bytes.NewReader(testEncode(t, tc.expect)))
			})
			t.Run("partial", func(t *testing.T) {
				expectEqual(t, tc.expect, &halfReader{
					rem: testEncode(t, tc.expect),
				})
			})
		})
	}
}

func decodeTest[T any](t *testing.T, value []byte, expect T) {
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
		for _, tc := range uint8TestCases() {
			t.Run("", func(t *testing.T) {
				decodeTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint16", func(t *testing.T) {
		for _, tc := range uint16TestCases() {
			t.Run("", func(t *testing.T) {
				decodeTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint32", func(t *testing.T) {
		for _, tc := range uint32TestCases() {
			t.Run("", func(t *testing.T) {
				decodeTest(t, tc.expect, tc.value)
			})
		}
	})
	t.Run("uint64", func(t *testing.T) {
		for _, tc := range uint64TestCases() {
			t.Run("", func(t *testing.T) {
				decodeTest(t, tc.expect, tc.value)
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

func TestRoundTrip(t *testing.T) {
	buf := make([]byte, 20)
	rand.Read(buf)
	output, err := compat.RoundTrip(buf)
	require.NoError(t, err)
	require.Equal(t, buf, output)
}
