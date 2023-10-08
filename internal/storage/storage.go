package storage

type KeyValuer[K any, V any] struct {
	Key   K
	Value V
}

type Storage[K any, V any] interface {
	Put(K, V)
	Get(K) (V, bool)
	Delete(K)
	List() []KeyValuer[K, V]
}
