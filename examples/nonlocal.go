package examples

//go:generate scalegen

type StructWithNonLocalField struct {
	Name   string `scale:"max=20"`
	SomeID string `scale:"nonlocal,max=20"`
}
