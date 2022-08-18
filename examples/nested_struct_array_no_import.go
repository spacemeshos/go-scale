package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedStructArray struct {
	Value [3]nested.Struct
}
