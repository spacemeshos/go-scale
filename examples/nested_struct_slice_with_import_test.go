package examples

import (
	"github.com/spacemeshos/go-scale/tester"
	"testing"
)

func FuzzNestedStructSliceConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStructSlice](f)
}

func FuzzNestedStructSliceSafety(f *testing.F) {
	tester.FuzzSafety[NestedStructSlice](f)
}
