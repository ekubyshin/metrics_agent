package reporter

import (
	"net/http"
	"testing"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/collector"
)

func TestAgentReporter_Report(t *testing.T) {
	type fields struct {
		reader   collector.Reader
		client   *http.Client
		interval time.Duration
	}
	type args struct {
		data Report
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AgentWriter{
				reader:   tt.fields.reader,
				client:   tt.fields.client,
				interval: tt.fields.interval,
			}
			if err := r.Report(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AgentReporter.Report() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
