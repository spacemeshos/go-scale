package scale

import (
	"errors"
	"fmt"
	"io"
	"math"
)

// ErrDecodeTooManyElements is returned when scale limit tag is used and collection has too many elements to decode.
var ErrDecodeTooManyElements = errors.New("too many elements to decode in collection with scale limit set")

type Decodable interface {
	DecodeScale(*Decoder) (int, error)
}

type DecodablePtr[B any] interface {
	Decodable
	*B
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Reset(r io.Reader) {
	d.r = r
}

type Decoder struct {
	r       io.Reader
	scratch [9]byte
}

func (d *Decoder) read(buf []byte) (int, error) {
	return io.ReadFull(d.r, buf)
}

func DecodeCompact8(d *Decoder) (uint8, int, error) {
	var (
		value uint8
		total int
	)
	n, err := d.read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint8(d.scratch[0]) >> 2
	case 1:
		n, err := d.read(d.scratch[1:2])
		if err != nil {
			return value, 0, err
		}
		total += n
		if rst := (uint16(d.scratch[0]) | uint16(d.scratch[1])<<8) >> 2; rst > math.MaxUint8 {
			return 0, 0, fmt.Errorf("value %d overflows uint8", rst)
		} else {
			value = uint8(rst)
		}
	default:
		return 0, 0, fmt.Errorf("value will overflow uint8")
	}
	return value, total, nil
}

func DecodeCompact16(d *Decoder) (uint16, int, error) {
	var (
		value uint16
		total int
	)
	n, err := d.read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint16(d.scratch[0]) >> 2
	case 1:
		n, err := d.read(d.scratch[1:2])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = (uint16(d.scratch[0]) | uint16(d.scratch[1])<<8) >> 2
	case 2:
		n, err := d.read(d.scratch[1:4])
		if err != nil {
			return value, 0, err
		}
		total += n
		if rst := (uint32(d.scratch[0]) |
			uint32(d.scratch[1])<<8 |
			uint32(d.scratch[2])<<16 |
			uint32(d.scratch[3])<<24) >> 2; rst > math.MaxUint16 {
			return 0, 0, fmt.Errorf("value %d overflows uint16", rst)
		} else {
			value = uint16(rst)
		}
	default:
		return 0, 0, fmt.Errorf("value will overflow uint16")
	}
	return value, total, nil
}

func DecodeCompact32(d *Decoder) (uint32, int, error) {
	var (
		value uint32
		total int
	)
	n, err := d.read(d.scratch[:1])
	if err != nil {
		return value, 0, err
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint32(d.scratch[0]) >> 2
	case 1:
		n, err := d.read(d.scratch[1:2])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = (uint32(d.scratch[0]) | uint32(d.scratch[1])<<8) >> 2
	case 2:
		n, err := d.read(d.scratch[1:4])
		if err != nil {
			return value, 0, err
		}
		total += n
		value = (uint32(d.scratch[0]) |
			uint32(d.scratch[1])<<8 |
			uint32(d.scratch[2])<<16 |
			uint32(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0])>>2 + 4
		if needed > 4 {
			return value, 0, fmt.Errorf("invalid compact32 needs %d bytes", needed)
		}
		_, err := d.read(d.scratch[:needed])
		if err != nil {
			return value, 0, err
		}
		total += int(needed)
		for i := 0; i < int(needed); i++ {
			value |= uint32(d.scratch[i]) << (8 * i)
		}
	}
	return value, total, nil
}

func DecodeLen(d *Decoder, limit uint32) (uint32, int, error) {
	v, n, err := DecodeCompact32(d)
	if err != nil {
		return v, n, err
	}
	if v > limit {
		return v, n, fmt.Errorf("%w: (%d > %d)", ErrDecodeTooManyElements, v, limit)
	}
	return v, n, err
}

