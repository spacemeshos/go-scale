package examples

import (
	"bytes"
	"github.com/spacemeshos/go-scale"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/spacemeshos/go-scale/tester"
)

func FuzzBytes20Consistency(f *testing.F) {
	tester.FuzzConsistency[Bytes20](f)
}

func FuzzBytes20Safety(f *testing.F) {
	tester.FuzzSafety[Bytes20](f)
}

func FuzzBytesSliceConsistency(f *testing.F) {
	tester.FuzzConsistency[Slice](f)
}

func FuzzBytesSliceSafety(f *testing.F) {
	tester.FuzzSafety[Slice](f)
}

func FuzzBytesSliceWithLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[SliceWithLimit](f)
}

func FuzzBytesSliceWithLimitSafety(f *testing.F) {
	tester.FuzzSafety[SliceWithLimit](f)
}

func FuzzSliceOfByteSliceWithLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[SliceOfByteSliceWithLimit](f)
}

func FuzzSliceOfByteSliceWithLimitSafety(f *testing.F) {
	tester.FuzzSafety[SliceOfByteSliceWithLimit](f)
}

func TestNilByteSlice(t *testing.T) {
	s := Slice{
		Value: nil,
	}
	buf := bytes.NewBuffer(nil)
	encoder := scale.NewEncoder(buf)
	_, err := s.EncodeScale(encoder)
	require.ErrorIs(t, err, scale.ErrNilSlice)
}
