package scale

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type newU16 uint16

func (newU16) EncodeScale(enc *Encoder) (int, error) {
	panic("uninmplemented")
}

type newU32 uint32

func (newU32) EncodeScale(enc *Encoder) (int, error) {
	panic("uninmplemented")
}

type newU64 uint64

func (newU64) EncodeScale(enc *Encoder) (int, error) {
	panic("uninmplemented")
}

func Test_getScaleType_Slices(t *testing.T) {
	t.Run("[]uint16", func(t *testing.T) {
		type Foo struct {
			Slice []uint16 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint16SliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint16 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU16 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint16 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint16
		type Foo struct {
			Slice []newT `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint16SliceWithLimit", scaleT.Name)
	})
	t.Run("[]uint32", func(t *testing.T) {
		type Foo struct {
			Slice []uint32 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint32SliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint32 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU32 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint32 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint32
		type Foo struct {
			Slice []newT `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint32SliceWithLimit", scaleT.Name)
	})
	t.Run("[]uint64", func(t *testing.T) {
		type Foo struct {
			Slice []uint64 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint64SliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint64 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU64 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
	})
	t.Run("[]newUint64 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint64
		type Foo struct {
			Slice []newT `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint64SliceWithLimit", scaleT.Name)
	})
}
