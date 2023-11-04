package agent

import (
	"github.com/ekubyshin/metrics_agent/internal/metrics"
)

type SystemInfo struct {
	Alloc         metrics.Gauge   `json:"alloc"`
	BuckHashSys   metrics.Gauge   `json:"buck_hash_sys"`
	Frees         metrics.Gauge   `json:"frees"`
	GCCPUFraction metrics.Gauge   `json:"gc_cpu_fraction"`
	GCSys         metrics.Gauge   `json:"gc_sys"`
	HeapAlloc     metrics.Gauge   `json:"heap_alloc"`
	HeapIdle      metrics.Gauge   `json:"heap_idle"`
	HeapInuse     metrics.Gauge   `json:"heap_inuse"`
	HeapObjects   metrics.Gauge   `json:"heap_objects"`
	HeapReleased  metrics.Gauge   `json:"heap_released"`
	HeapSys       metrics.Gauge   `json:"heap_sys"`
	LastGC        metrics.Gauge   `json:"last_gc"`
	Lookups       metrics.Gauge   `json:"lookups"`
	MCacheInuse   metrics.Gauge   `json:"mcache_in_use"`
	MSpanInuse    metrics.Gauge   `json:"mspan_in_use"`
	Mallocs       metrics.Gauge   `json:"mallocs"`
	NextGC        metrics.Gauge   `json:"next_gc"`
	NumForcedGC   metrics.Gauge   `json:"num_forced_gc"`
	NumGC         metrics.Gauge   `json:"num_gc"`
	OtherSys      metrics.Gauge   `json:"other_sys"`
	PauseTotalNs  metrics.Gauge   `json:"pause_total_ns"`
	StackInuse    metrics.Gauge   `json:"stack_in_use"`
	StackSys      metrics.Gauge   `json:"stack_sys"`
	Sys           metrics.Gauge   `json:"sys"`
	TotalAlloc    metrics.Gauge   `json:"total_alloc"`
	RandomValue   metrics.Gauge   `json:"random_value"`
	MCacheSys     metrics.Gauge   `json:"mcache_sys"`
	MSpanSys      metrics.Gauge   `json:"mspan_sys"`
	PollCount     metrics.Counter `json:"poll_counter"`
}
