// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

// nolint
package testdata

import (
	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-scale/examples/nested"
)

func (t *Data) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStringWithLimit(enc, string(t.Str), 20)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := t.NestedStruct.EncodeScale(enc)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeOption(enc, t.NestedStructPointer)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeStructSliceWithLimit(enc, t.NestedStructSlice, 5)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *Data) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStringWithLimit(dec, 20)
		if err != nil {
			return total, err
		}
		total += n
		t.Str = string(field)
	}
	{
		n, err := t.NestedStruct.DecodeScale(dec)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		field, n, err := scale.DecodeOption[nested.Struct](dec)
		if err != nil {
			return total, err
		}
		total += n
		t.NestedStructPointer = field
	}
	{
		field, n, err := scale.DecodeStructSliceWithLimit[nested.Struct](dec, 5)
		if err != nil {
			return total, err
		}
		total += n
		t.NestedStructSlice = field
	}
	return total, nil
}

func (t *MoreData) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStringWithLimit(enc, string(t.NestedAlias), 20)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeStringSliceWithLimit(enc, t.StrSlice, 5)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeByteArray(enc, t.ByteArray[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeByteSliceWithLimit(enc, t.ByteSlice, 20)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeSliceOfByteSliceWithLimit(enc, t.SliceOfByteSlices, 10)
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		n, err := scale.EncodeCompact64(enc, uint64(t.Uint64))
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *MoreData) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStringWithLimit(dec, 20)
		if err != nil {
			return total, err
		}
		total += n
		t.NestedAlias = nested.StringAlias(field)
	}
	{
		field, n, err := scale.DecodeStringSliceWithLimit(dec, 5)
		if err != nil {
			return total, err
		}
		total += n
		t.StrSlice = field
	}
	{
		n, err := scale.DecodeByteArray(dec, t.ByteArray[:])
		if err != nil {
			return total, err
		}
		total += n
	}
	{
		field, n, err := scale.DecodeByteSliceWithLimit(dec, 20)
		if err != nil {
			return total, err
		}
		total += n
		t.ByteSlice = field
	}
	{
		field, n, err := scale.DecodeSliceOfByteSliceWithLimit(dec, 10)
		if err != nil {
			return total, err
		}
		total += n
		t.SliceOfByteSlices = field
	}
	{
		field, n, err := scale.DecodeCompact64(dec)
		if err != nil {
			return total, err
		}
		total += n
		t.Uint64 = uint64(field)
	}
	return total, nil
}
