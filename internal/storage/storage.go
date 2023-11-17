package storage

import (
	"context"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
)

type KeyValuer[K any, V any] struct {
	Key   K
	Value V
}

type Storage[K any, V any] interface {
	Put(context.Context, K, V) (*V, error)
	PutBatch(context.Context, []KeyValuer[K, V]) ([]V, error)
	Get(context.Context, K) (V, bool)
	Delete(context.Context, K) error
	List(context.Context) ([]KeyValuer[K, V], error)
	Ping(context.Context) error
	Close() error
}

func AutoLoadStorage[K any, V metrics.Keyable[K]](ctx context.Context, cfg config.Config) (Storage[K, V], error) {
	var st Storage[K, V]
	var err error
	if cfg.DatabaseDSN != nil && *cfg.DatabaseDSN != "" {
		st, err = NewDBStorage[K, V](ctx, &cfg)
		if err != nil {
			return nil, err
		}
	} else {
		memSt := NewMemoryStorage[K, V]()
		if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
			st, err = NewFileStorage[K, V](ctx, memSt, *cfg.FileStoragePath, *cfg.Restore, cfg.StoreDuration())
			if err != nil {
				return nil, err
			}
		} else {
			st = memSt
		}
	}
	return st, nil
}
