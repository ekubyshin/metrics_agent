package handlers

// import (
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestMetrics_parsePath(t *testing.T) {
// 	type fields struct {
// 		route string
// 	}
// 	type args struct {
// 		url *url.URL
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *metricsHanlerPath
// 		wantErr bool
// 	}{
// 		{
// 			name: "without type",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/upload/",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "without name",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/update/type/",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "without value",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/update/type/name/",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "without value",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/update/counter/someMetric/527",
// 				},
// 			},
// 			want: &metricsHanlerPath{
// 				metricsType:  "counter",
// 				metricsName:  "someMetric",
// 				metricsValue: "527",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "without invalid val",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/update/counter/someMetric/none",
// 				},
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid type",
// 			fields: fields{
// 				route: "/update/",
// 			},
// 			args: args{
// 				url: &url.URL{
// 					Path: "/update/none/someMetric/none",
// 				},
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m := &Metrics{
// 				route: tt.fields.route,
// 			}
// 			got, err := m.parsePath(tt.args.url)
// 			if (err != nil) != tt.wantErr {
// 				assert.NotNil(t, err)
// 				return
// 			}
// 			assert.Equal(t, got, tt.want)
// 		})
// 	}
// }

// func TestMetrics_ServeHTTP(t *testing.T) {
// 	type fields struct {
// 		path   string
// 		method string
// 	}
// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   want
// 	}{
// 		{
// 			name: "400",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "",
// 				response:    ``,
// 			},
// 		},
// 		{
// 			name: "404",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/counter/",
// 			},
// 			want: want{
// 				code:        http.StatusNotFound,
// 				contentType: "",
// 				response:    ``,
// 			},
// 		},
// 		{
// 			name: "404",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/counter/name/",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "",
// 				response:    ``,
// 			},
// 		},
// 		{
// 			name: "200 counter",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/counter/name/1",
// 			},
// 			want: want{
// 				code:        http.StatusOK,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 		{
// 			name: "400 none",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/none/name/value/",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 		{
// 			name: "400 gauge none",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/gauge/name/none",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 		{
// 			name: "400 counter none",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/counter/name/none",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 		{
// 			name: "200 counter 1",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/counter/name/1",
// 			},
// 			want: want{
// 				code:        http.StatusOK,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 		{
// 			name: "200 gauge 1.0",
// 			fields: fields{
// 				method: http.MethodPost,
// 				path:   "/update/gauge/name/1.0",
// 			},
// 			want: want{
// 				code:        http.StatusOK,
// 				contentType: "application/json",
// 				response:    `{}`,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			request := httptest.NewRequest(tt.fields.method, tt.fields.path, nil)
// 			w := httptest.NewRecorder()
// 			m := NewMetricsHandler()
// 			m.ServeHTTP(w, request)
// 			res := w.Result()
// 			assert.Equal(t, tt.want.code, res.StatusCode)
// 			defer res.Body.Close()
// 			resBody, err := io.ReadAll(res.Body)
// 			require.NoError(t, err)
// 			if res.StatusCode == http.StatusOK {
// 				assert.JSONEq(t, tt.want.response, string(resBody))
// 				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
// 			}
// 		})
// 	}
// }
