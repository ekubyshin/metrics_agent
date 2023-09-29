package storage

type Storage interface {
	Put(_ any) error
	Get() (any, error)
	Delete() error
}

type MemStorage struct {
}

func (m *MemStorage) Put(_ any) error {
	panic("not implemented") // TODO: Implement
}

func (m *MemStorage) Get() (any, error) {
	panic("not implemented") // TODO: Implement
}

func (m *MemStorage) Delete() error {
	panic("not implemented") // TODO: Implement
}
