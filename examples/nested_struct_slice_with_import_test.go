package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzNestedStructSliceConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStructSlice](f)
}

func FuzzNestedStructSliceSafety(f *testing.F) {
	tester.FuzzSafety[NestedStructSlice](f)
}
