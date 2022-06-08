package examples

import (
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzOptionsConsistency(f *testing.F) {
	tester.FuzzConsistency[Options](f)
}

func FuzzOptionsSafety(f *testing.F) {
	tester.FuzzSafety[Options](f)
}
