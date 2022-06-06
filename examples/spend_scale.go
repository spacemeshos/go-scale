// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package examples

import (
	"github.com/spacemeshos/go-scale"
)

func (t *Spend) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeCompact8(enc, t.Type); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := t.Body.EncodeScale(enc); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *Spend) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeCompact8(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Type = field
	}
	if n, err := t.Body.DecodeScale(dec); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendBody) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteArray(enc, t.Adress[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeCompact8(enc, t.Selector); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := t.Payload.EncodeScale(enc); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendBody) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeByteArray(dec, t.Adress[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	if field, n, err := scale.DecodeCompact8(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Selector = field
	}
	if n, err := t.Payload.DecodeScale(dec); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendPayload) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := t.Arguments.EncodeScale(enc); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := t.Nonce.EncodeScale(enc); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeCompact32(enc, t.GasPrice); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeByteArray(enc, t.Signature[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendPayload) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := t.Arguments.DecodeScale(dec); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := t.Nonce.DecodeScale(dec); err != nil {
		return total, err
	} else {
		total += n
	}
	if field, n, err := scale.DecodeCompact32(dec); err != nil {
		return total, err
	} else {
		total += n
		t.GasPrice = field
	}
	if n, err := scale.DecodeByteArray(dec, t.Signature[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendArguments) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeByteArray(enc, t.Recipient[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeCompact64(enc, t.Amount); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendArguments) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if n, err := scale.DecodeByteArray(dec, t.Recipient[:]); err != nil {
		return total, err
	} else {
		total += n
	}
	if field, n, err := scale.DecodeCompact64(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Amount = field
	}
	return total, nil
}

func (t *SpendNonce) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeCompact32(enc, t.Counter); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeCompact64(enc, t.Bitfield); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *SpendNonce) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeCompact32(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Counter = field
	}
	if field, n, err := scale.DecodeCompact64(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Bitfield = field
	}
	return total, nil
}