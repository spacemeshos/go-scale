package test

import (
	"bytes"
	"testing"

	"github.com/spacemeshos/go-scale"
	"github.com/stretchr/testify/require"
)

func Test_DecodeDeepNestedModule(t *testing.T) {
	t.Run("default decoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedModule{})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf)
		_, _, err = scale.DecodeStruct[DeepNestedModule](dec)
		require.NoError(t, err)
	})

	t.Run("default decoder, with 5 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedModule{
			Value: &DeepNestedModule{
				Value: &DeepNestedModule{
					Value: &DeepNestedModule{
						Value: &DeepNestedModule{},
					},
				},
			},
		})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf)
		_, _, err = scale.DecodeStruct[DeepNestedModule](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})

	t.Run("decoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedModule{
			Value: &DeepNestedModule{
				Value: &DeepNestedModule{
					Value: &DeepNestedModule{},
				},
			},
		})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf, scale.WithDecodeMaxNested(2))
		_, _, err = scale.DecodeStruct[DeepNestedModule](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})
}

func Test_DecodeDeepNestedSliceModule(t *testing.T) {
	t.Run("default decoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf)
		_, _, err = scale.DecodeStruct[DeepNestedSliceModule](dec)
		require.NoError(t, err)
	})

	t.Run("default decoder, with 5 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{
			Value: []DeepNestedSliceModule{{
				Value: []DeepNestedSliceModule{{
					Value: []DeepNestedSliceModule{{
						Value: []DeepNestedSliceModule{{}},
					}},
				}},
			}},
		})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf, scale.WithDecodeMaxNested(2))
		_, _, err = scale.DecodeStruct[DeepNestedSliceModule](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})

	t.Run("decoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedSliceModule{
			Value: []DeepNestedSliceModule{{
				Value: []DeepNestedSliceModule{{
					Value: []DeepNestedSliceModule{{}},
				}},
			}},
		})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf, scale.WithDecodeMaxNested(2))
		_, _, err = scale.DecodeStruct[DeepNestedSliceModule](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})
}

func Test_DecodeDeepNestedArrayModule(t *testing.T) {
	t.Run("default decoder, no nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, DeepNestedArrayModule{})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf, scale.WithDecodeMaxNested(2))
		_, _, err = scale.DecodeStruct[DeepNestedArrayModule](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})

	t.Run("decoder with max nested 2, with 3 recursive nesting", func(t *testing.T) {
		var buf bytes.Buffer
		enc := scale.NewEncoder(&buf, scale.WithEncodeMaxNested(10))

		_, err := scale.EncodeStruct(enc, Level2{})
		require.NoError(t, err)

		dec := scale.NewDecoder(&buf, scale.WithDecodeMaxNested(2))
		_, _, err = scale.DecodeStruct[Level2](dec)
		require.ErrorIs(t, scale.ErrDecodeNestedTooDeep, err)
	})
}
