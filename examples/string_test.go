package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzStructWithStringConsistency(f *testing.F) {
	tester.FuzzConsistency[StructWithString](f)
}

func FuzzStructWithStringSafety(f *testing.F) {
	tester.FuzzSafety[StructWithString](f)
}
