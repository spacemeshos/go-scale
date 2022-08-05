// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

// nolint
package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *Ex2) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructSlice(enc, t.Slice)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeStructArray(enc, t.Array[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Ex2) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStructSlice[Ex2](dec)
		if err != nil {
			return total, err
		}
		total += n
		t.Slice = field
	}
	{
		n, err := scale.DecodeStructArray(dec, t.Array[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Smth) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeCompact32(enc, uint32(t.Val))
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Smth) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeCompact32(dec)
		if err != nil {
			return total, err
		}
		total += n
		t.Val = uint32(field)
	}
	return total, nil
}

func (t *StructSliceWithLimit) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructSliceWithLimit(enc, t.Slice, 2)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *StructSliceWithLimit) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStructSliceWithLimit[Smth](dec, 2)
		if err != nil {
			return total, err
		}
		total += n
		t.Slice = field
	}
	return total, nil
}
