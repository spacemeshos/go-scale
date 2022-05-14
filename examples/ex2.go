package examples

//go:generate scalegen -pkg examples -file ex2_scale.go -types Ex2,Smth -imports github.com/spacemeshos/go-scale/examples

type Ex2 struct {
	Slice []Ex2
	Array [5]Smth
}

type Smth struct {
	Val uint32
}
