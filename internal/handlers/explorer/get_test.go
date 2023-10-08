package explorer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExplorerHandler_ServeHTTP(t *testing.T) {
	type args struct {
		Type  string
		Name  string
		Value any
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args []args
		want want
	}{
		{
			"empty",
			[]args{},
			want{
				contentType: "application/json",
				code:        http.StatusOK,
				response:    "{}",
			},
		},
		{
			"one element",
			[]args{
				{
					Type:  "gauge",
					Name:  "someMetric",
					Value: types.Gauge(1.0),
				},
			},
			want{
				contentType: "application/json",
				code:        http.StatusOK,
				response:    `{"someMetric":1.0}`,
			},
		},
		{
			"several elements",
			[]args{
				{
					Type:  "gauge",
					Name:  "someMetric",
					Value: types.Gauge(1.0),
				},
				{
					Type:  "gauge",
					Name:  "someMetric2",
					Value: types.Gauge(123.0),
				},
				{
					Type:  "counter",
					Name:  "someMetric3",
					Value: types.Counter(1),
				},
			},
			want{
				contentType: "application/json",
				code:        http.StatusOK,
				response: `{
					"someMetric":1.0,
					"someMetric2":123.0,
					"someMetric3":1
				}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[handlers.Key, any]()
			mr := NewExplorerHandler(st)
			w := httptest.NewRecorder()
			for _, v := range tt.args {
				st.Put(handlers.Key{Type: v.Type, Name: v.Name}, v.Value)
			}
			router.Get("/", mr.ServeHTTP)
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if res.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.want.response, string(resBody))
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
