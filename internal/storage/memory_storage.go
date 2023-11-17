package storage

import (
	"context"
	"sync"
)

type MemStorage[K any, V any] struct {
	data sync.Map
}

func NewMemoryStorage[K any, V any]() *MemStorage[K, V] {
	return &MemStorage[K, V]{
		data: sync.Map{},
	}
}

func (m *MemStorage[K, V]) Put(ctx context.Context, key K, val V) (*V, error) {
	m.data.Swap(key, val)
	return nil, nil
}

func (m *MemStorage[K, V]) PutBatch(ctx context.Context, vals []KeyValuer[K, V]) ([]V, error) {
	for _, v := range vals {
		if _, err := m.Put(ctx, v.Key, v.Value); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (m *MemStorage[K, V]) Get(ctx context.Context, key K) (V, bool) {
	var res V
	if val, ok := m.data.Load(key); ok {
		if castVal, ok := val.(V); ok {
			return castVal, true
		}
		return res, false
	}
	return res, false
}

func (m *MemStorage[K, V]) Delete(ctx context.Context, key K) error {
	m.data.Delete(key)
	return nil
}

func (m *MemStorage[K, V]) Ping(ctx context.Context) error {
	return nil
}

func (m *MemStorage[K, V]) Close() error {
	return nil
}

func (m *MemStorage[K, V]) List(ctx context.Context) ([]KeyValuer[K, V], error) {
	arr := make([]KeyValuer[K, V], 0, 100)
	m.data.Range(func(key, value any) bool {
		if key == nil || key == "" || value == nil {
			return true
		}
		if catKey, ok := key.(K); ok {
			if castVal, ok := value.(V); ok {
				arr = append(arr, KeyValuer[K, V]{Key: catKey, Value: castVal})
			}
		}
		return true
	})
	return arr, nil
}
