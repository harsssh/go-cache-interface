package cache

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type redisCache[V any] struct {
	rdb redis.Client
}

func (c *redisCache[V]) Get(ctx context.Context, key string) (Maybe[V], error) {
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

func (c *redisCache[V]) Set(ctx context.Context, key string, value V) error {
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

func (c *redisCache[V]) Delete(ctx context.Context, key string) error {
	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisCache[V]) Clear(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func NewRedisCache[V any](rdb redis.Client) Cache[string, V] {
	return &redisCache[V]{rdb: rdb}
}
