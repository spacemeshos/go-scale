// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

// nolint
package test

import (
	"github.com/spacemeshos/go-scale"
)

func (t *DeepNestedModule) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeOption(enc, t.Value)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *DeepNestedModule) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeOption[DeepNestedModule](dec)
		if err != nil {
			return total, err
		}
		total += n
		t.Value = field
	}
	return total, nil
}

func (t *DeepNestedSliceModule) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructSliceWithLimit(enc, t.Value, 2)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *DeepNestedSliceModule) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStructSliceWithLimit[DeepNestedSliceModule](dec, 2)
		if err != nil {
			return total, err
		}
		total += n
		t.Value = field
	}
	return total, nil
}

func (t *DeepNestedArrayModule) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *DeepNestedArrayModule) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level1) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level1) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level2) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level2) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level3) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level3) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level4) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStructArray(enc, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level4) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		n, err := scale.DecodeStructArray(dec, t.Value[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level5) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStringWithLimit(enc, string(t.Value), 64)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Level5) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStringWithLimit(dec, 64)
		if err != nil {
			return total, err
		}
		total += n
		t.Value = string(field)
	}
	return total, nil
}
