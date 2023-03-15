// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

// nolint
package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *StructWithString) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStringWithLimit(enc, string(t.Value), 3)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *StructWithString) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStringWithLimit(dec, 3)
		if err != nil {
			return total, err
		}
		total += n
		t.Value = string(field)
	}
	return total, nil
}

func (t *StructWithStringAlias) EncodeScale(enc *scale.Encoder) (total int, err error) {
	{
		n, err := scale.EncodeStringWithLimit(enc, string(t.Value), 3)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func (t *StructWithStringAlias) DecodeScale(dec *scale.Decoder) (total int, err error) {
	{
		field, n, err := scale.DecodeStringWithLimit(dec, 3)
		if err != nil {
			return total, err
		}
		total += n
		t.Value = StringAlias(field)
	}
	return total, nil
}
