package testdata

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type Data struct {
	Str                 string
	NestedStruct        nested.Struct
	NestedStructPointer *nested.Struct
	NestedStructSlice   []nested.Struct
	NestedAlias         nested.StringAlias
	StrSlice            []string
	ByteArray           [20]byte
	ByteSlice           []byte
	SliceOfByteSlices   [][]byte
	Uint64              uint64
}
