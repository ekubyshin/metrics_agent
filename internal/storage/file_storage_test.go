package storage

import (
	"context"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestRestoreStorage(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"ok",
			args{
				"./test/test.json",
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			fs, err := NewFileStorage(ctx, db, tt.args.filename, true, 0)
			defer func() { _ = fs.Close() }()
			assert.NoError(t, err)
			elems, _ := fs.List(context.TODO())
			assert.Equal(t, tt.want, len(elems))
		})
	}
}
