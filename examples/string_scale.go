// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *StructWithString) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeString(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *StructWithString) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeString(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}
	return total, nil
}