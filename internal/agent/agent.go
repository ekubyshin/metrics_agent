package agent

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
	"github.com/ekubyshin/metrics_agent/internal/reporter"
	"github.com/go-resty/resty/v2"
)

type Agent interface {
	Start()
}

type MetricsAgent struct {
	reporter        reporter.Writer
	collector       collector.Reader
	sendInterval    time.Duration
	refreshInterval time.Duration
	queue           chan collector.SystemInfo
	done            chan bool
	batchSize       int64
	canceled        bool
}

func NewMetricsAgent(
	reportInterval time.Duration,
	pollInterval time.Duration,
	endpoint string,
) Agent {
	colector := collector.NewRuntimeReader(pollInterval)
	client := resty.New()
	reporter := reporter.NewAgentReporter(reportInterval, client, endpoint)
	return &MetricsAgent{
		reporter:        reporter,
		collector:       colector,
		sendInterval:    reportInterval,
		refreshInterval: pollInterval,
		queue:           make(chan collector.SystemInfo, 100),
		batchSize:       int64(math.Ceil(float64(reportInterval) / float64(pollInterval))),
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
		a.queue <- st
	}
}

func (a *MetricsAgent) report() {
	count := 0
	time.Sleep(a.sendInterval)
	for {
		select {
		case <-a.done:
			a.canceled = true
			return
		case st := <-a.queue:
			count++
			if count <= int(a.batchSize) {
				a.reporter.WriteBatch(convertSystemInfoToReport(st))
			} else {
				count = 0
				time.Sleep(a.sendInterval)
			}
		}
	}
}

func convertSystemInfoToReport(info collector.SystemInfo) []reporter.Report {
	v := reflect.ValueOf(info)
	reports := make([]reporter.Report, v.NumField())
	for f := 0; f < v.NumField(); f++ {
		field := v.Field(f)
		tName := field.Type().Name()
		fieldName := v.Type().Field(f).Name
		report := reporter.Report{
			Type: strings.ToLower(tName),
			Name: fieldName,
		}
		switch field.Kind() {
		case reflect.Float64:
			report.Value = strconv.FormatFloat(field.Float(), 'f', 1, 64)
		case reflect.Int64:
			report.Value = strconv.FormatInt(field.Int(), 10)
		}
		reports[f] = report
	}
	return reports
}
