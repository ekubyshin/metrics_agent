package reporter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_reportToMap(t *testing.T) {
	type args struct {
		data Report
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"check full",
			args{
				Report{
					Type:  "gauge",
					Name:  "some",
					Value: "1.0",
				},
			},
			map[string]string{
				"type":  "gauge",
				"name":  "some",
				"value": "1.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reportToMap(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reportToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentWriter_WriteBatch(t *testing.T) {
	type args struct {
		data []Report
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
				[]Report{
					{
						Type:  "gauge",
						Name:  "someCounter",
						Value: "1.0",
					},
				},
			},
			[]error{},
		},
		{
			"check several",
			args{
				[]Report{
					{
						Type:  "gauge",
						Name:  "someCounter",
						Value: "1.0",
					},
					{
						Type:  "counter",
						Name:  "someCounter",
						Value: "1",
					},
					{
						Type:  "gauge",
						Name:  "someCounter2",
						Value: "2.0",
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
			for _, r := range tt.args.data {
				httpmock.RegisterResponder(
					"POST",
					fmt.Sprintf("http://%v/update/%v/%v/%v", endPoint, strings.ToLower(r.Type), strings.ToLower(r.Name), r.Value),
					httpmock.NewStringResponder(200, ``),
				)
			}
			r := NewAgentReporter(client, endPoint)
			errs := r.WriteBatch(tt.args.data)
			assert.Len(t, errs, 0)
			total := httpmock.GetTotalCallCount()
			assert.Equal(t, len(tt.args.data), total)
		})
	}
}
