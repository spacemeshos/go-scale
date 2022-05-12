package scale

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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
