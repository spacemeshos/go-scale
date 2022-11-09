package scale

import (
	"reflect"
	"unsafe"
)

// stringToBytes converts a string to a byte slice without copying the underlying data.
// IMPORTANT: The returned byte slice must not be modified!
// This is a low-level function and should be used carefully.
func stringToBytes(s string) []byte {
	if s == "" {
		return nil
	}

	const max = 0x7fff0000
	if len(s) > max {
		panic("string too long")
	}
	return (*[max]byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))[:len(s):len(s)]
}

// bytesToString converts a byte slice to a string without copying the underlying data.
// This is a low-level function and should be used carefully.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&b))
}
