package examples

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

func TestNonLocal(t *testing.T) {
	s := StructWithNonLocalField{
		Name:   "foo",
		SomeID: "bar",
	}

	buf := bytes.NewBuffer(nil)
	encoder := scale.NewEncoder(buf)
	n, err := s.EncodeScale(encoder)
	require.NoError(t, err)
	require.Equal(t, 8, n)

	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	var s1 StructWithNonLocalField
	n, err = s1.DecodeScale(decoder)
	require.NoError(t, err)
	require.Equal(t, 8, n)
	require.Equal(t, s, s1)

	buf = bytes.NewBuffer(nil)
	encoder = scale.NewEncoder(buf, scale.WithEncodeLocal())
	n, err = s.EncodeScale(encoder)
	require.NoError(t, err)
	require.Equal(t, 4, n)

	decoder = scale.NewDecoder(bytes.NewReader(buf.Bytes()), scale.WithDecodeLocal())
	var s2 StructWithNonLocalField
	n, err = s2.DecodeScale(decoder)
	require.NoError(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, StructWithNonLocalField{Name: "foo"}, s2)
}
