package storage

import (
	"context"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStorage[K any, V any] struct {
	conn sqlx.DB
}

func NewDBStorage[K any, V any](ctx context.Context, cfg *config.Config) (storage *DBStorage[K, V], err error) {
	if cfg.DatabaseDSN == nil || *cfg.DatabaseDSN == "" {
		return
	}
	conn, err := sqlx.Open("postgres", *cfg.DatabaseDSN)
	if err != nil {
		return
	}
	storage = &DBStorage[K, V]{
		conn: *conn,
	}
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				conn.Close()
				return
			default:
				continue
			}
		}
	}(ctx)
	return
}

func (db *DBStorage[K, V]) Ping() error {
	return db.conn.Ping()
}

func (db *DBStorage[K, V]) Close() error {
	return db.conn.Close()
}

func (m *DBStorage[K, V]) Put(key K, val V) {
	panic("not implemented")
}

func (m *DBStorage[K, V]) Get(key K) (V, bool) {
	panic("not implemented")
}

func (m *DBStorage[K, V]) Delete(key K) {
	panic("not implemented")
}

func (m *DBStorage[K, V]) List() []KeyValuer[K, V] {
	panic("not implemented")
}
