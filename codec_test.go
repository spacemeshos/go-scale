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
	var rst uint32
	dec := NewDecoder(buf)
	_, err = DecodeCompact32(dec, &rst)
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
	var value uint64
	_, err = DecodeCompact64(dec, &value)
	require.NoError(t, err)

	require.EqualValues(t, encoded, value)
}
