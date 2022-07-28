package examples

//go:generate scalegen

type StructWithString struct {
	Value string
}

type StructWithStringLimit struct {
	Value string `scale:"max=3"`
}
