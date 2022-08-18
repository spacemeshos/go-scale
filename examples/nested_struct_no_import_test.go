package examples

import (
	"github.com/spacemeshos/go-scale/tester"
	"testing"
)

func FuzzNestedStructConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStruct](f)
}

func FuzzNestedStructSafety(f *testing.F) {
	tester.FuzzSafety[NestedStruct](f)
}
