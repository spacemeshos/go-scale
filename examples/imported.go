package examples

import "github.com/spacemeshos/go-scale/examples/alias"

//go:generate scalegen

type ImportedB struct {
	List []alias.B `scale:"type=StructArray"`
}
