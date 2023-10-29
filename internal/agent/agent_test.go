package agent

import (
	"reflect"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/collector"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

func Test_convertSystemInfoToReport(t *testing.T) {
	type args struct {
		info collector.SystemInfo
	}
	tests := []struct {
		name string
		args args
		want []types.Metrics
	}{
		{
			"check",
			args{
				info: collector.SystemInfo{
					Alloc:         1.0,
					BuckHashSys:   1.0,
					Frees:         1.0,
					GCCPUFraction: 1.0,
					GCSys:         1.0,
					HeapAlloc:     1.0,
					HeapIdle:      1.0,
					HeapInuse:     1.0,
					HeapObjects:   1.0,
					HeapReleased:  1.0,
					HeapSys:       1.0,
					LastGC:        1.0,
					Lookups:       1.0,
					MCacheInuse:   1.0,
					MSpanInuse:    1.0,
					Mallocs:       1.0,
					NextGC:        1.0,
					NumForcedGC:   1.0,
					NumGC:         1.0,
					OtherSys:      1.0,
					PauseTotalNs:  1.0,
					StackInuse:    1.0,
					StackSys:      1.0,
					Sys:           1.0,
					TotalAlloc:    1.0,
					RandomValue:   1.0,
					PollCount:     1,
				},
			},
			[]types.Metrics{
				{
					MType: "gauge",
					ID:    "Alloc",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "BuckHashSys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "Frees",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "GCCPUFraction",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "GCSys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapAlloc",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapIdle",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapInuse",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapObjects",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapReleased",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "HeapSys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "LastGC",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "Lookups",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "MCacheInuse",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "MSpanInuse",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "Mallocs",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "NextGC",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "NumForcedGC",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "NumGC",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "OtherSys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "PauseTotalNs",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "StackInuse",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "StackSys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "Sys",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "TotalAlloc",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "gauge",
					ID:    "RandomValue",
					Value: utils.ToPointer[float64](1.0),
				},
				{
					MType: "counter",
					ID:    "PollCount",
					Delta: utils.ToPointer[int64](1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSystemInfoToMetric(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRuntimeReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
