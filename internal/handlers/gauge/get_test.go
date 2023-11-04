package gauge

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugeGetHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route   string
		value   metrics.Gauge
		valName string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			"404",
			fields{
				route: "/gauge/testSetGet111",
			},
			want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    ``,
			},
		},
		{
			"test 200",
			fields{
				route:   "/gauge/testSetGet111",
				value:   1524.1,
				valName: "testSetGet111",
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    `1524.1`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tt.fields.route, nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
			mr := NewGaugeGetHandler(st)
			mw := NewGaugePostHandler(st)
			router.Get(mr.BaseURL(), mr.ServeHTTP)
			router.Post(mw.BaseURL(), mw.ServeHTTP)
			w := httptest.NewRecorder()
			if tt.fields.valName != "" {
				st.Put(
					metrics.MetricsKey{ID: tt.fields.valName, MType: handlers.GaugeActionKey},
					metrics.Metrics{ID: tt.fields.valName, MType: handlers.GaugeActionKey, Value: pointer.From[float64](float64(tt.fields.value))})
			}
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, tt.want.response, string(resBody))
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
