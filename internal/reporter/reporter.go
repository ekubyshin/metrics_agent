package reporter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
)

const rootUrl = "http://127.0.0.1:8080/update"

type Writer interface {
	Report(data Report) error
	ReportBatch(data []Report) []error
}

type Report struct {
	Type  string
	Name  string
	Value string
}

type AgentWriter struct {
	reader   collector.Reader
	client   *http.Client
	interval time.Duration
}

func NewAgentReporter(interval time.Duration, client *http.Client) Writer {
	return &AgentWriter{
		client:   client,
		interval: interval,
	}
}

func (r *AgentWriter) Report(data Report) error {
	url := fmt.Sprintf("%s/%s/%s/%s", rootUrl, data.Type, data.Name, data.Value)
	request, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		return err
	}

	_, err = r.client.Do(request)
	return err
}

func (r *AgentWriter) ReportBatch(data []Report) []error {
	var wg sync.WaitGroup
	resp := make([]error, 0, len(data))
	for _, v := range data {
		wg.Add(1)
		go func(v Report) {
			defer wg.Done()
			err := r.Report(v)
			if err != nil {
				resp = append(resp, err)
			}
		}(v)
	}
	wg.Wait()
	return resp
}
