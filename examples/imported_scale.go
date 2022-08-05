// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-scale/examples/alias"
)

func (t *ImportedA) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeStructArray(enc, t.ListA[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	if n, err := scale.EncodeStructSlice(enc, t.ListB); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *ImportedA) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeStructArray(dec, t.ListA[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	if field, n, err := scale.DecodeStructSlice[alias.A](dec); err != nil {
		return total, err
	} else { // nolint
		total += n
		t.ListB = field
	}
	return total, nil
}