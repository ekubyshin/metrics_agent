package storage

type KeyValuer struct {
	Key   any
	Value any
}

type Storage interface {
	Put(any, any)
	Get(any) (any, error)
	Delete(any)
	List() []KeyValuer
}
