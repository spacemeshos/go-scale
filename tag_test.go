package scale

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Name   string `scale:"max=64"`
	Items  []int  `scale:"max=4242"`
	Bad    []int  `scale:"max=4242x"`
	Bad1   []int  `scale:"abc=4242"`
	NoTag  []int
	Name1  string `scale:"nonlocal"`
	Items1 []int  `scale:"max=4243,nonlocal"`
	Items2 []int  `scale:"nonlocal,max=4244"`
}

func TestGetMaxElements(t *testing.T) {
	n, err := MaxScaleElements[Foo]("Name")
	require.NoError(t, err)
	require.Equal(t, uint32(64), n)
	n, err = MaxScaleElements[Foo]("Items")
	require.NoError(t, err)
	require.Equal(t, uint32(4242), n)
	n, err = MaxScaleElements[*Foo]("Name")
	require.NoError(t, err)
	require.Equal(t, uint32(64), n)
	n, err = MaxScaleElements[*Foo]("Items")
	require.NoError(t, err)
	require.Equal(t, uint32(4242), n)
	n, err = MaxScaleElements[*Foo]("Items1")
	require.NoError(t, err)
	require.Equal(t, uint32(4243), n)
	n, err = MaxScaleElements[*Foo]("Items2")
	require.NoError(t, err)
	require.Equal(t, uint32(4244), n)
	require.Equal(t, uint32(64), MustGetMaxElements[Foo]("Name"))
	require.Equal(t, uint32(4242), MustGetMaxElements[Foo]("Items"))
	require.Equal(t, uint32(4243), MustGetMaxElements[Foo]("Items1"))
	require.Equal(t, uint32(4244), MustGetMaxElements[Foo]("Items2"))

	_, err = MaxScaleElements[Foo]("NoSuchField")
	require.Error(t, err)
	_, err = MaxScaleElements[int]("Name")
	require.Error(t, err)
	_, err = MaxScaleElements[int]("Bad")
	require.Error(t, err)
	_, err = MaxScaleElements[int]("Bad1")
	require.Error(t, err)
	_, err = MaxScaleElements[int]("NoTag")
	require.Error(t, err)
}

func TestNonLocal(t *testing.T) {
	nameField, found := reflect.TypeOf(Foo{}).FieldByName("Name")
	require.True(t, found)
	nl, err := nonLocal(nameField.Tag)
	require.NoError(t, err)
	require.False(t, nl)

	name1Field, found := reflect.TypeOf(Foo{}).FieldByName("Name1")
	require.True(t, found)
	nl, err = nonLocal(name1Field.Tag)
	require.NoError(t, err)
	require.True(t, nl)

	itemsField, found := reflect.TypeOf(Foo{}).FieldByName("Items")
	require.True(t, found)
	nl, err = nonLocal(itemsField.Tag)
	require.NoError(t, err)
	require.False(t, nl)

	items1Field, found := reflect.TypeOf(Foo{}).FieldByName("Items1")
	require.True(t, found)
	nl, err = nonLocal(items1Field.Tag)
	require.NoError(t, err)
	require.True(t, nl)

	items2Field, found := reflect.TypeOf(Foo{}).FieldByName("Items2")
	require.True(t, found)
	nl, err = nonLocal(items2Field.Tag)
	require.NoError(t, err)
	require.True(t, nl)
}
