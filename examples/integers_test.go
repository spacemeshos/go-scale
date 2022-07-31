package examples

import (
	"path/filepath"
	"testing"

	"github.com/spacemeshos/go-scale/tester"
	"github.com/stretchr/testify/require"
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

func FuzzUConsistency(f *testing.F) {
	tester.FuzzConsistency[U](f)
}

func FuzzUSafety(f *testing.F) {
	tester.FuzzSafety[U](f)
}

func TestGoldenIntegers(t *testing.T) {
	golden, err := filepath.Abs("./golden")
	require.NoError(t, err)
	t.Run("U8", func(t *testing.T) {
		tester.GoldenTest[U8](t, filepath.Join(golden, "U8.json"))
	})
	t.Run("U16", func(t *testing.T) {
		tester.GoldenTest[U16](t, filepath.Join(golden, "U16.json"))
	})
	t.Run("U32", func(t *testing.T) {
		tester.GoldenTest[U32](t, filepath.Join(golden, "U32.json"))
	})
	t.Run("U64", func(t *testing.T) {
		tester.GoldenTest[U64](t, filepath.Join(golden, "U64.json"))
	})
}
