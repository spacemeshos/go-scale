package tester

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

func FuzzConsistency[T any, H scale.TypePtr[T]](f *testing.F, fuzzFuncs ...any) {
	f.Fuzz(func(t *testing.T, data []byte) {
		fuzzer := fuzz.NewFromGoFuzz(data).Funcs(fuzzFuncs...)
		var object T
		fuzzer.Fuzz(&object)

		buf := bytes.NewBuffer(nil)
		enc := scale.NewEncoder(buf)
		n1, err := H(&object).EncodeScale(enc)
		if errors.Is(err, scale.ErrEncodeTooManyElements) {
			return
		}
		require.NoError(t, err)
		require.Equal(t, n1, buf.Len())

		dec := scale.NewDecoder(buf)
		var decoded T
		n2, err := H(&decoded).DecodeScale(dec)
		require.NoError(t, err)

		require.Equal(t, n1, n2)

		if !cmp.Equal(object, decoded, cmpopts.EquateEmpty()) {
			t.Errorf("decoded didn't match original: %s", cmp.Diff(object, decoded))
		}
	})
}

func FuzzSafety[T any, H scale.TypePtr[T]](f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		dec := scale.NewDecoder(bytes.NewReader(data))
		var decoded T
		H(&decoded).DecodeScale(dec)
	})
}
