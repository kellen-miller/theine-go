package theine_test

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Yiling-J/theine-go"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	client, err := theine.New[string, string](1000)
	require.Nil(t, err)
	for i := 0; i < 20000; i++ {
		key := fmt.Sprintf("key:%d", rand.Intn(100000))
		client.Set(key, key, 1)
	}
	time.Sleep(300 * time.Millisecond)
	require.True(t, client.Len() < 1200)
	client.Close()
}

func TestSetParallel(t *testing.T) {
	client, err := theine.New[string, string](1000)
	require.Nil(t, err)
	var wg sync.WaitGroup
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key:%d", rand.Intn(100000))
				client.Set(key, key, 1)
			}

		}()
	}
	wg.Wait()
	time.Sleep(300 * time.Millisecond)
	require.True(t, client.Len() < 1200)
}

func TestGetSet(t *testing.T) {
	client, err := theine.New[string, string](1000)
	require.Nil(t, err)
	for i := 0; i < 20000; i++ {
		key := fmt.Sprintf("key:%d", rand.Intn(3000))
		v, ok := client.Get(key)
		if !ok {
			client.Set(key, key, 1)
		} else {
			require.Equal(t, key, v)
		}
	}
	time.Sleep(300 * time.Millisecond)
	require.True(t, client.Len() < 1200)
}

func TestDelete(t *testing.T) {
	client, err := theine.New[string, string](100)
	require.Nil(t, err)
	client.Set("foo", "foo", 1)
	v, ok := client.Get("foo")
	require.True(t, ok)
	require.Equal(t, "foo", v)
	client.Delete("foo")
	_, ok = client.Get("foo")
	require.False(t, ok)

	client.SetWithTTL("foo", "foo", 1, 10*time.Second)
	v, ok = client.Get("foo")
	require.True(t, ok)
	require.Equal(t, "foo", v)
	client.Delete("foo")
	_, ok = client.Get("foo")
	require.False(t, ok)
}

func TestGetSetParallel(t *testing.T) {
	client, err := theine.New[string, string](1000)
	require.Nil(t, err)
	var wg sync.WaitGroup
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				key := fmt.Sprintf("key:%d", rand.Intn(3000))
				v, ok := client.Get(key)
				if !ok {
					client.Set(key, key, 1)
				} else {
					require.Equal(t, key, v)
				}
			}
		}()
	}
	wg.Wait()
	time.Sleep(300 * time.Millisecond)
	require.True(t, client.Len() < 1200)
}

func TestSetWithTTL(t *testing.T) {
	client, err := theine.New[string, string](500)
	require.Nil(t, err)
	client.SetWithTTL("foo", "foo", 1, 3600*time.Second)
	require.Equal(t, 1, client.Len())
	time.Sleep(1 * time.Second)
	client.SetWithTTL("foo", "foo", 1, 1*time.Second)
	require.Equal(t, 1, client.Len())
	time.Sleep(2 * time.Second)
	_, ok := client.Get("foo")
	require.False(t, ok)
	require.Equal(t, 0, client.Len())
}

func TestSetWithTTLAutoExpire(t *testing.T) {
	client, err := theine.New[string, string](500)
	require.Nil(t, err)
	for i := 0; i < 30; i++ {
		key1 := fmt.Sprintf("key:%d", i)
		client.SetWithTTL(key1, key1, 1, time.Duration(i+1)*time.Second)
		key2 := fmt.Sprintf("key:%d:2", i)
		client.SetWithTTL(key2, key2, 1, time.Duration(i+100)*time.Second)
	}
	current := 60
	counter := 0
	for {
		time.Sleep(5 * time.Second)
		counter += 1
		require.True(t, client.Len() < current)
		current = client.Len()
		if current <= 30 {
			break
		}
	}
	require.True(t, counter < 10)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		_, ok := client.Get(key)
		require.False(t, ok)
	}
}

func TestGetSetDeleteNoRace(t *testing.T) {
	for _, size := range []int{500, 2000, 10000, 50000} {
		client, err := theine.New[string, string](int64(size))
		require.Nil(t, err)
		var wg sync.WaitGroup
		keys := []string{}
		for i := 0; i < 100000; i++ {
			keys = append(keys, fmt.Sprintf("%d", rand.Intn(1000000)))
		}
		for i := 1; i <= 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100000; i++ {
					key := keys[i]
					client.Get(key)
					if i%3 == 0 {
						client.SetWithTTL(key, key, 1, time.Second*time.Duration(i%25+5))
					}
					if i%5 == 0 {
						client.Delete(key)
					}
				}
			}()
		}
		wg.Wait()
		time.Sleep(300 * time.Millisecond)
		require.True(t, client.Len() < size+50)
	}
}

func TestCost(t *testing.T) {
	client, err := theine.New[string, string](500)
	require.Nil(t, err)
	success := client.Set("z", "z", 501)
	require.False(t, success)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		success = client.Set(key, key, 20)
		require.True(t, success)
	}
	time.Sleep(time.Second)
	require.True(t, client.Len() == 25)

	// test cost func
	client, err = theine.New[string, string](500)
	require.Nil(t, err)
	client.SetCost(func(v string) int64 {
		return int64(len(v))
	})
	success = client.Set("z", strings.Repeat("z", 501), 0)
	require.False(t, success)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		success = client.Set(key, strings.Repeat("z", 20), 0)
		require.True(t, success)
	}
	time.Sleep(time.Second)
	require.True(t, client.Len() == 25)
}

func TestCostUpdate(t *testing.T) {
	client, err := theine.New[string, string](500)
	require.Nil(t, err)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		success := client.Set(key, key, 20)
		require.True(t, success)
	}
	time.Sleep(time.Second)
	require.True(t, client.Len() == 25)
	// update cost
	success := client.Set("key:10", "", 200)
	require.True(t, success)
	time.Sleep(time.Second)
	// 15 * 20 + 200
	require.True(t, client.Len() == 16)
}

func TestDoorkeeper(t *testing.T) {
	client, err := theine.New[string, string](500)
	require.Nil(t, err)
	client.SetDoorkeeper(true)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		success := client.Set(key, key, 20)
		require.False(t, success)
	}
	require.True(t, client.Len() == 0)
	time.Sleep(time.Second)
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key:%d", i)
		success := client.Set(key, key, 20)
		require.True(t, success)
	}
	require.True(t, client.Len() > 0)
}
