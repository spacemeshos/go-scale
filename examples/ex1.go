package examples

//go:generate scalegen

type Ex1 struct {
	Option *Ex1
	Bool   bool
}
