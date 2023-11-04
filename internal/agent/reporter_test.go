package agent

import (
	"fmt"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestAgentWriter_WriteBatch(t *testing.T) {
	type args struct {
		data []types.Metrics
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
				[]types.Metrics{
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
				[]types.Metrics{
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
			resp, _ := httpmock.NewJsonResponder(200, types.Metrics{})
			httpmock.RegisterResponder(
				"POST",
				fmt.Sprintf("http://%v/update/", endPoint),
				resp,
			)
			r := NewAgentReporter(client, endPoint)
			errs := r.WriteBatch(tt.args.data)
			assert.Len(t, errs, 0)
			total := httpmock.GetTotalCallCount()
			assert.Equal(t, len(tt.args.data), total)
		})
	}
}
