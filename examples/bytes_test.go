package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzBytes20Consistency(f *testing.F) {
	tester.FuzzConsistency[Bytes20](f)
}

func FuzzBytes20Safety(f *testing.F) {
	tester.FuzzSafety[Bytes20](f)
}

func FuzzBytesSliceConsistency(f *testing.F) {
	tester.FuzzConsistency[Slice](f)
}

func FuzzBytesSliceSafety(f *testing.F) {
	tester.FuzzSafety[Slice](f)
}

func FuzzBytesSliceWithLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[SliceWithLimit](f)
}

func FuzzBytesSliceWithLimitSafety(f *testing.F) {
	tester.FuzzSafety[SliceWithLimit](f)
}
