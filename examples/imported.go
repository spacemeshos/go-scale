package examples

import "github.com/spacemeshos/go-scale/examples/alias"

//go:generate scalegen

type ImportedA struct {
	ListA []alias.A `scale:"type=StructArray"`
	ListB []alias.B
}
