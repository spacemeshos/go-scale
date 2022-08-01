package tester

import (
	"bytes"
	"errors"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

func FuzzConsistency[T any, H scale.TypePtr[T]](f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		fuzzer := fuzz.NewFromGoFuzz(data)
		var object T
		fuzzer.Fuzz(&object)

		buf := bytes.NewBuffer(nil)
		enc := scale.NewEncoder(buf)
		_, err := H(&object).EncodeScale(enc)
		if errors.Is(err, scale.ErrEncodeTooManyElements) {
			return
		}
		require.NoError(t, err)

		dec := scale.NewDecoder(buf)
		var decoded T
		_, err = H(&decoded).DecodeScale(dec)
		require.NoError(t, err)

		require.Equal(t, object, decoded)
	})
}

func FuzzSafety[T any, H scale.TypePtr[T]](f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		dec := scale.NewDecoder(bytes.NewReader(data))
		var decoded T
		H(&decoded).DecodeScale(dec)
	})
}
