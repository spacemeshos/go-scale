package scale

import (
	"fmt"
	"io"
)

type Decodable interface {
	DecodeScale(*Decoder) (int, error)
}

type DecodableHelper[B any] interface {
	Decodable
	*B
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

type Decoder struct {
	r       io.Reader
	scratch [9]byte
}

func DecodeCompact8(d *Decoder) (uint8, int, error) {
	var (
		value uint8
		total int
	)
	n, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint8(d.scratch[0]) >> 2
	case 1:
		n, err := d.r.Read(d.scratch[1:2])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = uint8((uint16(d.scratch[0]) | uint16(d.scratch[1])<<8) >> 2)
	}
	return value, total, nil
}

func DecodeCompact32(d *Decoder) (uint32, int, error) {
	var (
		value uint32
		total int
	)
	n, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint32(d.scratch[0]) >> 2
	case 1:
		n, err := d.r.Read(d.scratch[1:2])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = (uint32(d.scratch[0]) | uint32(d.scratch[1])<<8) >> 2
	case 2:
		n, err := d.r.Read(d.scratch[1:4])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = (uint32(d.scratch[0]) |
			uint32(d.scratch[1])<<8 |
			uint32(d.scratch[2])<<16 |
			uint32(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0]) >> 2
		_, err := d.r.Read(d.scratch[:needed])
		if err != nil {
			return value, 0, err
		}
		total += int(needed)
		if needed > 4 {
			return value, 0, fmt.Errorf("invalid compact32 encoding %x needs %d bytes", d.scratch[:needed], needed)
		}
		for i := 0; i < int(needed); i++ {
			value |= uint32(d.scratch[i]) << (8 * i)
		}
	}
	return value, total, nil
}

func DecodeCompact64(d *Decoder) (uint64, int, error) {
	var value uint64
	_, err := d.r.Read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	switch d.scratch[0] % 4 {
	case 0:
		value = uint64(d.scratch[0]) >> 2
	case 1:
		_, err := d.r.Read(d.scratch[1:2])
		if err != nil {
			return 0, 0, err
		}
		value = (uint64(d.scratch[0]) | uint64(d.scratch[1])<<8) >> 2
	case 2:
		_, err := d.r.Read(d.scratch[1:4])
		if err != nil {
			return 0, 0, err
		}
		value = (uint64(d.scratch[0]) |
			uint64(d.scratch[1])<<8 |
			uint64(d.scratch[2])<<16 |
			uint64(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0]) >> 2
		_, err := d.r.Read(d.scratch[:needed])
		if err != nil {
			return 0, 0, err
		}
		if needed > 8 {
			return value, 0, fmt.Errorf("invalid compact64 encoding %x needs %d bytes", d.scratch[:needed], needed)
		}
		for i := 0; i < int(needed); i++ {
			value |= uint64(d.scratch[i]) << (8 * i)
		}
	}
	return 0, 0, nil
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

func DecodeStruct[V any, H DecodableHelper[V]](d *Decoder) (V, int, error) {
	var empty V
	n, err := H(&empty).DecodeScale(d)
	return empty, n, err
}

func DecodeStructSlice[V any, H DecodableHelper[V]](d *Decoder) ([]V, int, error) {
	lth, total, err := DecodeCompact32(d)
	if err != nil {
		return nil, 0, err
	}
	if lth == 0 {
		return nil, 0, err
	}
	value := make([]V, 0, lth)

	for i := uint32(0); i < lth; i++ {
		val, n, err := DecodeStruct[V, H](d)
		if err != nil {
			return nil, 0, err
		}
		value = append(value, val)
		total += n
	}
	return value, total, nil
}

func DecodeStructArray[V Decodable](d *Decoder, value *[]V) (int, error) {
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

func DecodeOption[V any, H DecodableHelper[V]](d *Decoder) (*V, int, error) {
	var exists bool
	total, err := DecodeBool(d, &exists)
	if !exists || err != nil {
		return nil, total, err
	}
	var empty V
	n, err := H(&empty).DecodeScale(d)
	if err != nil {
		return nil, 0, err
	}
	return &empty, total + n, nil
}
