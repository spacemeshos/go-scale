// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *Ex2) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeStructSlice(enc, t.Slice); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeStructArray(enc, t.Array[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *Ex2) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeStructSlice[Ex2](dec); err != nil {
		return total, err
	} else {
		total += n
		t.Slice = field
	}
	if n, err := scale.DecodeStructArray(dec, t.Array[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *Smth) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeCompact32(enc, t.Val); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *Smth) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeCompact32(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Val = field
	}
	return total, nil
}
