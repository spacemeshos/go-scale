package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzNestedStructPointerConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStructPointer](f)
}

func FuzzNestedStructPointerSafety(f *testing.F) {
	tester.FuzzSafety[NestedStructPointer](f)
}
