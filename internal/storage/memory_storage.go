package storage

import (
	"errors"
	"sync"
)

type MemStorage struct {
	data sync.Map
}

func NewMemoryStorage() Storage {
	return &MemStorage{
		data: sync.Map{},
	}
}

func (m *MemStorage) Put(key any, val any) {
	m.data.Swap(key, val)
}

func (m *MemStorage) Get(key any) (any, error) {
	if val, ok := m.data.Load(key); ok {
		return val, nil
	}
	return nil, errors.New("NotFound")
}

func (m *MemStorage) Delete(key any) {
	m.data.Delete(key)
}

func (m *MemStorage) List() []KeyValuer {
	arr := make([]KeyValuer, 0, 100)
	m.data.Range(func(key, value any) bool {
		if key != nil && key != "" && value != nil {
			arr = append(arr, KeyValuer{Key: key, Value: value})
			return true
		}
		return false
	})
	return arr
}
