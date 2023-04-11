package theine

import (
	"errors"
	"time"

	"github.com/Yiling-J/theine-go/internal"
)

const (
	ZERO_TTL = 0 * time.Second
)

type Cache[K comparable, V any] struct {
	store *internal.Store[K, V]
}

func New[K comparable, V any](maxsize int64) (*Cache[K, V], error) {
	if maxsize <= 0 {
		return nil, errors.New("size must be positive")
	}

	return &Cache[K, V]{
		store: internal.NewStore[K, V](maxsize),
	}, nil
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.store.Get(key)
}

func (c *Cache[K, V]) SetCost(cost func(v V) int64) *Cache[K, V] {
	c.store.SetCost(cost)
	return c
}
func (c *Cache[K, V]) SetDoorkeeper(enabled bool) *Cache[K, V] {
	c.store.SetDoorkeeper(enabled)
	return c
}

func (c *Cache[K, V]) SetWithTTL(key K, value V, cost int64, ttl time.Duration) bool {
	return c.store.Set(key, value, cost, ttl)
}

func (c *Cache[K, V]) Set(key K, value V, cost int64) bool {
	return c.SetWithTTL(key, value, cost, ZERO_TTL)
}

func (c *Cache[K, V]) Delete(key K) {
	c.store.Delete(key)
}

func (c *Cache[K, V]) Len() int {
	return c.store.Len()
}

func (c *Cache[K, V]) Close() {
	c.store.Close()
}
