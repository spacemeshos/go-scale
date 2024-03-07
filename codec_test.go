package scale

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompact32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)

	const expected = 321
	_, err := EncodeCompact32(enc, expected)
	require.NoError(t, err)
	dec := NewDecoder(buf)
	rst, _, err := DecodeCompact32(dec)
	require.NoError(t, err)
	require.EqualValues(t, expected, rst)
}

func TestCompactBigInt(t *testing.T) {
	const encoded uint64 = math.MaxUint64
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	_, err := EncodeCompact64(enc, encoded)
	require.NoError(t, err)

	dec := NewDecoder(buf)
	rst, _, err := DecodeCompact64(dec)
	require.NoError(t, err)
	require.EqualValues(t, encoded, rst)
}

func TestLenSize(t *testing.T) {
	for _, l := range []uint32{
		0,
		1,
		16,
		64,
		65,
		1<<14 - 1,
		1 << 14,
		1 << 20,
		1 << 30,
		1 << 31,
		1<<32 - 1,
	} {
		var b bytes.Buffer
		n, err := EncodeLen(NewEncoder(&b), l, 1<<32-1)
		require.NoError(t, err)
		require.Equal(t, uint32(b.Len()), LenSize(l))
		require.Equal(t, uint32(n), LenSize(l))
		decoded, n, err := DecodeLen(NewDecoder(&b), 1<<32-1)
		require.NoError(t, err)
		require.Equal(t, l, decoded)
		require.Equal(t, uint32(n), LenSize(l))
	}
}
