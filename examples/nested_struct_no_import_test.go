package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzNestedStructConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStruct](f)
}

func FuzzNestedStructSafety(f *testing.F) {
	tester.FuzzSafety[NestedStruct](f)
}
