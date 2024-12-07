package isuutil

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
)

type Maybe[V any] struct {
	Value V
	Found bool
}

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (Maybe[V], error)
	Set(ctx context.Context, key K, value V) error
	Delete(ctx context.Context, key K) error
	Clear(ctx context.Context) error
}

type InMemoryLRUCache[K comparable, V any] struct {
	l *lru.Cache[K, V]
}

func (c *InMemoryLRUCache[K, V]) Get(ctx context.Context, key K) (Maybe[V], error) {
	v, ok := c.l.Get(key)
	if !ok {
		return Maybe[V]{Found: false}, nil
	}
	return Maybe[V]{Value: v, Found: true}, nil
}

func (c *InMemoryLRUCache[K, V]) Set(ctx context.Context, key K, value V) error {
	c.l.Add(key, value)
	return nil
}

func (c *InMemoryLRUCache[K, V]) Delete(ctx context.Context, key K) error {
	c.l.Remove(key)
	return nil
}

func (c *InMemoryLRUCache[K, V]) Clear(ctx context.Context) error {
	c.l.Purge()
	return nil
}

func NewInMemoryLRUCache[K comparable, V any](size int) (Cache[K, V], error) {
	l, err := lru.New[K, V](size)
	if err != nil {
		return nil, err
	}
	return &InMemoryLRUCache[K, V]{l: l}, nil
}

type RedisCache[V any] struct {
	rdb redis.Client
}

func (c *RedisCache[V]) Get(ctx context.Context, key string) (Maybe[V], error) {
	raw, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return Maybe[V]{Found: false}, nil
		}
		return Maybe[V]{Found: false}, err
	}

	var v V
	err = json.UnmarshalContext(ctx, []byte(raw), &v)
	if err != nil {
		return Maybe[V]{Found: false}, err
	}

	return Maybe[V]{Value: v, Found: true}, nil
}

func (c *RedisCache[V]) Set(ctx context.Context, key string, value V) error {
	b, err := json.MarshalContext(ctx, value)
	if err != nil {
		return err
	}

	err = c.rdb.Set(ctx, key, b, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache[V]) Delete(ctx context.Context, key string) error {
	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache[V]) Clear(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func NewRedisCache[V any](rdb redis.Client) Cache[string, V] {
	return &RedisCache[V]{rdb: rdb}
}
