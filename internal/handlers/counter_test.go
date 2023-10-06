package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route  string
		method string
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
				method: "POST",
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
				route:  "/none/someCounter/1",
				method: "POST",
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
				route:  "/counter/someCounter/a",
				method: "POST",
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
				route:  "/counter/someCounter/1.0",
				method: "POST",
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
				route:  "/counter/someCounter/1",
				method: "POST",
			},
			want{
				code:        http.StatusOK,
				contentType: "application/json",
				response:    `1`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.fields.method, tt.fields.route, nil)
			router := chi.NewMux()
			m := NewCounterHandler()
			w := httptest.NewRecorder()
			router.Post(m.BaseURL(), m.ServeHTTP)
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
