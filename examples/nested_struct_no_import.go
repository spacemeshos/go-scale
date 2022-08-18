package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedStruct struct {
	Value nested.Struct
}
