package gauge

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

func TestGaugeGetHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route   string
		method  string
		value   types.Gauge
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
			"test 404",
			fields{
				route:  "/none/",
				method: "GET",
			},
			want{
				code:        http.StatusNotFound,
				contentType: "",
				response:    ``,
			},
		},
		{
			"test 404",
			fields{
				route:  "/none/someCounter",
				method: "GET",
			},
			want{
				code:        http.StatusNotFound,
				contentType: "",
				response:    ``,
			},
		},
		{
			"test 404",
			fields{
				route:  "/gauge/someCounter",
				method: "POST",
			},
			want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
				response:    ``,
			},
		},
		{
			"test 200",
			fields{
				route:   "/gauge/testSetGet111",
				method:  "GET",
				value:   1524.0,
				valName: "testSetGet111",
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    `1524.0`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.fields.method, tt.fields.route, nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage()
			mr := NewGaugeGetHandler(st)
			router.Get(mr.BaseURL(), mr.ServeHTTP)
			w := httptest.NewRecorder()
			if tt.fields.valName != "" {
				st.Put(handlers.Key{Type: "gauge", Name: tt.fields.valName}, tt.fields.value)
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
