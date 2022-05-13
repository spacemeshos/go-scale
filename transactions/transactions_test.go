package transactions

import (
	"bytes"
	"testing"

	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-scale/gen"
	"github.com/spacemeshos/go-scale/transactions/types"

	xdr "github.com/nullstyle/go-xdr/xdr3"
	"github.com/stretchr/testify/require"
)

func TestGenTransactions(t *testing.T) {
	require.NoError(t, gen.Generate("types", "types/types_scale.go",
		types.SelfSpawn{}, types.SelfSpawnArguments{}, types.SelfSpawnBody{}, types.SelfSpawnPayload{}))
}

func TestSelfSpawn(t *testing.T) {
	tx := types.SelfSpawn{
		Type: 1,
		Body: types.SelfSpawnBody{
			Address:  scale.Address{1, 1, 1},
			Selector: 15,
			Payload: types.SelfSpawnPayload{
				Template: scale.Address{9, 9, 9},
				Arguments: types.SelfSpawnArguments{
					Key: scale.PublicKey{19, 19},
				},
				GasPrice:  100,
				Signature: scale.Signature{13, 23},
			},
		},
	}

	buf := bytes.NewBuffer(nil)
	enc := scale.NewEncoder(buf)
	en, err := tx.EncodeScale(enc)
	require.NoError(t, err)

	dec := scale.NewDecoder(buf)
	var decoded types.SelfSpawn
	dn, err := decoded.DecodeScale(dec)
	require.NoError(t, err)
	require.Equal(t, tx, decoded)
	require.Equal(t, en, dn)
}

func BenchmarkSelfSpawn(b *testing.B) {
	tx := types.SelfSpawn{
		Type: 1,
		Body: types.SelfSpawnBody{
			Address:  scale.Address{1, 1, 1},
			Selector: 15,
			Payload: types.SelfSpawnPayload{
				Template: scale.Address{9, 9, 9},
				Arguments: types.SelfSpawnArguments{
					Key: scale.PublicKey{19, 19},
				},
				GasPrice:  100,
				Signature: scale.Signature{13, 23},
			},
		},
	}

	b.Run("Encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(nil)
			enc := scale.NewEncoder(buf)
			_, err := tx.EncodeScale(enc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	buf := bytes.NewBuffer(nil)
	enc := scale.NewEncoder(buf)
	_, err := tx.EncodeScale(enc)
	if err != nil {
		b.Fatal(err)
	}
	byts := buf.Bytes()

	b.Run("Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dec := scale.NewDecoder(bytes.NewBuffer(byts))
			var decoded types.SelfSpawn
			_, err = decoded.DecodeScale(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("EncodeXDR", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(nil)
			_, err := xdr.Marshal(buf, tx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	buf = bytes.NewBuffer(nil)
	_, err = xdr.Marshal(buf, tx)
	if err != nil {
		b.Fatal(err)
	}
	byts = buf.Bytes()

	b.Run("DecodeXDR", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(byts)
			var decoded types.SelfSpawn
			_, err := xdr.Unmarshal(buf, &decoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
