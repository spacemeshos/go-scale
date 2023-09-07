package scale

import (
	"unsafe"
)

// stringToBytes converts a string to a byte slice without copying the underlying data.
// IMPORTANT: The returned byte slice must not be modified!
// This is a low-level function and should be used carefully.
func stringToBytes(s string) (b []byte) {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// bytesToString converts a byte slice to a string without copying the underlying data.
// This is a low-level function and should be used carefully.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
