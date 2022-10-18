package tester

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

type GoldenTestCase[T any, H scale.TypePtr[T]] struct {
	Name   string
	Object T
	Hex    string
	Error  string
}

func GoldenTest[T any, H scale.TypePtr[T]](t *testing.T, path string) {
	f, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })
	jcodec := json.NewDecoder(f)
	for {
		var tc GoldenTestCase[T, H]
		if err := jcodec.Decode(&tc); errors.Is(err, io.EOF) {
			return
		} else {
			require.NoError(t, err)
		}
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("Encode", func(t *testing.T) {
				if len(tc.Error) > 0 {
					t.Skip("skip encoding invalid entry")
				}
				buf := bytes.NewBuffer(nil)
				encoder := scale.NewEncoder(buf)
				_, err := H(&tc.Object).EncodeScale(encoder)
				require.NoError(t, err)
				hx := hex.EncodeToString(buf.Bytes())
				require.Equal(t, tc.Hex, hx)
			})
			t.Run("Decode", func(t *testing.T) {
				buf, err := hex.DecodeString(tc.Hex)
				require.NoError(t, err)
				decoder := scale.NewDecoder(bytes.NewReader(buf))
				var object T
				_, err = H(&object).DecodeScale(decoder)
				if len(tc.Error) > 0 {
					require.Equal(t, tc.Error, err.Error())
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.Object, object)
				}
			})
		})
	}
}
