package agent

import (
	"fmt"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestAgentWriter_WriteBatch(t *testing.T) {
	type args struct {
		data []metrics.Metrics
	}
	const endPoint = "localhost:8080"
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			"check ok",
			args{
				[]metrics.Metrics{
					{
						MType: "gauge",
						ID:    "someCounter",
						Value: pointer.From[float64](1.0),
					},
				},
			},
			[]error{},
		},
		{
			"check several",
			args{
				[]metrics.Metrics{
					{
						MType: "gauge",
						ID:    "someCounter",
						Value: pointer.From[float64](1.0),
					},
					{
						MType: "counter",
						ID:    "someCounter",
						Delta: pointer.From[int64](1),
					},
				},
			},
			[]error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			httpmock.ActivateNonDefault(client.GetClient())
			defer httpmock.DeactivateAndReset()
			resp, _ := httpmock.NewJsonResponder(200, metrics.Metrics{})
			httpmock.RegisterResponder(
				"POST",
				fmt.Sprintf("http://%s/update/", endPoint),
				resp,
			)
			httpmock.RegisterResponder(
				"POST",
				fmt.Sprintf("http://%s/updates/", endPoint),
				resp,
			)
			r := NewAgentReporter(client, endPoint, nil)
			err := r.WriteBatch(tt.args.data)
			assert.NoError(t, err)
		})
	}
}
