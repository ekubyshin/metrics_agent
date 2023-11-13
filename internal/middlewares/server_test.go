package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/agent"
	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/handlers/rest"
	"github.com/ekubyshin/metrics_agent/internal/metrics"
	"github.com/ekubyshin/metrics_agent/internal/pointer"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gzipReader(t *testing.T) {
	in := metrics.Metrics{MType: "gauge", ID: "test", Value: pointer.From[float64](1.0)}
	bSend, err := json.Marshal(in)
	require.NoError(t, err)
	compB, err := agent.Compress(bSend)
	require.NoError(t, err)
	reader := bytes.NewReader(compB)
	request := httptest.NewRequest("POST", "/update/", reader)
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	router := chi.NewMux()
	st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
	_ = st.Put(
		context.TODO(),
		metrics.MetricsKey{ID: "test", Type: handlers.GaugeActionKey},
		metrics.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: pointer.From[float64](1.0)})
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(st)
	router.Use(GzipReader)
	router.Post("/update/", m.Update)
	router.ServeHTTP(w, request)
	val, ok := st.Get(context.TODO(), metrics.MetricsKey{ID: "test", Type: handlers.GaugeActionKey})
	assert.True(t, ok)
	assert.Equal(t, metrics.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: pointer.From[float64](1.0)}, val)
	res := w.Result()
	defer res.Body.Close()
}

func Test_gzipWriter(t *testing.T) {
	in := metrics.Metrics{MType: "gauge", ID: "test", Value: pointer.From[float64](1.0)}
	bSend, err := json.MarshalIndent(in, "", "  ")
	require.NoError(t, err)
	reader := bytes.NewReader(bSend)
	request := httptest.NewRequest("POST", "/value/", reader)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	router := chi.NewMux()
	st := storage.NewMemoryStorage[metrics.MetricsKey, metrics.Metrics]()
	_ = st.Put(
		context.TODO(),
		metrics.MetricsKey{ID: "test", Type: handlers.GaugeActionKey},
		metrics.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: pointer.From[float64](1.0)})
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(st)
	router.Use(GzipHandler)
	router.Post("/value/", m.Value)
	router.ServeHTTP(w, request)
	val, ok := st.Get(context.TODO(), metrics.MetricsKey{ID: "test", Type: handlers.GaugeActionKey})
	assert.True(t, ok)
	assert.Equal(t, metrics.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: pointer.From[float64](1.0)}, val)
	res := w.Result()
	defer res.Body.Close()
	compB, err := agent.Compress(bSend)
	require.NoError(t, err)
	r, err := io.ReadAll(res.Body)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.NoError(t, err)
	assert.Equal(t, string(compB), string(r))
}
