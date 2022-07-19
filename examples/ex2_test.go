package examples

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-scale/tester"
)

func FuzzStructSliceWithLimitConsistency(f *testing.F) {
	tester.FuzzConsistency[StructSliceWithLimit](f)
}

func FuzzStructSliceWithLimitSafety(f *testing.F) {
	tester.FuzzSafety[StructSliceWithLimit](f)
}

func TestStructSliceWithLimitEncodeTooManyElements(t *testing.T) {
	s := StructSliceWithLimit{
		Slice: []Smth{
			{Val: 1},
			{Val: 2},
			{Val: 3},
		},
	}
	buf := bytes.NewBuffer(nil)
	encoder := scale.NewEncoder(buf)
	_, err := s.EncodeScale(encoder)
	require.ErrorIs(t, err, scale.ErrEncodeTooManyElements)
}

func TestStructSliceWithLimitDecodeTooManyElements(t *testing.T) {
	var structSliceWith3ElementsHexStr = "0c04080c"
	buf, err := hex.DecodeString(structSliceWith3ElementsHexStr)
	require.NoError(t, err)
	decoder := scale.NewDecoder(bytes.NewReader(buf))
	var s StructSliceWithLimit
	_, err = s.DecodeScale(decoder)
	require.ErrorIs(t, err, scale.ErrDecodeTooManyElements)
}
