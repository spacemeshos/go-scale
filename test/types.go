package test

//go:generate scalegen

type DeepNestedModule struct {
	Value *DeepNestedModule
}

type DeepNestedSliceModule struct {
	Value []DeepNestedSliceModule `scale:"max=2"`
}

type DeepNestedArrayModule struct {
	Value [2]Level1
}

type Level1 struct {
	Value [2]Level2
}

type Level2 struct {
	Value [2]Level3
}

type Level3 struct {
	Value [2]Level4
}

type Level4 struct {
	Value [2]Level5
}

type Level5 struct {
	Value string `scale:"max=64"`
}
