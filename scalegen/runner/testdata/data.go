package testdata

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type Data struct {
	Str                 string `scale:"max=20"`
	NestedStruct        nested.Struct
	NestedStructPointer *nested.Struct
	NestedStructSlice   []nested.Struct `scale:"max=5"`
}

type MoreData struct {
	NestedAlias       nested.StringAlias `scale:"max=20"`
	StrSlice          []string           `scale:"max=5"`
	ByteArray         [20]byte
	ByteSlice         []byte   `scale:"max=20"`
	SliceOfByteSlices [][]byte `scale:"max=10"`
	Uint64            uint64
}
