// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *Bytes20) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteArray(enc, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Bytes20) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeByteArray(dec, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Bytes32) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteArray(enc, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Bytes32) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeByteArray(dec, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Bytes64) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteArray(enc, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Bytes64) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeByteArray(dec, t.Value[:]); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Slice) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteSlice(enc, t.Value); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *Slice) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeByteSlice(dec); err != nil {
		return total, err
	} else { // nolint
		total += n
		t.Value = field
	}
	return total, nil
}

func (t *SliceWithLimit) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteSliceWithLimit(enc, t.Value, 10); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *SliceWithLimit) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeByteSliceWithLimit(dec, 10); err != nil {
		return total, err
	} else { // nolint
		total += n
		t.Value = field
	}
	return total, nil
}

func (t *SliceOfByteSliceWithLimit) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeSliceOfByteSlice(enc, t.Value); err != nil {
		return total, err
	} else { // nolint
		total += n
	}
	return total, nil
}

func (t *SliceOfByteSliceWithLimit) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeSliceOfByteSlice(dec); err != nil {
		return total, err
	} else { // nolint
		total += n
		t.Value = field
	}
	return total, nil
}
