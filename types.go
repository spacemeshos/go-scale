package scale

type Address [20]byte

func (a *Address) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Address) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type Hash20 [20]byte

func (a *Hash20) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Hash20) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type Hash32 [32]byte

func (a *Hash32) EncodeScale(e *Encoder) (int, error) {
	return e.w.Write(a[:])
}

func (a *Hash32) DecodeScale(d *Decoder) (int, error) {
	return d.r.Read(a[:])
}

type PublicKey = Hash32

type ByteSlice []byte

func (b *ByteSlice) EncodeScale(e *Encoder) (int, error) {
	total, err := EncodeCompact32(e, uint32(len(*b)))
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
	lth, total, err := DecodeCompact32(d)
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
