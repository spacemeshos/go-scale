package scale

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type newU8 uint8

func (newU8) EncodeScale(enc *Encoder) (int, error) {
	panic("uninmplemented")
}

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

func Test_getScaleType_NonLocal(t *testing.T) {
	type Foo struct {
		Name string `scale:"nonlocal,max=64"`
	}

	rtype := reflect.TypeOf(Foo{})
	scaleT, err := getScaleType(rtype, rtype.Field(0))
	require.NoError(t, err)
	require.Equal(t, "StringWithLimit", scaleT.Name)
	require.True(t, scaleT.NonLocal)
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
		require.False(t, scaleT.NonLocal)
	})
	t.Run("[]newUint16 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU16 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
		require.False(t, scaleT.NonLocal)
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
		require.False(t, scaleT.NonLocal)
	})
	t.Run("[]uint32", func(t *testing.T) {
		type Foo struct {
			Slice []uint32 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint32SliceWithLimit", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("[]newUint32 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU32 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
		require.False(t, scaleT.NonLocal)
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
		require.False(t, scaleT.NonLocal)
	})
	t.Run("[]uint64", func(t *testing.T) {
		type Foo struct {
			Slice []uint64 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Uint64SliceWithLimit", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("[]newUint64 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Slice []newU64 `scale:"max=2"`
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "StructSliceWithLimit", scaleT.Name)
		require.False(t, scaleT.NonLocal)
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
		require.False(t, scaleT.NonLocal)
	})
}

func Test_getScaleTypePtrs(t *testing.T) {
	t.Run("*uint8", func(t *testing.T) {
		type Foo struct {
			Ptr *uint8
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact8Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU8 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Ptr *newU8
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Option", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU8 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint8
		type Foo struct {
			Ptr *newT
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact8Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*uint16", func(t *testing.T) {
		type Foo struct {
			Ptr *uint16
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact16Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU16 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Ptr *newU16
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Option", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU16 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint16
		type Foo struct {
			Ptr *newT
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact16Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*uint32", func(t *testing.T) {
		type Foo struct {
			Ptr *uint32
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact32Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU32 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Ptr *newU32
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Option", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU32 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint32
		type Foo struct {
			Ptr *newT
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact32Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*uint64", func(t *testing.T) {
		type Foo struct {
			Ptr *uint64
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact64Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU64 (implements Encodable)", func(t *testing.T) {
		type Foo struct {
			Ptr *newU64
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Option", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
	t.Run("*newU64 (doesn't implement Encodable)", func(t *testing.T) {
		type newT uint64
		type Foo struct {
			Ptr *newT
		}

		rtype := reflect.TypeOf(Foo{})
		scaleT, err := getScaleType(rtype, rtype.Field(0))
		require.NoError(t, err)
		require.Equal(t, "Compact64Ptr", scaleT.Name)
		require.False(t, scaleT.NonLocal)
	})
}
