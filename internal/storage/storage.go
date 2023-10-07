package storage

type Storage interface {
	Put(any, any)
	Get(any) (any, error)
	Delete(any)
}
