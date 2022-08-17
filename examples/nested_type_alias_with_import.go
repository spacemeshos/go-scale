package examples

import "github.com/spacemeshos/go-scale/examples/nested"

//go:generate scalegen

type NestedTypeAliasWithImport struct {
	Value nested.StringAlias
}
