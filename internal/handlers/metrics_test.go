package handlers

import (
	"net/url"
	"reflect"
	"testing"
)

func TestMetrics_parsePath(t *testing.T) {
	type fields struct {
		route string
	}
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metricsHanlerPath
		wantErr bool
	}{
		{
			name: "without type",
			fields: fields{
				route: "/update/",
			},
			args: args{
				url: &url.URL{
					Path: "/upload/",
				},
			},
			wantErr: true,
		},
		{
			name: "without name",
			fields: fields{
				route: "/update/",
			},
			args: args{
				url: &url.URL{
					Path: "/update/type/",
				},
			},
			wantErr: true,
		},
		{
			name: "without value",
			fields: fields{
				route: "/update/",
			},
			args: args{
				url: &url.URL{
					Path: "/update/type/name/",
				},
			},
			wantErr: true,
		},
		{
			name: "without value",
			fields: fields{
				route: "/update/",
			},
			args: args{
				url: &url.URL{
					Path: "/upload/counter/someMetric/527",
				},
			},
			want: &metricsHanlerPath{
				metricsType:  "counter",
				metricsName:  "someMetric",
				metricsValue: "527",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				route: tt.fields.route,
			}
			got, err := m.parsePath(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metrics.parsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metrics.parsePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
