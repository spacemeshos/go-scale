package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzU8Consistency(f *testing.F) {
	tester.FuzzConsistency[U8](f)
}

func FuzzU8Safety(f *testing.F) {
	tester.FuzzSafety[U8](f)
}

func FuzzU16Consistency(f *testing.F) {
	tester.FuzzConsistency[U16](f)
}

func FuzzU16Safety(f *testing.F) {
	tester.FuzzSafety[U16](f)
}

func FuzzU32Consistency(f *testing.F) {
	tester.FuzzConsistency[U32](f)
}

func FuzzU32Safety(f *testing.F) {
	tester.FuzzSafety[U32](f)
}

func FuzzU64Consistency(f *testing.F) {
	tester.FuzzConsistency[U64](f)
}

func FuzzU64Safety(f *testing.F) {
	tester.FuzzSafety[U64](f)
}
