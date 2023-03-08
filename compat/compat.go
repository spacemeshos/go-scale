package compat

/*
#include <stdlib.h>
#include "rust/compat.h"
#cgo LDFLAGS: rust/target/release/libcompat.a -ldl
*/
import "C"

import (
	"errors"
	"unsafe"
)

func RoundTrip(input []byte) ([]byte, error) {
	ptr := C.CBytes(input)
	defer C.free(ptr)
	response := C.round_trip((*C.uchar)(ptr), C.size_t(len(input)))
	buf := make([]byte, response.len)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(response.ptr)), response.len))
	if response.code == 1 {
		return nil, errors.New(string(buf))
	}
	return buf, nil
}
