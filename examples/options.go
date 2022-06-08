package examples

import "github.com/spacemeshos/go-scale"

type array [29]byte

func (a *array) EncodeScale(encoder *scale.Encoder) (int, error) {
	return scale.EncodeByteArray(encoder, a[:])
}

func (a *array) DecodeScale(decoder *scale.Decoder) (int, error) {
	return scale.DecodeByteArray(decoder, a[:])
}

type slice []byte

func (a *slice) EncodeScale(encoder *scale.Encoder) (int, error) {
	return scale.EncodeByteSlice(encoder, *a)
}

func (a *slice) DecodeScale(decoder *scale.Decoder) (int, error) {
	field, n, err := scale.DecodeByteSlice(decoder)
	if err != nil {
		return n, err
	}
	*a = field
	return n, nil
}

//go:generate scalegen

type Options struct {
	ArrayPtr *array
	SlicePtr *slice
}
