package testdata

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type Data struct {
	Str                 string `scale:"max=20"`
	NestedStruct        nested.Struct
	NestedStructPointer *nested.Struct
	NestedStructSlice   []nested.Struct `scale:"max=5"`
}

type Name struct {
	Value string `scale:"max=20"`
}

type MoreData struct {
	NestedAlias nested.StringAlias `scale:"max=20"`
	StrSlice    []Name             `scale:"max=5"`
	ByteArray   [20]byte
	ByteSlice   []byte `scale:"max=20"`
	Uint64      uint64
}
