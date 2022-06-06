// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package nested

import (
	"github.com/spacemeshos/go-scale"
)

func (t *NestedModule) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteSlice(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *NestedModule) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeByteSlice(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}
	return total, nil
}
