package reporter

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
	"github.com/go-resty/resty/v2"
)

const rootURL = "http://localhost:8080/update/{Type}/{Name}/{Value}"

type Writer interface {
	Write(data Report) error
	WriteBatch(data []Report) []error
}

type Report struct {
	Type  string
	Name  string
	Value string
}

type AgentWriter struct {
	reader   collector.Reader
	client   *resty.Client
	interval time.Duration
	queue    chan Report
}

func NewAgentReporter(interval time.Duration, client *resty.Client) Writer {
	return &AgentWriter{
		client:   client,
		interval: interval,
		queue:    make(chan Report, 100),
	}
}

func (r *AgentWriter) send(data Report) error {

	_, err := r.client.R().SetPathParams(reportToMap(data)).Get(rootURL)

	if err != nil {
		return err
	}

	return err
}

func (r *AgentWriter) Write(data Report) error {
	return r.send(data)
}

func reportToMap(data Report) map[string]string {
	res := make(map[string]string)
	v := reflect.ValueOf(data)
	for f := 0; f < v.NumField(); f++ {
		field := v.Field(f)
		fieldName := strings.ToLower(v.Type().Field(f).Name)
		res[fieldName] = strings.ToLower(field.String())
	}
	return res
}

func (r *AgentWriter) WriteBatch(data []Report) []error {
	var wg sync.WaitGroup
	resp := make([]error, 0, len(data))
	for _, v := range data {
		wg.Add(1)
		go func(v Report) {
			defer wg.Done()
			err := r.Write(v)
			if err != nil {
				resp = append(resp, err)
			}
		}(v)
	}
	wg.Wait()
	return resp
}