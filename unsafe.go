package scale

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}

	const max = 0x7fff0000
	if len(s) > max {
		panic("string too long")
	}
	return (*[max]byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))[:len(s):len(s)]
}

func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&b))
}
