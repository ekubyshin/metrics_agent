package collector

import (
	"errors"
	"math/rand"
	"runtime"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/types"
)

type Reader interface {
	Run()
	Stop()
	Read() (*SystemInfo, error)
}

type RuntimeReader struct {
	stats        chan runtime.MemStats
	done         chan bool
	closed       bool
	pollInterval time.Duration
	pollCounter  types.Counter
}

func NewRuntimeReader(pollInterval time.Duration) *RuntimeReader {
	return &RuntimeReader{
		stats:        make(chan runtime.MemStats, 100),
		done:         make(chan bool),
		pollInterval: pollInterval,
	}
}

func (r *RuntimeReader) Run() {
	go (func() {
		stats := runtime.MemStats{}
		for {
			select {
			case r.stats <- stats:
				runtime.ReadMemStats(&stats)
				r.pollCounter += 1
				time.Sleep(r.pollInterval)
			case d := <-r.done:
				if d {
					r.closed = true
					close(r.stats)
					close(r.done)
					return
				}
			}
		}
	})()
}

func generateRandom() float64 {
	return rand.Float64()
}

func (r *RuntimeReader) Read() (*SystemInfo, error) {
	if r.closed {
		return nil, errors.New("monitor stopped")
	}
	return r.convertToStat(<-r.stats), nil
}

func (r *RuntimeReader) Stop() {
	r.done <- true
}

func (r *RuntimeReader) convertToStat(st runtime.MemStats) *SystemInfo {
	return &SystemInfo{
		Alloc:         types.Gauge(st.Alloc),
		BuckHashSys:   types.Gauge(st.BuckHashSys),
		Frees:         types.Gauge(st.Frees),
		GCCPUFraction: types.Gauge(st.GCCPUFraction),
		GCSys:         types.Gauge(st.GCSys),
		HeapAlloc:     types.Gauge(st.HeapAlloc),
		HeapIdle:      types.Gauge(st.HeapIdle),
		HeapInuse:     types.Gauge(st.HeapInuse),
		HeapObjects:   types.Gauge(st.HeapObjects),
		HeapReleased:  types.Gauge(st.HeapReleased),
		HeapSys:       types.Gauge(st.HeapSys),
		LastGC:        types.Gauge(st.LastGC),
		Lookups:       types.Gauge(st.Lookups),
		MCacheInuse:   types.Gauge(st.MCacheInuse),
		MSpanInuse:    types.Gauge(st.MSpanInuse),
		Mallocs:       types.Gauge(st.Mallocs),
		NextGC:        types.Gauge(st.NextGC),
		NumForcedGC:   types.Gauge(st.NumForcedGC),
		NumGC:         types.Gauge(st.NumGC),
		OtherSys:      types.Gauge(st.OtherSys),
		PauseTotalNs:  types.Gauge(st.PauseTotalNs),
		StackInuse:    types.Gauge(st.StackInuse),
		StackSys:      types.Gauge(st.StackSys),
		Sys:           types.Gauge(st.Sys),
		TotalAlloc:    types.Gauge(st.TotalAlloc),
		PollCount:     r.pollCounter,
		RandomValue:   types.Gauge(generateRandom()),
	}
}
