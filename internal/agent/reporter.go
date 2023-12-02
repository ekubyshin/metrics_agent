package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/crypto"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/go-resty/resty/v2"
)

const (
	contentEncoding = "Content-Encoding"
)

type Writer interface {
	WriteBatch([]metrics.Metrics) error
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
	secret   *string
}

func NewAgentReporter(client *resty.Client, endpoint string, secret *string) *AgentWriter {
	client.Header.Add("Content-Type", "application/json")
	client.Header.Add("Accept", "application/json")
	client.SetTimeout(1 * time.Second)
	return &AgentWriter{
		client:   client,
		endpoint: endpoint,
		queue:    make(chan Report, 100),
		secret:   secret,
	}
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

func (r *AgentWriter) WriteBatch(data []metrics.Metrics) error {
	bSend, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req := r.client.R()
	hash, err := r.hashData(bSend)
	if err == nil {
		req = req.SetHeader(crypto.HashHeader, string(hash))
	}
	compB, err := Compress(bSend)
	if err == nil {
		bSend = compB
	}
	req = req.SetBody(bSend)
	if err == nil {
		req = req.SetHeader(contentEncoding, "gzip")
	}

	return r.send(req, 1)
}

func (r *AgentWriter) hashData(d []byte) ([]byte, error) {
	if r.secret == nil {
		return nil, errors.New("empty key")
	}
	return crypto.HashData(d, *r.secret)
}

func (r *AgentWriter) send(req *resty.Request, i int64) error {
	_, err := req.Post(fmt.Sprintf("http://%s/updates/", r.endpoint))
	if err != nil && errors.Is(err, syscall.ECONNREFUSED) && i <= 5 {
		time.Sleep(time.Duration(i) * time.Second)
		return r.send(req, i+2)
	}

	return err
}
