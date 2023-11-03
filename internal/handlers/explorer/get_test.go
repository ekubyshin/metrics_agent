package explorer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/ekubyshin/metrics_agent/internal/utils"
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
				contentType: "text/html",
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
				contentType: "text/html",
				code:        http.StatusOK,
				response:    `{"gauge_someMetric":1}`,
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
				contentType: "text/html",
				code:        http.StatusOK,
				response: `{
					"gauge_someMetric":1,
					"gauge_someMetric2":123,
					"counter_someMetric3":1
				}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
			mr := NewExplorerHandler(st)
			w := httptest.NewRecorder()
			for _, v := range tt.args {
				if v.Type == handlers.GaugeActionKey {
					if val, ok := v.Value.(types.Gauge); ok {
						st.Put(types.MetricsKey{ID: v.Name, MType: v.Type}, types.Metrics{ID: v.Name, MType: v.Type, Value: utils.ToPointer[float64](float64(val))})
					}
				} else {
					if val, ok := v.Value.(types.Counter); ok {
						st.Put(types.MetricsKey{ID: v.Name, MType: v.Type}, types.Metrics{ID: v.Name, MType: v.Type, Delta: utils.ToPointer[int64](int64(val))})
					}
				}
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
