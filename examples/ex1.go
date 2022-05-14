package examples

//go:generate scalegen -pkg examples -file ex1_scale.go -types Ex1 -imports github.com/spacemeshos/go-scale/examples

type Ex1 struct {
	Option *Ex1
	Bool   bool
}
