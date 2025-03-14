package theine_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kellen-miller/theine-go"
	"github.com/stretchr/testify/require"
)

func TestPersist_Basic(t *testing.T) {
	client, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	for i := 0; i < 1000; i++ {
		client.Set(i, i, 1)
	}
	f, err := os.Create("ptest")
	defer os.Remove("ptest")
	require.Nil(t, err)
	err = client.SaveCache(0, f)
	require.Nil(t, err)
	f.Close()

	f, err = os.Open("ptest")
	require.Nil(t, err)
	new, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	err = new.LoadCache(0, f)
	require.Nil(t, err)
	f.Close()
	m := map[int]int{}
	new.Range(func(key, value int) bool {
		m[key] = value
		return true
	})
	require.Equal(t, 100, len(m))
	for k, v := range m {
		require.Equal(t, k, v)
	}

}

func TestPersist_LoadingBasic(t *testing.T) {
	client, err := theine.NewBuilder[int, int](100).BuildWithLoader(func(
		ctx context.Context,
		key int,
	) (theine.Loaded[int], error) {
		return theine.Loaded[int]{Value: key, Cost: 1, TTL: 0}, nil
	})
	require.Nil(t, err)
	for i := 0; i < 1000; i++ {
		client.Set(i, i, 1)
	}
	f, err := os.Create("ptest")
	defer os.Remove("ptest")
	require.Nil(t, err)
	err = client.SaveCache(0, f)
	require.Nil(t, err)
	f.Close()

	f, err = os.Open("ptest")
	require.Nil(t, err)
	new, err := theine.NewBuilder[int, int](100).BuildWithLoader(func(
		ctx context.Context,
		key int,
	) (theine.Loaded[int], error) {
		return theine.Loaded[int]{Value: key, Cost: 1, TTL: 0}, nil
	})
	require.Nil(t, err)
	err = new.LoadCache(0, f)
	require.Nil(t, err)
	f.Close()
	m := map[int]int{}
	new.Range(func(key, value int) bool {
		m[key] = value
		return true
	})
	require.Equal(t, 100, len(m))
	for k, v := range m {
		require.Equal(t, k, v)
	}

}

func TestPersist_TestVersionMismatch(t *testing.T) {
	client, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	f, err := os.Create("ptest")
	defer os.Remove("ptest")
	require.Nil(t, err)
	err = client.SaveCache(0, f)
	require.Nil(t, err)
	f.Close()

	f, err = os.Open("ptest")
	require.Nil(t, err)
	new, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	err = new.LoadCache(1, f)
	require.Equal(t, theine.VersionMismatch, err)
}

func TestPersist_TestChecksumMismatch(t *testing.T) {
	client, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	f, err := os.Create("ptest")
	defer os.Remove("ptest")
	require.Nil(t, err)
	err = client.SaveCache(1, f)
	require.Nil(t, err)
	// change file content
	for _, i := range []int64{222} {
		_, err = f.WriteAt([]byte{1, 0, 1, 1}, i)
		require.Nil(t, err)
	}
	f.Close()

	f, err = os.Open("ptest")
	require.Nil(t, err)
	new, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	err = new.LoadCache(1, f)
	require.Equal(t, "checksum mismatch", err.Error())
}

type PStruct struct {
	Id   int
	Name string
	Data []byte
}

func TestPersist_Large(t *testing.T) {
	client, err := theine.NewBuilder[int, PStruct](100000).Build()
	require.Nil(t, err)
	for i := 0; i < 100000; i++ {
		client.Set(i, PStruct{
			Id:   i,
			Name: fmt.Sprintf("struct-%d", i),
			Data: make([]byte, i%1000),
		}, 1)
	}
	require.Equal(t, 100000, client.Len())
	f, err := os.Create("ptest")
	defer os.Remove("ptest")
	require.Nil(t, err)
	time.Sleep(time.Second)
	err = client.SaveCache(0, f)
	require.Nil(t, err)
	f.Close()

	f, err = os.Open("ptest")
	require.Nil(t, err)
	new, err := theine.NewBuilder[int, PStruct](100000).Build()
	require.Nil(t, err)
	err = new.LoadCache(0, f)
	require.Nil(t, err)
	f.Close()
	m := map[int]PStruct{}
	new.Range(func(key int, value PStruct) bool {
		m[key] = value
		return true
	})
	require.Equal(t, 100000, len(m))
	for k, v := range m {
		require.Equal(t, k, v.Id)
		require.Equal(t, fmt.Sprintf("struct-%d", k), v.Name)
		require.Equal(t, k%1000, len(v.Data))
	}
}

func TestPersist_OS(t *testing.T) {
	f, err := os.Open("otest")
	require.Nil(t, err)
	client, err := theine.NewBuilder[int, int](100).Build()
	require.Nil(t, err)
	err = client.LoadCache(0, f)
	require.Nil(t, err)
	f.Close()
	m := map[int]int{}
	client.Range(func(key, value int) bool {
		m[key] = value
		return true
	})
	require.Equal(t, 100, len(m))
	for k, v := range m {
		require.Equal(t, k, v)
	}
}
