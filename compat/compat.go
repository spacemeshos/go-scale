package compat

/*
#include <stdlib.h>
#include "./rust/compat.h"
#cgo LDFLAGS: ${SRCDIR}/rust/target/release/libcompat.a
*/
import "C"

import (
	"errors"
	"unsafe"
)

//go:generate scalegen

type Struct struct {
	Field1 uint16
	Field2 [3]byte
}

type Byte32 struct {
	Array [32]byte
}

type Compat struct {
	Field1  uint8
	Field2  uint16
	Field3  uint32
	Field4  uint64
	Field5  [8]byte
	Field6  bool
	Field7  *Struct
	Field8  Struct
	Field9  [4]Struct
	Field10 []byte
	Field11 []Byte32
	Field12 [][]byte
	Field13 []Struct
}

func RoundTrip(input []byte) ([]byte, error) {
	ptr := C.CBytes(input)
	defer C.free(ptr)
	response := C.round_trip((*C.uchar)(ptr), C.size_t(len(input)))
	defer C.free_response(response)
	buf := make([]byte, response.len)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(response.ptr)), response.len))
	if response.code == 1 {
		return nil, errors.New(string(buf))
	}
	return buf, nil
}
