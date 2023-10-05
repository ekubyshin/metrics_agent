package agent

import (
	"reflect"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/collector"
	"github.com/ekubyshin/metrics_agent/internal/reporter"
)

func Test_convertSystemInfoToReport(t *testing.T) {
	type args struct {
		info collector.SystemInfo
	}
	tests := []struct {
		name string
		args args
		want []reporter.Report
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
				},
			},
			[]reporter.Report{
				{
					Type:  "gauge",
					Name:  "Alloc",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "BuckHashSys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "Frees",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "GCCPUFraction",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "GCSys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapAlloc",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapIdle",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapInuse",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapObjects",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapReleased",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "HeapSys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "LastGC",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "Lookups",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "MCacheInuse",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "MSpanInuse",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "Mallocs",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "NextGC",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "NumForcedGC",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "NumGC",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "OtherSys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "PauseTotalNs",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "StackInuse",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "StackSys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "Sys",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "TotalAlloc",
					Value: "1.0",
				},
				{
					Type:  "gauge",
					Name:  "RandomValue",
					Value: "1.0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSystemInfoToReport(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRuntimeReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
