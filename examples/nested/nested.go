package nested

//go:generate scalegen

type NestedModule struct {
	Value []byte
}

type Struct struct {
	A uint64
}

type StringAlias string
