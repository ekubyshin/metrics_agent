package storage

import (
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/types"
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
			db := NewMemoryStorage[types.MetricsKey, types.Metrics]()
			fs, err := NewFileStorage(db, tt.args.filename, true, 0)
			defer func() { _ = fs.Close() }()
			assert.NoError(t, err)
			elems := fs.List()
			assert.Equal(t, tt.want, len(elems))
		})
	}
}
