package nested

//go:generate scalegen

type NestedModule struct {
	Value []byte `scale:"max=32"`
}

type Struct struct {
	A uint64
}

type StringAlias string
