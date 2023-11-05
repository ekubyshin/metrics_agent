package storage

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	conn sqlx.DB
}

func NewDBStorage(cfg *config.Config) (storage *DBStorage, err error) {
	if cfg.DatabaseDSN == nil || *cfg.DatabaseDSN == "" {
		return
	}
	conn, err := sqlx.Open("postgres", *cfg.DatabaseDSN)
	if err != nil {
		return
	}
	storage = &DBStorage{
		conn: *conn,
	}
	return
}

func (db *DBStorage) Ping() error {
	return db.conn.Ping()
}

func (db *DBStorage) Close() error {
	return db.conn.Close()
}
