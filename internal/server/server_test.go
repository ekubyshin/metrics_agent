package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ekubyshin/metrics_agent/internal/handlers"
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
	st := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
	st.Put(
		types.MetricsKey{ID: "test", MType: handlers.GaugeActionKey},
		types.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: utils.ToPointer[float64](1.0)})
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(st)
	router.Use(gzipReader)
	router.Post("/update/", m.Update)
	router.ServeHTTP(w, request)
	val, ok := st.Get(types.MetricsKey{ID: "test", MType: handlers.GaugeActionKey})
	assert.True(t, ok)
	assert.Equal(t, types.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: utils.ToPointer[float64](1.0)}, val)
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
	st := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
	st.Put(
		types.MetricsKey{ID: "test", MType: handlers.GaugeActionKey},
		types.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: utils.ToPointer[float64](1.0)})
	w := httptest.NewRecorder()
	m := rest.NewRestHandler(st)
	router.Use(gzipHandle)
	router.Post("/value/", m.Value)
	router.ServeHTTP(w, request)
	val, ok := st.Get(types.MetricsKey{ID: "test", MType: handlers.GaugeActionKey})
	assert.True(t, ok)
	assert.Equal(t, types.Metrics{ID: "test", MType: handlers.GaugeActionKey, Value: utils.ToPointer[float64](1.0)}, val)
	res := w.Result()
	defer res.Body.Close()
	compB, err := reporter.Compress(bSend)
	require.NoError(t, err)
	r, err := io.ReadAll(res.Body)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.NoError(t, err)
	assert.Equal(t, string(compB), string(r))
}

// func TestRestoreStorage(t *testing.T) {
// 	type args struct {
// 		filename string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			"ok",
// 			args{
// 				"./test/test.json",
// 			},
// 			false,
// 		},
// 		{
// 			"false",
// 			args{
// 				"./test/test2.json",
// 			},
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			st := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
// 			err := RestoreStorage(st, tt.args.filename)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			if !tt.wantErr {
// 				elems := st.List()
// 				assert.Equal(t, 2, len(elems))
// 			}
// 		})
// 	}
// }
