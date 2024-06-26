package scale

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/bits"
)

const (
	// MaxElements is the maximum number of elements allowed in a collection if not set explicitly during
	// encoding/decoding.
	MaxElements uint32 = 1 << 20

	// MaxNested is the maximum nested level allowed if not set explicitly during encoding/decoding.
	MaxNested uint = 4
)

var (
	// ErrEncodeTooManyElements is returned when scale limit tag is used and collection has too many elements to encode.
	ErrEncodeTooManyElements = errors.New("too many elements to encode in collection with scale limit set")

	// ErrEncodeNestedTooDeep is returned when the depth of nested types exceeds the limit.
	ErrEncodeNestedTooDeep = errors.New("nested level is too deep")
)

type Encodable interface {
	EncodeScale(enc *Encoder) (int, error)
}

type EncodablePtr[B any] interface {
	Encodable
	*B
}

type encoderOpts func(*Encoder)

// WithEncodeMaxNested sets the nested level of the encoder.
// A value of 0 means no nesting is allowed. The default value is 4.
func WithEncodeMaxNested(nested uint) encoderOpts {
	return func(e *Encoder) {
		e.maxNested = nested
	}
}

// WithEncodeMaxElements sets the maximum number of elements allowed in a collection.
// The default value is 1 << 20.
func WithEncodeMaxElements(elements uint32) encoderOpts {
	return func(e *Encoder) {
		e.maxElements = elements
	}
}

