package counter

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterGetHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route   string
		method  string
		value   metrics.Counter
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
			"test 200",
			fields{
				route:   "/counter/testSetGet110",
				method:  "GET",
				value:   1,
				valName: "testSetGet110",
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    `1`,
			},
		},
		{
			"test 200",
			fields{
				route:   "/counter/testSetGet111",
				method:  "GET",
				value:   1524,
				valName: "testSetGet111",
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				response:    `1524`,
			},
		},
		{
			"test 404",
			fields{
				route:   "/counter/testSetGet111",
				method:  "GET",
				value:   1524,
				valName: "testSetGet112",
			},
			want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.fields.method, tt.fields.route, nil)
			router := chi.NewMux()
			st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
			mr := NewCounterHandler(st)
			w := httptest.NewRecorder()
			router.Get(GetURL, mr.Get)
			router.Post(PostURL, mr.Post)
			if tt.fields.valName != "" {
				request := httptest.NewRequest("POST", "/counter/"+tt.fields.valName+"/"+strconv.Itoa(int(tt.fields.value)), nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, request)
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

func TestCounterHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route  string
		method string
		value  int64
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
			"string instead of number",
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
			"float instead of int",
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
			st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
			m := NewCounterHandler(st)
			w := httptest.NewRecorder()
			router.Post(PostURL, m.Post)
			router.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				v, err := st.Get(context.TODO(), metrics.MetricsKey{ID: tt.fields.key, Type: handlers.CounterActionKey})
				assert.True(t, err)
				assert.Equal(t, tt.fields.value, *v.Delta)
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
