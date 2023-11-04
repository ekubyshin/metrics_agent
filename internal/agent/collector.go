package agent

import (
	"errors"
	"math/rand"
	"runtime"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
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
	pollCounter  metrics.Counter
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
		Alloc:         metrics.Gauge(st.Alloc),
		BuckHashSys:   metrics.Gauge(st.BuckHashSys),
		Frees:         metrics.Gauge(st.Frees),
		GCCPUFraction: metrics.Gauge(st.GCCPUFraction),
		GCSys:         metrics.Gauge(st.GCSys),
		HeapAlloc:     metrics.Gauge(st.HeapAlloc),
		HeapIdle:      metrics.Gauge(st.HeapIdle),
		HeapInuse:     metrics.Gauge(st.HeapInuse),
		HeapObjects:   metrics.Gauge(st.HeapObjects),
		HeapReleased:  metrics.Gauge(st.HeapReleased),
		HeapSys:       metrics.Gauge(st.HeapSys),
		LastGC:        metrics.Gauge(st.LastGC),
		Lookups:       metrics.Gauge(st.Lookups),
		MCacheInuse:   metrics.Gauge(st.MCacheInuse),
		MSpanInuse:    metrics.Gauge(st.MSpanInuse),
		Mallocs:       metrics.Gauge(st.Mallocs),
		NextGC:        metrics.Gauge(st.NextGC),
		NumForcedGC:   metrics.Gauge(st.NumForcedGC),
		NumGC:         metrics.Gauge(st.NumGC),
		OtherSys:      metrics.Gauge(st.OtherSys),
		PauseTotalNs:  metrics.Gauge(st.PauseTotalNs),
		StackInuse:    metrics.Gauge(st.StackInuse),
		StackSys:      metrics.Gauge(st.StackSys),
		Sys:           metrics.Gauge(st.Sys),
		TotalAlloc:    metrics.Gauge(st.TotalAlloc),
		PollCount:     r.pollCounter,
		MCacheSys:     metrics.Gauge(st.MCacheSys),
		RandomValue:   metrics.Gauge(generateRandom()),
	}
}
