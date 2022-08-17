package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedStructNoImport struct {
	Value nested.Struct
}
