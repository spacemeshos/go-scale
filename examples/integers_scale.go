// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *U8) EncodeScale(enc *scale.Encoder) (total int, err error) {
	// field Value (0)
	if n, err := scale.EncodeCompact8(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}

	return total, nil
}

func (t *U8) DecodeScale(dec *scale.Decoder) (total int, err error) {
	// field Value (0)
	if field, n, err := scale.DecodeCompact8(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}

	return total, nil
}

func (t *U16) EncodeScale(enc *scale.Encoder) (total int, err error) {
	// field Value (0)
	if n, err := scale.EncodeCompact16(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}

	return total, nil
}

func (t *U16) DecodeScale(dec *scale.Decoder) (total int, err error) {
	// field Value (0)
	if field, n, err := scale.DecodeCompact16(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}

	return total, nil
}

func (t *U32) EncodeScale(enc *scale.Encoder) (total int, err error) {
	// field Value (0)
	if n, err := scale.EncodeCompact32(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}

	return total, nil
}

func (t *U32) DecodeScale(dec *scale.Decoder) (total int, err error) {
	// field Value (0)
	if field, n, err := scale.DecodeCompact32(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}

	return total, nil
}

func (t *U64) EncodeScale(enc *scale.Encoder) (total int, err error) {
	// field Value (0)
	if n, err := scale.EncodeCompact64(enc, t.Value); err != nil {
		return total, err
	} else {
		total += n
	}

	return total, nil
}

func (t *U64) DecodeScale(dec *scale.Decoder) (total int, err error) {
	// field Value (0)
	if field, n, err := scale.DecodeCompact64(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Value = field
	}

	return total, nil
}
