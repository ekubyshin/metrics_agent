package reporter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"sync"

	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-resty/resty/v2"
)

const (
	path            = "/update/"
	contentEncoding = "Content-Encoding"
)

type Writer interface {
	Write(data types.Metrics) error
	WriteBatch(data []types.Metrics) []error
}

type Report struct {
	Type  string
	Name  string
	Value string
}

type AgentWriter struct {
	client   *resty.Client
	queue    chan Report
	endpoint string
}

func NewAgentReporter(client *resty.Client, endpoint string) *AgentWriter {
	client.Header.Add("Content-Type", "application/json")
	return &AgentWriter{
		client:   client,
		endpoint: endpoint,
		queue:    make(chan Report, 100),
	}
}

func (r *AgentWriter) send(data types.Metrics) error {
	bSend, err := json.Marshal(data)
	if err != nil {
		return err
	}
	compB, err := Compress(bSend)
	if err == nil {
		bSend = compB
	}
	req := r.client.R().SetBody(bSend)
	if err == nil {
		req.SetHeader(contentEncoding, "gzip")
	}
	_, err = req.Post("http://" + r.endpoint + path)
	if err != nil {
		return err
	}
	return err
}

func (r *AgentWriter) Write(data types.Metrics) error {
	return r.send(data)
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	zw, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = zw.Write(data)
	if err != nil {
		return nil, err
	}
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (r *AgentWriter) WriteBatch(data []types.Metrics) []error {
	var wg sync.WaitGroup
	resp := make([]error, 0, len(data))
	for _, v := range data {
		wg.Add(1)
		go func(v types.Metrics) {
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
