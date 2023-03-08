package compat

import (
	"bytes"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

func FuzzRoundTrip(f *testing.F) {
	f.Add([]byte("0"))
	f.Add([]byte("321321dssae12312asada3421"))
	f.Fuzz(func(t *testing.T, data []byte) {
		fuzzer := fuzz.NewFromGoFuzz(data)

		compat := Compat{}
		fuzzer.Fuzz(&compat)
		buf := bytes.NewBuffer(nil)
		_, err := compat.EncodeScale(scale.NewEncoder(buf))
		require.NoError(t, err)

		output, err := RoundTrip(buf.Bytes())
		require.NoError(t, err)
		require.Equal(t, buf.Bytes(), output)

		result := Compat{}
		_, err = result.DecodeScale(scale.NewDecoder(bytes.NewReader(output)))
		require.NoError(t, err)
		require.Equal(t, compat, result)
	})
}
