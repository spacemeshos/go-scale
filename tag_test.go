package scale

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Name  string `scale:"max=64"`
	Items []int  `scale:"max=4242"`
	Bad   []int  `scale:"max=4242x"`
	Bad1  []int  `scale:"abc=4242"`
	NoTag []int
}

func TestGetMaxElements(t *testing.T) {
	n, err := GetMaxElements[Foo]("Name")
	require.NoError(t, err)
	require.Equal(t, uint32(64), n)
	n, err = GetMaxElements[Foo]("Items")
	require.NoError(t, err)
	require.Equal(t, uint32(4242), n)
	n, err = GetMaxElements[*Foo]("Name")
	require.NoError(t, err)
	require.Equal(t, uint32(64), n)
	n, err = GetMaxElements[*Foo]("Items")
	require.NoError(t, err)
	require.Equal(t, uint32(4242), n)
	require.Equal(t, uint32(64), MustGetMaxElements[Foo]("Name"))
	require.Equal(t, uint32(4242), MustGetMaxElements[Foo]("Items"))

	_, err = GetMaxElements[Foo]("NoSuchField")
	require.Error(t, err)
	_, err = GetMaxElements[int]("Name")
	require.Error(t, err)
	_, err = GetMaxElements[int]("Bad")
	require.Error(t, err)
	_, err = GetMaxElements[int]("Bad1")
	require.Error(t, err)
	_, err = GetMaxElements[int]("NoTag")
	require.Error(t, err)
}
