package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

//go:generate scalegen

func FuzzNestedStructArrayConsistency(f *testing.F) {
	tester.FuzzConsistency[NestedStructArray](f)
}

func FuzzNestedStructArraySafety(f *testing.F) {
	tester.FuzzSafety[NestedStructArray](f)
}
