package test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-scale"
)

func Test_EncodeDeepNestedModule(t *testing.T) {
	t.Run("default encoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf)

		_, err := scale.EncodeStruct(enc, DeepNestedModule{})
		require.NoError(t, err)
	})

	t.Run("default encoder, with 5 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf)

		_, err := scale.EncodeStruct(enc, DeepNestedModule{
			Value: &DeepNestedModule{
				Value: &DeepNestedModule{
					Value: &DeepNestedModule{
						Value: &DeepNestedModule{},
					},
				},
			},
		})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})

	t.Run("encoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(2))

		_, err := scale.EncodeStruct(enc, DeepNestedModule{
			Value: &DeepNestedModule{
				Value: &DeepNestedModule{
					Value: &DeepNestedModule{},
				},
			},
		})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})
}

func Test_EncodeDeepNestedSliceModule(t *testing.T) {
	t.Run("default encoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf)

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{})
		require.NoError(t, err)
	})

	t.Run("default encoder, with 5 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf)

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{
			Value: []DeepNestedSliceModule{{
				Value: []DeepNestedSliceModule{{
					Value: []DeepNestedSliceModule{{
						Value: []DeepNestedSliceModule{{}},
					}},
				}},
			}},
		})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})

	t.Run("encoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(2))

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{
			Value: []DeepNestedSliceModule{{
				Value: []DeepNestedSliceModule{{
					Value: []DeepNestedSliceModule{{}},
				}},
			}},
		})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})
}

func Test_EncodeDeepNestedArrayModule(t *testing.T) {
	t.Run("default encoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf)

		_, err := scale.EncodeStruct(enc, DeepNestedArrayModule{})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})

	t.Run("encoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(2))

		_, err := scale.EncodeStruct(enc, Level2{})
		require.ErrorIs(t, err, scale.ErrEncodeNestedTooDeep)
	})
}
