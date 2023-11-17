package agent

import (
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/go-resty/resty/v2"
)

type Agent interface {
	Start()
}

type MetricsAgent struct {
	reporter       Writer
	collector      Reader
	reportInterval time.Duration
	pollInterval   time.Duration
	queue          chan SystemInfo
	done           chan bool
	batchSize      int64
	canceled       bool
}

func NewMetricsAgent(
	cfg config.Config,
) *MetricsAgent {
	collector := NewRuntimeReader(cfg.PollDuration())
	client := resty.New()
	reporter := NewAgentReporter(client, cfg.Address.ToString())
	return &MetricsAgent{
		reporter:       reporter,
		collector:      collector,
		reportInterval: cfg.ReportDuration(),
		pollInterval:   cfg.PollDuration(),
		queue:          make(chan SystemInfo, 100),
		batchSize:      int64(math.Ceil(float64(cfg.ReportInterval) / float64(cfg.PollInterval))),
	}
}

func (a *MetricsAgent) Start() {
	go a.collect()
	go a.report()
	<-a.done
}

func (a *MetricsAgent) collect() {
	a.collector.Run()
	for {
		st, err := a.collector.Read()
		if err != nil {
			a.done <- true
			a.collector.Stop()
			return
		}
		if st == nil {
			return
		}
		a.queue <- *st
	}
}

func (a *MetricsAgent) report() {
	count := 0
	time.Sleep(a.reportInterval)
	for {
		select {
		case <-a.done:
			a.canceled = true
			return
		case st := <-a.queue:
			count++
			if count <= int(a.batchSize) {
				_ = a.reporter.WriteBatch(convertSystemInfoToMetric(st))
			} else {
				count = 0
				time.Sleep(a.reportInterval)
			}
		}
	}
}

func convertSystemInfoToMetric(info SystemInfo) []metrics.Metrics {
	v := reflect.ValueOf(info)
	reports := make([]metrics.Metrics, v.NumField())
	for f := 0; f < v.NumField(); f++ {
		field := v.Field(f)
		tName := field.Type().Name()
		fieldName := v.Type().Field(f).Name
		metric := metrics.Metrics{
			ID:    fieldName,
			MType: strings.ToLower(tName),
		}
		switch field.Kind() {
		case reflect.Float64:
			metric.Value = pointer.From[float64](field.Float())
		case reflect.Int64:
			metric.Delta = pointer.From[int64](field.Int())
		}
		reports[f] = metric
	}
	return reports
}
