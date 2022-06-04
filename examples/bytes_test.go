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
