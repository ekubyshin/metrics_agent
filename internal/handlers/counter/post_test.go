package counter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCounterHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route  string
		method string
		value  types.Counter
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
			"test 400",
			fields{
				route:  "/counter/someCounter/a",
				method: "POST",
			},
			want{
				code:        http.StatusBadRequest,
				contentType: "",
				response:    ``,
			},
		},
		{
			"test 404",
			fields{
				route:  "/counter/someCounter/1.0",
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
				route:  "/counter/someCounter/1",
				method: "POST",
				value:  1,
				key:    "someCounter",
			},
			want{
				code:        http.StatusOK,
				contentType: "",
				response:    `1`,
			},
		},
		{
			"test 200",
			fields{
				route:  "/counter/someCounter/1234",
				method: "POST",
				value:  1234,
				key:    "someCounter",
			},
			want{
				code:        http.StatusOK,
				contentType: "",
				response:    `1234`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.fields.method, tt.fields.route, nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[handlers.Key, any]()
			m := NewCounterPostHandler(st)
			w := httptest.NewRecorder()
			router.Post(m.BaseURL(), m.ServeHTTP)
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				v, err := st.Get(handlers.Key{Type: handlers.CounterActionKey, Name: tt.fields.key})
				assert.True(t, err)
				assert.Equal(t, tt.fields.value, v)
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
