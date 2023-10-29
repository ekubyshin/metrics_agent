package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/ekubyshin/metrics_agent/internal/config"
	"github.com/ekubyshin/metrics_agent/internal/handlers"
	"github.com/ekubyshin/metrics_agent/internal/handlers/counter"
	"github.com/ekubyshin/metrics_agent/internal/handlers/explorer"
	"github.com/ekubyshin/metrics_agent/internal/handlers/gauge"
	"github.com/ekubyshin/metrics_agent/internal/handlers/rest"
	l "github.com/ekubyshin/metrics_agent/internal/logger"
	"github.com/ekubyshin/metrics_agent/internal/storage"
	"github.com/ekubyshin/metrics_agent/internal/types"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	Run() error
}

type ChiServer struct {
	router   *chi.Mux
	endpoint config.Address
}

func NewServer(endpoint config.Address, logger l.Logger) *ChiServer {
	dbCounter := storage.NewMemoryStorage[string, types.Counter]()
	dbGauge := storage.NewMemoryStorage[string, types.Gauge]()
	router := chi.NewRouter()
	router.Use(l.NewRequestLogger(logger))
	router.Use(l.NewResponseLogger(logger))
	router.Use(gzipReader)
	router.Use(gzipHandle)
	registerRoutes(router, dbCounter, dbGauge)
	return &ChiServer{
		router:   router,
		endpoint: endpoint,
	}
}

func registerRoutes(
	router *chi.Mux,
	dbCounter storage.Storage[string, types.Counter],
	dbGauge storage.Storage[string, types.Gauge]) {
	gaugePostHandler := gauge.NewGaugePostHandler(dbGauge)
	counterPostHandler := counter.NewCounterPostHandler(dbCounter)
	gaugeGetHandler := gauge.NewGaugeGetHandler(dbGauge)
	counterGetHanlder := counter.NewCounterGetHandler(dbCounter)
	listHanlder := explorer.NewExplorerHandler(dbCounter, dbGauge)
	restHandler := rest.NewRestHandler(dbCounter, dbGauge)
	router.Get(listHanlder.BaseURL(), listHanlder.ServeHTTP)
	router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, handlers.ParamTypeKey)
		switch t {
		case handlers.GaugeActionKey:
			gaugePostHandler.ServeHTTP(w, r)
		case handlers.CounterActionKey:
			counterPostHandler.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})
	router.Post("/update/{type}/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(GetErrorStatusCode(r))
	})
	router.Route("/value", func(r chi.Router) {
		r.Get(gaugeGetHandler.BaseURL(), gaugeGetHandler.ServeHTTP)
		r.Get(counterGetHanlder.BaseURL(), counterGetHanlder.ServeHTTP)
	})
	router.Post("/update/", restHandler.Update)
	router.Post("/value/", restHandler.Value)
}

func (s *ChiServer) Run() error {
	return http.ListenAndServe(s.endpoint.ToString(), s.router)
}

func GetErrorStatusCode(r *http.Request) int {
	t := chi.URLParam(r, handlers.ParamTypeKey)
	switch t {
	case handlers.GaugeActionKey:
		return http.StatusNotFound
	case handlers.CounterActionKey:
		return http.StatusNotFound
	default:
		return http.StatusNotImplemented
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		oldBody := r.Body
		defer oldBody.Close()
		zr, err := gzip.NewReader(oldBody)
		if err != nil {
			io.WriteString(w, err.Error()) //nolint
			return
		}
		r.Body = zr
		next.ServeHTTP(w, r)
	})
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error()) //nolint
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
