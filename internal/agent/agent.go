package agent

import (
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
	"github.com/ekubyshin/metrics_agent/internal/reporter"
)

type Agent interface {
	Start()
}

type MetricsAgent struct {
	reporter        reporter.Writer
	collector       collector.Reader
	pollInterval    time.Duration
	refreshInterval time.Duration
	counter         int64
}

func NewMetricsAgent(
	pollInterval time.Duration,
	refreshInterval time.Duration) Agent {
	colector := collector.NewRuntimeReader(refreshInterval)
	client := reporter.NewReporterClient()
	reporter := reporter.NewAgentReporter(pollInterval, client)
	return &MetricsAgent{
		reporter:        reporter,
		collector:       colector,
		pollInterval:    pollInterval,
		refreshInterval: refreshInterval,
	}
}

func (a *MetricsAgent) Start() {
	a.collector.Run()
	for {
		st, err := a.collector.Read()
		a.counter++
		if err != nil {
			panic(err)
		}
		time.Sleep(10)
		batch := convertSystemInfoToReport(st)
		batch = append(
			batch,
			reporter.Report{
				Type:  "counter",
				Name:  "PollCount",
				Value: strconv.FormatInt(a.counter, 10),
			},
			reporter.Report{
				Type:  "gauge",
				Name:  "RandomValue",
				Value: generateRandom(),
			},
		)
		a.reporter.ReportBatch(convertSystemInfoToReport(st))
	}
}

func generateRandom() string {
	return strconv.FormatFloat(rand.Float64(), 'f', 10, 64)
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
