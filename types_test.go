package scale

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type Complex struct {
	Field1 ByteSlice
	Uint64 uint64
	Nested []Complex
}

func (c Complex) EncodeScale(enc *Encoder) (int, error) {
	total := 0
	if n, err := c.Field1.EncodeScale(enc); err != nil {
		return 0, err
	} else {
		total += n
	}
	if n, err := EncodeCompact64(enc, c.Uint64); err != nil {
		return 0, err
	} else {
		total += n
	}
	if n, err := EncodeStructSlice(enc, c.Nested); err != nil {
		return 0, err
	} else {
		total += n
	}
	return total, nil
}

func (c *Complex) DecodeScale(dec *Decoder) (int, error) {
	var (
		total, n int
		err      error
	)
	val, n, err := DecodeStruct[ByteSlice](dec)
	if err != nil {
		return 0, err
	} else {
		total += n
	}
	c.Field1 = val

	if val, n, err := DecodeCompact64(dec); err != nil {
		return 0, err
	} else {
		total += n
		c.Uint64 = val
	}
	if val, n, err := DecodeStructSlice[Complex](dec); err != nil {
		return 0, err
	} else {
		total += n
		c.Nested = val
	}
	return total, nil
}

func TestStruct(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	c := Complex{
		Field1: ByteSlice("dsasda"),
		Uint64: 321,
		Nested: []Complex{
			{Uint64: 11},
			{Field1: ByteSlice("321313")},
		},
	}
	_, err := c.EncodeScale(enc)
	require.NoError(t, err)

	var rst Complex
	rst.Nested = []Complex{}
	dec := NewDecoder(buf)
	_, err = rst.DecodeScale(dec)
	require.NoError(t, err)
	require.Equal(t, c, rst)
}

func TestBytes(t *testing.T) {
	bsl := ByteSlice("dasdk13213dsaasda")
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	_, err := bsl.EncodeScale(enc)
	require.NoError(t, err)

	dec := NewDecoder(buf)
	var decoded ByteSlice
	_, err = decoded.DecodeScale(dec)
	require.NoError(t, err)
	require.Equal(t, bsl, decoded)
}
