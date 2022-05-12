package scale

import "io"

type Decodable interface {
	DecodeScale(*Decoder) (int, error)
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

type Decoder struct {
	r       io.Reader
	scratch [9]byte
}

func DecodeCompact32(d *Decoder, value *uint32) (int, error) {
	_, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return 0, err
	}
	switch d.scratch[0] % 4 {
	case 0:
		*value = uint32(d.scratch[0]) >> 2
	case 1:
		_, err := d.r.Read(d.scratch[1:2])
		if err != nil {
			return 0, err
		}
		*value = (uint32(d.scratch[0]) | uint32(d.scratch[1])<<8) >> 2
	case 2:
		_, err := d.r.Read(d.scratch[1:4])
		if err != nil {
			return 0, err
		}
		*value = (uint32(d.scratch[0]) |
			uint32(d.scratch[1])<<8 |
			uint32(d.scratch[2])<<16 |
			uint32(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0]) >> 2
		_, err := d.r.Read(d.scratch[:needed])
		if err != nil {
			return 0, err
		}
		for i := 0; i < int(needed); i++ {
			*value |= uint32(d.scratch[i]) << (8 * i)
		}
	}
	return 0, nil
}

func DecodeCompact64(d *Decoder, value *uint64) (int, error) {
	_, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return 0, err
	}
	switch d.scratch[0] % 4 {
	case 0:
		*value = uint64(d.scratch[0]) >> 2
	case 1:
		_, err := d.r.Read(d.scratch[1:2])
		if err != nil {
			return 0, err
		}
		*value = (uint64(d.scratch[0]) | uint64(d.scratch[1])<<8) >> 2
	case 2:
		_, err := d.r.Read(d.scratch[1:4])
		if err != nil {
			return 0, err
		}
		*value = (uint64(d.scratch[0]) |
			uint64(d.scratch[1])<<8 |
			uint64(d.scratch[2])<<16 |
			uint64(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0]) >> 2
		_, err := d.r.Read(d.scratch[:needed])
		if err != nil {
			return 0, err
		}
		for i := 0; i < int(needed); i++ {
			*value |= uint64(d.scratch[i]) << (8 * i)
		}
	}
	return 0, nil
}

func DecodeBool(d *Decoder, value *bool) (int, error) {
	n, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return 0, err
	}
	if d.scratch[0] == 1 {
		*value = true
	} else {
		*value = false
	}
	return n, nil
}

func DecodeStructSlice[V *[]Decodable](d *Decoder, value V) (int, error) {
	var lth uint32
	total, err := DecodeCompact32(d, &lth)
	if err != nil {
		return 0, err
	}
	if uint32(len(*value)) < lth {
		*value = make([]Decodable, lth)
	}
	for i := range *value {
		n, err := (*value)[i].DecodeScale(d)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

func DecodeStructArray[V *[]Decodable](d *Decoder, value V) (int, error) {
	total := 0
	for i := range *value {
		n, err := (*value)[i].DecodeScale(d)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

func DecodeOption[V *Decodable](d *Decoder, value V) (int, error) {
	var exists bool
	total, err := DecodeBool(d, &exists)
	if !exists || err != nil {
		return total, err
	}
	var empty Decodable
	n, err := empty.DecodeScale(d)
	if err != nil {
		return 0, err
	}
	*value = empty
	return total + n, nil
}
