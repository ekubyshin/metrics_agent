package collector

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/types"
)

func TestRuntimeReader_convertToStat(t *testing.T) {
	type fields struct {
		stats        chan runtime.MemStats
		done         chan bool
		closed       bool
		pollInterval time.Duration
		pollCounter  types.Counter
	}
	type args struct {
		st runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   SystemInfo
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RuntimeReader{
				stats:        tt.fields.stats,
				done:         tt.fields.done,
				closed:       tt.fields.closed,
				pollInterval: tt.fields.pollInterval,
				pollCounter:  tt.fields.pollCounter,
			}
			if got := r.convertToStat(tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuntimeReader.convertToStat() = %v, want %v", got, tt.want)
			}
		})
	}
}
