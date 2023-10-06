package reporter

import (
	"reflect"
	"testing"
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
