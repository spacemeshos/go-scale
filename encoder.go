package scale

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/bits"
)

var (
	MaxElements uint32 = 1 << 20
)

var (
	// ErrEncodeTooManyElements is returned when scale limit tag is used and collection has too many elements to encode.
	ErrEncodeTooManyElements = errors.New("too many elements to encode in collection with scale limit set")
)

const (
	// 0b00 | value
	zerozero = 63
	// 0b01 | value << 2
	zeroone = 16383
	// 0b10 | value << 2
	onezero = 1073741823
	// oneone is a big integer mode
)

type Encodable interface {
	EncodeScale(*Encoder) (int, error)
}

type EncodablePtr[B any] interface {
	Encodable
	*B
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

type Encoder struct {
	w       io.Writer
	scratch [9]byte
}

func EncodeByteSlice(e *Encoder, value []byte) (int, error) {
	return EncodeByteSliceWithLimit(e, value, MaxElements)
}

func EncodeByteSliceWithLimit(e *Encoder, value []byte, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	n, err := EncodeByteArray(e, value)
	if err != nil {
		return 0, err
	}
	return total + n, nil
}

func EncodeByteArray(e *Encoder, value []byte) (int, error) {
	return e.w.Write(value)
}

func EncodeString(e *Encoder, value string) (int, error) {
	return EncodeStringWithLimit(e, value, MaxElements)
}

func EncodeStringWithLimit(e *Encoder, value string, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	n, err := io.WriteString(e.w, value)
	if err != nil {
		return 0, err
	}
	return total + n, nil
}

func EncodeStructSlice[V any, H EncodablePtr[V]](e *Encoder, value []V) (int, error) {
	return EncodeStructSliceWithLimit[V, H](e, value, MaxElements)
}

func EncodeStructSliceWithLimit[V any, H EncodablePtr[V]](e *Encoder, value []V, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	for i := range value {
		n, err := H(&value[i]).EncodeScale(e)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

func EncodeSliceOfByteSlice(e *Encoder, value [][]byte) (int, error) {
	return EncodeSliceOfByteSliceWithLimit(e, value, MaxElements)
}

func EncodeSliceOfByteSliceWithLimit(e *Encoder, value [][]byte, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, fmt.Errorf("failed encoding len for a slice of byte slices: %w", err)
	}
	for _, byteSlice := range value {
		n, err := EncodeLen(e, uint32(len(byteSlice)), MaxElements)
		if err != nil {
			return 0, fmt.Errorf("failed encoding len for a byte slice: %w", err)
		}
		total += n
		n, err = EncodeByteSliceWithLimit(e, byteSlice, MaxElements)
		if err != nil {
			return 0, fmt.Errorf("failed encoding byte slice: %w", err)
		}
		total += n
	}
	return total, nil
}

func EncodeStructArray[V any, H EncodablePtr[V]](e *Encoder, value []V) (int, error) {
	total := 0
	for i := range value {
		n, err := H(&value[i]).EncodeScale(e)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

func EncodeBool(e *Encoder, value bool) (int, error) {
	if value {
		e.scratch[0] = 1
	} else {
		e.scratch[0] = 0
	}
	return e.w.Write(e.scratch[:1])
}

func EncodeByte(e *Encoder, value byte) (int, error) {
	e.scratch[0] = value
	return e.w.Write(e.scratch[:1])
}

func EncodeUint16(e *Encoder, value uint16) (int, error) {
	binary.LittleEndian.PutUint16(e.scratch[:2], value)
	return e.w.Write(e.scratch[:2])
}

func EncodeUint32(e *Encoder, value uint32) (int, error) {
	binary.LittleEndian.PutUint32(e.scratch[:4], value)
	return e.w.Write(e.scratch[:4])
}

func EncodeUint64(e *Encoder, value uint64) (int, error) {
	binary.LittleEndian.PutUint64(e.scratch[:8], value)
	return e.w.Write(e.scratch[:8])
}

func encodeZeroZero[V ~uint8 | ~uint16 | ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	return e.w.Write(e.scratch[:1])
}

func encodeZeroOne[V ~uint16 | ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	e.scratch[1] = byte(v >> 8)
	return e.w.Write(e.scratch[:2])
}

func encodeOneZero[V ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	e.scratch[1] = byte(v >> 8)
	e.scratch[2] = byte(v >> 16)
	e.scratch[3] = byte(v >> 24)
	return e.w.Write(e.scratch[:4])
}

func encodeOneOne(e *Encoder, v uint64) (int, error) {
	needed := 8 - bits.LeadingZeros64(v)/8
	e.scratch[0] = byte(needed)<<2 | 0b11
	for i := 1; i <= needed; i++ {
		e.scratch[i] = byte(v)
		v >>= 8
	}
	return e.w.Write(e.scratch[:needed+1])
}

func EncodeCompact8(e *Encoder, v uint8) (int, error) {
	if v <= zerozero {
		return encodeZeroZero(e, v<<2)
	}
	return encodeZeroOne(e, uint16(v)<<2|0b01)
}

func EncodeCompact16(e *Encoder, v uint16) (int, error) {
	if v <= zerozero {
		return encodeZeroZero(e, v<<2)
	} else if v <= zeroone {
		return encodeZeroOne(e, v<<2|0b01)
	}
	return encodeOneZero(e, uint32(v)<<2|0b01)
}

func EncodeCompact32(e *Encoder, v uint32) (int, error) {
	if v <= zerozero {
		return encodeZeroZero(e, v<<2)
	} else if v <= zeroone {
		return encodeZeroOne(e, v<<2|0b01)
	} else if v <= onezero {
		return encodeOneZero(e, v<<2|0b10)
	}
	return encodeOneOne(e, uint64(v))
}

func EncodeCompact64(e *Encoder, v uint64) (int, error) {
	if v <= zerozero {
		return encodeZeroZero(e, v<<2)
	} else if v <= zeroone {
		return encodeZeroOne(e, v<<2|0b01)
	} else if v <= onezero {
		return encodeOneZero(e, v<<2|0b10)
	}
	return encodeOneOne(e, uint64(v))
}

func EncodeLen(e *Encoder, v uint32, limit uint32) (int, error) {
	if v > limit {
		return 0, fmt.Errorf("%w: %d", ErrEncodeTooManyElements, limit)
	}
	return EncodeCompact32(e, v)
}

func EncodeOption[V any, H EncodablePtr[V]](e *Encoder, value *V) (int, error) {
	if value == nil {
		return EncodeBool(e, false)
	}
	total, err := EncodeBool(e, true)
	if err != nil {
		return 0, err
	}
	n, err := H(value).EncodeScale(e)
	if err != nil {
		return 0, err
	}
	return total + n, nil
}

func EncodeStruct[V any, H EncodablePtr[V]](e *Encoder, value V) (int, error) {
	n, err := H(&value).EncodeScale(e)
	if err != nil {
		return 0, err
	}
	return n, nil
}