// NewEncoder returns a new encoder that writes to w.
// If w implements io.StringWriter, the returned encoder will be more efficient in encoding strings.
func NewEncoder(w io.Writer, opts ...encoderOpts) *Encoder {
	e := &Encoder{
		w:           w,
		maxNested:   MaxNested,
		maxElements: MaxElements,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

type Encoder struct {
	w           io.Writer
	scratch     [9]byte
	maxNested   uint
	maxElements uint32
}

func (e *Encoder) enterNested() error {
	if e.maxNested == 0 {
		return ErrEncodeNestedTooDeep
	}
	e.maxNested--
	return nil
}

func (e *Encoder) leaveNested() {
	e.maxNested++
}

func EncodeByteSlice(e *Encoder, value []byte) (int, error) {
	return EncodeByteSliceWithLimit(e, value, e.maxElements)
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

func EncodeUint16Slice(e *Encoder, value []uint16) (int, error) {
	return EncodeUint16SliceWithLimit(e, value, e.maxElements)
}

func EncodeUint16SliceWithLimit(e *Encoder, value []uint16, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	for _, v := range value {
		n, err := EncodeCompact16(e, v)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}

func EncodeUint32Slice(e *Encoder, value []uint32) (int, error) {
	return EncodeUint32SliceWithLimit(e, value, e.maxElements)
}

func EncodeUint32SliceWithLimit(e *Encoder, value []uint32, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	for _, v := range value {
		n, err := EncodeCompact32(e, v)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}

func EncodeUint64Slice(e *Encoder, value []uint64) (int, error) {
	return EncodeUint64SliceWithLimit(e, value, e.maxElements)
}

func EncodeUint64SliceWithLimit(e *Encoder, value []uint64, limit uint32) (int, error) {
	total, err := EncodeLen(e, uint32(len(value)), limit)
	if err != nil {
		return 0, err
	}
	for _, v := range value {
		n, err := EncodeCompact64(e, v)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}

func EncodeString(e *Encoder, value string) (int, error) {
	return EncodeStringWithLimit(e, value, e.maxElements)
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
	return EncodeStructSliceWithLimit[V, H](e, value, e.maxElements)
}

func EncodeStructSliceWithLimit[V any, H EncodablePtr[V]](e *Encoder, value []V, limit uint32) (int, error) {
	if err := e.enterNested(); err != nil {
		return 0, err
	}
	defer e.leaveNested()
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

func EncodeStringSlice(e *Encoder, value []string) (int, error) {
	return EncodeStringSliceWithLimit(e, value, e.maxElements)
}

func EncodeStringSliceWithLimit(e *Encoder, value []string, limit uint32) (int, error) {
	valueToBytes := make([][]byte, 0, len(value))
	for i := range value {
		valueToBytes = append(valueToBytes, stringToBytes(value[i]))
	}
	total, err := EncodeLen(e, uint32(len(valueToBytes)), limit)
	if err != nil {
		return 0, fmt.Errorf("EncodeLen failed: %w", err)
	}
	for _, byteSlice := range valueToBytes {
		n, err := EncodeByteSliceWithLimit(e, byteSlice, e.maxElements)
		if err != nil {
			return 0, fmt.Errorf("EncodeByteSliceWithLimit failed: %w", err)
		}
		total += n
	}
	return total, nil
}

func EncodeStructArray[V any, H EncodablePtr[V]](e *Encoder, value []V) (int, error) {
	if err := e.enterNested(); err != nil {
		return 0, err
	}
	defer e.leaveNested()
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

func encodeUint8[V ~uint8 | ~uint16 | ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	return e.w.Write(e.scratch[:1])
}

func encodeUint16[V ~uint16 | ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	e.scratch[1] = byte(v >> 8)
	return e.w.Write(e.scratch[:2])
}

func encodeUint32[V ~uint32 | ~uint64](e *Encoder, v V) (int, error) {
	e.scratch[0] = byte(v)
	e.scratch[1] = byte(v >> 8)
	e.scratch[2] = byte(v >> 16)
	e.scratch[3] = byte(v >> 24)
	return e.w.Write(e.scratch[:4])
}

func encodeBigUint(e *Encoder, v uint64) (int, error) {
	needed := 8 - bits.LeadingZeros64(v)/8
	e.scratch[0] = byte(needed-4)<<2 | 0b11
	for i := 1; i <= needed; i++ {
		e.scratch[i] = byte(v)
		v >>= 8
	}
	return e.w.Write(e.scratch[:needed+1])
}

func EncodeCompact8(e *Encoder, v uint8) (int, error) {
	if v <= maxUint6 {
		return encodeUint8(e, v<<2)
	}
	return encodeUint16(e, uint16(v)<<2|0b01)
}

func EncodeCompact16(e *Encoder, v uint16) (int, error) {
	if v <= maxUint6 {
		return encodeUint8(e, v<<2)
	} else if v <= maxUint14 {
		return encodeUint16(e, v<<2|0b01)
	}
	return encodeUint32(e, uint32(v)<<2|0b10)
}

func EncodeCompact32(e *Encoder, v uint32) (int, error) {
	if v <= maxUint6 {
		return encodeUint8(e, v<<2)
	} else if v <= maxUint14 {
		return encodeUint16(e, v<<2|0b01)
	} else if v <= maxUint30 {
		return encodeUint32(e, v<<2|0b10)
	}
	return encodeBigUint(e, uint64(v))
}

func EncodeCompact64(e *Encoder, v uint64) (int, error) {
	if v <= maxUint6 {
		return encodeUint8(e, v<<2)
	} else if v <= maxUint14 {
		return encodeUint16(e, v<<2|0b01)
	} else if v <= maxUint30 {
		return encodeUint32(e, v<<2|0b10)
	}
	return encodeBigUint(e, uint64(v))
}

func EncodeLen(e *Encoder, v, limit uint32) (int, error) {
	if v > limit {
		return 0, fmt.Errorf("%w: %d", ErrEncodeTooManyElements, limit)
	}
	return EncodeCompact32(e, v)
}

func EncodeOption[V any, H EncodablePtr[V]](e *Encoder, value *V) (int, error) {
	if err := e.enterNested(); err != nil {
		return 0, err
	}
	defer e.leaveNested()
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

func EncodeCompact8Ptr(e *Encoder, value *uint8) (int, error) {
	if value == nil {
		return EncodeBool(e, false)
	}
	total, err := EncodeBool(e, true)
	if err != nil {
		return 0, err
	}
	n, err := EncodeCompact8(e, *value)
	if err != nil {
		return 0, err
	}
	return total + n, err
}

func EncodeCompact16Ptr(e *Encoder, value *uint16) (int, error) {
	if value == nil {
		return EncodeBool(e, false)
	}
	total, err := EncodeBool(e, true)
	if err != nil {
		return 0, err
	}
	n, err := EncodeCompact16(e, *value)
	if err != nil {
		return 0, err
	}
	return total + n, err
}

func EncodeCompact32Ptr(e *Encoder, value *uint32) (int, error) {
	if value == nil {
		return EncodeBool(e, false)
	}
	total, err := EncodeBool(e, true)
	if err != nil {
		return 0, err
	}
	n, err := EncodeCompact32(e, *value)
	if err != nil {
		return 0, err
	}
	return total + n, err
}

func EncodeCompact64Ptr(e *Encoder, value *uint64) (int, error) {
	if value == nil {
		return EncodeBool(e, false)
	}
	total, err := EncodeBool(e, true)
	if err != nil {
		return 0, err
	}
	n, err := EncodeCompact64(e, *value)
	if err != nil {
		return 0, err
	}
	return total + n, err
}

func EncodeStruct[V any, H EncodablePtr[V]](e *Encoder, value V) (int, error) {
	n, err := H(&value).EncodeScale(e)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// LenSize returns the size in bytes required to encode a length value.
func LenSize(v uint32) uint32 {
	switch {
	case v <= maxUint6:
		return 1
	case v <= maxUint14:
		return 2
	case v <= maxUint30:
		return 4
	default:
		return 5
	}
}
