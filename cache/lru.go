package cache

import (
	"context"
	lru "github.com/hashicorp/golang-lru/v2"
)

type inMemoryLRUCache[K comparable, V any] struct {
	l *lru.Cache[K, V]
}

func (c *inMemoryLRUCache[K, V]) Get(ctx context.Context, key K) (Maybe[V], error) {
	v, ok := c.l.Get(key)
	if !ok {
		return Maybe[V]{Found: false}, nil
	}
	return Maybe[V]{Value: v, Found: true}, nil
}

func (c *inMemoryLRUCache[K, V]) Set(ctx context.Context, key K, value V) error {
	c.l.Add(key, value)
	return nil
}

func (c *inMemoryLRUCache[K, V]) Delete(ctx context.Context, key K) error {
	c.l.Remove(key)
	return nil
}

func (c *inMemoryLRUCache[K, V]) Clear(ctx context.Context) error {
	c.l.Purge()
	return nil
}

func NewInMemoryLRUCache[K comparable, V any](size int) (Cache[K, V], error) {
	l, err := lru.New[K, V](size)
	if err != nil {
		return nil, err
	}
	return &inMemoryLRUCache[K, V]{l: l}, nil
}
