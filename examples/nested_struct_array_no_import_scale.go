// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

// nolint
package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *NestedStructArray) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *NestedStructArray) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}
