package scale

import (
	"bytes"
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
