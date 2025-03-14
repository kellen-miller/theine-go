//go:build go1.24
// +build go1.24

package hasher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar string
}

func TestHasher_StringKey(t *testing.T) {
	hasher := NewHasher[string](nil)
	h := hasher.Hash(strconv.Itoa(123456))
	for i := 0; i < 10; i++ {
		require.Equal(t, h, hasher.Hash(strconv.Itoa(123456)))
	}
}

func TestHasher_StructStringKey(t *testing.T) {
	hasher1 := NewHasher[Foo](nil)
	hasher2 := NewHasher[Foo](func(k Foo) string {
		return k.Bar
	})
	h1 := uint64(0)
	h2 := uint64(0)
	for i := 0; i < 10; i++ {
		foo := Foo{Bar: strconv.Itoa(123456)}
		if h1 == 0 {
			h1 = hasher1.Hash(foo)
		} else {
			require.Equal(t, h1, hasher1.Hash(foo))
		}
	}
	for i := 0; i < 10; i++ {
		foo := Foo{Bar: strconv.Itoa(123456)}
		if h2 == 0 {
			h2 = hasher2.Hash(foo)
		} else {
			require.Equal(t, h2, hasher2.Hash(foo))
		}
	}
}
