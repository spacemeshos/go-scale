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

func FuzzStructWithStringSliceAndLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[StructWithStringSliceAndLimit](f)
}

func FuzzStructWithStringSliceAndLimitSafety(f *testing.F) {
	tester.FuzzSafety[StructWithStringSliceAndLimit](f)
}

func FuzzStructWithStringAliasConsistency(f *testing.F) {
	tester.FuzzConsistency[StructWithStringAlias](f)
}

func FuzzStructWithStringAliasSafety(f *testing.F) {
	tester.FuzzSafety[StructWithStringAlias](f)
}

func FuzzStructWithStringAliasAndLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[StructWithStringAliasAndLimit](f)
}

func FuzzStructWithStringAliasAndLimitSafety(f *testing.F) {
	tester.FuzzSafety[StructWithStringAliasAndLimit](f)
}
