package examples

import (
	"github.com/spacemeshos/go-scale/tester"
	"testing"
)

func FuzzNestedStructPointerConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStructPointer](f)
}

func FuzzNestedStructPointerSafety(f *testing.F) {
	tester.FuzzSafety[NestedStructPointer](f)
}
