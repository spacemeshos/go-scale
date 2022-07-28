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

func FuzzStructWithStringLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[StructWithStringLimit](f)
}

func FuzzStructWithStringLimitSafety(f *testing.F) {
	tester.FuzzSafety[StructWithStringLimit](f)
}
