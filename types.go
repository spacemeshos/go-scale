package scale

type Type interface {
	Encodable
	Decodable
}

type TypePtr[T any] interface {
	Type
	*T
}

type U8 uint8

func (u *U8) EncodeScale(e *Encoder) (int, error) {
	return EncodeCompact8(e, uint8(*u))
}

func (u *U8) DecodeScale(d *Decoder) (int, error) {
	r, n, err := DecodeCompact8(d)
	if err != nil {
		return n, err
	}
	*u = U8(r)
	return n, nil
}
