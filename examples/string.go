package examples

//go:generate scalegen

type StructWithString struct {
	Value string
}

type StructWithStringLimit struct {
	Value string `scale:"max=3"`
}

type StructWithStringSliceAndLimit struct {
	Value []string `scale:"max=3"`
}

type StringAlias string

type StructWithStringAlias struct {
	Value StringAlias
}

type StructWithStringAliasAndLimit struct {
	Value StringAlias `scale:"max=3"`
}
