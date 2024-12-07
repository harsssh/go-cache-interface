package isucon_utility

import "context"

type Option[V any] struct {
	Value V
	Found bool
}

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (Option[V], error)
	Set(ctx context.Context, key K, value V) error
	Delete(ctx context.Context, key K) error
	Clear(ctx context.Context) error
	Len(ctx context.Context) (uint, error)
}
