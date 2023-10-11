package gauge

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// nolint
func TestGaugeHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route  string
		method string
		value  types.Gauge
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
			st := storage.NewMemoryStorage[string, types.Gauge]()
			m := NewGaugePostHandler(st)
			w := httptest.NewRecorder()
			router.Post(m.BaseURL(), m.ServeHTTP)
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
