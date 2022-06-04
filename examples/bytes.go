package examples

//go:generate scalegen

type Bytes20 struct {
	Value [20]byte
}

type Bytes32 struct {
	Value [32]byte
}

type Bytes64 struct {
	Value [64]byte
}

type Slice struct {
	Value []byte
}
