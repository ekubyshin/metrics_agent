package agent

import (
	"github.com/ekubyshin/metrics_agent/internal/types"
)

type SystemInfo struct {
	Alloc         types.Gauge   `json:"alloc"`
	BuckHashSys   types.Gauge   `json:"buck_hash_sys"`
	Frees         types.Gauge   `json:"frees"`
	GCCPUFraction types.Gauge   `json:"gc_cpu_fraction"`
	GCSys         types.Gauge   `json:"gc_sys"`
	HeapAlloc     types.Gauge   `json:"heap_alloc"`
	HeapIdle      types.Gauge   `json:"heap_idle"`
	HeapInuse     types.Gauge   `json:"heap_inuse"`
	HeapObjects   types.Gauge   `json:"heap_objects"`
	HeapReleased  types.Gauge   `json:"heap_released"`
	HeapSys       types.Gauge   `json:"heap_sys"`
	LastGC        types.Gauge   `json:"last_gc"`
	Lookups       types.Gauge   `json:"lookups"`
	MCacheInuse   types.Gauge   `json:"mcache_in_use"`
	MSpanInuse    types.Gauge   `json:"mspan_in_use"`
	Mallocs       types.Gauge   `json:"mallocs"`
	NextGC        types.Gauge   `json:"next_gc"`
	NumForcedGC   types.Gauge   `json:"num_forced_gc"`
	NumGC         types.Gauge   `json:"num_gc"`
	OtherSys      types.Gauge   `json:"other_sys"`
	PauseTotalNs  types.Gauge   `json:"pause_total_ns"`
	StackInuse    types.Gauge   `json:"stack_in_use"`
	StackSys      types.Gauge   `json:"stack_sys"`
	Sys           types.Gauge   `json:"sys"`
	TotalAlloc    types.Gauge   `json:"total_alloc"`
	RandomValue   types.Gauge   `json:"random_value"`
	MCacheSys     types.Gauge   `json:"mcache_sys"`
	MSpanSys      types.Gauge   `json:"mspan_sys"`
	PollCount     types.Counter `json:"poll_counter"`
}
