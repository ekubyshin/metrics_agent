package reporter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
)

const rootUrl = "http://localhost:8080/update"

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
	client   *http.Client
	interval time.Duration
	queue    chan Report
}

func NewAgentReporter(interval time.Duration, client *http.Client) Writer {
	return &AgentWriter{
		client:   client,
		interval: interval,
		queue:    make(chan Report, 100),
	}
}

func (r *AgentWriter) send(data Report) error {
	url := fmt.Sprintf("%s/%s/%s/%s", rootUrl, data.Type, data.Name, data.Value)
	request, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		return err
	}

	_, err = r.client.Do(request)
	return err
}

func (r *AgentWriter) Write(data Report) error {
	return r.send(data)
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
