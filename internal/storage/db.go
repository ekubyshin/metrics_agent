package storage

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type DBStorage struct {
	conn sqlx.DB
}

func NewDBStorage(cfg *config.Config) (storage *DBStorage, err error) {
	conn, err := sqlx.Open("pgx", *cfg.DBDSN)
	if err != nil {
		return
	}
	storage = &DBStorage{
		conn: *conn,
	}
	return
}
