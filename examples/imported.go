package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type Imported struct {
	B []nested.Struct
}
