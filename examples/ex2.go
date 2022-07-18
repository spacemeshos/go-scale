package examples

//go:generate scalegen

type Ex2 struct {
	Slice []Ex2
	Array [5]Smth
}

type Smth struct {
	Val uint32
}

type StructSliceWithLimit struct {
	Slice []Smth `scale:"max=10"`
}
