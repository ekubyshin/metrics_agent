package gauge

import (
	"context"
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
			mr := NewGaugeHandler(st)
			router.Get(GetURL, mr.Get)
			router.Post(PostURL, mr.Post)
			w := httptest.NewRecorder()
			if tt.fields.valName != "" {
				_, _ = st.Put(
					context.TODO(),
					metrics.MetricsKey{ID: tt.fields.valName, Type: handlers.GaugeActionKey},
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

// nolint
func TestGaugeHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route  string
		method string
		value  metrics.Gauge
		key    string
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
			"test 404",
			fields{
				route:  "/gauge/someCounter/a",
				method: "POST",
			},
			want{
				code:        http.StatusBadRequest,
				contentType: "",
				response:    ``,
			},
		},
		{
			"test 200",
			fields{
				route:  "/gauge/testSetGet185/117067.144",
				method: "POST",
				value:  117067.144,
				key:    "testSetGet185",
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain",
				response:    `117067.144`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.fields.method, tt.fields.route, nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
			m := NewGaugeHandler(st)
			w := httptest.NewRecorder()
			router.Post(PostURL, m.Post)
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
