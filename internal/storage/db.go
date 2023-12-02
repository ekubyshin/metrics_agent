package storage

import (
	"context"
	"embed"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

func (m *DBStorage[K, V]) Put(ctx context.Context, key K, val V) (*V, error) {
	r, err := m.conn.NamedQueryContext(
		ctx,
		`
		INSERT INTO metrics(id, type, delta, value)
		VALUES (:id, :type, :delta, :value)
		ON CONFLICT (id, type) DO
		UPDATE 
			SET delta = EXCLUDED.delta,
				value = EXCLUDED.value
		RETURNING metrics.id, metrics.type, metrics.delta, metrics.value;
		`,
		val,
	)
	if err != nil {
		return nil, err
	}
	v := new(V)
	for r.Next() {
		if r.Err() != nil {
			continue
		}
		err = r.StructScan(v)
	}
	return v, err
}

func (m *DBStorage[K, V]) PutBatch(ctx context.Context, vals []KeyValuer[K, V]) ([]V, error) {
	tx, err := m.conn.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() //nolint
		}
	}()
	out := make([]V, 0, len(vals))
	for _, v := range vals {
		r, err := sqlx.NamedQueryContext(
			ctx,
			tx,
			`
				INSERT INTO metrics(id, type, delta, value)
				VALUES (:id, :type, :delta, :value)
				ON CONFLICT (id, type) DO
				UPDATE 
					SET delta = metrics.delta + EXCLUDED.delta,
						value = EXCLUDED.value
				RETURNING metrics.id, metrics.type, metrics.delta, metrics.value;
			`,
			v.Value,
		)
		if err != nil {
			return nil, err
		}
		for r.Next() {
			if r.Err() != nil {
				continue
			}
			v := new(V)
			err := r.StructScan(v)
			if err == nil {
				out = append(out, *v)
			}
		}

	}
	err = tx.Commit()
	return out, err
}

func (m *DBStorage[K, V]) Get(ctx context.Context, key K) (V, bool) {
	var r V
	rows, err := m.conn.NamedQueryContext(
		ctx,
		`
			SELECT * FROM metrics WHERE id=:id AND type=:type;
		`,
		key,
	)
	if err != nil {
		return r, false
	}
	for rows.Next() {
		if rows.Err() != nil {
			continue
		}
		err := rows.StructScan(&r)
		if err != nil {
			return r, false
		}
	}
	return r, true
}

func (m *DBStorage[K, V]) Delete(ctx context.Context, key K) error {
	_, err := m.conn.NamedExecContext(
		ctx,
		`
			DELETE FROM metrics WHERE id=:id AND type=:type
		`,
		key,
	)
	return err
}

func (m *DBStorage[K, V]) List(ctx context.Context) ([]KeyValuer[K, V], error) {
	rows, err := m.conn.NamedQueryContext(
		ctx,
		`
			SELECT * FROM metrics;
		`,
		map[string]any{},
	)
	if err != nil {
		return nil, err
	}
	res := make([]KeyValuer[K, V], 0)
	for rows.Next() {
		if rows.Err() != nil {
			continue
		}
		r := new(V)
		err := rows.StructScan(r)
		if err != nil {
			continue
		}
		res = append(res, KeyValuer[K, V]{Key: (*r).Key(), Value: *r})
	}
	return res, err
}
