package storage

import (
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

func (m *MemStorage[K, V]) Put(key K, val V) {
	m.data.Swap(key, val)
}

func (m *MemStorage[K, V]) Get(key K) (V, bool) {
	var res V
	if val, ok := m.data.Load(key); ok {
		if castVal, ok := val.(V); ok {
			return castVal, true
		}
		return res, false
	}
	return res, false
}

func (m *MemStorage[K, V]) Delete(key K) {
	m.data.Delete(key)
}

func (m *MemStorage[K, V]) Ping() error {
	return nil
}

func (m *MemStorage[K, V]) List() []KeyValuer[K, V] {
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
	return arr
}
