package logger

import (
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func NewResponseLogger(l Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			h.ServeHTTP(&lw, r)
			l.Info(
				"Response",
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"size", responseData.size,
			)
		}
		return http.HandlerFunc(logFn)
	}
}

func NewRequestLogger(l Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			h.ServeHTTP(w, r)

			duration := time.Since(start)

			l.Info(
				"Request",
				"uri", r.RequestURI,
				"method", r.Method,
				"duration", duration,
			)
		}
		return http.HandlerFunc(logFn)
	}
}
