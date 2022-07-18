package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzStructSliceWithLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[StructSliceWithLimit](f)
}

func FuzzStructSliceWithLimitSafety(f *testing.F) {
	tester.FuzzSafety[StructSliceWithLimit](f)
}
