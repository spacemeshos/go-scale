package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedStructSlice struct {
	Value []nested.Struct `scale:"max=5"`
}
