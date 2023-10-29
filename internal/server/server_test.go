package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers/rest"
	"github.com/ekubyshin/metrics_agent/internal/reporter"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/ekubyshin/metrics_agent/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gzipReader(t *testing.T) {
	in := types.Metrics{MType: "gauge", ID: "test", Value: utils.ToPointer[float64](1.0)}
	bSend, err := json.Marshal(in)
	require.NoError(t, err)
	compB, err := reporter.Compress(bSend)
	require.NoError(t, err)
	reader := bytes.NewReader(compB)
	request := httptest.NewRequest("POST", "/update/", reader)
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	router := chi.NewMux()
	stg := storage.NewMemoryStorage[string, types.Gauge]()
	stc := storage.NewMemoryStorage[string, types.Counter]()
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(stc, stg)
	router.Use(gzipReader)
	router.Post("/update/", m.Update)
	router.ServeHTTP(w, request)
	val, ok := stg.Get("test")
	assert.True(t, ok)
	assert.Equal(t, types.Gauge(1.0), val)
	res := w.Result()
	defer res.Body.Close()
}

func Test_gzipWriter(t *testing.T) {
	in := types.Metrics{MType: "gauge", ID: "test", Value: utils.ToPointer[float64](1.0)}
	bSend, err := json.MarshalIndent(in, "", "  ")
	require.NoError(t, err)
	reader := bytes.NewReader(bSend)
	request := httptest.NewRequest("POST", "/value/", reader)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	router := chi.NewMux()
	stg := storage.NewMemoryStorage[string, types.Gauge]()
	stc := storage.NewMemoryStorage[string, types.Counter]()
	stg.Put("test", types.Gauge(1.0))
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(stc, stg)
	router.Use(gzipHandle)
	router.Post("/value/", m.Value)
	router.ServeHTTP(w, request)
	val, ok := stg.Get("test")
	assert.True(t, ok)
	assert.Equal(t, types.Gauge(1.0), val)
	res := w.Result()
	defer res.Body.Close()
	compB, err := reporter.Compress(bSend)
	require.NoError(t, err)
	r, err := io.ReadAll(res.Body)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.NoError(t, err)
	assert.Equal(t, string(compB), string(r))
}