func DecodeCompact64(d *Decoder) (uint64, int, error) {
	var (
		value uint64
		total int
	)
	n, err := d.read(d.scratch[:1])
	if err != nil {
		return value, 0, fmt.Errorf("read compact header: %w", err)
	}
	total += n
	switch d.scratch[0] % 4 {
	case 0:
		value = uint64(d.scratch[0]) >> 2
	case 1:
		n, err := d.read(d.scratch[1:2])
		if err != nil {
			return 0, 0, err
		}
		total += n
		value = (uint64(d.scratch[0]) | uint64(d.scratch[1])<<8) >> 2
	case 2:
		n, err := d.read(d.scratch[1:4])
		if err != nil {
			return 0, 0, err
		}
		total += n
		value = (uint64(d.scratch[0]) |
			uint64(d.scratch[1])<<8 |
			uint64(d.scratch[2])<<16 |
			uint64(d.scratch[3])<<24) >> 2
	case 3:
		needed := byte(d.scratch[0]) >> 2
		if needed > 8 {
			return value, 0, fmt.Errorf("invalid compact64 needs %d bytes", needed)
		}
		// add back 4 bytes deducted in encoder
		needed += 4
		n, err := d.read(d.scratch[:needed])
		if err != nil {
			return 0, 0, err
		}
		total += n
		for i := 0; i < int(needed); i++ {
			value |= uint64(d.scratch[i]) << (8 * i)
		}
	}
	return value, total, nil
}

func DecodeBool(d *Decoder) (bool, int, error) {
	n, err := d.read(d.scratch[:1])
	if err != nil {
		return false, 0, err
	}
	if d.scratch[0] == 1 {
		return true, n, nil
	}
	return false, n, nil
}

func DecodeStruct[V any, H DecodablePtr[V]](d *Decoder) (V, int, error) {
	var empty V
	n, err := H(&empty).DecodeScale(d)
	return empty, n, err
}

func DecodeByteSlice(d *Decoder) ([]byte, int, error) {
	return DecodeByteSliceWithLimit(d, MaxElements)
}

func DecodeByteSliceWithLimit(d *Decoder, limit uint32) ([]byte, int, error) {
	lth, total, err := DecodeLen(d, limit)
	if err != nil {
		return nil, 0, err
	}
	if lth == 0 {
		return nil, total, nil
	}
	value := make([]byte, lth)
	n, err := DecodeByteArray(d, value)
	if err != nil {
		return value, 0, err
	}
	return value, total + n, nil
}

func DecodeByteArray(d *Decoder, value []byte) (int, error) {
	return d.read(value)
}

func DecodeString(d *Decoder) (string, int, error) {
	return DecodeStringWithLimit(d, MaxElements)
}

func DecodeStringWithLimit(d *Decoder, limit uint32) (string, int, error) {
	bytes, total, err := DecodeByteSliceWithLimit(d, limit)
	return string(bytes), total, err
}

func DecodeStructSlice[V any, H DecodablePtr[V]](d *Decoder) ([]V, int, error) {
	return DecodeStructSliceWithLimit[V, H](d, MaxElements)
}

func DecodeStructSliceWithLimit[V any, H DecodablePtr[V]](d *Decoder, limit uint32) ([]V, int, error) {
	lth, total, err := DecodeLen(d, limit)
	if err != nil {
		return nil, 0, err
	}
	if lth == 0 {
		return nil, 0, nil
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

func DecodeSliceOfByteSlice(d *Decoder) ([][]byte, int, error) {
	return DecodeSliceOfByteSliceWithLimit(d, MaxElements)
}

func DecodeSliceOfByteSliceWithLimit(d *Decoder, limit uint32) ([][]byte, int, error) {
	resultLen, total, err := DecodeLen(d, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("DecodeLen failed: %w", err)
	}
	if resultLen == 0 {
		return nil, 0, nil
	}
	result := make([][]byte, 0, resultLen)

	for i := uint32(0); i < resultLen; i++ {
		val, n, err := DecodeByteSlice(d)
		if err != nil {
			return nil, 0, fmt.Errorf("DecodeByteSlice failed: %w", err)
		}
		result = append(result, val)
		total += n
	}

	return result, total, nil
}

func DecodeStringSlice(d *Decoder) ([]string, int, error) {
	return DecodeStringSliceWithLimit(d, MaxElements)
}

func DecodeStringSliceWithLimit(d *Decoder, limit uint32) ([]string, int, error) {
	sliceOfByteSlices, n, err := DecodeSliceOfByteSliceWithLimit(d, limit)
	if err != nil {
		return nil, 0, err
	}
	if sliceOfByteSlices == nil {
		return nil, 0, nil
	}
	result := make([]string, 0, len(sliceOfByteSlices))
	for i := range sliceOfByteSlices {
		result = append(result, string(sliceOfByteSlices[i]))
	}

	return result, n, nil
}

func DecodeStructArray[V any, H DecodablePtr[V]](d *Decoder, value []V) (int, error) {
	total := 0
	for i := range value {
		n, err := H(&value[i]).DecodeScale(d)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

func DecodeOption[V any, H DecodablePtr[V]](d *Decoder) (*V, int, error) {
	exists, total, err := DecodeBool(d)
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
