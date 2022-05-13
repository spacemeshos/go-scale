package types

import "github.com/spacemeshos/go-scale"

type Tx[T any, H scale.TypeHelper[T]] struct {
	Type uint8
	Body struct {
		CommonBody
		Payload T
	}
}

type CommonBody struct {
	Address  scale.Address
	Selector uint8
}

func Encode[T any, H scale.TypeHelper[T]](enc *scale.Encoder, tx Tx[T, H]) (int, error) {
	total := 0
	if n, err := scale.EncodeCompact8(enc, tx.Type); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeStruct(enc, tx.Body.CommonBody); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeStruct[T, H](enc, tx.Body.Payload); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}
