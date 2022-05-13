package scale

type Type interface {
	Encodable
	Decodable
}

type TypeHelper[T any] interface {
	Type
	*T
}

type Bytes20 [20]byte

func (a *Bytes20) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Bytes20) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type Bytes32 [32]byte

func (a *Bytes32) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Bytes32) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type Hash20 = Bytes20

type Address = Bytes20

type PublicKey = Bytes32

type Hash32 = Bytes32

type Signature [64]byte

func (a *Signature) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Signature) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type ByteSlice []byte

func (b *ByteSlice) EncodeScale(e *Encoder) (int, error) {
	total, err := EncodeLen(e, uint32(len(*b)))
	if err != nil {
		return 0, err
	}
	n, err := e.w.Write(*b)
	if err != nil {
		return 0, err
	}
	return total + n, nil
}

func (b *ByteSlice) DecodeScale(d *Decoder) (int, error) {
	lth, total, err := DecodeLen(d)
	if err != nil {
		return 0, err
	}
	if uint32(len(*b)) < lth {
		*b = make([]byte, lth)
	}
	n, err := d.r.Read(*b)
	if err != nil {
		return 0, err
	}
	return total + n, nil
}

type String = ByteSlice

func (s String) String() string {
	return string(s)
}
