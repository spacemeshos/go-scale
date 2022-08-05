package alias

import "github.com/spacemeshos/go-scale"

type A [20]byte

func (a *A) EncodeScale(enc *scale.Encoder) (int, error) {
	return scale.EncodeByteArray(enc, a[:])
}

func (a *A) DecodeScale(dec *scale.Decoder) (int, error) {
	return scale.DecodeByteArray(dec, a[:])
}

type B = A
