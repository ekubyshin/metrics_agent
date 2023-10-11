package collector

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeReader_convertToStat(t *testing.T) {
	type args struct {
		st runtime.MemStats
	}
	tests := []struct {
		name string
		args args
		want *SystemInfo
	}{
		{
			"test",
			args{
				st: runtime.MemStats{
					Alloc: 1.0,
				},
			},
			&SystemInfo{Alloc: 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RuntimeReader{
				stats:        make(chan runtime.MemStats, 100),
				done:         make(chan bool),
				pollInterval: 1 * time.Second,
			}
			got := r.convertToStat(tt.args.st)
			got.RandomValue = 0
			assert.NotNil(t, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuntimeReader.convertToStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRandom(t *testing.T) {
	got1 := generateRandom()
	got2 := generateRandom()
	assert.NotEqual(t, got1, got2)
}
