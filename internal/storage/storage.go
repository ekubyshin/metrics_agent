package storage

import "context"

type KeyValuer[K any, V any] struct {
	Key   K
	Value V
}

type Storage[K any, V any] interface {
	Put(context.Context, K, V) (*V, error)
	PutBatch(context.Context, []KeyValuer[K, V]) ([]V, error)
	Get(context.Context, K) (V, bool)
	Delete(context.Context, K) error
	List(context.Context) ([]KeyValuer[K, V], error)
	Ping(context.Context) error
	Close() error
}
