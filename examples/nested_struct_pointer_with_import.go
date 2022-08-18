package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedStructPointer struct {
	Value *nested.Struct
}
