package storage

import (
	"context"
	"embed"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const timeout = 3 * time.Second

type DBStorage[K any, V metrics.Keyable[K]] struct {
	conn sqlx.DB
}

func NewDBStorage[K any, V metrics.Keyable[K]](ctx context.Context, cfg *config.Config) (storage *DBStorage[K, V], err error) {
	if cfg.DatabaseDSN == nil || *cfg.DatabaseDSN == "" {
		return
	}
	conn, err := sqlx.Open("postgres", *cfg.DatabaseDSN)
	if err != nil {
		return
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn.DB, "migrations"); err != nil {
		panic(err)
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

func (m *DBStorage[K, V]) Ping(ctx context.Context) error {
	return m.conn.PingContext(ctx)
}

func (m *DBStorage[K, V]) Close() error {
	return m.conn.Close()
}

func (m *DBStorage[K, V]) Put(ctx context.Context, key K, val V) error {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := m.conn.NamedExecContext(
		c,
		`
			INSERT INTO metrics(id, type, delta, value)
			VALUES (:id, :type, :delta, :value)
			ON CONFLICT (id, type) DO
			UPDATE 
				SET delta = EXCLUDED.delta,
					value = EXCLUDED.value;
		`,
		val.Serialize(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBStorage[K, V]) Get(ctx context.Context, key K) (V, bool) {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r V
	rows, err := m.conn.NamedQueryContext(
		c,
		`
			SELECT * FROM metrics WHERE id=:id AND type=:type;
		`,
		key,
	)
	if err != nil {
		return r, false
	}
	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return r, false
		}
	}
	return r, true
}

func (m *DBStorage[K, V]) Delete(ctx context.Context, key K) error {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := m.conn.NamedExecContext(
		c,
		`
			DELETE FROM metrics WHERE id=:id AND type=:type
		`,
		key,
	)
	return err
}

func (m *DBStorage[K, V]) List(ctx context.Context) ([]KeyValuer[K, V], error) {
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r V
	res := make([]KeyValuer[K, V], 0)
	rows, err := m.conn.NamedQueryContext(
		c,
		`
			SELECT * FROM metrics;
		`,
		map[string]any{},
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.StructScan(r)
		if err != nil {
			continue
		}
		res = append(res, KeyValuer[K, V]{Key: r.Key(), Value: r})
	}
	return res, err
}
