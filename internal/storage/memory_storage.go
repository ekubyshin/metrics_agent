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
